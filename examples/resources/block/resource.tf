terraform {
  required_providers {
    mageai = {
      source = "komminarlabs/mageai"
    }
  }
}

provider "mageai" {}

resource "mageai_block" "default" {
  name          = "example_block"
  pipeline_uuid = "example_pipeline"
  type          = "data_loader"
  content       = file("${path.module}/script.py")

  configuration = {
    data_provider = "example_data_provider"
  }
}

output "default_block" {
  value = mageai_block.default
}
