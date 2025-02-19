package scaleway

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/scaleway/scaleway-sdk-go/api/k8s/v1"
)

func TestAccScalewayK8SCluster_PoolBasic(t *testing.T) {
	tt := NewTestTools(t)
	defer tt.Cleanup()

	latestK8SVersion := testAccScalewayK8SClusterGetLatestK8SVersion(tt)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      testAccCheckScalewayK8SClusterDestroy(tt),
		Steps: []resource.TestStep{
			{
				Config: testAccCheckScalewayK8SPoolConfigMinimal(latestK8SVersion, false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckScalewayK8SClusterExists(tt, "scaleway_k8s_cluster.minimal"),
					testAccCheckScalewayK8SPoolExists(tt, "scaleway_k8s_pool.default"),
					resource.TestCheckResourceAttr("scaleway_k8s_pool.default", "node_type", "gp1_xs"),
					resource.TestCheckResourceAttr("scaleway_k8s_pool.default", "size", "1"),
					resource.TestCheckResourceAttr("scaleway_k8s_pool.default", "autohealing", "true"),
					resource.TestCheckResourceAttr("scaleway_k8s_pool.default", "autoscaling", "true"),
					resource.TestCheckResourceAttr("scaleway_k8s_pool.default", "version", latestK8SVersion),
					resource.TestCheckResourceAttrSet("scaleway_k8s_pool.default", "id"),
					resource.TestCheckResourceAttr("scaleway_k8s_pool.default", "tags.0", "terraform-test"),
					resource.TestCheckResourceAttr("scaleway_k8s_pool.default", "tags.1", "scaleway_k8s_cluster"),
					resource.TestCheckResourceAttr("scaleway_k8s_pool.default", "tags.2", "default"),
				),
			},
			{
				Config: testAccCheckScalewayK8SPoolConfigMinimal(latestK8SVersion, true),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckScalewayK8SClusterExists(tt, "scaleway_k8s_cluster.minimal"),
					testAccCheckScalewayK8SPoolExists(tt, "scaleway_k8s_pool.minimal"),
					resource.TestCheckResourceAttr("scaleway_k8s_pool.minimal", "node_type", "gp1_xs"),
					resource.TestCheckResourceAttr("scaleway_k8s_pool.minimal", "size", "1"),
					resource.TestCheckResourceAttr("scaleway_k8s_pool.minimal", "autohealing", "true"),
					resource.TestCheckResourceAttr("scaleway_k8s_pool.minimal", "autoscaling", "true"),
					resource.TestCheckResourceAttr("scaleway_k8s_pool.minimal", "version", latestK8SVersion),
					resource.TestCheckResourceAttrSet("scaleway_k8s_pool.minimal", "id"),
					resource.TestCheckResourceAttr("scaleway_k8s_pool.minimal", "tags.0", "terraform-test"),
					resource.TestCheckResourceAttr("scaleway_k8s_pool.minimal", "tags.1", "scaleway_k8s_cluster"),
					resource.TestCheckResourceAttr("scaleway_k8s_pool.minimal", "tags.2", "minimal"),
				),
			},
			{
				Config: testAccCheckScalewayK8SPoolConfigMinimal(latestK8SVersion, false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckScalewayK8SClusterExists(tt, "scaleway_k8s_cluster.minimal"),
					testAccCheckScalewayK8SPoolDestroy(tt, "scaleway_k8s_pool.minimal"),
				),
			},
		},
	})
}

