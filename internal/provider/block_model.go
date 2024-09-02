package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/komminarlabs/terraform-provider-mageai/internal/sdk/mageai"
)

type BlockDataSourceModel struct {
	PipelineUUID types.String `tfsdk:"pipeline_uuid"`
	BlockModel
}

type BlockResourceModel struct {
	PipelineUUID types.String `tfsdk:"pipeline_uuid"`
	BlockModel
}

type BlocksDataSourceModel struct {
	PipelineUUID types.String `tfsdk:"pipeline_uuid"`
	Blocks       []BlockModel `tfsdk:"blocks"`
}

type BlockModel struct {
	AllUpstreamBlocksExecuted types.Bool   `tfsdk:"all_upstream_blocks_executed"`
	Configuration             types.Object `tfsdk:"configuration"`
	Content                   types.String `tfsdk:"content"`
	DownstreamBlocks          types.Set    `tfsdk:"downstream_blocks"`
	ExecutorType              types.String `tfsdk:"executor_type"`
	ExtensionUUID             types.String `tfsdk:"extension_uuid"`
	HasCallback               types.Bool   `tfsdk:"has_callback"`
	Language                  types.String `tfsdk:"language"`
	Name                      types.String `tfsdk:"name"`
	Priority                  types.Int32  `tfsdk:"priority"`
	RetryConfig               types.Object `tfsdk:"retry_config"`
	Status                    types.String `tfsdk:"status"`
	Timeout                   types.Int64  `tfsdk:"timeout"`
	Type                      types.String `tfsdk:"type"`
	UpstreamBlocks            types.Set    `tfsdk:"upstream_blocks"`
	UUID                      types.String `tfsdk:"uuid"`
}

type BlockConfigurationModel struct {
	DataProvider         types.String `tfsdk:"data_provider"`
	DataProviderDatabase types.String `tfsdk:"data_provider_database"`
	DataProviderProfile  types.String `tfsdk:"data_provider_profile"`
	DataProviderSchema   types.String `tfsdk:"data_provider_schema"`
	DataProviderTable    types.String `tfsdk:"data_provider_table"`
	ExportWritePolicy    types.String `tfsdk:"export_write_policy"`
	UseRawSql            types.String `tfsdk:"use_raw_sql"`
}

func (b BlockModel) GetAttrType() attr.Type {
	return types.ObjectType{AttrTypes: map[string]attr.Type{
		"all_upstream_blocks_executed": types.BoolType,
		"configuration":                types.ObjectType{AttrTypes: BlockConfigurationModel{}.GetAttrType()},
		"content":                      types.StringType,
		"downstream_blocks":            types.SetType{ElemType: types.StringType},
		"executor_type":                types.StringType,
		"extension_uuid":               types.StringType,
		"has_callback":                 types.BoolType,
		"language":                     types.StringType,
		"name":                         types.StringType,
		"priority":                     types.Int32Type,
		"retry_config":                 types.ObjectType{AttrTypes: RetryConfigModel{}.GetAttrType()},
		"status":                       types.StringType,
		"timeout":                      types.Int64Type,
		"type":                         types.StringType,
		"upstream_blocks":              types.SetType{ElemType: types.StringType},
		"uuid":                         types.StringType,
	}}
}

func (b BlockConfigurationModel) GetAttrType() map[string]attr.Type {
	return map[string]attr.Type{
		"data_provider":          types.StringType,
		"data_provider_database": types.StringType,
		"data_provider_profile":  types.StringType,
		"data_provider_schema":   types.StringType,
		"data_provider_table":    types.StringType,
		"export_write_policy":    types.StringType,
		"use_raw_sql":            types.StringType,
	}
}

