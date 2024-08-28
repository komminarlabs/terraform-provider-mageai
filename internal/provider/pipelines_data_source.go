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
	_ datasource.DataSource              = &PipelinesDataSource{}
	_ datasource.DataSourceWithConfigure = &PipelinesDataSource{}
)

// NewPipelinesDataSource is a helper function to simplify the provider implementation.
func NewPipelinesDataSource() datasource.DataSource {
	return &PipelinesDataSource{}
}

// PipelinesDataSource is the data source implementation.
type PipelinesDataSource struct {
	client mageai.Client
}

// PipelinesDataSourceModel describes the data source data model.
type PipelinesDataSourceModel struct {
	Pipelines []PipelineModel `tfsdk:"pipelines"`
}

// Metadata returns the data source type name.
func (d *PipelinesDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_pipelines"
}

// Schema defines the schema for the data source.
func (d *PipelinesDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		Description: "To retrieve all pipelines.",

		Attributes: map[string]schema.Attribute{
			"pipelines": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
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
										Description: "Block file contents.",
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
										Description: "The blocks objects of a pipeline.",
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
						"cache_block_output_in_memory": schema.BoolAttribute{
							Computed:    true,
							Description: "The cache_block_output_in_memory.",
						},
						"created_at": schema.StringAttribute{
							Computed:    true,
							Description: "The created_at.",
						},
						"description": schema.StringAttribute{
							Computed:    true,
							Description: "The description.",
						},
						"executor_count": schema.Int32Attribute{
							Computed:    true,
							Description: "The executor count.",
						},
						"name": schema.StringAttribute{
							Computed:    true,
							Description: "Human readable name of the pipeline.",
						},
						"retry_config": schema.SingleNestedAttribute{
							Computed:    true,
							Description: "The blocks objects of a pipeline.",
							Attributes: map[string]schema.Attribute{
								"delay": schema.Int32Attribute{
									Computed:    true,
									Description: "The type of template part.",
								},
								"exponential_backoff": schema.BoolAttribute{
									Computed:    true,
									Description: "The type of template part.",
								},
								"max_delay": schema.Int32Attribute{
									Computed:    true,
									Description: "The type of template part.",
								},
								"retries": schema.Int32Attribute{
									Computed:    true,
									Description: "The type of template part.",
								},
							},
						},
						"run_pipeline_in_one_process": schema.BoolAttribute{
							Computed:    true,
							Description: "The bool value for run_pipeline_in_one_process.",
						},
						"tags": schema.SetAttribute{
							Computed:    true,
							Description: "The tags.",
							ElementType: types.StringType,
						},
						"type": schema.StringAttribute{
							Computed:    true,
							Description: "The type of the pipeline: `integration`, `pyspark`, `python`, `streaming`. **Note:** that `python` is a standard (batch) pipeline with a python backend, while `pyspark` is a batch pipeline with a spark backend.",
						},
						"uuid": schema.StringAttribute{
							Computed:    true,
							Description: "The uuid.",
						},
						"updated_at": schema.StringAttribute{
							Computed:    true,
							Description: "The updated_at value.",
						},
						"variables_dir": schema.StringAttribute{
							Computed:    true,
							Description: "The data directory path.",
						},
					},
				},
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (d *PipelinesDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *PipelinesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state PipelinesDataSourceModel

	readDatabasesResponse, err := d.client.PipelineAPI().ReadPipelines(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error getting pipelines",
			err.Error(),
		)
		return
	}

	// Map response body to model
	for _, pipeline := range readDatabasesResponse.Pipelines {
		pipelineState, err := getPipelineModel(ctx, pipeline)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error getting pipelines",
				err.Error(),
			)
			return
		}
		state.Pipelines = append(state.Pipelines, *pipelineState)
	}

	// Set state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
