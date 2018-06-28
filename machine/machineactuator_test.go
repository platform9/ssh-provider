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

	sshCreds, err := util.SecretFromFile("./testdata/sshcredentials-secret.yaml")
	if err != nil {
		t.Fatal(err)
	}

	clusterCA, err := util.SecretFromFile("./testdata/clusterca-secret.yaml")
	if err != nil {
		t.Fatal(err)
	}

	sa, err := machine.NewActuator(cms, sshCreds, clusterCA)
	if err != nil {
		t.Fatal(err)
	}

	err = sa.Create(c, ms[0])
	if err != nil {
		t.Fatal(err)
	}
}
