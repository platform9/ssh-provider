/*
Copyright 2018 Platform 9 Systems, Inc.
*/

package machine

import (
	"fmt"
	"log"
	"net"

	"github.com/platform9/ssh-provider/provisionedmachine"

	"golang.org/x/crypto/ssh"

	sshconfigv1 "github.com/platform9/ssh-provider/sshproviderconfig/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	clusterv1 "sigs.k8s.io/cluster-api/pkg/apis/cluster/v1alpha1"
	clusterutil "sigs.k8s.io/cluster-api/pkg/util"
)

type SSHActuator struct {
	InsecureIgnoreHostKey bool
	sshProviderCodec      *sshconfigv1.SSHProviderCodec

	provisionedMachineConfigMaps []*corev1.ConfigMap
	sshCredentials               *corev1.Secret
	clusterCA                    *corev1.Secret
}

func NewActuator(provisionedMachineConfigMaps []*corev1.ConfigMap, sshCredentials *corev1.Secret, clusterCA *corev1.Secret) (*SSHActuator, error) {
	codec, err := sshconfigv1.NewCodec()
	if err != nil {
		return nil, err
	}
	return &SSHActuator{
		sshProviderCodec:             codec,
		provisionedMachineConfigMaps: provisionedMachineConfigMaps,
		sshCredentials:               sshCredentials,
		clusterCA:                    clusterCA,
	}, nil
}

func (sa *SSHActuator) Create(cluster *clusterv1.Cluster, machine *clusterv1.Machine) error {
	// caCert, ok := clusterCA.Data["ca.crt"]
	// if !ok {
	// 	return fmt.Errorf("error reading cluster CA certificate: %s", ok)
	// }
	// caKey, ok := clusterCA.Data["ca.key"]
	// if !ok {
	// 	return fmt.Errorf("error reading cluster CA private key: %s", ok)
	// }

	cm, err := sa.selectProvisionedMachine(machine)
	if err != nil {
		return fmt.Errorf("error finding a compatible ProvisionedMachine for Machine %q: %s", machine.Name, err)
	}
	// err = sa.linkProvisionedMachineWithMachine(cm, machine)
	// if err != nil {
	// 	return fmt.Errorf("error linking ProvisionedMachine ConfigMap %q and Machine %q: %s", cm.Name, machine.Name, err)
	// }

	client, err := sshClient(cm, sa.sshCredentials, sa.InsecureIgnoreHostKey)
	if err != nil {
		return fmt.Errorf("error creating SSH client for machine %q: %s", machine.Name, err)
	}
	defer client.Close()

	if clusterutil.IsMaster(machine) {
		return sa.createMaster(cluster, machine, client)
	} else {
		return sa.createNode(cluster, machine, client)
	}
}

// TODO(dlipovetsky) Find a compatible ProvisionedMachine, or return an error
func (sa *SSHActuator) selectProvisionedMachine(machine *clusterv1.Machine) (*corev1.ConfigMap, error) {
	return sa.provisionedMachineConfigMaps[0], nil
}

// TODO(dlipovetsky) Persist changes
func (sa *SSHActuator) linkProvisionedMachineWithMachine(cm *corev1.ConfigMap, machine *clusterv1.Machine) error {
	pm, err := provisionedmachine.NewFromConfigMap(cm)
	if err != nil {
		return fmt.Errorf("error parsing ProvisionedMachine from ConfigMap %q: %s", cm.Name, err)
	}
	// Update ProvisionedMachine annotations
	cm.Annotations["sshprovider.platform9.com/machine-name"] = machine.Name
	// Update Machine annotations
	machine.Annotations["sshprovider.platform9.com/provisionedmachine-name"] = cm.Name
	// Update Machine.Status.ProviderStatus
	sshProviderStatus := &sshconfigv1.SSHMachineProviderStatus{
		SSHConfig: pm.SSHConfig,
	}
	if providerStatus, err := sa.sshProviderCodec.EncodeToProviderStatus(sshProviderStatus); err != nil {
		return fmt.Errorf("error creating machine ProviderStatus: %s", err)
	} else {
		machine.Status.ProviderStatus = *providerStatus
	}
	return nil
}

func (sa *SSHActuator) createMaster(cluster *clusterv1.Cluster, machine *clusterv1.Machine, client *ssh.Client) error {
	session, err := client.NewSession()
	if err != nil {
		log.Fatalf("error creating new SSH session for machine %q: %s", machine.Name, err)
	}
	out, err := session.CombinedOutput("echo writing ca cert and key")
	if err != nil {
		return fmt.Errorf("error invoking ssh command %s", err)
	} else {
		log.Println(string(out))
	}

	session, err = client.NewSession()
	if err != nil {
		log.Fatalf("error creating new SSH session for machine %q: %s", machine.Name, err)
	}
	out, err = session.CombinedOutput("echo running nodeadm init")
	if err != nil {
		return fmt.Errorf("error invoking ssh command %s", err)
	} else {
		log.Println(string(out))
	}

	session, err = client.NewSession()
	if err != nil {
		log.Fatalf("error creating new SSH session for machine %q: %s", machine.Name, err)
	}
	out, err = session.CombinedOutput("echo running etcdadm init")
	if err != nil {
		return fmt.Errorf("error invoking ssh command %s", err)
	} else {
		log.Println(string(out))
	}

	// TODO(dlipovetsky) Update cluster CA Secret with actual CA

	return nil
}

func (sa *SSHActuator) createNode(cluster *clusterv1.Cluster, machine *clusterv1.Machine, client *ssh.Client) error {
	session, err := client.NewSession()
	if err != nil {
		log.Fatalf("error creating new SSH session for machine %q: %s", machine.Name, err)
	}
	out, err := session.CombinedOutput("echo running nodeadm join")
	if err != nil {
		return fmt.Errorf("error invoking ssh command %s", err)
	} else {
		log.Println(out)
	}
	return nil
}

func (sa *SSHActuator) Delete(cluster *clusterv1.Cluster, machine *clusterv1.Machine) error {
	return nil
}

func (sa *SSHActuator) Update(cluster *clusterv1.Cluster, machine *clusterv1.Machine) error {
	return nil
}

func (sa *SSHActuator) Exists(cluster *clusterv1.Cluster, machine *clusterv1.Machine) (bool, error) {
	return false, nil
}

func (sa *SSHActuator) machineproviderconfig(providerConfig clusterv1.ProviderConfig) (*sshconfigv1.SSHMachineProviderConfig, error) {
	var config sshconfigv1.SSHMachineProviderConfig
	err := sa.sshProviderCodec.DecodeFromProviderConfig(providerConfig, &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}

func (sa *SSHActuator) clusterproviderconfig(providerConfig clusterv1.ProviderConfig) (*sshconfigv1.SSHClusterProviderConfig, error) {
	var config sshconfigv1.SSHClusterProviderConfig
	err := sa.sshProviderCodec.DecodeFromProviderConfig(providerConfig, &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}

// FixedHostKeys is a version of ssh.FixedHostKey that checks a list of SSH public keys
func FixedHostKeys(keys []ssh.PublicKey) ssh.HostKeyCallback {
	callbacks := make([]ssh.HostKeyCallback, len(keys))
	for i, expectedKey := range keys {
		callbacks[i] = ssh.FixedHostKey(expectedKey)
	}

	return func(hostname string, remote net.Addr, actualKey ssh.PublicKey) error {
		for _, callback := range callbacks {
			err := callback(hostname, remote, actualKey)
			if err == nil {
				return nil
			}
		}
		return fmt.Errorf("host key does not match any expected keys")
	}
}
