package nodeadm_test

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"

	spv1 "github.com/platform9/ssh-provider/pkg/apis/sshprovider/v1alpha1"
	"github.com/platform9/ssh-provider/pkg/controller"
	"github.com/platform9/ssh-provider/pkg/nodeadm"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	clusterv1 "sigs.k8s.io/cluster-api/pkg/apis/cluster/v1alpha1"
)

var (
	testClusterName     = "test-cluster"
	testNamespace       = "test-namespace"
	testMachineName     = "test-machine"
	testProvMachineName = "test-pm"

	testServicesCIDRBlocks = []string{"172.0.0.0/24"}
	testPodsCIDRBlocks     = []string{"10.0.0.0/16"}
	testServiceDomain      = "cluster.local"

	testVIP      = "192.168.1.1"
	testRouterID = 100

	testVIPNetworkInterface = "eth0"

	testToken                    = "foo"
	testDiscoveryTokenAPIServers = []string{"example.com:1234"}
	testDiscoveryTokenCAHashes   = []string{"bar"}
)

func newCluster() (*clusterv1.Cluster, error) {
	cluster := clusterv1.Cluster{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Cluster",
			APIVersion: "cluster.k8s.io/v1alpha1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      testClusterName,
			Namespace: testNamespace,
		},
		Spec: clusterv1.ClusterSpec{
			ClusterNetwork: clusterv1.ClusterNetworkingConfig{
				Services: clusterv1.NetworkRanges{
					CIDRBlocks: testServicesCIDRBlocks,
				},
				Pods: clusterv1.NetworkRanges{
					CIDRBlocks: testPodsCIDRBlocks,
				},
				ServiceDomain: testServiceDomain,
			},
		},
	}
	cps := spv1.ClusterSpec{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ClusterSpec",
			APIVersion: "sshprovider.platform9.com/v1alpha1",
		},
		ClusterConfig: &spv1.ClusterConfig{
			KubeAPIServer: map[string]string{
				"service-node-port-range": "8000-32767",
				"secure-port":             "6445",
			},
			KubeControllerManager: map[string]string{
				"pod-eviction-timeout": "42s",
			},
			KubeScheduler: map[string]string{
				"log-dir": "/var/log/scheduler",
			},
			KubeProxy: &spv1.KubeProxyConfiguration{
				Mode: spv1.ProxyMode("iptables"),
			},
			Kubelet: &spv1.KubeletConfiguration{
				KubeAPIBurst: 84,
			},
		},
		VIPConfiguration: &spv1.VIPConfiguration{
			IP:       testVIP,
			RouterID: testRouterID,
		},
	}
	if err := controller.PutClusterSpec(cps, &cluster); err != nil {
		return nil, fmt.Errorf("unable to serialize provider cluster spec to cluster object: %s", err)
	}
	return &cluster, nil
}

func newMachine() (*clusterv1.Machine, error) {
	machine := clusterv1.Machine{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Machine",
			APIVersion: "cluster.k8s.io/v1alpha1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      testMachineName,
			Namespace: testNamespace,
		},
		Spec: clusterv1.MachineSpec{
			Taints: []v1.Taint{},
		},
	}
	mps := spv1.MachineSpec{
		TypeMeta: metav1.TypeMeta{
			Kind:       "MachineSpec",
			APIVersion: "sshprovider.platform9.com/v1alpha1",
		},
		ComponentVersions: &spv1.MachineComponentVersions{
			KubernetesVersion: "v1.11.7",
		},
	}
	if err := controller.PutMachineSpec(mps, &machine); err != nil {
		return nil, fmt.Errorf("unable to serialize provider machine spec to machine object: %s", err)
	}
	return &machine, nil
}

func newProvMachine() *spv1.ProvisionedMachine {
	pm := spv1.ProvisionedMachine{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ProvisionedMachine",
			APIVersion: "sshprovider.platform9.com/v1alpha1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      testProvMachineName,
			Namespace: testNamespace,
		},
		Spec: spv1.ProvisionedMachineSpec{
			VIPNetworkInterface: testVIPNetworkInterface,
		},
	}
	return &pm
}

