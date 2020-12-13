output "vpc" {
  value = google_compute_network.vpc
}

output "subnets" {
  value = {
    asia_northeast1 = google_compute_subnetwork.asia_northeast1
  }
}
