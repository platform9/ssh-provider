/*
Copyright 2018 Platform 9 Systems, Inc.
*/

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type SSHMachineProviderConfig struct {
	metav1.TypeMeta `json:",inline"`

	// maxPods is the number of pods that can run on this Kubelet.
	// Dynamic Kubelet Config (beta): If dynamically updating this field, consider that
	// changes may cause Pods to fail admission on Kubelet restart, and may change
	// the value reported in Node.Status.Capacity[v1.ResourcePods], thus affecting
	// future scheduling decisions. Increasing this value may also decrease performance,
	// as more Pods can be packed into a single node.
	// Default: 110
	// +optional
	MaxPods int32 `json:"maxPods,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type SSHMachineProviderStatus struct {
	metav1.TypeMeta `json:",inline"`

	Host string `json:"host"`
	Port int    `json:"port"`

	// The host's known SSH public keys
	PublicKeys    []string `json:"publicKeys"`
	SSHSecretName string   `json:"sshSecretName"`
}
