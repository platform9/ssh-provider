/*
Copyright 2018 Platform 9 Systems, Inc.
*/

package v1alpha1

import (
	"net"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type SSHClusterProviderConfig struct {
	metav1.TypeMeta `json:",inline"`

	// CASecretName is the name of the Secret with the cluster CA certificate and
	// private key. If it is not specified, the default name is derived from the
	// cluster name. If the Secret is not present, the provider generates a
	// self-signed one and creates the Secret.
	// +optional
	CASecretName string `json:"caSecretName"`

	// VIPConfiguration is the configuration of the VIP for the API. If it is not
	// specified, the VIP is not created.
	// +optional
	VIPConfiguration *VIPConfiguraton `json:"vipConfiguration,omitempty"`
}

// VIPConfiguration specifies the parameters used to provision a virtual IP
// which API servers advertise and accept requests on.
type VIPConfiguration struct {
	// The virtual IP.
	IP net.IP `json:"ip"`
	// The virtual router ID. Must be in the range [0, 254]. Must be unique within
	// a single L2 network domain.
	RouterID string `json:"routerID"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type SSHClusterProviderStatus struct {
	metav1.TypeMeta `json:",inline"`
}
