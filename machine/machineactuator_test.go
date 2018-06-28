package machine_test

import (
	"testing"

	"github.com/platform9/ssh-provider/machine"
	"github.com/platform9/ssh-provider/util"
)

func TestCreateMaster(t *testing.T) {
	c, err := util.ClusterFromFile("./testdata/cluster.yaml")
	if err != nil {
		t.Fatal(err)
	}

	ms, err := util.MachinesFromFile("./testdata/machines.yaml")
	if err != nil {
		t.Fatal(err)
	}

	cms, err := util.ConfigMapsFromFile("./testdata/provisionedmachine-configmaps.yaml")
	if err != nil {
		t.Fatal(err)
	}

	sshCreds, err := util.SecretFromFile("./testdata/ssh-credentials.yaml")
	if err != nil {
		t.Fatal(err)
	}

	etcdCA, err := util.SecretFromFile("./testdata/etcd-ca.yaml")
	if err != nil {
		t.Fatal(err)
	}

	apiServerCA, err := util.SecretFromFile("./testdata/apiserver-ca.yaml")
	if err != nil {
		t.Fatal(err)
	}

	frontProxyCA, err := util.SecretFromFile("./testdata/front-proxy-ca.yaml")
	if err != nil {
		t.Fatal(err)
	}

	serviceAccountKey, err := util.SecretFromFile("./testdata/serviceaccount-key.yaml")
	if err != nil {
		t.Fatal(err)
	}

	sa, err := machine.NewActuator(cms, sshCreds, etcdCA, apiServerCA, frontProxyCA, serviceAccountKey)
	if err != nil {
		t.Fatal(err)
	}

	err = sa.Create(c, ms[0])
	if err != nil {
		t.Fatal(err)
	}
}
