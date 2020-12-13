resource "google_compute_network" "vpc" {
  name                            = local.project
  description                     = "VPC for ${local.project}"
  auto_create_subnetworks         = false
  routing_mode                    = "REGIONAL"
  mtu                             = 1460
  delete_default_routes_on_create = false
}

resource "google_compute_subnetwork" "asia_northeast1" {
  ip_cidr_range = "10.128.0.0/20"
  name          = "${google_compute_network.vpc.name}-${local.region}"
  region        = local.region
  network       = google_compute_network.vpc.id
}