func TestAccScalewayK8SCluster_PoolWait(t *testing.T) {
	tt := NewTestTools(t)
	defer tt.Cleanup()
	latestK8SVersion := testAccScalewayK8SClusterGetLatestK8SVersion(tt)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      testAccCheckScalewayK8SClusterDestroy(tt),
		Steps: []resource.TestStep{
			{
				Config: testAccCheckScalewayK8SPoolConfigWait(latestK8SVersion, false, 0),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckScalewayK8SClusterExists(tt, "scaleway_k8s_cluster.minimal"),
					testAccCheckScalewayK8SPoolExists(tt, "scaleway_k8s_pool.default"),
					resource.TestCheckResourceAttr("scaleway_k8s_pool.default", "size", "1"),
					resource.TestCheckResourceAttr("scaleway_k8s_pool.default", "min_size", "1"),
					resource.TestCheckResourceAttr("scaleway_k8s_pool.default", "max_size", "1"),
					resource.TestCheckResourceAttr("scaleway_k8s_pool.default", "status", k8s.PoolStatusReady.String()),
					resource.TestCheckResourceAttr("scaleway_k8s_pool.default", "nodes.0.status", k8s.NodeStatusReady.String()),
					resource.TestCheckResourceAttr("scaleway_k8s_pool.default", "wait_for_pool_ready", "true"),
				),
			},
			{
				Config: testAccCheckScalewayK8SPoolConfigWait(latestK8SVersion, true, 1),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckScalewayK8SClusterExists(tt, "scaleway_k8s_cluster.minimal"),
					testAccCheckScalewayK8SPoolExists(tt, "scaleway_k8s_pool.minimal"),
					resource.TestCheckResourceAttr("scaleway_k8s_pool.minimal", "size", "1"),
					resource.TestCheckResourceAttr("scaleway_k8s_pool.minimal", "min_size", "1"),
					resource.TestCheckResourceAttr("scaleway_k8s_pool.minimal", "max_size", "1"),
					resource.TestCheckResourceAttr("scaleway_k8s_pool.minimal", "status", k8s.PoolStatusReady.String()),
					resource.TestCheckResourceAttr("scaleway_k8s_pool.minimal", "nodes.0.status", k8s.NodeStatusReady.String()),
					resource.TestCheckResourceAttr("scaleway_k8s_pool.minimal", "wait_for_pool_ready", "true"),
				),
			},
			{
				Config: testAccCheckScalewayK8SPoolConfigWait(latestK8SVersion, true, 2),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckScalewayK8SClusterExists(tt, "scaleway_k8s_cluster.minimal"),
					testAccCheckScalewayK8SPoolExists(tt, "scaleway_k8s_pool.minimal"),
					resource.TestCheckResourceAttr("scaleway_k8s_pool.minimal", "size", "2"),
					resource.TestCheckResourceAttr("scaleway_k8s_pool.minimal", "min_size", "1"),
					resource.TestCheckResourceAttr("scaleway_k8s_pool.minimal", "max_size", "2"),
					resource.TestCheckResourceAttr("scaleway_k8s_pool.minimal", "status", k8s.PoolStatusReady.String()),
					resource.TestCheckResourceAttr("scaleway_k8s_pool.minimal", "nodes.0.status", k8s.NodeStatusReady.String()),
					resource.TestCheckResourceAttr("scaleway_k8s_pool.minimal", "nodes.1.status", k8s.NodeStatusReady.String()), // check that the new node has the "ready" status
					resource.TestCheckResourceAttr("scaleway_k8s_pool.minimal", "wait_for_pool_ready", "true"),
				),
			},
			{
				Config: testAccCheckScalewayK8SPoolConfigWait(latestK8SVersion, true, 1),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckScalewayK8SClusterExists(tt, "scaleway_k8s_cluster.minimal"),
					testAccCheckScalewayK8SPoolExists(tt, "scaleway_k8s_pool.minimal"),
					resource.TestCheckResourceAttr("scaleway_k8s_pool.minimal", "size", "1"),
					resource.TestCheckResourceAttr("scaleway_k8s_pool.minimal", "min_size", "1"),
					resource.TestCheckResourceAttr("scaleway_k8s_pool.minimal", "max_size", "1"),
					resource.TestCheckResourceAttr("scaleway_k8s_pool.minimal", "status", k8s.PoolStatusReady.String()),
					testAccCheckScalewayK8SPoolNodesOneOfIsDeleting("scaleway_k8s_pool.minimal"), // check that one of the nodes is deleting (nodes are not ordered)
					resource.TestCheckResourceAttr("scaleway_k8s_pool.minimal", "nodes.#", "2"),  // the node that is deleting should still exist
					resource.TestCheckResourceAttr("scaleway_k8s_pool.minimal", "wait_for_pool_ready", "true"),
				),
			},
			{
				Config: testAccCheckScalewayK8SPoolConfigWait(latestK8SVersion, false, 0),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckScalewayK8SClusterExists(tt, "scaleway_k8s_cluster.minimal"),
					testAccCheckScalewayK8SPoolDestroy(tt, "scaleway_k8s_pool.minimal"),
				),
			},
		},
	})
}

