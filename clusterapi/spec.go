package clusterapi

import (
	cmav1alpha1 "sigs.k8s.io/cluster-api/pkg/apis/cluster/v1alpha1"
	"sigs.k8s.io/cluster-api/pkg/client/clientset_generated/clientset"
)

type Cluster interface {
	AddNodePool(name string) NodePoolBuilder
	Delete()

	GetLiveObject() (*cmav1alpha1.Cluster, error)
}

type ClusterBuilder interface {
	InRegion(r string) ClusterBuilder
	WithMaster() MasterBuilder

	Create() (Cluster, error)
}

type MasterBuilder interface {
	AvailabilityZone(az string) ClusterBuilder
	InstanceType(t string) ClusterBuilder
}

type NodePool interface {
	Delete()
	GetLiveObject() (*cmav1alpha1.MachineDeployment, error)
}

type NodePoolBuilder interface {
	WithAvailabilityZone(az string) NodePoolBuilder
	WithInstanceType(t string) NodePoolBuilder
	WithMinWorkers(n int) NodePoolBuilder
	WithMaxWorkers(n int) NodePoolBuilder
	Create() (NodePool, error)
}

type TestCase interface {
	NewCluster(name string) ClusterBuilder

	CMAClient() clientset.Interface
}
