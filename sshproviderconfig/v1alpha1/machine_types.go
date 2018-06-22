/*
Copyright 2018 Platform 9 Systems, Inc.
*/

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kubeadmv1 "k8s.io/kubernetes/cmd/kubeadm/app/apis/kubeadm/v1alpha2"
)

// SSHMachineProviderConfig defines the desired provider-specific state of the
// machine
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type SSHMachineProviderConfig struct {
	metav1.TypeMeta `json:",inline"`

	// MasterConfiguraton holds optional kubeadm configuration for a master.
	// Values here override values in the kubeadm configuration used by the
	// provider.
	// +optional

	MasterConfiguration *kubeadmv1.MasterConfiguration `json:"masterConfiguration,omitempty"`
	// NodeConfiguraton holds optional kubeadm configuration for a node.
	// Values here override values in the kubeadm configuration used by the
	// provider.
	// +optional
	NodeConfiguration *kubeadmv1.NodeConfiguration `json:"nodeConfiguration,omitEmpty"`
}

// SSHMachineProviderStatus defines the observed provider-specific state of the
// machine
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type SSHMachineProviderStatus struct {
	metav1.TypeMeta `json:",inline"`

	SSHConfig SSHConfig `json:"sshConfig"`

	// EtcdMember defines the observed etcd configuration of the machine.
	// This field is populated for masters only.
	// +optional
	EtcdMember *EtcdMember `json:"etcdMember,omitempty"`
}

// SSHConfig specifies everything needed to ssh to a host
type SSHConfig {
	// The IP or hostname used to SSH to the machine
	Host string `json:"host"`
	// The used to SSH to the machine
	Port int `json:"port"`
	// The SSH public keys of the machine
	PublicKeys []string `json:"publicKeys"`
	// The Secret with the username and private key used to SSH to the machine
	SecretName string `json:"secretName"`
}

// EtcdMember defines the configuration of an etcd member.
type EtcdMember struct {
	// ID is the member ID for this member.
	ID uint64 `json:"ID"`
	// Name is the human-readable name of the member.
	Name string `json:"name"`
	// PeerURLs is the list of URLs the member exposes to the cluster for communication.
	PeerURLs []string `json:"peerURLs"`
	// ClientURLs is the list of URLs the member exposes to clients for communication.
	ClientURLs []string `json:"clientURLs"`
}