func TestAccScalewayK8SCluster_PoolPlacementGroup(t *testing.T) {
	tt := NewTestTools(t)
	defer tt.Cleanup()
	latestK8SVersion := testAccScalewayK8SClusterGetLatestK8SVersion(tt)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      testAccCheckScalewayK8SClusterDestroy(tt),
		Steps: []resource.TestStep{
			{
				Config: testAccCheckScalewayK8SPoolConfigPlacementGroup(latestK8SVersion),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckScalewayK8SClusterExists(tt, "scaleway_k8s_cluster.placement_group"),
					testAccCheckScalewayK8SPoolExists(tt, "scaleway_k8s_pool.placement_group"),
					resource.TestCheckResourceAttr("scaleway_k8s_pool.placement_group", "node_type", "gp1_xs"),
					resource.TestCheckResourceAttr("scaleway_k8s_pool.placement_group", "size", "1"),
					resource.TestCheckResourceAttrSet("scaleway_k8s_pool.placement_group", "id"),
					resource.TestCheckResourceAttrSet("scaleway_k8s_pool.placement_group", "placement_group_id"),
				),
			},
			{
				Config: testAccCheckScalewayK8SPoolConfigPlacementGroupWithCustomZone(latestK8SVersion),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckScalewayK8SClusterExists(tt, "scaleway_k8s_cluster.cluster"),
					testAccCheckScalewayK8SPoolExists(tt, "scaleway_k8s_pool.pool"),
					resource.TestCheckResourceAttrPair("scaleway_k8s_pool.pool", "placement_group_id", "scaleway_instance_placement_group.placement_group", "id"),
					resource.TestCheckResourceAttr("scaleway_k8s_pool.pool", "zone", "nl-ams-2"),
					resource.TestCheckResourceAttr("scaleway_k8s_pool.pool", "node_type", "gp1_xs"),
					resource.TestCheckResourceAttr("scaleway_k8s_pool.pool", "size", "1"),
					resource.TestCheckResourceAttrSet("scaleway_k8s_pool.pool", "id"),
					resource.TestCheckResourceAttrSet("scaleway_k8s_pool.pool", "placement_group_id"),
				),
			},
			{
				Config: testAccCheckScalewayK8SPoolConfigPlacementGroupWithMultiZone(latestK8SVersion),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckScalewayK8SClusterExists(tt, "scaleway_k8s_cluster.cluster"),
					testAccCheckScalewayK8SPoolExists(tt, "scaleway_k8s_pool.pool"),
					resource.TestCheckResourceAttrPair("scaleway_k8s_pool.pool", "placement_group_id", "scaleway_instance_placement_group.placement_group", "id"),
					resource.TestCheckResourceAttr("scaleway_k8s_pool.pool", "zone", "nl-ams-1"),
					resource.TestCheckResourceAttr("scaleway_k8s_pool.pool", "node_type", "gp1_xs"),
					resource.TestCheckResourceAttr("scaleway_k8s_pool.pool", "size", "1"),
					resource.TestCheckResourceAttrSet("scaleway_k8s_pool.pool", "id"),
					resource.TestCheckResourceAttrSet("scaleway_k8s_pool.pool", "placement_group_id"),
				),
			},
		},
	})
}

