/*
Copyright 2018 Platform 9 Systems, Inc.
*/

package machine

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/platform9/ssh-provider/pkg/controller"
	"github.com/platform9/ssh-provider/pkg/machine"

	"k8s.io/client-go/kubernetes"

	sshproviderv1 "github.com/platform9/ssh-provider/pkg/apis/sshprovider/v1alpha1"
	spclient "github.com/platform9/ssh-provider/pkg/client/clientset_generated/clientset"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	clusterv1 "sigs.k8s.io/cluster-api/pkg/apis/cluster/v1alpha1"
	clusterclient "sigs.k8s.io/cluster-api/pkg/client/clientset_generated/clientset"
	clusterutil "sigs.k8s.io/cluster-api/pkg/util"
)

const (
	EtcdadmPath       = "/opt/bin/etcdadm"
	NodeadmPath       = "/opt/bin/nodeadm"
	NodeadmConfigPath = "/etc/nodeadm.yaml"
	CachePath         = "/var/cache/ssh-provider"
)

type machineClientBuilder func(host string, port int, username string, privateKey string, publicKeys []string, insecureIgnoreHostKey bool) (machine.Client, error)

type Actuator struct {
	InsecureIgnoreHostKey bool

	kubeClient           kubernetes.Interface
	clusterClient        clusterclient.Interface
	spClient             spclient.Interface
	machineClientBuilder machineClientBuilder
}

func NewActuator(kubeClient kubernetes.Interface, clusterClient clusterclient.Interface, spClient spclient.Interface, machineClientBuilder machineClientBuilder, insecureIgnoreHostKey bool) *Actuator {
	return &Actuator{
		InsecureIgnoreHostKey: insecureIgnoreHostKey,

		kubeClient:           kubeClient,
		clusterClient:        clusterClient,
		spClient:             spClient,
		machineClientBuilder: machineClientBuilder,
	}
}

//version could include commit if shas are different for last tag and last commit.
//given version vMajor.Minor.Patch[.Commit] getScrubbedVersion returns vMajor.Minor.Patch
func getScrubbedVersion(version string) (string, error) {
	split := strings.Split(version, ".")
	if len(split) < 3 {
		return "", fmt.Errorf("unable to parse version of %s", version)
	}
	return fmt.Sprintf("%s.%s.%s", split[0], split[1], split[2]), nil
}

func ensureComponentVersion(componentInstallPath, expectedVersion, componentName string, machineClient machine.Client) error {
	exists, err := machineClient.Exists(componentInstallPath)
	if err != nil {
		return fmt.Errorf("unable to check if %s already exists: %v", componentInstallPath, err)
	}
	if exists {
		existingVersionBytes, _, err := machineClient.RunCommand(fmt.Sprintf("%s version --short", componentInstallPath))
		if err != nil {
			return fmt.Errorf("unable to check version of %s: %v", componentInstallPath, err)
		}
		existingVersion, err := getScrubbedVersion(strings.TrimSpace(string(existingVersionBytes)))
		if err != nil {
			return fmt.Errorf("unable to get scrubbed version: %v", err)
		}
		fmt.Printf("doing version check for %s existing version %s expected version %s\n", componentName, existingVersion, expectedVersion)
		if existingVersion == expectedVersion {
			fmt.Printf("found expected verision for %s\n", componentName)
			return nil
		}
	}
	fmt.Printf("looking for expected version %s for %s in cache\n", expectedVersion, componentName)
	componentCachePath := filepath.Join(CachePath, componentName, expectedVersion, componentName)
	exists, err = machineClient.Exists(componentCachePath)
	if err != nil {
		return fmt.Errorf("unable to check if %s already exists: %v", componentCachePath, err)
	}
	if exists {
		fmt.Printf("installing %s binary from %s\n", componentName, componentCachePath)
		machineClient.CopyFile(componentCachePath, componentInstallPath)
		return nil
	} else {
		//TODO Try download from a hosted location
	}
	return fmt.Errorf("unable to copy component binary from %s to %s, source exists %t", componentCachePath, componentInstallPath, exists)
}

func ensureComponentsVersion(machineSpec *sshproviderv1.MachineSpec, machineClient machine.Client) error {
	if err := ensureComponentVersion(EtcdadmPath, machineSpec.ComponentVersions.EtcdadmVersion, "etcdadm", machineClient); err != nil {
		return err 
	}
	return ensureComponentVersion(NodeadmPath, machineSpec.ComponentVersions.NodeadmVersion, "nodeadm", machineClient)
}

