/*
Copyright 2020 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package testutils

import (
	"fmt"

	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/kops/pkg/apis/kops"
	"k8s.io/kops/upup/pkg/fi"
)

// BuildMinimalCluster a generic minimal cluster
func BuildMinimalCluster(clusterName string) *kops.Cluster {
	c := &kops.Cluster{}
	c.ObjectMeta.Name = clusterName
	c.Spec.KubernetesVersion = "1.23.2"
	c.Spec.Networking.Subnets = []kops.ClusterSubnetSpec{
		{Name: "subnet-us-test-1a", Zone: "us-test-1a", CIDR: "172.20.1.0/24", Type: kops.SubnetTypePrivate},
	}

	c.Spec.ContainerRuntime = "containerd"
	c.Spec.Containerd = &kops.ContainerdConfig{}

	c.Spec.API.PublicName = fmt.Sprintf("api.%v", clusterName)
	c.Spec.API.Access = []string{"0.0.0.0/0"}
	c.Spec.SSHAccess = []string{"0.0.0.0/0"}

	// Default to public topology
	c.Spec.Networking.Topology = &kops.TopologySpec{
		DNS: kops.DNSTypePublic,
	}

	c.Spec.Networking.NetworkCIDR = "172.20.0.0/16"
	c.Spec.Networking.Subnets = []kops.ClusterSubnetSpec{
		{Name: "subnet-us-test-1a", Zone: "us-test-1a", CIDR: "172.20.1.0/24", Type: kops.SubnetTypePublic},
		{Name: "subnet-us-test-1b", Zone: "us-test-1b", CIDR: "172.20.2.0/24", Type: kops.SubnetTypePublic},
		{Name: "subnet-us-test-1c", Zone: "us-test-1c", CIDR: "172.20.3.0/24", Type: kops.SubnetTypePublic},
	}

	c.Spec.Networking.NonMasqueradeCIDR = "100.64.0.0/10"
	c.Spec.CloudProvider.AWS = &kops.AWSSpec{}

	c.Spec.ConfigBase = "memfs://unittest-bucket/" + clusterName

	c.Spec.DNSZone = "test.com"

	c.Spec.SSHKeyName = fi.PtrTo("test")

	addEtcdClusters(c)

	return c
}

func addEtcdClusters(c *kops.Cluster) {
	subnetNames := sets.NewString()
	for _, z := range c.Spec.Networking.Subnets {
		subnetNames.Insert(z.Name)
	}
	etcdZones := subnetNames.List()

	for _, etcdCluster := range []string{"main", "events"} {
		etcd := kops.EtcdClusterSpec{}
		etcd.Name = etcdCluster
		for _, zone := range etcdZones {
			m := kops.EtcdMemberSpec{}
			m.Name = zone
			m.InstanceGroup = fi.PtrTo("master-" + zone)
			etcd.Members = append(etcd.Members, m)
		}
		c.Spec.EtcdClusters = append(c.Spec.EtcdClusters, etcd)
	}
}

func BuildMinimalNodeInstanceGroup(name string, subnets ...string) kops.InstanceGroup {
	g := kops.InstanceGroup{}
	g.ObjectMeta.Name = name
	g.Spec.Role = kops.InstanceGroupRoleNode
	g.Spec.Image = "ubuntu/images/hvm-ssd/ubuntu-focal-20.04-amd64-server-20220404"
	g.Spec.Subnets = subnets

	return g
}

func BuildMinimalBastionInstanceGroup(name string, subnets ...string) kops.InstanceGroup {
	g := kops.InstanceGroup{}
	g.ObjectMeta.Name = name
	g.Spec.Role = kops.InstanceGroupRoleNode
	g.Spec.Image = "ubuntu/images/hvm-ssd/ubuntu-focal-20.04-amd64-server-20220404"
	g.Spec.Subnets = subnets

	return g
}

func BuildMinimalMasterInstanceGroup(subnet string) kops.InstanceGroup {
	g := kops.InstanceGroup{}
	g.ObjectMeta.Name = "master-" + subnet
	g.Spec.Role = kops.InstanceGroupRoleControlPlane
	g.Spec.Subnets = []string{subnet}
	g.Spec.Image = "ami-1234abcd"

	return g
}
