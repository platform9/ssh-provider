package machine_test

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net"
	"os"
	"testing"

	"github.com/ghodss/yaml"
	"github.com/platform9/ssh-provider/pkg/nodeadm"

	"github.com/apparentlymart/go-cidr/cidr"
	"github.com/golang/mock/gomock"

	spv1 "github.com/platform9/ssh-provider/pkg/apis/sshprovider/v1alpha1"
	spclientfake "github.com/platform9/ssh-provider/pkg/client/clientset_generated/clientset/fake"
	"github.com/platform9/ssh-provider/pkg/clusterapi/machine"
	"github.com/platform9/ssh-provider/pkg/controller"
	spmachine "github.com/platform9/ssh-provider/pkg/machine"
	mockmachine "github.com/platform9/ssh-provider/pkg/machine/mock"

	kubeclientfake "k8s.io/client-go/kubernetes/fake"
	clusterclientfake "sigs.k8s.io/cluster-api/pkg/client/clientset_generated/clientset/fake"

	capicommon "sigs.k8s.io/cluster-api/pkg/apis/cluster/common"
	clusterv1 "sigs.k8s.io/cluster-api/pkg/apis/cluster/v1alpha1"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/equality"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/diff"
)

const (
	testNamespace = "test-namespace"

	testServiceDomain = "cluster.local"

	testVIPNetworkInterface = "testiface"
	testClusterName         = "test-cluster"
	testSSHPort             = 22

	testSSHCredentialSecretName = "test-ssh-credential"

	// Corresponding public key: `ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQC7L2PL1Jh2dKEQjaqjIyYkxKjp+roLafIOsOS+9GQXqwjWDBajaXUVVOqPGiRYfGZD1/JGtaLZmcnPh4ROkWDanWEiatWtWOpnfESGo0CCYIpOkOZ+hhH5xDv7Iw//M1ESMCLuQDtCr5LzMeM/tP7BtLTDHWagu50Phrqp2iDauq8TDKK7R1BMnuH9MuoiszeTweeidKBRUO3xCbJ/h/GaqwoHvZiVjb/y8GbSymN68mljcPWCjm1BFlX9qPnKiisJILGy2DUogzfTFWUxqfVIJY4F2Nw2QzRPr5c3NUjVBx9UU8m3yujePxlSdiLWSmwuAoUj3fGDlwQo76cOiUul test-ssh-username@localmachine`
	testSSHPrivateKey = `-----BEGIN RSA PRIVATE KEY-----
MIIEpAIBAAKCAQEAuy9jy9SYdnShEI2qoyMmJMSo6fq6C2nyDrDkvvRkF6sI1gwW
o2l1FVTqjxokWHxmQ9fyRrWi2ZnJz4eETpFg2p1hImrVrVjqZ3xEhqNAgmCKTpDm
foYR+cQ7+yMP/zNREjAi7kA7Qq+S8zHjP7T+wbS0wx1moLudD4a6qdog2rqvEwyi
u0dQTJ7h/TLqIrM3k8HnonSgUVDt8Qmyf4fxmqsKB72YlY2/8vBm0spjevJpY3D1
go5tQRZV/aj5yoorCSCxstg1KIM30xVlMan1SCWOBdjcNkM0T6+XNzVI1QcfVFPJ
t8ro3j8ZUnYi1kpsLgKFI93xg5cEKO+nDolLpQIDAQABAoIBAGo09FnV0X/8otqi
lvwjWGQqVXEU+vS480fUpKWYQyaOu9+/UiT8FBu7Z680dQTj6J5765KlQrZWpQQk
bzSpFhxUiaWemojV14JKJxOBk3umTTNQ0gdeCNr/tczD0zLIqio4t8kZFsa6mhV0
6+zcxGOiJaJDj1SJvw7tMgJnqpaPtsoJK+Y8PRpkw5FafPyxH50Dqf0NvvayxQTX
OKjBgxMqDkPcTeVIykiuI4apRPUh2wqFIBrLwXhkFYxlg6eV4mgPM0AflmQlC4c6
sHzOLtDh9mcTJ7385JMkAMb1oKfxIi7HQcJ1afxuMRiVT0nQa8dEqTKahuEZIYTo
xCwDMC0CgYEA7TJDE3MyRfqEUrBgZNPyJLxCY9UOdVK6QQW0lCc/92TTeWXKDpKX
iFge0s7BCatIbSerqkkWEztOzi+SnfQ8YuxcA/EN5qsfju5ZDOEGlOXVEhkMtV+J
BRo+pasNKTTY/cFuXXrNj0gwWQyEW0bBRuTseWxEkG5Q2sjK3yAmhOcCgYEAygYv
J5dkAM8ScXY9NRodVKSSUsJWaMJWWUL06wOild03jzcJvkgQSrE2BYWXI8BY58bt
BR0CtSGIrOqBkFLnEmrJPbPkPmxY372KcPOByiZtfItZuAgljtrWVQPTv/Kr84tv
72dT/sKUFFkmP0Y7SquiePcm3bc5kpBB6KJSzZMCgYEA12IKkoDKJ80hlsxd23Cx
sjXYuzpeFJ74Tc7WeKljZkwB03xgi/cE7mPGKXpgw1zrOWMKeGhBSVlTZ9q+8fsz
Ukh6DYI4Mzs0Pt+jsRQsm8xPPE3OwmtrylxYgqreGorIdzPL+cpnGjJc5J9/GDsZ
ijyZlnB/mf7yIJivzwgssgUCgYAYQP3QRwCeiD2ymXtnsKbszoUyolo1YX90X/F/
dYRDcbeF3hmxWn16OiQ+LXejtyf1S5CRGJaGoGwENnMvnIRZVbCrU7mHNQLGeb7x
KIrgdhSW6zPuQCwiJmus8TSuyreSGZ9qooEXCM02VC2fUxMzN7/dve5Ql/q1edQv
1a0TOQKBgQCbxFgxeOEAIm7x9AlHy3LRWXo2vP6FcmSH2aQrsHZ3LBAH/qH1fMxK
xZTdlwxDuehPVIko3zwy5JVaLD7BpnrYifWFk9ijJInsp5dD3wvW3UzHox4ZEINT
BVf3lXQbxrHFMPHuAZ3lSAmarx8AzGYPP00rgFJoTlBC5HGl/1nQLg==
-----END RSA PRIVATE KEY-----`
	testSSHUsername = "test-ssh-username"

	testEtcdCASecretName            = "test-etcd-ca-secret"
	testAPIServerCASecretName       = "test-apiserver-ca-secret"
	testFrontProxyCASecretName      = "test-front-proxy-ca-secret"
	testServiceAccountKeySecretName = "test-service-account-key-secret"

	testVIP      = "192.168.1.1"
	testRouterID = 100

	testKubeletVersion      = "1.10.6"
	testControlPlaneVersion = "1.10.6"
)