func TestAccScalewayK8SCluster_PoolUpgradePolicy(t *testing.T) {
	tt := NewTestTools(t)
	defer tt.Cleanup()

	latestK8SVersion := testAccScalewayK8SClusterGetLatestK8SVersion(tt)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      testAccCheckScalewayK8SClusterDestroy(tt),
		Steps: []resource.TestStep{
			{
				Config: testAccCheckScalewayK8SPoolConfigUpgradePolicy(latestK8SVersion, 2, 3),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckScalewayK8SClusterExists(tt, "scaleway_k8s_cluster.upgrade_policy"),
					testAccCheckScalewayK8SPoolExists(tt, "scaleway_k8s_pool.default"),
					resource.TestCheckResourceAttr("scaleway_k8s_pool.default", "node_type", "gp1_xs"),
					resource.TestCheckResourceAttr("scaleway_k8s_pool.default", "size", "1"),
					resource.TestCheckResourceAttr("scaleway_k8s_pool.default", "autohealing", "true"),
					resource.TestCheckResourceAttr("scaleway_k8s_pool.default", "autoscaling", "true"),
					resource.TestCheckResourceAttr("scaleway_k8s_pool.default", "version", latestK8SVersion),
					resource.TestCheckResourceAttrSet("scaleway_k8s_pool.default", "id"),
					resource.TestCheckResourceAttr("scaleway_k8s_pool.default", "tags.0", "terraform-test"),
					resource.TestCheckResourceAttr("scaleway_k8s_pool.default", "tags.1", "scaleway_k8s_cluster"),
					resource.TestCheckResourceAttr("scaleway_k8s_pool.default", "tags.2", "upgrade_policy"),
					resource.TestCheckResourceAttr("scaleway_k8s_pool.default", "upgrade_policy.0.max_surge", "2"),
					resource.TestCheckResourceAttr("scaleway_k8s_pool.default", "upgrade_policy.0.max_unavailable", "3"),
				),
			},
			{
				Config: testAccCheckScalewayK8SPoolConfigUpgradePolicy(latestK8SVersion, 0, 1),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckScalewayK8SClusterExists(tt, "scaleway_k8s_cluster.upgrade_policy"),
					testAccCheckScalewayK8SPoolExists(tt, "scaleway_k8s_pool.default"),
					resource.TestCheckResourceAttr("scaleway_k8s_pool.default", "node_type", "gp1_xs"),
					resource.TestCheckResourceAttr("scaleway_k8s_pool.default", "size", "1"),
					resource.TestCheckResourceAttr("scaleway_k8s_pool.default", "autohealing", "true"),
					resource.TestCheckResourceAttr("scaleway_k8s_pool.default", "autoscaling", "true"),
					resource.TestCheckResourceAttr("scaleway_k8s_pool.default", "version", latestK8SVersion),
					resource.TestCheckResourceAttrSet("scaleway_k8s_pool.default", "id"),
					resource.TestCheckResourceAttr("scaleway_k8s_pool.default", "tags.0", "terraform-test"),
					resource.TestCheckResourceAttr("scaleway_k8s_pool.default", "tags.1", "scaleway_k8s_cluster"),
					resource.TestCheckResourceAttr("scaleway_k8s_pool.default", "tags.2", "upgrade_policy"),
					resource.TestCheckResourceAttr("scaleway_k8s_pool.default", "upgrade_policy.0.max_surge", "0"),
					resource.TestCheckResourceAttr("scaleway_k8s_pool.default", "upgrade_policy.0.max_unavailable", "1"),
				),
			},
		},
	})
}

func TestAccScalewayK8SCluster_PoolKubeletArgs(t *testing.T) {
	tt := NewTestTools(t)
	defer tt.Cleanup()

	latestK8SVersion := testAccScalewayK8SClusterGetLatestK8SVersion(tt)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      testAccCheckScalewayK8SClusterDestroy(tt),
		Steps: []resource.TestStep{
			{
				Config: testAccCheckScalewayK8SPoolConfigKubeletArgs(latestK8SVersion, 1337),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckScalewayK8SClusterExists(tt, "scaleway_k8s_cluster.kubelet_args"),
					testAccCheckScalewayK8SPoolExists(tt, "scaleway_k8s_pool.default"),
					resource.TestCheckResourceAttr("scaleway_k8s_pool.default", "version", latestK8SVersion),
					resource.TestCheckResourceAttrSet("scaleway_k8s_pool.default", "id"),
					resource.TestCheckResourceAttr("scaleway_k8s_pool.default", "kubelet_args.maxPods", "1337"),
				),
			},
			{
				Config: testAccCheckScalewayK8SPoolConfigKubeletArgs(latestK8SVersion, 50),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckScalewayK8SClusterExists(tt, "scaleway_k8s_cluster.kubelet_args"),
					testAccCheckScalewayK8SPoolExists(tt, "scaleway_k8s_pool.default"),
					resource.TestCheckResourceAttr("scaleway_k8s_pool.default", "version", latestK8SVersion),
					resource.TestCheckResourceAttrSet("scaleway_k8s_pool.default", "id"),
					resource.TestCheckResourceAttr("scaleway_k8s_pool.default", "kubelet_args.maxPods", "50"),
				),
			},
		},
	})
}

