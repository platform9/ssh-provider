package util

import (
	"io/ioutil"

	"github.com/ghodss/yaml"
	corev1 "k8s.io/api/core/v1"
	clusterv1 "sigs.k8s.io/cluster-api/pkg/apis/cluster/v1alpha1"
)

func MachinesFromFile(file string) ([]*clusterv1.Machine, error) {
	bytes, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	list := &clusterv1.MachineList{}
	err = yaml.Unmarshal(bytes, &list)
	if err != nil {
		return nil, err
	}

	if list == nil {
		return []*clusterv1.Machine{}, nil
	}

	return machineP(list.Items), nil
}

func ClusterFromFile(file string) (*clusterv1.Cluster, error) {
	bytes, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	cluster := &clusterv1.Cluster{}
	err = yaml.Unmarshal(bytes, &cluster)
	if err != nil {
		return nil, err
	}

	return cluster, nil
}

func ConfigMapsFromFile(file string) ([]*corev1.ConfigMap, error) {
	bytes, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	list := &corev1.ConfigMapList{}
	err = yaml.Unmarshal(bytes, &list)
	if err != nil {
		return nil, err
	}

	if list == nil {
		return []*corev1.ConfigMap{}, nil
	}

	return configMapP(list.Items), nil
}

func SecretFromFile(file string) (*corev1.Secret, error) {
	bytes, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	secret := &corev1.Secret{}
	err = yaml.Unmarshal(bytes, &secret)
	if err != nil {
		return nil, err
	}

	return secret, nil
}

// Convert to list of pointers
func configMapP(configMaps []corev1.ConfigMap) []*corev1.ConfigMap {
	var ret []*corev1.ConfigMap
	for _, cm := range configMaps {
		ret = append(ret, cm.DeepCopy())
	}
	return ret
}

// Convert to list of pointers
func machineP(machines []clusterv1.Machine) []*clusterv1.Machine {
	var ret []*clusterv1.Machine
	for _, machine := range machines {
		ret = append(ret, machine.DeepCopy())
	}
	return ret
}
