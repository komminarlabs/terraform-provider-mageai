terraform {
  required_providers {
    mageai = {
      source = "komminarlabs/mageai"
    }
  }
}

provider "mageai" {}

data "mageai_block" "default" {
  pipeline_uuid = "example_pipeline"
  uuid          = "daring_butterfly"
}

output "default_pipeline_block" {
  value = data.mageai_block.default
}
