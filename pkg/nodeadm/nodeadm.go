package nodeadm

import (
	"fmt"
	"strconv"

	"k8s.io/api/core/v1"

	spconstants "github.com/platform9/ssh-provider/constants"
	spv1 "github.com/platform9/ssh-provider/pkg/apis/sshprovider/v1alpha1"
	"github.com/platform9/ssh-provider/pkg/controller"
	clusterv1 "sigs.k8s.io/cluster-api/pkg/apis/cluster/v1alpha1"
)

type InitConfiguration struct {
	MasterConfiguration KubeadmMasterConfiguration `json:"masterConfiguration,omitempty"`

	VIPConfiguration VIPConfiguration  `json:"vipConfiguration,omitempty"`
	NetworkBackend   map[string]string `json:"networkBackend,omitempty"`
	KeepAlived       map[string]string `json:"keepAlived,omitempty"`
}

type JoinConfiguration struct {
	NodeConfiguration KubeadmNodeConfiguration `json:"nodeConfiguration"`
}

type VIPConfiguration struct {
	// The virtual IP.
	IP string `json:"ip,omitempty"`
	// The virtual router ID. Must be in the range [0, 254]. Must be unique within
	// a single L2 network domain.
	RouterID int `json:"routerID,omitempty"`
	// Network interface chosen to create the virtual IP. If it is not specified,
	// the interface of the default gateway is chosen.
	NetworkInterface string `json:"networkInterface,omitempty"`
}

// Networking contains elements describing cluster's networking configuration
type Networking struct {
	// ServiceSubnet is the subnet used by k8s services. Defaults to "10.96.0.0/12".
	ServiceSubnet string `json:"serviceSubnet,omitempty"`
	// PodSubnet is the subnet used by pods.
	PodSubnet string `json:"podSubnet,omitempty"`
	// DNSDomain is the dns domain used by k8s services. Defaults to "cluster.local".
	DNSDomain string `json:"dnsDomain,omitempty"`
}

type KubeadmMasterConfiguration struct {
	API                        API                         `json:"api,omitempty"`
	APIServerCertSANs          []string                    `json:"apiServerCertSANs,omitempty"`
	Etcd                       Etcd                        `json:"etcd,omitempty"`
	KubernetesVersion          string                      `json:"kubernetesVersion,omitempty"`
	Networking                 Networking                  `json:"networking,omitempty"`
	KubeletConfiguration       spv1.KubeletConfiguration   `json:"kubeletConfiguration,omitempty"`
	KubeProxy                  spv1.KubeProxyConfiguration `json:"kubeProxy,omitempty"`
	APIServerExtraArgs         map[string]string           `json:"apiServerExtraArgs,omitempty"`
	ControllerManagerExtraArgs map[string]string           `json:"controllerManagerExtraArgs,omitempty"`
	SchedulerExtraArgs         map[string]string           `json:"schedulerExtraArgs,omitempty"`

	NodeRegistration NodeRegistrationOptions `json:"nodeRegistration"`
}

type KubeadmNodeConfiguration struct {
	NodeRegistration NodeRegistrationOptions `json:"nodeRegistration"`
}

// NodeRegistrationOptions holds fields that relate to registering a new master or node to the cluster, either via "kubeadm init" or "kubeadm join"
type NodeRegistrationOptions struct {
	// Name is the `.Metadata.Name` field of the Node API object that will be created in this `kubeadm init` or `kubeadm joiń` operation.
	// This field is also used in the CommonName field of the kubelet's client certificate to the API server.
	// Defaults to the hostname of the node if not provided.
	Name string `json:"name,omitempty"`

	// Taints specifies the taints the Node API object should be registered with. If this field is unset, i.e. nil, in the `kubeadm init` process
	// it will be defaulted to []v1.Taint{'node-role.kubernetes.io/master=""'}. If you don't want to taint your master node, set this field to an
	// empty slice, i.e. `taints: {}` in the YAML file. This field is solely used for Node registration.
	Taints []v1.Taint `json:"taints,omitempty"`

	// KubeletExtraArgs passes through extra arguments to the kubelet. The arguments here are passed to the kubelet command line via the environment file
	// kubeadm writes at runtime for the kubelet to source. This overrides the generic base-level configuration in the kubelet-config-1.X ConfigMap
	// Flags have higher higher priority when parsing. These values are local and specific to the node kubeadm is executing on.
	KubeletExtraArgs map[string]string `json:"kubeletExtraArgs,omitempty"`
}

type API struct {
	AdvertiseAddress     string `json:"advertiseAddress,omitempty"`
	BindPort             int32  `json:"bindPort,omitempty"`
	ControlPlaneEndpoint string `json:"controlPlaneEndpoint"`
}

type Etcd struct {
	Endpoints []string `json:"endpoints,omitempty"`
	CAFile    string   `json:"caFile,omitempty"`
	CertFile  string   `json:"certFile,omitempty"`
	KeyFile   string   `json:"keyFile,omitempty"`
}