func TestAccScalewayK8SCluster_PoolZone(t *testing.T) {
	tt := NewTestTools(t)
	defer tt.Cleanup()

	latestK8SVersion := testAccScalewayK8SClusterGetLatestK8SVersion(tt)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      testAccCheckScalewayK8SClusterDestroy(tt),
		Steps: []resource.TestStep{
			{
				Config: testAccCheckScalewayK8SPoolConfigZone(latestK8SVersion, "fr-par-2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckScalewayK8SClusterExists(tt, "scaleway_k8s_cluster.zone"),
					testAccCheckScalewayK8SPoolExists(tt, "scaleway_k8s_pool.default"),
					resource.TestCheckResourceAttr("scaleway_k8s_pool.default", "version", latestK8SVersion),
					resource.TestCheckResourceAttrSet("scaleway_k8s_pool.default", "id"),
					resource.TestCheckResourceAttr("scaleway_k8s_pool.default", "zone", "fr-par-2"),
				),
			},
		},
	})
}

func TestAccScalewayK8SCluster_PoolSize(t *testing.T) {
	tt := NewTestTools(t)
	defer tt.Cleanup()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      testAccCheckScalewayK8SClusterDestroy(tt),
		Steps: []resource.TestStep{
			{
				Config: `
				resource "scaleway_k8s_cluster" "cluster" {
				  name    = "cluster"
				  version = "1.23"
				  cni     = "cilium"
				  delete_additional_resources = true
				  auto_upgrade {
				    enable = true
				    maintenance_window_start_hour = 12
				    maintenance_window_day = "monday"
				  }
				}
				
				resource "scaleway_k8s_pool" "pool" {
				  cluster_id          = scaleway_k8s_cluster.cluster.id
				  name                = "pool"
				  node_type           = "gp1_xs"
				  size                = 1
				  autoscaling         = false
				  autohealing         = true
				  wait_for_pool_ready = true
				}`,
			},
			{
				Config: `
				resource "scaleway_k8s_cluster" "cluster" {
				  name    = "cluster"
				  version = "1.23"
				  cni     = "cilium"
				  auto_upgrade {
				  enable = true
				  maintenance_window_start_hour = 12
				  maintenance_window_day = "monday"
				  }
				  delete_additional_resources = true
				}
				
				resource "scaleway_k8s_pool" "pool" {
				  cluster_id          = scaleway_k8s_cluster.cluster.id
				  name                = "pool"
				  node_type           = "gp1_xs"
				  size                = 2
				  autoscaling         = false
				  autohealing         = true
				  wait_for_pool_ready = true
				}`,
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccCheckScalewayK8SPoolDestroy(tt *TestTools, n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return nil
		}

		k8sAPI, region, poolID, err := k8sAPIWithRegionAndID(tt.Meta, rs.Primary.ID)
		if err != nil {
			return err
		}

		_, err = k8sAPI.GetPool(&k8s.GetPoolRequest{
			Region: region,
			PoolID: poolID,
		})
		// If no error resource still exist
		if err == nil {
			return fmt.Errorf("pool (%s) still exists", rs.Primary.ID)
		}
		// Unexpected api error we return it
		if !is404Error(err) {
			return err
		}

		return nil
	}
}

func testAccCheckScalewayK8SPoolExists(tt *TestTools, n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("resource not found: %s", n)
		}

		k8sAPI, region, poolID, err := k8sAPIWithRegionAndID(tt.Meta, rs.Primary.ID)
		if err != nil {
			return err
		}

		_, err = k8sAPI.GetPool(&k8s.GetPoolRequest{
			Region: region,
			PoolID: poolID,
		})
		if err != nil {
			return err
		}

		return nil
	}
}

