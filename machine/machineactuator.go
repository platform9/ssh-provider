/*
Copyright 2018 Platform 9 Systems, Inc.
*/

package machine

import (
	"fmt"
	"io/ioutil"
	"log"

	"golang.org/x/crypto/ssh"

	sshconfigv1 "github.com/platform9/ssh-provider/sshproviderconfig/v1alpha1"
	clusterv1 "sigs.k8s.io/cluster-api/pkg/apis/cluster/v1alpha1"
)

type SSHActuator struct {
	InsecureIgnoreHostKey  bool
	sshProviderConfigCodec *sshconfigv1.SSHProviderConfigCodec
}

func (sa *SSHActuator) Create(cluster *clusterv1.Cluster, machine *clusterv1.Machine) error {
	machineConfig, err := sa.machineproviderconfig(machine.Spec.ProviderConfig)
	if err != nil {
		return err
	}

	// get username and ssh private key from Secret "sshcreds-machine-name"
	username := "root"
	key, err := ioutil.ReadFile("/home/daniel/.ssh/pf9_dev")
	if err != nil {
		log.Fatalf("unable to read private key: %v", err)
	}
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		log.Fatalf("unable to parse private key: %v", err)
	}

	// get host address and ssh fingerprint from SSHMachineProviderConfig
	sshConfig := &ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	connection, err := ssh.Dial("tcp", fmt.Sprintf("%s:%s", machineConfig.Host, machineConfig.Port), sshConfig)
	defer connection.Close()
	if err != nil {
		log.Fatalf("unable to dial: %s", err)
	}

	session, err := connection.NewSession()
	if err != nil {
		log.Fatalf("unable to create session: %s", err)
	}
	out, err := session.CombinedOutput("ls -al")
	if err != nil {
		log.Fatalf("unable to run ls -al")
	}
	fmt.Println(string(out))

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
	err := sa.sshProviderConfigCodec.DecodeFromProviderConfig(providerConfig, &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}