func InitConfigurationForMachine(cluster clusterv1.Cluster, machine clusterv1.Machine, pm spv1.ProvisionedMachine) (*InitConfiguration, error) {
	cfg := &InitConfiguration{}

	cpc, err := controller.GetClusterSpec(cluster)
	if err != nil {
		return nil, fmt.Errorf("unable to decode cluster spec: %v", err)
	}

	// NodeRegistrationOptions
	cfg.MasterConfiguration.NodeRegistration.Name = machine.Name
	cfg.MasterConfiguration.NodeRegistration.Taints = machine.Spec.Taints

	// MasterConfiguration
	cfg.MasterConfiguration.KubernetesVersion = machine.Spec.Versions.ControlPlane
	cfg.MasterConfiguration.Etcd.Endpoints = []string{"https://127.0.0.1:2379"}
	cfg.MasterConfiguration.Etcd.CAFile = "/etc/etcd/pki/ca.crt"
	cfg.MasterConfiguration.Etcd.CertFile = "/etc/etcd/pki/apiserver-etcd-client.crt"
	cfg.MasterConfiguration.Etcd.KeyFile = "/etc/etcd/pki/apiserver-etcd-client.key"
	if err := validateClusterNetworkingConfiguration(cluster); err != nil {
		return nil, fmt.Errorf("invalid cluster networking configuration: %v", err)
	}
	cfg.MasterConfiguration.Networking.DNSDomain = cluster.Spec.ClusterNetwork.ServiceDomain
	cfg.MasterConfiguration.Networking.PodSubnet = cluster.Spec.ClusterNetwork.Pods.CIDRBlocks[0]
	cfg.MasterConfiguration.Networking.ServiceSubnet = cluster.Spec.ClusterNetwork.Services.CIDRBlocks[0]

	// VIPConfiguration (optional)
	if cpc.VIPConfiguration != nil {
		cfg.VIPConfiguration.IP = cpc.VIPConfiguration.IP
		cfg.VIPConfiguration.RouterID = cpc.VIPConfiguration.RouterID
		cfg.VIPConfiguration.NetworkInterface = pm.Spec.VIPNetworkInterface

		cfg.MasterConfiguration.API.ControlPlaneEndpoint = cpc.VIPConfiguration.IP
		cfg.MasterConfiguration.APIServerCertSANs = []string{cpc.VIPConfiguration.IP}
	}

	// ClusterConfig (optional)
	if cpc.ClusterConfig != nil {
		setInitConfigFromClusterConfig(cfg, cpc.ClusterConfig)
	}

	return cfg, nil
}

// SetKubeAPIServerConfig sets configuration for API Server.
// Depending on the parameter name this function sets
// the MasterConfiguration fields or APIServerExtraArgs
func setKubeAPIServerConfig(cfg *InitConfiguration, clusterConfig *spv1.ClusterConfig) error {
	if clusterConfig.KubeAPIServer != nil {
		// Set fields for API server manually as there is no upstream type yet.
		// BindPort
		bindPortStr, ok := clusterConfig.KubeAPIServer[spconstants.KubeAPIServerSecurePortKey]
		if ok {
			bindPort, err := strconv.ParseInt(bindPortStr, 10, 32)
			if err != nil {
				return fmt.Errorf("unable to parse port value: %s", bindPortStr)
			}
			cfg.MasterConfiguration.API.BindPort = int32(bindPort)
			// delete as it should not be considered as an extra arg
			delete(clusterConfig.KubeAPIServer, spconstants.KubeAPIServerSecurePortKey)
		}
	}
	return nil
}

func setInitConfigFromClusterConfig(cfg *InitConfiguration, clusterConfig *spv1.ClusterConfig) error {
	if err := setKubeAPIServerConfig(cfg, clusterConfig); err != nil {
		return fmt.Errorf("unable to set configurable parameters for api-server: %v", err)
	}
	cfg.MasterConfiguration.ControllerManagerExtraArgs = clusterConfig.KubeControllerManager
	if clusterConfig.KubeProxy != nil {
		cfg.MasterConfiguration.KubeProxy = *clusterConfig.KubeProxy
	}
	cfg.MasterConfiguration.SchedulerExtraArgs = clusterConfig.KubeScheduler
	if clusterConfig.Kubelet != nil {
		cfg.MasterConfiguration.KubeletConfiguration = *clusterConfig.Kubelet
	}
	cfg.NetworkBackend = clusterConfig.NetworkBackend
	cfg.KeepAlived = clusterConfig.KeepAlived
	return nil
}

func JoinConfigurationForMachine(cluster *clusterv1.Cluster, machine *clusterv1.Machine) (*JoinConfiguration, error) {
	cfg := &JoinConfiguration{}

	// NodeRegistrationOptions
	cfg.NodeConfiguration.NodeRegistration.Name = machine.Name
	cfg.NodeConfiguration.NodeRegistration.Taints = machine.Spec.Taints

	return cfg, nil
}

func validateClusterNetworkingConfiguration(cluster clusterv1.Cluster) error {
	switch cbl := len(cluster.Spec.ClusterNetwork.Pods.CIDRBlocks); {
	case cbl < 1:
		return fmt.Errorf("cluster %q spec.clusterNetwork.pods.cidrBlocks must contain at least one block", cluster.Name)
	case cbl > 1:
		return fmt.Errorf("cluster %q spec.clusterNetwork.pods.cidrBlocks must contain at most one block", cluster.Name)
	}
	switch cbl := len(cluster.Spec.ClusterNetwork.Services.CIDRBlocks); {
	case cbl < 1:
		return fmt.Errorf("cluster %q spec.clusterNetwork.pods.cidrBlocks must contain at least one block", cluster.Name)
	case cbl > 1:
		return fmt.Errorf("cluster %q spec.clusterNetwork.pods.cidrBlocks must contain at most one block", cluster.Name)
	}
	return nil
}