func (a *Actuator) Create(cluster *clusterv1.Cluster, machine *clusterv1.Machine) error {
	machineSpec, err := controller.GetMachineSpec(*machine)
	if err != nil {
		return fmt.Errorf("unable to decode spec of machine %q: %v", machine.Name, err)
	}
	pm, err := a.spClient.SshproviderV1alpha1().ProvisionedMachines(machine.Namespace).Get(machineSpec.ProvisionedMachineName, metav1.GetOptions{})
	if err != nil {
		return fmt.Errorf("unable to get provider machine %q bound to machine %q: %v", machineSpec.ProvisionedMachineName, machine.Name, err)
	}
	// TODO(dlipovetsky) validate machine-provisioned machine binding
	credentialSecret, err := a.kubeClient.CoreV1().Secrets(machine.Namespace).Get(pm.Spec.SSHConfig.CredentialSecret.Name, metav1.GetOptions{})
	if err != nil {
		return fmt.Errorf("unable to get ssh credential of provisioned machine %q bound to machine %q: %v", pm.Name, machine.Name, err)
	}
	username, privateKey, err := controller.UsernameAndKeyFromSecret(credentialSecret)
	if err != nil {
		return fmt.Errorf("unable to get ssh credential for machine %q: %v", machine.Name, err)
	}
	machineClient, err := a.machineClientBuilder(
		pm.Spec.SSHConfig.Host,
		pm.Spec.SSHConfig.Port,
		username,
		privateKey,
		pm.Spec.SSHConfig.PublicKeys,
		a.InsecureIgnoreHostKey,
	)
	if err != nil {
		return fmt.Errorf("error creating client for machine %q: %s", machine.Name, err)
	}
	if err := ensureComponentsVersion(machineSpec, machineClient); err != nil {
		return fmt.Errorf("enable to ensure correct versions of tools: %v", err)
	}
	if clusterutil.IsMaster(machine) {
		if err := a.createMaster(cluster, machine, pm, machineClient); err != nil {
			return fmt.Errorf("error creating machine %q: %s", machine.Name, err)
		}
	} else {
		if err := a.createNode(cluster, machine, machineClient); err != nil {
			return fmt.Errorf("error creating machine %q: %s", machine.Name, err)
		}
	}

	if _, err := a.clusterClient.ClusterV1alpha1().Machines(machine.Namespace).UpdateStatus(machine); err != nil {
		return fmt.Errorf("error updating machine %q: %v", machine.Name, err)
	}
	return nil
}

func (a *Actuator) Delete(cluster *clusterv1.Cluster, machine *clusterv1.Machine) error {
	machineSpec, err := controller.GetMachineSpec(*machine)
	if err != nil {
		return fmt.Errorf("unable to decode spec of machine %q: %v", machine.Name, err)
	}
	pm, err := a.spClient.SshproviderV1alpha1().ProvisionedMachines(machine.Namespace).Get(machineSpec.ProvisionedMachineName, metav1.GetOptions{})
	if err != nil {
		return fmt.Errorf("unable to get provider machine %q bound to machine %q: %v", machineSpec.ProvisionedMachineName, machine.Name, err)
	}
	// TODO(dlipovetsky) validate machine-provisioned machine binding
	credentialSecret, err := a.kubeClient.CoreV1().Secrets(machine.Namespace).Get(pm.Spec.SSHConfig.CredentialSecret.Name, metav1.GetOptions{})
	if err != nil {
		return fmt.Errorf("unable to get ssh credential of provisioned machine %q bound to machine %q: %v", pm.Name, machine.Name, err)
	}
	username, privateKey, err := controller.UsernameAndKeyFromSecret(credentialSecret)
	if err != nil {
		return fmt.Errorf("unable to get ssh credential for machine %q: %v", machine.Name, err)
	}
	machineClient, err := a.machineClientBuilder(
		pm.Spec.SSHConfig.Host,
		pm.Spec.SSHConfig.Port,
		username,
		privateKey,
		pm.Spec.SSHConfig.PublicKeys,
		a.InsecureIgnoreHostKey,
	)
	if err != nil {
		return fmt.Errorf("error creating client for machine %q: %s", machine.Name, err)
	}
	if err := ensureComponentsVersion(machineSpec, machineClient); err != nil {
		return fmt.Errorf("enable to ensure correct versions of tools: %v", err)
	}
	if clusterutil.IsMaster(machine) {
		if err := a.deleteMaster(machine, machineClient); err != nil {
			return fmt.Errorf("error deleting machine %q: %s", machine.Name, err)
		}
	} else {
		if err := a.deleteNode(machine, machineClient); err != nil {
			return fmt.Errorf("error deleting machine %q: %s", machine.Name, err)
		}
	}

	if _, err := a.clusterClient.ClusterV1alpha1().Machines(machine.Namespace).UpdateStatus(machine); err != nil {
		return fmt.Errorf("error updating machine %q: %v", machine.Name, err)
	}
	return nil
}

func (a *Actuator) Update(cluster *clusterv1.Cluster, machine *clusterv1.Machine) error {
	return nil
}

func (a *Actuator) Exists(cluster *clusterv1.Cluster, machine *clusterv1.Machine) (bool, error) {
	return false, nil
}