var (
	testProvisionedMachinesCIDR = "192.168.0.0/24"
	testServicesCIDRBlocks      = []string{"172.0.0.0/24"}
	testPodsCIDRBlocks          = []string{"10.0.0.0/16"}
	testVIPConfiguration        = spv1.VIPConfiguration{
		IP:       testVIP,
		RouterID: testRouterID,
	}
	testSSHCredentialSecret = corev1.Secret{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Secret",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      testSSHCredentialSecretName,
			Namespace: testNamespace,
		},
		Data: map[string][]byte{
			spv1.CredentialSecretSSHPrivateKeyKey: []byte(testSSHPrivateKey),
			spv1.CredentialSecretUsernameKey:      []byte(testSSHUsername),
		},
	}

	testNodeadmInitCmd = fmt.Sprintf("%s init --cfg %s", machine.NodeadmPath, machine.NodeadmConfigPath)
	testEtcdadmInitCmd = fmt.Sprintf("%s init", machine.EtcdadmPath)
	testEtcdadmInfoCmd = fmt.Sprintf("%s info", machine.EtcdadmPath)
)

type sequentialIPAllocator struct {
	t        *testing.T
	ipNet    *net.IPNet
	maxHosts int
	nextHost int
}

func newSequentialIPAllocator(t *testing.T, CIDR string) *sequentialIPAllocator {
	_, ipNet, err := net.ParseCIDR(CIDR)
	if err != nil {
		t.Fatalf("unable to create new IP allocator: %v", err)
	}
	return &sequentialIPAllocator{
		t:        t,
		ipNet:    ipNet,
		maxHosts: int(cidr.AddressCount(ipNet)) - 1, // do not allocate broadcast address
		nextHost: 1,                                 // do not allocate network identifier
	}
}