func getBlockModel(ctx context.Context, block mageai.Block) (*BlockModel, error) {
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
		Delay:              types.Int32Value(block.RetryConfig.Delay),
		ExponentialBackoff: types.BoolValue(block.RetryConfig.ExponentialBackoff),
		MaxDelay:           types.Int32Value(block.RetryConfig.MaxDelay),
		Retries:            types.Int32Value(block.RetryConfig.Retries),
	}

	blockRetryConfigObjectValue, diags := types.ObjectValueFrom(ctx, blockRetryConfigValue.GetAttrType(), blockRetryConfigValue)
	if diags.HasError() {
		return nil, fmt.Errorf("error getting block retry_config")
	}

	upstreamBlocks, diags := types.SetValueFrom(ctx, types.StringType, block.UpstreamBlocks)
	if diags.HasError() {
		return nil, fmt.Errorf("error getting upstream_blocks")
	}

	blockState := BlockModel{
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
	return &blockState, nil
}

func convertUpstreamBlocksSetToStringSlice(upstreamBlocks basetypes.SetValue) []string {
	upstreamBlocksSlice := []string{}
	for _, block := range upstreamBlocks.Elements() {
		upstreamBlocksSlice = append(upstreamBlocksSlice, block.String())
	}
	return upstreamBlocksSlice
}

func convertBlockConfigurationObjectToModel(ctx context.Context, blockConfigurationObject basetypes.ObjectValue) (*mageai.BlockConfiguration, error) {
	configurationModel := BlockConfigurationModel{}
	diags := blockConfigurationObject.As(ctx, &configurationModel, basetypes.ObjectAsOptions{})
	if diags.HasError() {
		return nil, fmt.Errorf("could not get configuration, unexpected error: %v", diags.Errors())
	}

	configuration := &mageai.BlockConfiguration{
		DataProvider:         configurationModel.DataProvider.ValueString(),
		DataProviderDatabase: configurationModel.DataProviderDatabase.ValueString(),
		DataProviderProfile:  configurationModel.DataProviderProfile.ValueString(),
		DataProviderSchema:   configurationModel.DataProviderSchema.ValueString(),
		DataProviderTable:    configurationModel.DataProviderTable.ValueString(),
		ExportWritePolicy:    configurationModel.ExportWritePolicy.ValueString(),
		UseRawSql:            configurationModel.UseRawSql.ValueString(),
	}
	return configuration, nil
}

func makeCreateBlockRequestFromModel(ctx context.Context, b BlockResourceModel) (*mageai.CreateBlockRequest, error) {
	upstreamBlocks := convertUpstreamBlocksSetToStringSlice(b.UpstreamBlocks)

	configuration, err := convertBlockConfigurationObjectToModel(ctx, b.Configuration)
	if err != nil {
		return nil, fmt.Errorf("error converting block configuration: %v", err)
	}

	return &mageai.CreateBlockRequest{
		Block: mageai.BlockRequest{
			Configuration:  *configuration,
			Content:        b.Content.ValueString(),
			ExtensionUUID:  b.ExtensionUUID.ValueString(),
			Language:       b.Language.ValueString(),
			Name:           b.Name.ValueString(),
			Priority:       b.Priority.ValueInt32(),
			Type:           mageai.BlockType(b.Type.ValueString()),
			UpstreamBlocks: upstreamBlocks,
		},
	}, nil
}

func makeUpdateBlockRequestFromModel(ctx context.Context, b BlockResourceModel) (*mageai.UpdateBlockRequest, error) {
	configuration, err := convertBlockConfigurationObjectToModel(ctx, b.Configuration)
	if err != nil {
		return nil, fmt.Errorf("error converting block configuration")
	}

	return &mageai.UpdateBlockRequest{
		Block: mageai.BlockRequest{
			Configuration:  *configuration,
			Content:        b.Content.ValueString(),
			ExtensionUUID:  b.ExtensionUUID.ValueString(),
			Language:       b.Language.ValueString(),
			Name:           b.Name.ValueString(),
			Priority:       b.Priority.ValueInt32(),
			Type:           mageai.BlockType(b.Type.ValueString()),
			UpstreamBlocks: convertUpstreamBlocksSetToStringSlice(b.UpstreamBlocks),
		},
	}, nil
}
