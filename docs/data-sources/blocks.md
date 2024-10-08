---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "mageai_blocks Data Source - terraform-provider-mageai"
subcategory: ""
description: |-
  Fetch and return the contents of all blocks in a pipeline.
---

# mageai_blocks (Data Source)

Fetch and return the contents of all blocks in a pipeline.



<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `pipeline_uuid` (String) The UUID of the pipeline to fetch the blocks from.

### Read-Only

- `blocks` (Attributes List) The blocks objects of a pipeline. (see [below for nested schema](#nestedatt--blocks))

<a id="nestedatt--blocks"></a>
### Nested Schema for `blocks`

Read-Only:

- `all_upstream_blocks_executed` (Boolean) Whether or not all upstream blocks have been successfully executed.
- `configuration` (Attributes) Miscellaneous configuration settings for the block. (see [below for nested schema](#nestedatt--blocks--configuration))
- `content` (String) Blocks file contents.
- `downstream_blocks` (Set of String) The block UUIDs that depend on this block.
- `executor_type` (String) The type of executor to use for the block: `ecs`, `gcp_cloud_run`, `azure_container_instance`, `k8s`, `local_python`, `pyspark`. See the [Kubernetes config](https://docs.mage.ai/production/configuring-production-settings/compute-resource#2-set-executor-type-and-customize-the-compute-resource-of-the-mage-executor) page for more details.
- `extension_uuid` (String) The extension uuid.
- `has_callback` (Boolean) The has_callback boolean.
- `language` (String) The language.
- `name` (String) Human readable name of block.
- `priority` (Number) The priority.
- `retry_config` (Attributes) The blocks objects of a block. (see [below for nested schema](#nestedatt--blocks--retry_config))
- `status` (String) Status of block: `executed`, `failed`, `not_executed`, `updated`.
- `timeout` (Number) The timeout.
- `type` (String) Type of block: `callback`, `chart`, `conditional`, `custom`, `data_exporter`, `data_loader`, `dbt`, `extension`, `global_data_product`, `markdown`, `scratchpad`, `sensor`, `transformer`.
- `upstream_blocks` (Set of String) The block UUIDs that this block depends on.
- `uuid` (String) Unique identifier for the block.

<a id="nestedatt--blocks--configuration"></a>
### Nested Schema for `blocks.configuration`

Read-Only:

- `data_provider` (String) Database or data warehouse for the SQL block to connect to.
- `data_provider_database` (String) Database name to use when saving the output of the SQL block.
- `data_provider_profile` (String) Profile target for the dbt block.
- `data_provider_schema` (String) Schema name to use when saving the output of the SQL block.
- `data_provider_table` (String) Table name to use when saving the output of the SQL block.
- `export_write_policy` (String) Whether to `replace` the existing table of the SQL block output, `append`, or raise an error and `fail`.
- `use_raw_sql` (String) Toggle writing raw SQL in the block. Read more [here](https://docs.mage.ai/guides/blocks/sql-blocks#using-raw-sql).


<a id="nestedatt--blocks--retry_config"></a>
### Nested Schema for `blocks.retry_config`

Read-Only:

- `delay` (Number) Initial delay (in seconds) before retry. If exponential_backoff is true, the delay time is multiplied by 2 for the next retry.
- `exponential_backoff` (Boolean) Whether to use exponential backoff retry.
- `max_delay` (Number) Maximum time between the first attempt and the last retry.
- `retries` (Number) Number of retry times.
