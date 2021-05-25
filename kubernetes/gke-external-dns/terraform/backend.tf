terraform {
  backend "gcs" {
    bucket = "tf-state-sandbox-230212"
    prefix = "gke-external-dns"
  }
}
