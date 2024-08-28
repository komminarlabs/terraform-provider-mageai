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
		blockConfigurationValue := BlockConfigurationModel{
			DataProvider:         types.StringValue(block.Configuration.DataProvider),
			DataProviderDatabase: types.StringValue(block.Configuration.DataProviderDatabase),
			DataProviderProfile:  types.StringValue(block.Configuration.DataProviderProfile),
			DataProviderSchema:   types.StringValue(block.Configuration.DataProviderSchema),
			DataProviderTable:    types.StringValue(block.Configuration.DataProviderTable),
			ExportWritePolicy:    types.StringValue(block.Configuration.ExportWritePolicy),
			UseRawSql:            types.StringValue(block.Configuration.UseRawSql),
		}

		blockConfigurationObjectValue, diags := types.ObjectValueFrom(ctx, blockConfigurationValue.GetAttrType(), blockConfigurationValue)
		if diags.HasError() {
			return nil, fmt.Errorf("error getting block configuration")
		}

		downstreamBlocks, diags := types.SetValueFrom(ctx, types.StringType, block.DownstreamBlocks)
		if diags.HasError() {
			return nil, fmt.Errorf("error getting downstream_blocks")
		}

		blockRetryConfigValue := RetryConfigModel{
			Delay:              types.Int32Value(pipeline.RetryConfig.Delay),
			ExponentialBackoff: types.BoolValue(pipeline.RetryConfig.ExponentialBackoff),
			MaxDelay:           types.Int32Value(pipeline.RetryConfig.MaxDelay),
			Retries:            types.Int32Value(pipeline.RetryConfig.Retries),
		}

		blockRetryConfigObjectValue, diags := types.ObjectValueFrom(ctx, blockRetryConfigValue.GetAttrType(), blockRetryConfigValue)
		if diags.HasError() {
			return nil, fmt.Errorf("error getting block retry_config")
		}

		upstreamBlocks, diags := types.SetValueFrom(ctx, types.StringType, block.UpstreamBlocks)
		if diags.HasError() {
			return nil, fmt.Errorf("error getting upstream_blocks")
		}

		block := BlockModel{
			AllUpstreamBlocksExecuted: types.BoolValue(block.AllUpstreamBlocksExecuted),
			Configuration:             blockConfigurationObjectValue,
			Content:                   types.StringValue(block.Content),
			DownstreamBlocks:          downstreamBlocks,
			ExecutorType:              types.StringValue(block.ExecutorType),
			ExtensionUUID:             types.StringValue(block.ExtensionUUID),
			HasCallback:               types.BoolValue(block.HasCallback),
			Language:                  types.StringValue(block.Language),
			Name:                      types.StringValue(block.Name),
			Priority:                  types.Int32Value(block.Priority),
			RetryConfig:               blockRetryConfigObjectValue,
			Status:                    types.StringValue(block.Status),
			Timeout:                   types.Int64Value(block.Timeout),
			Type:                      types.StringValue(block.Type),
			UpstreamBlocks:            upstreamBlocks,
			UUID:                      types.StringValue(block.UUID),
		}
		blocks = append(blocks, block)
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