func (a *sequentialIPAllocator) NextIP() net.IP {
	ip, err := cidr.Host(a.ipNet, a.nextHost)
	if err != nil {
		a.t.Fatalf("unable to allocate IP: %v", err)
	}
	a.nextHost++
	return ip
}

func testCluster(t *testing.T, vipConfiguration spv1.VIPConfiguration) *clusterv1.Cluster {
	clusterSpec := spv1.ClusterSpec{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ClusterSpec",
			APIVersion: "sshprovider.platform9.com/v1alpha1",
		},
		EtcdCASecret: &corev1.LocalObjectReference{
			Name: testEtcdCASecretName,
		},
		APIServerCASecret: &corev1.LocalObjectReference{
			Name: testAPIServerCASecretName,
		},
		FrontProxyCASecret: &corev1.LocalObjectReference{
			Name: testFrontProxyCASecretName,
		},
		ServiceAccountKeySecret: &corev1.LocalObjectReference{
			Name: testServiceAccountKeySecretName,
		},
		VIPConfiguration: &vipConfiguration,
	}
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
	if err := controller.PutClusterSpec(clusterSpec, &cluster); err != nil {
		t.Fatalf("Unable to create test cluster: %v", err)
	}
	return &cluster
}

func newSSHConfig(name string) *spv1.SSHConfig {
	return &spv1.SSHConfig{
		Host: name,
		Port: testSSHPort,
		PublicKeys: []string{
			"ssh-dss AAAAB3NzaC1kc3MAAACBAPkFvDBlzDXDt/R1vMyASCyKohERRWVq2KHtz500N6yy6H7VEgXwEA8bDmAN9Xye0GBqFc/WmfV757+y6vuPul/t8Re5APTLCFH/Q4XBp2kxchVGtyGB6ihVwMoGar2IMMQuneVkz0+/fn92fC6wOZtT+YIAMpWghieqJhUjtgz7AAAAFQDsH/tN4UZIpP0sfsAxwh/o76LB6QAAAIA8OAwT8ufxdXB7EW9T7oP3AEKuvJ95q+eSyBEPJsOfaWWzGLsvy1lpbpuTm/Ffz77mSKKQnFO4PiZgCX+OTdTXjgc1JN9v+GPEmlXA2FaA4cDsiKMyXdgN1ncNnlfa1ZVRwzg9mCqsPrX2Zt2r0o0N9LSB+aQ62WOPEdtRqNVp+AAAAIAZUUOUNg0Jd8vTtj0p2zY5AVcBBlB+v41keodjbgxkT1gwX8FvRGyHzyqIJr0Vbyj4i7Yj7hlLI6JgkxGLlDcZSO8lqLll133tE75XY70yMTP4ff8Dp18II3PDUg/fqrDSmC9ZfsXY3dYQG7pfBnz9EOziWp3EzkO1CSe0n6RKtg== root@localhost",
			"ecdsa-sha2-nistp256 AAAAE2VjZHNhLXNoYTItbmlzdHAyNTYAAAAIbmlzdHAyNTYAAABBBLd8oe2tsyRTuctnUbIRAQMzbe00+E9IlPDi7znwXf5uOmIPKiG641tWEfCw1HIJBz10jkxljqNG+4zbOjFnhqg= root@localhost",
			"ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIMlvMz8+NI/SxOkQCmNZXFdDYuJ+sRgBLWCEWGPTloEn root@localhost",
			"ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDc2MCBwSo9MONFrdmNV+ouzAZIY83f+S/SUQUuWuIHy5lmtXVkBJFSK/3meYZQVoxVIQKyUshRG+6lQ2lwOrPY3UoMiX7ezYX3d7Zb594cHJ7lxj38I0oRJpllYF6xjIebHG7Sl8OEBw/eWu94+ISGmTGokxviBBJ9Rhkic3/NE8Kf0K62WoXtiHiOdhCPw9vSz3j1w/7KdN0hKxzEmRsBxBiu4aOnBKTBAYyys0Mc6F6pfFkzEG/q14Qn2aje325mnS+FiKKZcxEZiHtrqe220jYbvsOxEsYw3Iz6jWsj5kz4EMYYFj0cSiLW2b7ioK6nsIYwBgLn1Wii9SiW5+ib root@localhost",
		},
		CredentialSecret: corev1.LocalObjectReference{
			Name: testSSHCredentialSecretName,
		},
	}
}

