terraform {
  required_version = "0.14.2"

  required_providers {
    # https://registry.terraform.io/providers/hashicorp/google/latest
    google = {
      source  = "hashicorp/google"
      version = "3.50.0"
    }

    # https://registry.terraform.io/providers/hashicorp/google-beta/latest
    google-beta = {
      source = "hashicorp/google-beta"
      version = "3.50.0"
    }
  }
}

provider "google" {
  project = local.gcp_project
  region  = local.region
}

provider "google-beta" {
  project = local.gcp_project
  region  = local.region
}