package mageai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"path"
)

const (
	BlockAPIPath                         = "block"
	BlocksAPIPath                        = "blocks"
	callbackBlockType          BlockType = "callback"
	chartBlockType             BlockType = "chart"
	conditionalBlockType       BlockType = "conditional"
	customBlockType            BlockType = "custom"
	dataExporterBlockType      BlockType = "data_exporter"
	dataLoaderblocktype        BlockType = "data_loader"
	dbtBlockType               BlockType = "dbt"
	globalDataProductBlockType BlockType = "global_data_product"
	markdownBlockType          BlockType = "markdown"
	scratchpadBlockType        BlockType = "scratchpad"
	sensorBlockType            BlockType = "sensor"
	transformerBlockType       BlockType = "transformer"
)

type BlockType string

type blockResponse struct {
	Block Block `json:"block"`
}

type blocksResponse struct {
	Blocks []Block `json:"blocks"`
}

type CreateBlockRequest struct {
	Block BlockRequest `json:"block"`
}

type UpdateBlockRequest struct {
	Block BlockRequest `json:"block"`
}

type BlockRequest struct {
	Color          string             `json:"color"`
	Configuration  BlockConfiguration `json:"configuration"`
	Content        string             `json:"content"`
	ExtensionUUID  string             `json:"extension_uuid"`
	Language       string             `json:"language"`
	Name           string             `json:"name"`
	Priority       int32              `json:"priority"`
	RetryConfig    RetryConfig        `json:"retry_config"`
	Type           BlockType          `json:"type"`
	UpstreamBlocks []string           `json:"upstream_blocks"`
}

type BlockAPI interface {
	CreateBlock(ctx context.Context, pipelineUUID *string, blockRequest *CreateBlockRequest) (*blockResponse, error)
	DeleteBlock(ctx context.Context, pipelineUUID *string, blockUUID *string) error
	ReadBlock(ctx context.Context, pipelineUUID *string, blockUUID *string) (*blockResponse, error)
	ReadBlocks(ctx context.Context, pipelineUUID *string) (*blocksResponse, error)
	UpdateBlock(ctx context.Context, pipelineUUID *string, blockUUID *string, blockRequest *UpdateBlockRequest) (*blockResponse, error)
}

func (bt BlockType) IsValid() bool {
	switch bt {
	case callbackBlockType, chartBlockType, conditionalBlockType, customBlockType, dataExporterBlockType, dataLoaderblocktype, dbtBlockType, globalDataProductBlockType, markdownBlockType, scratchpadBlockType, sensorBlockType, transformerBlockType:
		return true
	}
	return false
}

func (c *client) CreateBlock(ctx context.Context, pipelineUUID *string, blockRequest *CreateBlockRequest) (*blockResponse, error) {
	if !blockRequest.Block.Type.IsValid() {
		return nil, fmt.Errorf("invalid block type: %s", blockRequest.Block.Type)
	}

	reqBody, err := json.Marshal(blockRequest)
	if err != nil {
		return nil, err
	}

	respBody, err := c.makeAPICall(http.MethodPost, path.Join(PipelinesAPIPath, *pipelineUUID, BlocksAPIPath), bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, err
	}

	createBlockResponse := blockResponse{}
	err = json.Unmarshal(respBody, &createBlockResponse)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling JSON: %w", err)
	}

	if createBlockResponse.Block.UUID == "" {
		errRes := errorResponse{}
		err = json.Unmarshal(respBody, &errRes)
		if err != nil {
			return nil, fmt.Errorf("error unmarshalling JSON: %w", err)
		}
		return nil, fmt.Errorf("error creating block for the pipeline: %s, Status code: %d", errRes.Error.Exception, errRes.Error.Code)
	}
	return &createBlockResponse, nil
}

func (c *client) DeleteBlock(ctx context.Context, pipelineUUID *string, blockUUID *string) error {
	respBody, err := c.makeAPICall(http.MethodDelete, path.Join(PipelinesAPIPath, *pipelineUUID, BlockAPIPath, *blockUUID), nil)
	if err != nil {
		return err
	}

	deleteBlockResponse := blockResponse{}
	err = json.Unmarshal(respBody, &deleteBlockResponse)
	if err != nil {
		return fmt.Errorf("error unmarshalling JSON: %w", err)
	}

	if deleteBlockResponse.Block.UUID == "" {
		errRes := errorResponse{}
		err = json.Unmarshal(respBody, &errRes)
		if err != nil {
			return fmt.Errorf("error unmarshalling JSON: %w", err)
		}
		return fmt.Errorf("error deleting block: %s, Status code: %d", errRes.Error.Exception, errRes.Error.Code)
	}
	return nil
}

func (c *client) ReadBlock(ctx context.Context, pipelineUUID *string, blockUUID *string) (*blockResponse, error) {
	readBlockResponse := blockResponse{}
	body, err := c.makeAPICall(http.MethodGet, path.Join(PipelinesAPIPath, *pipelineUUID, BlockAPIPath, *blockUUID), nil)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &readBlockResponse)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling JSON: %w", err)
	}

	if readBlockResponse.Block.UUID == "" {
		errRes := errorResponse{}
		err = json.Unmarshal(body, &errRes)
		if err != nil {
			return nil, fmt.Errorf("error unmarshalling JSON: %w", err)
		}
		return nil, fmt.Errorf("error getting the block for the pipeline: %s, Status code: %d", errRes.Error.Exception, errRes.Error.Code)
	}
	return &readBlockResponse, nil
}

func (c *client) ReadBlocks(ctx context.Context, pipelineUUID *string) (*blocksResponse, error) {
	readBlocksResponse := blocksResponse{}
	body, err := c.makeAPICall(http.MethodGet, path.Join(PipelinesAPIPath, *pipelineUUID, BlocksAPIPath), nil)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &readBlocksResponse)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling JSON: %w", err)
	}

	if readBlocksResponse.Blocks == nil {
		errRes := errorResponse{}
		err = json.Unmarshal(body, &errRes)
		if err != nil {
			return nil, fmt.Errorf("error unmarshalling JSON: %w", err)
		}
		return nil, fmt.Errorf("error getting blocks for the pipeline: %s, Status code: %d", errRes.Error.Exception, errRes.Error.Code)
	}
	return &readBlocksResponse, nil
}

func (c *client) UpdateBlock(ctx context.Context, pipelineUUID *string, blockUUID *string, blockRequest *UpdateBlockRequest) (*blockResponse, error) {
	if !blockRequest.Block.Type.IsValid() {
		return nil, fmt.Errorf("invalid block type: %s", blockRequest.Block.Type)
	}

	reqBody, err := json.Marshal(&blockRequest)
	if err != nil {
		return nil, err
	}

	respBody, err := c.makeAPICall(http.MethodPut, path.Join(PipelinesAPIPath, *pipelineUUID, BlocksAPIPath, *blockUUID), bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, err
	}

	updateBlockResponse := blockResponse{}
	err = json.Unmarshal(respBody, &updateBlockResponse)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling JSON: %w", err)
	}

	if updateBlockResponse.Block.UUID == "" {
		errRes := errorResponse{}
		err = json.Unmarshal(respBody, &errRes)
		if err != nil {
			return nil, fmt.Errorf("error unmarshalling JSON: %w", err)
		}
		return nil, fmt.Errorf("error updating block for the pipeline: %s, Status code: %d", errRes.Error.Exception, errRes.Error.Code)
	}
	return &updateBlockResponse, nil
}
