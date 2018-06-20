/*
Copyright 2018 Platform 9 Systems, Inc.
*/

package sshproviderconfig

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type SSHClusterProviderConfig struct {
	metav1.TypeMeta `json:",inline"`

	CASecretName       string      `json:"caSecretName"`
	VirtualAPIEndpoint APIEndpoint `json:"virtualAPIEndpoint"`
}

// APIEndpoint represents a reachable Kubernetes API endpoint.
type APIEndpoint struct {
	// The hostname on which the API server is serving.
	Host string `json:"host"`

	// The port on which the API server is serving.
	Port int `json:"port"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type SSHClusterProviderStatus struct {
	metav1.TypeMeta `json:",inline"`

	EtcdStatus EtcdStatus `json:"etcdStatus"`
}

type EtcdStatus struct {
	Members []EtcdMember `json:"members"`
	Token   string       `json:"token"`
}

type EtcdMember struct {
	// ID is the member ID for this member.
	ID uint64 `json:"ID,omitempty"`
	// name is the human-readable name of the member. If the member is not started, the name will be an empty string.
	Name string `json:"name,omitempty"`
	// peerURLs is the list of URLs the member exposes to the cluster for communication.
	PeerURLs []string `json:"peerURLs,omitempty"`
	// clientURLs is the list of URLs the member exposes to clients for communication. If the member is not started, clientURLs will be empty.
	ClientURLs []string `json:"clientURLs,omitempty"`
}
