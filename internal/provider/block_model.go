package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

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