package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/komminarlabs/terraform-provider-mageai/internal/sdk/mageai"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                = &PipelineResource{}
	_ resource.ResourceWithImportState = &PipelineResource{}
	_ resource.ResourceWithImportState = &PipelineResource{}
)

// NewPipelineResource is a helper function to simplify the provider implementation.
func NewPipelineResource() resource.Resource {
	return &PipelineResource{}
}

// PipelineResource defines the resource implementation.
type PipelineResource struct {
	client mageai.Client
}

// Metadata returns the resource type name.
func (r *PipelineResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_pipeline"
}

// Schema defines the schema for the resource.
func (r *PipelineResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		Description: "To create a pipeline.",
		Attributes: map[string]schema.Attribute{
			"blocks": schema.ListNestedAttribute{
				Computed:    true,
				Description: "The blocks objects of a pipeline.",
				Default:     listdefault.StaticValue(types.ListValueMust(BlockModel{}.GetAttrType(), []attr.Value{})),
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
				Required:    true,
				Description: "Human readable name of the pipeline.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
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
				Optional:    true,
				Description: "The type of the pipeline: `integration`, `pyspark`, `python`, `streaming`. **Note:** that `python` is a standard (batch) pipeline with a python backend, while `pyspark` is a batch pipeline with a spark backend.",
				Default:     stringdefault.StaticString("python"),
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"integration", "pyspark", "python", "streaming"}...),
				},
			},
			"uuid": schema.StringAttribute{
				Computed:    true,
				Description: "The uuid.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
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
	}
}

// Create creates the resource and sets the initial Terraform state.
func (r *PipelineResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan PipelineModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	createPipelineRequest := &mageai.CreatePipelineRequest{
		Pipeline: mageai.PipelineRequest{
			Name: plan.Name.ValueString(),
			Type: mageai.PipelineType(plan.Type.ValueString()),
		},
	}

	createPipelineResponse, err := r.client.PipelineAPI().CreatePipeline(ctx, createPipelineRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating pipeline",
			"Could not create pipeline, unexpected error: "+err.Error(),
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
	pipelineModel, err := getPipelineModel(ctx, createPipelineResponse.Pipeline)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error getting pipeline model",
			err.Error(),
		)
		return
	}
	plan = *pipelineModel

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *PipelineResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state PipelineModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed pipeline value from Mage AI
	readDatabaseResponse, err := r.client.PipelineAPI().ReadPipeline(ctx, state.UUID.ValueStringPointer())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error getting pipeline",
			err.Error(),
		)
		return
	}

	// Overwrite items with refreshed state
	pipelineModel, err := getPipelineModel(ctx, readDatabaseResponse.Pipeline)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error getting pipeline model",
			err.Error(),
		)
		return
	}
	state = *pipelineModel

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *PipelineResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan PipelineModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	updatePipelineRequest := &mageai.UpdatePipelineRequest{
		Pipeline: mageai.PipelineRequest{
			Name: plan.Name.ValueString(),
			Type: mageai.PipelineType(plan.Type.ValueString()),
		}}

	// Update existing pipeline
	updatePipelineResponse, err := r.client.PipelineAPI().UpdatePipeline(ctx, plan.UUID.ValueStringPointer(), updatePipelineRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating pipeline",
			"Could not update pipeline, unexpected error: "+err.Error(),
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
	pipelineModel, err := getPipelineModel(ctx, updatePipelineResponse.Pipeline)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error getting pipeline model",
			err.Error(),
		)
		return
	}
	plan = *pipelineModel

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *PipelineResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state PipelineModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete existing pipeline
	err := r.client.PipelineAPI().DeletePipeline(ctx, state.UUID.ValueStringPointer())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting pipeline",
			"Could not delete pipeline, unexpected error: "+err.Error(),
		)
		return
	}
}

// Configure adds the provider configured client to the resource.
func (r *PipelineResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	pd, ok := req.ProviderData.(providerData)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected mageai.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}
	r.client = pd.client
}

func (r *PipelineResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("uuid"), req, resp)
}
