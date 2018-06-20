/*
Copyright 2018 Platform 9 Systems, Inc.
*/

package sshproviderconfig

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type SSHMachineProviderConfig struct {
	metav1.TypeMeta `json:",inline"`
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
