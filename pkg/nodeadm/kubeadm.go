package nodeadm

import (
	spv1 "github.com/platform9/ssh-provider/pkg/apis/sshprovider/v1alpha1"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

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
	metav1.TypeMeta `json:",inline"`

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
	metav1.TypeMeta `json:",inline"`

	// NodeRegistration holds fields that relate to registering the new master node to the cluster
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