func newProvisionedMachine(name string) *spv1.ProvisionedMachine {
	provMachine := &spv1.ProvisionedMachine{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ProvisionedMachine",
			APIVersion: "sshprovider.platform9.com/v1alpha1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: testNamespace,
		},
		Spec: spv1.ProvisionedMachineSpec{
			SSHConfig:           newSSHConfig(name),
			VIPNetworkInterface: testVIPNetworkInterface,
		},
	}
	return provMachine
}

func newMachine(t *testing.T, name string) *clusterv1.Machine {
	machineSpec := spv1.MachineSpec{
		TypeMeta: metav1.TypeMeta{
			Kind:       "MachineSpec",
			APIVersion: "sshprovider.platform9.com/v1alpha1",
		},
	}
	machine := &clusterv1.Machine{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Machine",
			APIVersion: "cluster.k8s.io/v1alpha1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: testNamespace,
		},
		Spec: clusterv1.MachineSpec{
			Roles: []capicommon.MachineRole{
				capicommon.MasterRole,
			},
			Versions: clusterv1.MachineVersionInfo{
				Kubelet:      testKubeletVersion,
				ControlPlane: testControlPlaneVersion,
			},
		},
	}
	if err := controller.PutMachineSpec(machineSpec, machine); err != nil {
		t.Fatalf("Unable to create new machine: %v", err)
	}
	return machine
}

