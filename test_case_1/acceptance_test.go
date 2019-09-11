package main

import (
	"testing"

	"github.com/giantswarm/clusterfunk/clusterapi"
)

func Test_SimpleNodePoolCluster(t *testing.T) {
	test := clusterapi.NewTestCase("test case")

	cluster, err := test.NewCluster("cluster-1").
		WithMaster().AvailabilityZone("eu-central-1b").
		WithMaster().InstanceType("m4.large").
		InRegion("eu-central-1").
		Create()

	if err != nil {
		t.Fatal(err)
	}

	np, err := cluster.AddNodePool("np-1").
		WithAvailabilityZone("eu-central-1a").
		WithAvailabilityZone("eu-central-1b").
		WithInstanceType("m5.large").
		WithMinWorkers(2).
		WithMaxWorkers(5).
		Create()

	if err != nil {
		t.Fatal(err)
	}

	_, err := cluster.GetLiveObject()
	if err != nil {
		t.Fatal(err)
	}

	_, err := np.GetLiveObject()
	if err != nil {
		t.Fatal(err)
	}

	np.Delete()
	cluster.Delete()
}
