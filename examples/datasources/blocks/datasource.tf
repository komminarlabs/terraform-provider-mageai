terraform {
  required_providers {
    mageai = {
      source = "komminarlabs/mageai"
    }
  }
}

provider "mageai" {}

data "mageai_blocks" "default" {
  pipeline_uuid = "example_pipeline"
}

output "default_pipeline_blocks" {
  value = data.mageai_blocks.default
}
