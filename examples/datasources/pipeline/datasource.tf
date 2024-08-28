terraform {
  required_providers {
    mageai = {
      source = "komminarlabs/mageai"
    }
  }
}

provider "mageai" {}

data "mageai_pipeline" "default" {
  uuid = "example_pipeline"
}

output "default_pipeline" {
  value = data.mageai_pipeline.default
}
