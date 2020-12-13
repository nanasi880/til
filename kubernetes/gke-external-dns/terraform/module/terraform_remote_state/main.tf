variable "vpc" {
  type        = bool
  description = "VPCデータを読み取る場合はtrue"
  default     = false
}

data "terraform_remote_state" "vpc" {
  count   = var.vpc ? 1 : 0
  backend = "gcs"
  config = {
    bucket = "tf-state-sandbox-230212"
    prefix = "gke-external-dns/vpc"
  }
}

output "vpc" {
  value = var.vpc ? data.terraform_remote_state.vpc[0] : null
}