func TestCreateMaster(t *testing.T) {
	provMachineIPAllocator := newSequentialIPAllocator(t, testProvisionedMachinesCIDR)

	kubeClient := kubeclientfake.NewSimpleClientset()
	clusterClient := clusterclientfake.NewSimpleClientset()
	spClient := spclientfake.NewSimpleClientset()

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockMachineClient := mockmachine.NewMockClient(mockCtrl)
	mockMachineClientBuilder := func(host string, port int, username string, privateKey string, publicKeys []string, insecureIgnoreHostKey bool) (spmachine.Client, error) {
		return mockMachineClient, nil
	}
	a := machine.NewActuator(kubeClient, clusterClient, spClient, mockMachineClientBuilder, false)

	c := testCluster(t, testVIPConfiguration)
	pm := newProvisionedMachine("pm1")
	pmIP := provMachineIPAllocator.NextIP()
	m := newMachine(t, "master1")
	if err := controller.BindMachineAndProvisionedMachine(m, pm); err != nil {
		t.Fatalf("Unable to bind machine and provisioned machine: %v", err)
	}

	if _, err := kubeClient.CoreV1().Secrets(testNamespace).Create(&testSSHCredentialSecret); err != nil {
		t.Fatalf("Unable to store Secret %q: %v", testSSHCredentialSecret.Name, err)
	}
	if _, err := clusterClient.ClusterV1alpha1().Machines(testNamespace).Create(m); err != nil {
		t.Fatalf("Unable to store Machine %q: %v", m.Name, err)
	}
	if _, err := spClient.SshproviderV1alpha1().ProvisionedMachines(testNamespace).Create(pm); err != nil {
		t.Fatalf("Unable to store ProvisionedMachine %q: %v", pm.Name, err)
	}

	// Prepare expected command output
	expectedEtcdMember := &spv1.EtcdMember{
		ID:         rand.Uint64(),
		Name:       pm.Spec.SSHConfig.Host,
		ClientURLs: []string{fmt.Sprintf("https://%s:2379", pmIP.String())},
		PeerURLs:   []string{fmt.Sprintf("https://%s:2380", pmIP.String())},
	}
	expectedEtcdInfoStdout, err := json.Marshal(expectedEtcdMember)
	if err != nil {
		t.Fatalf("unable to create expected etcd info output: %v", err)
	}

	expectedNodeadmInitConfig, err := nodeadm.InitConfigurationForMachine(*c, *m, *pm)
	if err != nil {
		t.Fatalf("unable to create nodedam init configuration: %v", err)
	}
	expectedNodeadmInitConfigBytes, err := yaml.Marshal(expectedNodeadmInitConfig)
	if err != nil {
		t.Fatalf("unable to marshal nodeadm init configuration to YAML: %v", err)
	}

	// Prepare expected operations, in sequence
	gomock.InOrder(
		mockMachineClient.EXPECT().
			RunCommand(testEtcdadmInitCmd).
			Return([]byte{}, []byte{}, nil),
		mockMachineClient.EXPECT().
			RunCommand(testEtcdadmInfoCmd).
			Return(expectedEtcdInfoStdout, nil, nil),
		mockMachineClient.EXPECT().
			WriteFile(machine.NodeadmConfigPath, os.FileMode(0600), expectedNodeadmInitConfigBytes).
			Return(nil),
		mockMachineClient.EXPECT().
			RunCommand(testNodeadmInitCmd).
			Return([]byte{}, []byte{}, nil),
	)

	if err := a.Create(c, m); err != nil {
		t.Fatalf("Actuator Create() failed: %v", err)
	}

	// Verify actuator posted expected machine status
	m, err = clusterClient.ClusterV1alpha1().Machines(m.Namespace).Get(m.Name, metav1.GetOptions{})
	if err != nil {
		t.Fatalf("unable to get machine %q: %v", m.Name, err)
	}
	machineStatus, err := controller.GetMachineStatus(*m)
	if err != nil {
		t.Fatalf("unable to decode machine %q status: %v", m.Name, err)
	}
	if !equality.Semantic.DeepEqual(expectedEtcdMember, machineStatus.EtcdMember) {
		t.Fatalf("actual etcd member is different from the expected one: %v", diff.ObjectDiff(expectedEtcdMember, machineStatus.EtcdMember))
	}
	if !equality.Semantic.DeepEqual(pm.Spec.SSHConfig, machineStatus.SSHConfig) {
		t.Fatalf("actual ssh config is different from the expected one: %v", diff.ObjectDiff(pm.Spec.SSHConfig, machineStatus.SSHConfig))
	}
}

func TestCreateTwoMasters(t *testing.T) {
	// create cluster
	// create cluster secrets (etcd, apiserver, frontproxy, serviceaccount)
	// create machine m1
	// create provisionedmachine pm1
	// bind machine to provisionedmachine
	// update machinestatus etcdmember
	// update machinestatus sshconfig
	// update clusterstatus etcdmember
	// --
	// create machine m2
	// create provisionedmachine pm2
	// bind machine to provisionedmachine
	// a.Create(c, m2)
}

func TestCreateThreeMasters(t *testing.T) {

}
func TestCreateNode(t *testing.T) {

}

func TestCreateThreeMastersThreeNodes(t *testing.T) {

}
