variable "name" {
  type        = string
  description = "クラスタ名"
}

variable "location" {
  type        = string
  description = "GKEロケーション Region or Zone"
}

variable "node_locations" {
  type        = list(string)
  description = "GKEノードのロケーション一覧"
}

variable "release_channel" {
  type        = string
  description = "GKEリリースチャンネル"
}

variable "min_master_version" {
  type        = string
  description = "GKEマスターのマスターバージョン"
}

variable "network_id" {
  type        = string
  description = "VPCネットワークID"
}

variable "subnet_id" {
  type        = string
  description = "VPCサブネットID"
}

variable "pod_ipv4_cidr" {
  type        = string
  description = "PodのIPレンジ"
}

variable "service_ipv4_cidr" {
  type        = string
  description = "ServiceのIPレンジ"
}
