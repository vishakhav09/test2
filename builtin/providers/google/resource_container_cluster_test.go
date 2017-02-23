package google

import (
	"fmt"
	"testing"

	"strconv"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccContainerCluster_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckContainerClusterDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccContainerCluster_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckContainerCluster(
						"google_container_cluster.primary"),
				),
			},
		},
	})
}

func TestAccContainerCluster_withAdditionalZones(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckContainerClusterDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccContainerCluster_withAdditionalZones,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckContainerCluster(
						"google_container_cluster.with_additional_zones"),
				),
			},
		},
	})
}

func TestAccContainerCluster_withVersion(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckContainerClusterDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccContainerCluster_withVersion,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckContainerCluster(
						"google_container_cluster.with_version"),
				),
			},
		},
	})
}

func TestAccContainerCluster_withNodeConfig(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckContainerClusterDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccContainerCluster_withNodeConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckContainerCluster(
						"google_container_cluster.with_node_config"),
				),
			},
		},
	})
}

func TestAccContainerCluster_withNodeConfigScopeAlias(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckContainerClusterDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccContainerCluster_withNodeConfigScopeAlias,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckContainerCluster(
						"google_container_cluster.with_node_config_scope_alias"),
				),
			},
		},
	})
}

func TestAccContainerCluster_network(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckContainerClusterDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccContainerCluster_networkRef,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckContainerCluster(
						"google_container_cluster.with_net_ref_by_url"),
					testAccCheckContainerCluster(
						"google_container_cluster.with_net_ref_by_name"),
				),
			},
		},
	})
}

func TestAccContainerCluster_backend(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckContainerClusterDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccContainerCluster_backendRef,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckContainerClusterExists(
						"google_container_cluster.primary"),
				),
			},
		},
	})
}

func testAccCheckContainerClusterDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "google_container_cluster" {
			continue
		}

		attributes := rs.Primary.Attributes
		_, err := config.clientContainer.Projects.Zones.Clusters.Get(
			config.Project, attributes["zone"], attributes["name"]).Do()
		if err == nil {
			return fmt.Errorf("Cluster still exists")
		}
	}

	return nil
}

