package util

import (
	sshconfigv1 "github.com/platform9/ssh-provider/sshproviderconfig/v1alpha1"
	clusterv1 "sigs.k8s.io/cluster-api/pkg/apis/cluster/v1alpha1"
)

func MachineProviderConfig(providerConfig clusterv1.ProviderConfig) (*sshconfigv1.SSHMachineProviderConfig, error) {
	codec, err := sshconfigv1.NewCodec()
	if err != nil {
		return nil, err
	}

	var config sshconfigv1.SSHMachineProviderConfig
	err = codec.DecodeFromProviderConfig(providerConfig, &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}

func ClusterProviderConfig(providerConfig clusterv1.ProviderConfig) (*sshconfigv1.SSHClusterProviderConfig, error) {
	codec, err := sshconfigv1.NewCodec()
	if err != nil {
		return nil, err
	}

	var config sshconfigv1.SSHClusterProviderConfig
	err = codec.DecodeFromProviderConfig(providerConfig, &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}
