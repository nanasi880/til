resource "google_container_cluster" "this" {
  name           = var.name
  location       = var.location
  node_locations = var.node_locations
  release_channel {
    channel = var.release_channel
  }
  min_master_version = var.min_master_version
  master_auth {
    username = ""
    password = ""
    client_certificate_config {
      issue_client_certificate = false
    }
  }

  # Automation
  enable_autopilot = true
  vertical_pod_autoscaling {
    enabled = true
  }

  # Networking
  network         = var.network_id
  subnetwork      = var.subnet_id
  networking_mode = "VPC_NATIVE"
  ip_allocation_policy {
    cluster_ipv4_cidr_block  = var.pod_ipv4_cidr
    services_ipv4_cidr_block = var.service_ipv4_cidr
  }
}
