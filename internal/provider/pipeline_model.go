package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/komminarlabs/terraform-provider-mageai/internal/sdk/mageai"
)

type PipelineModel struct {
	Blocks                   []BlockModel `tfsdk:"blocks"`
	CacheBlockOutputInMemory types.Bool   `tfsdk:"cache_block_output_in_memory"`
	CreatedAt                types.String `tfsdk:"created_at"`
	Description              types.String `tfsdk:"description"`
	ExecutorCount            types.Int32  `tfsdk:"executor_count"`
	Name                     types.String `tfsdk:"name"`
	RetryConfig              types.Object `tfsdk:"retry_config"`
	RunPipelineInOneProcess  types.Bool   `tfsdk:"run_pipeline_in_one_process"`
	Tags                     types.Set    `tfsdk:"tags"`
	Type                     types.String `tfsdk:"type"`
	UUID                     types.String `tfsdk:"uuid"`
	UpdatedAt                types.String `tfsdk:"updated_at"`
	VariablesDir             types.String `tfsdk:"variables_dir"`
}

type RetryConfigModel struct {
	Delay              types.Int32 `tfsdk:"delay"`
	ExponentialBackoff types.Bool  `tfsdk:"exponential_backoff"`
	MaxDelay           types.Int32 `tfsdk:"max_delay"`
	Retries            types.Int32 `tfsdk:"retries"`
}

func (r RetryConfigModel) GetAttrType() map[string]attr.Type {
	return map[string]attr.Type{
		"delay":               types.Int32Type,
		"exponential_backoff": types.BoolType,
		"max_delay":           types.Int32Type,
		"retries":             types.Int32Type,
	}
}

func getPipelineModel(ctx context.Context, pipeline mageai.Pipeline) (*PipelineModel, error) {
	blocks := make([]BlockModel, 0)
	for _, block := range pipeline.Blocks {
		blockState, err := getBlockModel(ctx, block)
		if err != nil {
			return nil, fmt.Errorf("error getting blocks %s", err)
		}
		blocks = append(blocks, *blockState)
	}

	pipelineRetryConfigValue := RetryConfigModel{
		Delay:              types.Int32Value(pipeline.RetryConfig.Delay),
		ExponentialBackoff: types.BoolValue(pipeline.RetryConfig.ExponentialBackoff),
		MaxDelay:           types.Int32Value(pipeline.RetryConfig.MaxDelay),
		Retries:            types.Int32Value(pipeline.RetryConfig.Retries),
	}

	pipelineRetryConfigObjectValue, diags := types.ObjectValueFrom(ctx, pipelineRetryConfigValue.GetAttrType(), pipelineRetryConfigValue)
	if diags.HasError() {
		return nil, fmt.Errorf("error getting pipeline retry_config")
	}

	tags, diags := types.SetValueFrom(ctx, types.StringType, pipeline.Tags)
	if diags.HasError() {
		return nil, fmt.Errorf("error getting tags")
	}

	pipelineState := PipelineModel{
		Blocks:                   blocks,
		CacheBlockOutputInMemory: types.BoolValue(pipeline.CacheBlockOutputInMemory),
		CreatedAt:                types.StringValue(pipeline.CreatedAt),
		Description:              types.StringValue(pipeline.Description),
		ExecutorCount:            types.Int32Value(pipeline.ExecutorCount),
		Name:                     types.StringValue(pipeline.Name),
		RetryConfig:              pipelineRetryConfigObjectValue,
		RunPipelineInOneProcess:  types.BoolValue(pipeline.RunPipelineInOneProcess),
		Tags:                     tags,
		Type:                     types.StringValue(pipeline.Type),
		UUID:                     types.StringValue(pipeline.UUID),
		UpdatedAt:                types.StringValue(pipeline.UpdatedAt),
		VariablesDir:             types.StringValue(pipeline.VariablesDir),
	}
	return &pipelineState, nil
}
