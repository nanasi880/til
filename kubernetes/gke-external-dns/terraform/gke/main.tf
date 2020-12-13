module "terraform_remote_sate" {
  source = "../module/terraform_remote_state"
  vpc    = true
}

resource "google_container_cluster" "gke" {
  provider = google-beta

  # general
  name       = local.project
  location   = local.zone
  network    = module.terraform_remote_sate.vpc.outputs.vpc.id
  subnetwork = module.terraform_remote_sate.vpc.outputs.subnets.asia_northeast1.id

  # default node pool
  remove_default_node_pool = true
  node_locations           = [local.zone]
  initial_node_count       = 1
  node_config {
    preemptible  = true
    disk_size_gb = 10
    machine_type = "n1-standard-1"
    image_type   = "cos_containerd"
  }

  # networking mode
  networking_mode = "VPC_NATIVE"
  ip_allocation_policy {
    cluster_ipv4_cidr_block  = ""
    services_ipv4_cidr_block = ""
  }

  # cluster version
  # gcloud container get-server-config --zone asia-northeast1-a
  release_channel {
    channel = "RAPID"
  }
  min_master_version = "1.18.12-gke.1200"
  node_version       = "1.18.12-gke.1200"

  # cluster auth
  master_auth {
    username = ""
    password = ""
    client_certificate_config {
      issue_client_certificate = false
    }
  }

  # cluster autoscaling
  cluster_autoscaling {
    enabled = false
  }

  # cluster addons
  addons_config {
    dns_cache_config {
      enabled = true
    }
  }
}

resource "google_container_node_pool" "default" {
  cluster    = google_container_cluster.gke.name
  name       = "default"
  location   = local.zone
  node_count = 1
  version    = google_container_cluster.gke.master_version

  node_config {
    preemptible  = true
    disk_size_gb = 20
    machine_type = "n1-standard-2"
    image_type   = "cos_containerd"
    oauth_scopes = [
      "https://www.googleapis.com/auth/cloud-platform"
    ]
  }
}

output "test" {
  value = google_container_node_pool.default
}

