/*
Copyright 2018 Platform 9 Systems, Inc.
*/

package actuator

import (
	clusterv1 "sigs.k8s.io/cluster-api/pkg/apis/cluster/v1alpha1"
)

type SSHActuator struct {
	InsecureIgnoreHostKey bool
}

func (sa *SSHActuator) Create(cluster *clusterv1.Cluster, machine *clusterv1.Machine) error {
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
