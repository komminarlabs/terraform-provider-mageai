package mageai

type errorResponse struct {
	Error struct {
		Code      int    `json:"code"`
		Exception string `json:"exception"`
		Message   string `json:"message"`
	} `json:"error"`
}

type pipelineResponse struct {
	Pipeline Pipeline `json:"pipeline"`
}

type pipelinesResponse struct {
	Pipelines []Pipeline `json:"pipelines"`
}

type Pipeline struct {
	Blocks                   []Block     `json:"blocks"`
	CacheBlockOutputInMemory bool        `json:"cache_block_output_in_memory"`
	CreatedAt                string      `json:"created_at"`
	Description              string      `json:"description"`
	ExecutorCount            int32       `json:"executor_count"`
	Name                     string      `json:"name"`
	RetryConfig              RetryConfig `json:"retry_config"`
	RunPipelineInOneProcess  bool        `json:"run_pipeline_in_one_process"`
	Tags                     []string    `json:"tags"`
	Type                     string      `json:"type"`
	UUID                     string      `json:"uuid"`
	UpdatedAt                string      `json:"updated_at"`
	VariablesDir             string      `json:"variables_dir"`
}

type Block struct {
	AllUpstreamBlocksExecuted bool               `json:"all_upstream_blocks_executed"`
	Color                     string             `json:"color"`
	Configuration             BlockConfiguration `json:"configuration"`
	Content                   string             `json:"content"`
	DownstreamBlocks          []string           `json:"downstream_blocks"`
	ExecutorType              string             `json:"executor_type"`
	ExtensionUUID             string             `json:"extension_uuid"`
	HasCallback               bool               `json:"has_callback"`
	Language                  string             `json:"language"`
	Name                      string             `json:"name"`
	Pipelines                 []string           `json:"pipelines"`
	Priority                  int32              `json:"priority"`
	RetryConfig               RetryConfig        `json:"retry_config"`
	Status                    string             `json:"status"`
	Timeout                   int64              `json:"timeout"`
	Type                      string             `json:"type"`
	UpstreamBlocks            []string           `json:"upstream_blocks"`
	UUID                      string             `json:"uuid"`
}

type BlockConfiguration struct {
	DataProvider         string `json:"data_provider"`
	DataProviderDatabase string `json:"data_provider_database"`
	DataProviderProfile  string `json:"data_provider_profile"`
	DataProviderSchema   string `json:"data_provider_schema"`
	DataProviderTable    string `json:"data_provider_table"`
	ExportWritePolicy    string `json:"export_write_policy"`
	UseRawSql            string `json:"use_raw_sql"`
}

type RetryConfig struct {
	Delay              int32 `json:"delay"`
	ExponentialBackoff bool  `json:"exponential_backoff"`
	MaxDelay           int32 `json:"max_delay"`
	Retries            int32 `json:"retries"`
}
