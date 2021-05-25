module "gke_cluster" {
  source             = "./module/gke_cluster/autopilot"
  name               = local.project
  location           = local.region
  node_locations     = local.zones
  release_channel    = "RAPID"
  min_master_version = "1.20.6-gke.1000" # gcloud container get-server-config --zone asia-northeast1-a
  network_id         = google_compute_network.vpc.id
  subnet_id          = google_compute_subnetwork.asia_northeast1.id
  pod_ipv4_cidr      = "10.129.0.0/17"
  service_ipv4_cidr  = "10.129.128.0/22"
}

output "gke_cluster_kubeconfig" {
  sensitive = true
  value     = module.gke_cluster.kubeconfig
}
