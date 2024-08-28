terraform {
  required_providers {
    mageai = {
      source = "komminarlabs/mageai"
    }
  }
}

provider "mageai" {}

data "mageai_pipelines" "all" {}

output "pipelines" {
  value = data.mageai_pipelines.all
}