func testAccCheckScalewayK8SPoolConfigMinimal(version string, otherPool bool) string {
	pool := ""
	if otherPool {
		pool += `
resource "scaleway_k8s_pool" "minimal" {
    name = "minimal"
	cluster_id = "${scaleway_k8s_cluster.minimal.id}"
	node_type = "gp1_xs"
	autohealing = true
	autoscaling = true
	size = 1
	tags = [ "terraform-test", "scaleway_k8s_cluster", "minimal" ]
}`
	}

	return fmt.Sprintf(`
resource "scaleway_k8s_pool" "default" {
    name = "default"
	cluster_id = "${scaleway_k8s_cluster.minimal.id}"
	node_type = "gp1_xs"
	autohealing = true
	autoscaling = true
	size = 1
	tags = [ "terraform-test", "scaleway_k8s_cluster", "default" ]
}
resource "scaleway_k8s_cluster" "minimal" {
    name = "K8SPoolConfigMinimal"
	cni = "calico"
	version = "%s"
	tags = [ "terraform-test", "scaleway_k8s_cluster", "minimal" ]
	delete_additional_resources = true
}
%s`, version, pool)
}

func testAccCheckScalewayK8SPoolConfigWait(version string, otherPool bool, otherPoolSize int) string {
	pool := ""
	if otherPool {
		pool += fmt.Sprintf(`
resource "scaleway_k8s_pool" "minimal" {
    name = "minimal"
	cluster_id = scaleway_k8s_cluster.minimal.id
	node_type = "gp1_xs"
	size = %d
	min_size = 1
	max_size = %d

	wait_for_pool_ready = true
}`, otherPoolSize, otherPoolSize)
	}

	return fmt.Sprintf(`
resource "scaleway_k8s_pool" "default" {
    name = "default"
	cluster_id = scaleway_k8s_cluster.minimal.id
	node_type = "gp1_xs"
	size = 1
	min_size = 1
	max_size = 1
	wait_for_pool_ready = true
}

resource "scaleway_k8s_cluster" "minimal" {
    name = "PoolConfigWait"
	cni = "calico"
	version = "%s"
	tags = [ "terraform-test", "scaleway_k8s_cluster", "minimal" ]
	delete_additional_resources = true
}
%s`, version, pool)
}

func testAccCheckScalewayK8SPoolConfigPlacementGroup(version string) string {
	return fmt.Sprintf(`
resource "scaleway_instance_placement_group" "placement_group" {
  name        = "pool-placement-group"
  policy_type = "max_availability"
  policy_mode = "optional"
}

resource "scaleway_k8s_pool" "placement_group" {
    name               = "placement_group"
    cluster_id         = scaleway_k8s_cluster.placement_group.id
    node_type          = "gp1_xs"
    placement_group_id = scaleway_instance_placement_group.placement_group.id
    size               = 1
}

resource "scaleway_k8s_cluster" "placement_group" {
  name	  = "placement_group"
  cni	  = "calico"
  version = "%s"
  tags	  = [ "terraform-test", "scaleway_k8s_cluster", "placement_group" ]
  delete_additional_resources = true
}`, version)
}

func testAccCheckScalewayK8SPoolConfigPlacementGroupWithCustomZone(version string) string {
	return fmt.Sprintf(`
resource "scaleway_instance_placement_group" "placement_group" {
  name        = "pool-placement-group"
  policy_type = "max_availability"
  policy_mode = "optional"
  zone        = "nl-ams-2"
}

resource "scaleway_k8s_pool" "pool" {
  name               = "placement_group"
  cluster_id         = scaleway_k8s_cluster.cluster.id
  node_type          = "gp1_xs"
  placement_group_id = scaleway_instance_placement_group.placement_group.id
  size               = 1
  region             = scaleway_k8s_cluster.cluster.region
  zone               = scaleway_instance_placement_group.placement_group.zone
}

resource "scaleway_k8s_cluster" "cluster" {
  name	  = "placement_group"
  cni	  = "calico"
  version = "%s"
  tags	  = [ "terraform-test", "scaleway_k8s_cluster", "placement_group" ]
  region  = "nl-ams"
  delete_additional_resources = true
}`, version)
}

