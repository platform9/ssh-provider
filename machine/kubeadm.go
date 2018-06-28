package machine

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/Jeffail/gabs"
	"github.com/ghodss/yaml"
	sshconfigv1 "github.com/platform9/ssh-provider/sshproviderconfig/v1alpha1"
	kubeadmv1 "k8s.io/kubernetes/cmd/kubeadm/app/apis/kubeadm/v1alpha1"
	kubeletconfigv1alpha1 "k8s.io/kubernetes/pkg/kubelet/apis/kubeletconfig/v1alpha1"
	kubeproxyconfigv1alpha1 "k8s.io/kubernetes/pkg/proxy/apis/kubeproxyconfig/v1alpha1"
	clusterv1 "sigs.k8s.io/cluster-api/pkg/apis/cluster/v1alpha1"
)

func (sa *SSHActuator) NewKubeletConfiguration(cfg *sshconfigv1.SSHMachineProviderConfig) *kubeletconfigv1alpha1.KubeletConfiguration {
	kubeletConfiguration := cfg.KubeletConfiguration.DeepCopy()
	return kubeletConfiguration
}

func (sa *SSHActuator) NewKubeProxyConfiguration(cfg *sshconfigv1.SSHMachineProviderConfig) *kubeproxyconfigv1alpha1.KubeProxyConfiguration {
	kubeproxyConfiguration := cfg.KubeProxyConfiguration.DeepCopy()
	return kubeproxyConfiguration
}

func (sa *SSHActuator) NewMasterConfiguration(cluster *clusterv1.Cluster, machine *clusterv1.Machine) (*kubeadmv1.MasterConfiguration, error) {
	sshMachineProviderConfig, err := sa.machineproviderconfig(machine.Spec.ProviderConfig)
	if err != nil {
		return nil, fmt.Errorf("error creating MasterConfiguration: %s", err)
	}

	masterConfiguration := &kubeadmv1.MasterConfiguration{}

	masterConfiguration.KubernetesVersion = machine.Spec.Versions.ControlPlane

	masterConfiguration.Etcd.Endpoints = []string{"https://127.0.0.1:2379"}
	masterConfiguration.Etcd.CAFile = "/etc/etcd/pki/ca.crt"
	masterConfiguration.Etcd.CertFile = "/etc/etcd/pki/apiserver-etcd-client.crt"
	masterConfiguration.Etcd.KeyFile = "/etc/etcd/pki/apiserver-etcd-client.key"

	if sshMachineProviderConfig.KubeletConfiguration != nil {
		masterConfiguration.KubeletConfiguration.BaseConfig = sa.NewKubeletConfiguration(sshMachineProviderConfig)
	}
	if sshMachineProviderConfig.KubeProxyConfiguration != nil {
		masterConfiguration.KubeProxy.Config = sa.NewKubeProxyConfiguration(sshMachineProviderConfig)
	}

	switch len(cluster.Spec.ClusterNetwork.Pods.CIDRBlocks) {
	case 0:
		// Do nothing
	case 1:
		masterConfiguration.Networking.PodSubnet = cluster.Spec.ClusterNetwork.Pods.CIDRBlocks[0]
	case 2:
		return nil, fmt.Errorf("error creating MasterConfiguration: cluster %q spec.clusterNetwork.pods.cidrBlocks must contain at most one block", cluster.Name)
	}

	switch len(cluster.Spec.ClusterNetwork.Services.CIDRBlocks) {
	case 0:
		// Do nothing
	case 1:
		masterConfiguration.Networking.ServiceSubnet = cluster.Spec.ClusterNetwork.Services.CIDRBlocks[0]
	case 2:
		return nil, fmt.Errorf("error creating MasterConfiguration: cluster %q spec.clusterNetwork.pods.cidrBlocks must contain at most one block", cluster.Name)
	}

	kubeadmv1.SetDefaults_MasterConfiguration(masterConfiguration)
	return masterConfiguration, nil
}

func MarshalToYAMLWithFixedKubeProxyFeatureGates(masterConfiguration *kubeadmv1.MasterConfiguration) ([]byte, error) {
	j, err := json.Marshal(&masterConfiguration)
	if err != nil {
		return nil, fmt.Errorf("error marshalling kubeadm configuration: %s", err)
	}
	p, err := gabs.ParseJSON(j)
	fgString, ok := p.Path("kubeProxy.config.featureGates").Data().(string)
	if !ok {
		return nil, fmt.Errorf("error marshalling kubeadm configuration: error parsing kubeProxy.config.featureGates: %s", err)
	}
	p.ObjectP("kubeProxy.config.featureGates")
	if strings.Contains(fgString, ",") {
		for _, gate := range strings.Split(fgString, ",") {
			p.SetP(true, fmt.Sprintf("kubeProxy.config.featureGates.%s", gate))
		}
	}
	y, err := yaml.JSONToYAML(p.Bytes())
	if err != nil {
		return nil, fmt.Errorf("error marshalling kubeadm configuration: %s", err)
	}
	return y, nil
}
