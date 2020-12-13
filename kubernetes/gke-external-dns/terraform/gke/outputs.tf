output "gcloud" {
  value = "gcloud container clusters get-credentials ${google_container_cluster.gke.name} --project ${local.gcp_project} --zone ${google_container_cluster.gke.zone}"
}
