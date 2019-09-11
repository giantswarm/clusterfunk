package clusterapi

import (
	"github.com/giantswarm/microerror"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	cmav1alpha1 "sigs.k8s.io/cluster-api/pkg/apis/cluster/v1alpha1"
	"sigs.k8s.io/cluster-api/pkg/client/clientset_generated/clientset"
)

const (
	defaultAvailabilityZone   = "eu-central-1a"
	defaultMasterInstanceType = "m5.large"
	defaultRegion             = "eu-central-1"
)

type cluster struct {
	name      string
	namespace string

	testCase TestCase
}

type clusterBuilder struct {
	name   string
	master MasterBuilder
	region string

	testCase TestCase
}

type masterBuilder struct {
	availabilityZone string
	instanceType     string
	clusterBuilder   ClusterBuilder
}

type nodePoolBuilder struct {
	availabilityZones []string
	instanceType      string
	maxWorkers        int
	minWorkers        int
	name              string
}

type nodePool struct {
	name string

	testCase TestCase
}

type testCase struct {
	name string

	cmaClient clientset.Interface
}

func NewTestCase(name string) TestCase {
	return &testCase{name: name}
}

func (tc *testCase) NewCluster(name string) ClusterBuilder {
	cb := &clusterBuilder{
		name: name,

		testCase: tc,
	}

	cb.master = &masterBuilder{
		availabilityZone: defaultAvailabilityZone,
		instanceType:     defaultMasterInstanceType,
		clusterBuilder:   cb,
	}

	return cb
}

func (tc *testCase) CMAClient() clientset.Interface {
	return tc.cmaClient
}

func (cb *clusterBuilder) InRegion(r string) ClusterBuilder {
	cb.region = r
	return cb
}

func (cb *clusterBuilder) WithMaster() MasterBuilder {
	return cb.master
}

func (cb *clusterBuilder) Create() (Cluster, error) {
	cr := cmav1alpha1.Cluster{
		TypeMeta: metav1.TypeMeta{
			Kind: "Cluster",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:   cb.name,
			Labels: []string{},
		},
	}

	_, err := npb.testCase.CMAClient().ClusterV1alpha1().Clusters("").Create(&cr)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	c := &cluster{
		name:     cb.name,
		testCase: cb.testCase,
	}

	return c, nil
}

func (mb *masterBuilder) AvailabilityZone(az string) ClusterBuilder {
	mb.availabilityZone = az
	return mb.clusterBuilder
}

func (mb *masterBuilder) InstanceType(t string) ClusterBuilder {
	mb.instanceType = t
	return mb.clusterBuilder
}

func (c *cluster) AddNodePool(name string) NodePoolBuilder {
	npb := &nodePoolBuilder{
		name: name,
	}

	return npb
}

func (c *cluster) GetLiveObject() (*cmav1alpha1.Cluster, error) {
	cluster, err := c.testCase.CMAClient().ClusterV1alpha1().Clusters(c.namespace).Get(c.name, metav1.GetOptions{})
	if err != nil {
		return nil, microerror.Mask(err)
	}
	return cluster, nil
}

func (npb *nodePoolBuilder) WithAvailabilityZone(az string) NodePoolBuilder {
	npb.availabilityZones = append(npb.availabilityZones, az)
	return npb
}

func (npb *nodePoolBuilder) WithInstanceType(t string) NodePoolBuilder {
	npb.instanceType = t
	return npb
}

func (npb *nodePoolBuilder) WithMinWorkers(n int) NodePoolBuilder {
	npb.minWorkers = n
	return npb
}

func (npb *nodePoolBuilder) WithMaxWorkers(n int) NodePoolBuilder {
	npb.maxWorkers = n
	return npb
}

func (npb *nodePoolBuilder) Create() (NodePool, error) {
	cr := cmav1alpha1.MachineDeployment{
		TypeMeta: metav1.TypeMeta{
			Kind: "MachineDeployment",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:   npb.name,
			Labels: []string{},
		},
	}

	_, err := npb.testCase.CMAClient().ClusterV1alpha1().MachineDeployments("").Create(&cr)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	np := &nodePool{
		name:     npb.name,
		testCase: npb.testCase,
	}

	return np, nil
}

func (np *nodePool) GetLiveObject() (*cmav1alpha1.MachineDeployment, error) {
	md, err := np.testCase.CMAClient().ClusterV1alpha1().MachineDeployments(c.namespace).Get(np.name, metav1.GetOptions{})
	if err != nil {
		return nil, microerror.Mask(err)
	}
	return md, nil
}