// TestInitConfigurationForMachine
func TestInitConfigurationForMachine(t *testing.T) {
	cluster, err := newCluster()
	if err != nil {
		t.Fatal("unable to create cluster:", err)
	}
	machine, err := newMachine()
	if err != nil {
		t.Fatal("unable to create machine:", err)
	}
	pm := newProvMachine()

	expectedCfg := nodeadm.InitConfiguration{
		MasterConfiguration: nodeadm.KubeadmMasterConfiguration{
			TypeMeta: metav1.TypeMeta{
				Kind:       "MasterConfiguration",
				APIVersion: "kubeadm.k8s.io/v1alpha2",
			},
			API: nodeadm.API{
				AdvertiseAddress:     "",
				BindPort:             6445,
				ControlPlaneEndpoint: testVIP,
			},
			APIServerCertSANs: []string{testVIP},
			Etcd: nodeadm.Etcd{
				External: &nodeadm.ExternalEtcd{
					Endpoints: []string{"https://127.0.0.1:2379"},
					CAFile:    "/etc/etcd/pki/ca.crt",
					CertFile:  "/etc/etcd/pki/apiserver-etcd-client.crt",
					KeyFile:   "/etc/etcd/pki/apiserver-etcd-client.key",
				},
			},
			KubernetesVersion: "v1.11.7",
			Networking: nodeadm.Networking{
				ServiceSubnet: testServicesCIDRBlocks[0],
				PodSubnet:     testPodsCIDRBlocks[0],
				DNSDomain:     testServiceDomain},
			KubeletConfiguration: nodeadm.KubeletConfiguration{
				BaseConfig: &spv1.KubeletConfiguration{
					KubeAPIBurst: 84,
				},
			},
			KubeProxy: nodeadm.KubeProxy{
				Config: &spv1.KubeProxyConfiguration{
					Mode: spv1.ProxyMode("iptables"),
				},
			},
			APIServerExtraArgs: map[string]string{
				"service-node-port-range": "8000-32767",
			},
			ControllerManagerExtraArgs: map[string]string{
				"pod-eviction-timeout": "42s",
			},
			SchedulerExtraArgs: map[string]string{
				"log-dir": "/var/log/scheduler",
			},
			NodeRegistration: nodeadm.NodeRegistrationOptions{
				Name:             testMachineName,
				Taints:           []v1.Taint{},
				KubeletExtraArgs: nil,
			},
		},
		VIPConfiguration: &nodeadm.VIPConfiguration{
			IP:               testVIP,
			RouterID:         testRouterID,
			NetworkInterface: testVIPNetworkInterface,
		},
		NetworkBackend: nil,
		KeepAlived:     nil,
	}

	actualCfg, err := nodeadm.InitConfigurationForMachine(*cluster, *machine, *pm)
	if err != nil {
		t.Fatal("unable to generate kubeadm InitConfiguration:", err)
	}

	if diff := cmp.Diff(expectedCfg, *actualCfg); diff != "" {
		t.Errorf("InitConfigurationForMachine() mismatch (-want +got):\n%s", diff)
	}
}

// TestJoinConfigurationForMachine
func TestJoinConfigurationForMachine(t *testing.T) {
	cluster, err := newCluster()
	if err != nil {
		t.Fatal("unable to create cluster:", err)
	}
	machine, err := newMachine()
	if err != nil {
		t.Fatal("unable to create machine:", err)
	}

	expectedCfg := &nodeadm.JoinConfiguration{
		NodeConfiguration: nodeadm.KubeadmNodeConfiguration{
			TypeMeta: metav1.TypeMeta{
				Kind:       "NodeConfiguration",
				APIVersion: "kubeadm.k8s.io/v1alpha2",
			},
			Token:                      testToken,
			DiscoveryTokenAPIServers:   testDiscoveryTokenAPIServers,
			DiscoveryTokenCACertHashes: testDiscoveryTokenCAHashes,
			NodeRegistration: nodeadm.NodeRegistrationOptions{
				Name:             testMachineName,
				Taints:           []v1.Taint{},
				KubeletExtraArgs: nil,
			},
		},
	}

	actualCfg, err := nodeadm.JoinConfigurationForMachine(cluster, machine, testDiscoveryTokenAPIServers, testDiscoveryTokenCAHashes, testToken)
	if err != nil {
		t.Fatal("unable to generate kubeadm JoinConfiguration:", err)
	}

	if diff := cmp.Diff(expectedCfg, actualCfg); diff != "" {
		t.Errorf("JoinConfigurationForMachine() mismatch (-want +got):\n%s", diff)
	}

}
