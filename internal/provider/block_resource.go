package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/komminarlabs/terraform-provider-mageai/internal/sdk/mageai"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                = &BlockResource{}
	_ resource.ResourceWithImportState = &BlockResource{}
	_ resource.ResourceWithImportState = &BlockResource{}
)

// NewBlockResource is a helper function to simplify the provider implementation.
func NewBlockResource() resource.Resource {
	return &BlockResource{}
}

// BlockResource defines the resource implementation.
type BlockResource struct {
	client mageai.Client
}

// Metadata returns the resource type name.
func (r *BlockResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_block"
}

// Schema defines the schema for the resource.
func (r *BlockResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		Description: "Create a new block in a pipeline.",
		Attributes: map[string]schema.Attribute{
			"pipeline_uuid": schema.StringAttribute{
				Required:    true,
				Description: "The UUID of the pipeline to create the block.",
			},
			"all_upstream_blocks_executed": schema.BoolAttribute{
				Computed:    true,
				Description: "Whether or not all upstream blocks have been successfully executed.",
			},
			"configuration": schema.SingleNestedAttribute{
				Computed:    true,
				Optional:    true,
				Description: "Miscellaneous configuration settings for the block.",
				Attributes: map[string]schema.Attribute{
					"data_provider": schema.StringAttribute{
						Computed:    true,
						Optional:    true,
						Description: "Database or data warehouse for the SQL block to connect to.",
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
					},
					"data_provider_database": schema.StringAttribute{
						Computed:    true,
						Optional:    true,
						Description: "Database name to use when saving the output of the SQL block.",
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
					},
					"data_provider_profile": schema.StringAttribute{
						Computed:    true,
						Optional:    true,
						Description: "Profile target for the dbt block.",
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
					},
					"data_provider_schema": schema.StringAttribute{
						Computed:    true,
						Optional:    true,
						Description: "Schema name to use when saving the output of the SQL block.",
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
					},
					"data_provider_table": schema.StringAttribute{
						Computed:    true,
						Optional:    true,
						Description: "Table name to use when saving the output of the SQL block.",
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
					},
					"export_write_policy": schema.StringAttribute{
						Computed:    true,
						Optional:    true,
						Description: "Whether to `replace` the existing table of the SQL block output, `append`, or raise an error and `fail`.",
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
					},
					"use_raw_sql": schema.StringAttribute{
						Computed:    true,
						Optional:    true,
						Description: "Toggle writing raw SQL in the block. Read more [here](https://docs.mage.ai/guides/blocks/sql-blocks#using-raw-sql).",
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
					},
				},
			},
			"content": schema.StringAttribute{
				Computed:    true,
				Optional:    true,
				Description: "Block file contents.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
					stringplanmodifier.UseStateForUnknown(),
				},
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
				Optional:    true,
				Description: "The extension uuid.",
			},
			"has_callback": schema.BoolAttribute{
				Computed:    true,
				Description: "The has_callback boolean.",
			},
			"language": schema.StringAttribute{
				Computed:    true,
				Optional:    true,
				Description: "The language.",
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "Human readable name of block.",
			},
			"priority": schema.Int32Attribute{
				Computed:    true,
				Optional:    true,
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
				Required:    true,
				Description: "Type of block: `callback`, `chart`, `conditional`, `custom`, `data_exporter`, `data_loader`, `dbt`, `extension`, `global_data_product`, `markdown`, `scratchpad`, `sensor`, `transformer`.",
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"callback", "chart", "conditional", "custom", "data_exporter", "data_loader", "dbt", "extension", "global_data_product", "markdown", "scratchpad", "sensor", "transformer"}...),
				},
			},
			"upstream_blocks": schema.SetAttribute{
				Computed:    true,
				Description: "The block UUIDs that this block depends on.",
				Default:     setdefault.StaticValue(types.SetUnknown(types.StringType)),
				ElementType: types.StringType,
			},
			"uuid": schema.StringAttribute{
				Computed:    true,
				Description: "Unique identifier for the block.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

// Create creates the resource and sets the initial Terraform state.
func (r *BlockResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan BlockResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	createBlockRequest, err := makeCreateBlockRequestFromModel(ctx, plan)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating block",
			err.Error(),
		)
		return
	}

	createBlockResponse, err := r.client.BlockAPI().CreateBlock(ctx, plan.PipelineUUID.ValueStringPointer(), createBlockRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating block",
			"Could not create block, unexpected error: "+err.Error(),
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
	blockModel, err := getBlockModel(ctx, createBlockResponse.Block)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error getting block model",
			err.Error(),
		)
		return
	}
	plan.BlockModel = *blockModel

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *BlockResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state BlockResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed block value from Mage AI
	readDatabaseResponse, err := r.client.BlockAPI().ReadBlock(ctx, state.PipelineUUID.ValueStringPointer(), state.UUID.ValueStringPointer())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error getting block",
			err.Error(),
		)
		return
	}

	// Overwrite items with refreshed state
	blockState, err := getBlockModel(ctx, readDatabaseResponse.Block)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error getting blocks",
			err.Error(),
		)
		return
	}
	state.BlockModel = *blockState

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *BlockResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan BlockResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	updateBlockRequest, err := makeUpdateBlockRequestFromModel(ctx, plan)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating block",
			err.Error(),
		)
		return
	}

	// Update existing block
	updateBlockResponse, err := r.client.BlockAPI().UpdateBlock(ctx, plan.PipelineUUID.ValueStringPointer(), plan.UUID.ValueStringPointer(), updateBlockRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating block",
			"Could not update block, unexpected error: "+err.Error(),
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
	blockModel, err := getBlockModel(ctx, updateBlockResponse.Block)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error getting block model",
			err.Error(),
		)
		return
	}
	plan.BlockModel = *blockModel

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *BlockResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state BlockResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete existing block
	err := r.client.BlockAPI().DeleteBlock(ctx, state.PipelineUUID.ValueStringPointer(), state.UUID.ValueStringPointer())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting block",
			"Could not delete block, unexpected error: "+err.Error(),
		)
		return
	}
}

// Configure adds the provider configured client to the resource.
func (r *BlockResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	pd, ok := req.ProviderData.(providerData)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected mageai.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}
	r.client = pd.client
}

func (r *BlockResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("uuid"), req, resp)
}
