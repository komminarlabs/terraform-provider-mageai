terraform {
  required_providers {
    mageai = {
      source = "komminarlabs/mageai"
    }
  }
}

provider "mageai" {}

resource "mageai_pipeline" "default" {
  name = "example_pipeline"
  type = "python"
}

output "default_pipeline" {
  value = mageai_pipeline.default
}