func testAccCheckContainerCluster(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		attributes, err := getResourceAttributes(n, s)
		if err != nil {
			return err
		}

		config := testAccProvider.Meta().(*Config)
		cluster, err := config.clientContainer.Projects.Zones.Clusters.Get(
			config.Project, attributes["zone"], attributes["name"]).Do()
		if err != nil {
			return err
		}

		if cluster.Name != attributes["name"] {
			return fmt.Errorf("Cluster %s not found, found %s instead", attributes["name"], cluster.Name)
		}

		type clusterTestField struct {
			tf_attr  string
			gcp_attr interface{}
		}

		clusterTests := []clusterTestField{
			{"initial_node_count", strconv.FormatInt(cluster.InitialNodeCount, 10)},
			{"master_auth.0.client_certificate", cluster.MasterAuth.ClientCertificate},
			{"master_auth.0.client_key", cluster.MasterAuth.ClientKey},
			{"master_auth.0.cluster_ca_certificate", cluster.MasterAuth.ClusterCaCertificate},
			{"master_auth.0.password", cluster.MasterAuth.Password},
			{"master_auth.0.username", cluster.MasterAuth.Username},
			{"zone", cluster.Zone},
			{"cluster_ipv4_cidr", cluster.ClusterIpv4Cidr},
			{"description", cluster.Description},
			{"endpoint", cluster.Endpoint},
			{"instance_group_urls", cluster.InstanceGroupUrls},
			{"logging_service", cluster.LoggingService},
			{"monitoring_service", cluster.MonitoringService},
			// TODO(danawillow): Add this back in. Currently this field is saved via the config instead of from the API response,
			// and the config may contain the network name or self_link, whereas the API only returns the self_link.
			// {"network", cluster.Network},
			{"subnetwork", cluster.Subnetwork},
			{"node_config.0.machine_type", cluster.NodeConfig.MachineType},
			{"node_config.0.disk_size_gb", strconv.FormatInt(cluster.NodeConfig.DiskSizeGb, 10)},
			{"node_config.0.oauth_scopes", cluster.NodeConfig.OauthScopes},
			{"node_version", cluster.CurrentNodeVersion},
		}

		// Remove Zone from additional_zones since that's what the resource writes in state
		additionalZones := []string{}
		for _, location := range cluster.Locations {
			if location != cluster.Zone {
				additionalZones = append(additionalZones, location)
			}
		}
		clusterTests = append(clusterTests, clusterTestField{"additional_zones", additionalZones})

		// AddonsConfig is neither Required or Computed, so the API may return nil for it
		if cluster.AddonsConfig != nil {
			if cluster.AddonsConfig.HttpLoadBalancing != nil {
				clusterTests = append(clusterTests, clusterTestField{"addons_config.0.http_load_balancing.0.disabled", strconv.FormatBool(cluster.AddonsConfig.HttpLoadBalancing.Disabled)})
			}
			if cluster.AddonsConfig.HorizontalPodAutoscaling != nil {
				clusterTests = append(clusterTests, clusterTestField{"addons_config.0.horizontal_pod_autoscaling.0.disabled", strconv.FormatBool(cluster.AddonsConfig.HorizontalPodAutoscaling.Disabled)})
			}
		}

		for _, attrs := range clusterTests {
			if c := checkMatch(attributes, attrs.tf_attr, attrs.gcp_attr); c != "" {
				return fmt.Errorf(c)
			}
		}

		return nil
	}
}

func getResourceAttributes(n string, s *terraform.State) (map[string]string, error) {
	rs, ok := s.RootModule().Resources[n]
	if !ok {
		return nil, fmt.Errorf("Not found: %s", n)
	}

	if rs.Primary.ID == "" {
		return nil, fmt.Errorf("No ID is set")
	}

	return rs.Primary.Attributes, nil
}

func checkMatch(attributes map[string]string, attr string, gcp interface{}) string {
	if gcpList, ok := gcp.([]string); ok {
		return checkListMatch(attributes, attr, gcpList)
	}
	tf := attributes[attr]
	if tf != gcp {
		return fmt.Sprintf("Cluster has mismatched %s.\nTF State: %+v\nGCP State: %+v", attr, tf, gcp)
	}
	return ""
}

func checkListMatch(attributes map[string]string, attr string, gcpList []string) string {
	num, err := strconv.Atoi(attributes[attr+".#"])
	if err != nil {
		return fmt.Sprintf("error in number conversion for attribute %s", attr)
	}
	if num != len(gcpList) {
		return fmt.Sprintf("Cluster has mismatched %s size.\nTF Size: %d\nGCP Size: %d", attr, num, len(gcpList))
	}

	for i, gcp := range gcpList {
		if tf := attributes[fmt.Sprintf("%s.%d", attr, i)]; tf != gcp {
			return fmt.Sprintf("Cluster has mismatched %s[%d].\nTF State: %+v\nGCP State: %+v", attr, i, tf, gcp)
		}
	}

	return ""
}

var testAccContainerCluster_basic = fmt.Sprintf(`
resource "google_container_cluster" "primary" {
	name = "cluster-test-%s"
	zone = "us-central1-a"
	initial_node_count = 3

	master_auth {
		username = "mr.yoda"
		password = "adoy.rm"
	}
}`, acctest.RandString(10))

var testAccContainerCluster_withAdditionalZones = fmt.Sprintf(`
resource "google_container_cluster" "with_additional_zones" {
	name = "cluster-test-%s"
	zone = "us-central1-a"
	initial_node_count = 1

	additional_zones = [
		"us-central1-b",
		"us-central1-c"
	]

	master_auth {
		username = "mr.yoda"
		password = "adoy.rm"
	}
}`, acctest.RandString(10))