func testAccCheckScalewayK8SPoolConfigPlacementGroupWithMultiZone(version string) string {
	return fmt.Sprintf(`
resource "scaleway_instance_placement_group" "placement_group" {
  name        = "pool-placement-group"
  policy_type = "max_availability"
  policy_mode = "optional"
  zone        = "nl-ams-1"
}

resource "scaleway_k8s_pool" "pool" {
  name               = "placement_group"
  cluster_id         = scaleway_k8s_cluster.cluster.id
  node_type          = "gp1_xs"
  placement_group_id = scaleway_instance_placement_group.placement_group.id
  size               = 1
  region             = scaleway_k8s_cluster.cluster.region
  zone               = scaleway_instance_placement_group.placement_group.zone
}

resource "scaleway_k8s_cluster" "cluster" {
  name		= "placement_group"
  cni		= "kilo"
  version	= "%s"
  tags		= [ "terraform-test", "scaleway_k8s_cluster", "placement_group" ]
  region	= "fr-par"
  type		= "multicloud"
  delete_additional_resources = true
}`, version)
}

func testAccCheckScalewayK8SPoolConfigUpgradePolicy(version string, maxSurge, maxUnavailable int) string {
	return fmt.Sprintf(`
resource "scaleway_k8s_pool" "default" {
    name = "default"
	cluster_id = "${scaleway_k8s_cluster.upgrade_policy.id}"
	node_type = "gp1_xs"
	autohealing = true
	autoscaling = true
	size = 1
	tags = [ "terraform-test", "scaleway_k8s_cluster", "upgrade_policy" ]
	upgrade_policy {
		max_surge = %d
		max_unavailable = %d
	}
}
resource "scaleway_k8s_cluster" "upgrade_policy" {
    name = "K8SPoolConfigUpgradePolicy"
	cni = "cilium"
	version = "%s"
	tags = [ "terraform-test", "scaleway_k8s_cluster", "upgrade_policy" ]
	delete_additional_resources = true
}`, maxSurge, maxUnavailable, version)
}

func testAccCheckScalewayK8SPoolConfigKubeletArgs(version string, maxPods int) string {
	return fmt.Sprintf(`
resource "scaleway_k8s_pool" "default" {
    name = "default"
	cluster_id = "${scaleway_k8s_cluster.kubelet_args.id}"
	node_type = "gp1_xs"
	autohealing = true
	autoscaling = true
	size = 1
	tags = [ "terraform-test", "scaleway_k8s_cluster", "kubelet_args" ]
	kubelet_args = {
		maxPods = %d
	}
}
resource "scaleway_k8s_cluster" "kubelet_args" {
    name = "K8SPoolConfigKubeletArgs"
	cni = "cilium"
	version = "%s"
	tags = [ "terraform-test", "scaleway_k8s_cluster", "kubelet_args" ]
	delete_additional_resources = true
}`, maxPods, version)
}

func testAccCheckScalewayK8SPoolConfigZone(version string, zone string) string {
	return fmt.Sprintf(`
resource "scaleway_k8s_pool" "default" {
    name = "default"
	cluster_id = "${scaleway_k8s_cluster.zone.id}"
	node_type = "gp1_xs"
	autohealing = true
	autoscaling = true
	size = 1
	tags = [ "terraform-test", "scaleway_k8s_cluster", "zone" ]
	zone = "%s"
}
resource "scaleway_k8s_cluster" "zone" {
    name = "K8SPoolConfigZone"
	cni = "cilium"
	version = "%s"
	tags = [ "terraform-test", "scaleway_k8s_cluster", "zone" ]
	delete_additional_resources = true
}`, zone, version)
}

func testAccCheckScalewayK8SPoolNodesOneOfIsDeleting(name string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		rs, ok := state.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("resource not found: %s", name)
		}
		nodesZeroStatus, ok := rs.Primary.Attributes["nodes.0.status"]
		if !ok {
			return fmt.Errorf("attribute \"nodes.0.status\" was not set")
		}
		nodesOneStatus, ok := rs.Primary.Attributes["nodes.1.status"]
		if !ok {
			return fmt.Errorf("attribute \"nodes.1.status\" was not set")
		}
		if nodesZeroStatus == "ready" && nodesOneStatus == "deleting" ||
			nodesZeroStatus == "deleting" && nodesOneStatus == "ready" {
			return nil
		}
		return fmt.Errorf("nodes status were not as expected: got %q for nodes.0 and %q for nodes.1", nodesZeroStatus, nodesOneStatus)
	}
}
