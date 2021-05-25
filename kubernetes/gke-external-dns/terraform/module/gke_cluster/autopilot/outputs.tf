output "cluster" {
  value = google_container_cluster.this
}

output "kubeconfig" {
  value = templatefile("${path.module}/template/kubeconfig.yaml", {
    name                       = var.name
    certificate-authority-data = google_container_cluster.this.master_auth.0.cluster_ca_certificate
    server                     = "https://${google_container_cluster.this.endpoint}"
  })
}
