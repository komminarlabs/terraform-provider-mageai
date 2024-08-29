package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/komminarlabs/terraform-provider-mageai/internal/sdk/mageai"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &BlocksDataSource{}
	_ datasource.DataSourceWithConfigure = &BlocksDataSource{}
)

// NewBlocksDataSource is a helper function to simplify the provider implementation.
func NewBlocksDataSource() datasource.DataSource {
	return &BlocksDataSource{}
}

// BlocksDataSource is the data source implementation.
type BlocksDataSource struct {
	client mageai.Client
}

// Metadata returns the data source type name.
func (d *BlocksDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_blocks"
}

// Schema defines the schema for the data source.
func (d *BlocksDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		Description: "Fetch and return the contents of all blocks in a pipeline.",
		Attributes: map[string]schema.Attribute{
			"pipeline_uuid": schema.StringAttribute{
				Required:    true,
				Description: "The UUID of the pipeline to fetch the blocks from.",
			},
			"blocks": schema.ListNestedAttribute{
				Computed:    true,
				Description: "The blocks objects of a pipeline.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{

						"all_upstream_blocks_executed": schema.BoolAttribute{
							Computed:    true,
							Description: "Whether or not all upstream blocks have been successfully executed.",
						},
						"configuration": schema.SingleNestedAttribute{
							Computed:    true,
							Description: "Miscellaneous configuration settings for the block.",
							Attributes: map[string]schema.Attribute{
								"data_provider": schema.StringAttribute{
									Computed:    true,
									Description: "Database or data warehouse for the SQL block to connect to.",
								},
								"data_provider_database": schema.StringAttribute{
									Computed:    true,
									Description: "Database name to use when saving the output of the SQL block.",
								},
								"data_provider_profile": schema.StringAttribute{
									Computed:    true,
									Description: "Profile target for the dbt block.",
								},
								"data_provider_schema": schema.StringAttribute{
									Computed:    true,
									Description: "Schema name to use when saving the output of the SQL block.",
								},
								"data_provider_table": schema.StringAttribute{
									Computed:    true,
									Description: "Table name to use when saving the output of the SQL block.",
								},
								"export_write_policy": schema.StringAttribute{
									Computed:    true,
									Description: "Whether to `replace` the existing table of the SQL block output, `append`, or raise an error and `fail`.",
								},
								"use_raw_sql": schema.StringAttribute{
									Computed:    true,
									Description: "Toggle writing raw SQL in the block. Read more [here](https://docs.mage.ai/guides/blocks/sql-blocks#using-raw-sql).",
								},
							},
						},
						"content": schema.StringAttribute{
							Computed:    true,
							Description: "Blocks file contents.",
						},
						"downstream_blocks": schema.SetAttribute{
							Computed:    true,
							Description: "The block UUIDs that depend on this block.",
							ElementType: types.StringType,
						},
						"executor_type": schema.StringAttribute{
							Computed:    true,
							Description: "The type of executor to use for the block: `ecs`, `gcp_cloud_run`, `azure_container_instance`, `k8s`, `local_python`, `pyspark`. See the [Kubernetes config](https://docs.mage.ai/production/configuring-production-settings/compute-resource#2-set-executor-type-and-customize-the-compute-resource-of-the-mage-executor) page for more details.",
						},
						"extension_uuid": schema.StringAttribute{
							Computed:    true,
							Description: "The extension uuid.",
						},
						"has_callback": schema.BoolAttribute{
							Computed:    true,
							Description: "The has_callback boolean.",
						},
						"language": schema.StringAttribute{
							Computed:    true,
							Description: "The language.",
						},
						"name": schema.StringAttribute{
							Computed:    true,
							Description: "Human readable name of block.",
						},
						"priority": schema.Int32Attribute{
							Computed:    true,
							Description: "The priority.",
						},
						"retry_config": schema.SingleNestedAttribute{
							Computed:    true,
							Description: "The blocks objects of a block.",
							Attributes: map[string]schema.Attribute{
								"delay": schema.Int32Attribute{
									Computed:    true,
									Description: "Initial delay (in seconds) before retry. If exponential_backoff is true, the delay time is multiplied by 2 for the next retry.",
								},
								"exponential_backoff": schema.BoolAttribute{
									Computed:    true,
									Description: "Whether to use exponential backoff retry.",
								},
								"max_delay": schema.Int32Attribute{
									Computed:    true,
									Description: "Maximum time between the first attempt and the last retry.",
								},
								"retries": schema.Int32Attribute{
									Computed:    true,
									Description: "Number of retry times.",
								},
							},
						},
						"status": schema.StringAttribute{
							Computed:    true,
							Description: "Status of block: `executed`, `failed`, `not_executed`, `updated`.",
						},
						"timeout": schema.Int64Attribute{
							Computed:    true,
							Description: "The timeout.",
						},
						"type": schema.StringAttribute{
							Computed:    true,
							Description: "Type of block: `callback`, `chart`, `conditional`, `custom`, `data_exporter`, `data_loader`, `dbt`, `extension`, `global_data_product`, `markdown`, `scratchpad`, `sensor`, `transformer`.",
						},
						"upstream_blocks": schema.SetAttribute{
							Computed:    true,
							Description: "The block UUIDs that this block depends on.",
							ElementType: types.StringType,
						},
						"uuid": schema.StringAttribute{
							Computed:    true,
							Description: "Unique identifier for the block.",
						},
					},
				},
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (d *BlocksDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	pd, ok := req.ProviderData.(providerData)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected mageai.client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	d.client = pd.client
}

// Read refreshes the Terraform state with the latest data.
func (d *BlocksDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state BlocksDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	readDatabaseResponse, err := d.client.BlockAPI().ReadBlocks(ctx, state.PipelineUUID.ValueStringPointer())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error getting block",
			err.Error(),
		)
		return
	}

	// Map response body to model
	blocks := make([]BlockModel, 0)
	for _, block := range readDatabaseResponse.Blocks {
		blockState, err := getBlockModel(ctx, block)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error getting blocks",
				err.Error(),
			)
			return
		}
		blocks = append(blocks, *blockState)
	}
	state.Blocks = blocks

	// Set state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