var testAccContainerCluster_withVersion = fmt.Sprintf(`
resource "google_container_cluster" "with_version" {
	name = "cluster-test-%s"
	zone = "us-central1-a"
	node_version = "1.5.2"
	initial_node_count = 1

	master_auth {
		username = "mr.yoda"
		password = "adoy.rm"
	}
}`, acctest.RandString(10))

var testAccContainerCluster_withNodeConfig = fmt.Sprintf(`
resource "google_container_cluster" "with_node_config" {
	name = "cluster-test-%s"
	zone = "us-central1-f"
	initial_node_count = 1

	master_auth {
		username = "mr.yoda"
		password = "adoy.rm"
	}

	node_config {
		machine_type = "g1-small"
		disk_size_gb = 15
		oauth_scopes = [
			"https://www.googleapis.com/auth/compute",
			"https://www.googleapis.com/auth/devstorage.read_only",
			"https://www.googleapis.com/auth/logging.write",
			"https://www.googleapis.com/auth/monitoring"
		]
	}
}`, acctest.RandString(10))

var testAccContainerCluster_withNodeConfigScopeAlias = fmt.Sprintf(`
resource "google_container_cluster" "with_node_config_scope_alias" {
	name = "cluster-test-%s"
	zone = "us-central1-f"
	initial_node_count = 1

	master_auth {
		username = "mr.yoda"
		password = "adoy.rm"
	}

	node_config {
		machine_type = "g1-small"
		disk_size_gb = 15
		oauth_scopes = [ "compute-rw", "storage-ro", "logging-write", "monitoring" ]
	}
}`, acctest.RandString(10))

var testAccContainerCluster_networkRef = fmt.Sprintf(`
resource "google_compute_network" "container_network" {
	name = "container-net-%s"
	auto_create_subnetworks = true
}

resource "google_container_cluster" "with_net_ref_by_url" {
	name = "cluster-test-%s"
	zone = "us-central1-a"
	initial_node_count = 1

	master_auth {
		username = "mr.yoda"
		password = "adoy.rm"
	}

	network = "${google_compute_network.container_network.self_link}"
}

resource "google_container_cluster" "with_net_ref_by_name" {
	name = "cluster-test-%s"
	zone = "us-central1-a"
	initial_node_count = 1

	master_auth {
		username = "mr.yoda"
		password = "adoy.rm"
	}

	network = "${google_compute_network.container_network.name}"
}`, acctest.RandString(10), acctest.RandString(10), acctest.RandString(10))

var testAccContainerCluster_backendRef = fmt.Sprintf(`
resource "google_compute_backend_service" "my-backend-service" {
  name      = "terraform-test-%s"
  port_name = "http"
  protocol  = "HTTP"

  backend {
    group = "${element(google_container_cluster.primary.instance_group_urls, 1)}"
  }

  health_checks = ["${google_compute_http_health_check.default.self_link}"]
}

resource "google_compute_http_health_check" "default" {
  name               = "terraform-test-%s"
  request_path       = "/"
  check_interval_sec = 1
  timeout_sec        = 1
}

resource "google_container_cluster" "primary" {
  name               = "terraform-test-%s"
  zone               = "us-central1-a"
  initial_node_count = 3

  additional_zones = [
    "us-central1-b",
    "us-central1-c",
  ]

  master_auth {
    username = "mr.yoda"
    password = "adoy.rm"
  }

  node_config {
    oauth_scopes = [
      "https://www.googleapis.com/auth/compute",
      "https://www.googleapis.com/auth/devstorage.read_only",
      "https://www.googleapis.com/auth/logging.write",
      "https://www.googleapis.com/auth/monitoring",
    ]
  }
}
`, acctest.RandString(10), acctest.RandString(10), acctest.RandString(10))
