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
	PipelinesAPIPath                     = "pipelines"
	integrationPipelineType pipelineType = "integration"
	pysparkPipelineType     pipelineType = "pyspark"
	pythonPipelineType      pipelineType = "python"
	streamingPipelineType   pipelineType = "streaming"
)

type PipelineAPI interface {
	CreatePipeline(ctx context.Context, pipelineParams *CreatePipelineRequest) (*pipelineResponse, error)
	DeletePipeline(ctx context.Context, uuid *string) error
	ReadPipeline(ctx context.Context, uuid *string) (*pipelineResponse, error)
	ReadPipelines(ctx context.Context) (*pipelinesResponse, error)
	UpdatePipeline(ctx context.Context, uuid *string, pipelineParams *UpdatePipelineRequest) (*pipelineResponse, error)
}

type pipelineType string

type CreatePipelineRequest struct {
	Pipeline PipelineRequest `json:"pipeline"`
}

type UpdatePipelineRequest struct {
	Pipeline PipelineRequest `json:"pipeline"`
}

type PipelineRequest struct {
	Name string       `json:"name"`
	Type pipelineType `json:"type"`
}

func (pt pipelineType) IsValid() bool {
	switch pt {
	case integrationPipelineType, pysparkPipelineType, pythonPipelineType, streamingPipelineType:
		return true
	}
	return false
}

func (c *client) CreatePipeline(ctx context.Context, pipelineRequest *CreatePipelineRequest) (*pipelineResponse, error) {
	if !pipelineRequest.Pipeline.Type.IsValid() {
		return nil, fmt.Errorf("invalid pipeline type: %s", pipelineRequest.Pipeline.Type)
	}

	reqBody, err := json.Marshal(pipelineRequest)
	if err != nil {
		return nil, err
	}

	respBody, err := c.makeAPICall(http.MethodPost, PipelinesAPIPath, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, err
	}

	createPipelineResponse := pipelineResponse{}
	err = json.Unmarshal(respBody, &createPipelineResponse)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling JSON: %w", err)
	}

	if createPipelineResponse.Pipeline.UUID == "" {
		errRes := errorResponse{}
		err = json.Unmarshal(respBody, &errRes)
		if err != nil {
			return &createPipelineResponse, fmt.Errorf("error unmarshalling JSON: %w", err)
		}
		return &createPipelineResponse, fmt.Errorf("error creating pipeline: %s, Status code: %d", errRes.Error.Exception, errRes.Error.Code)
	}
	return &createPipelineResponse, nil
}

func (c *client) DeletePipeline(ctx context.Context, uuid *string) error {
	respBody, err := c.makeAPICall(http.MethodDelete, path.Join(PipelinesAPIPath, *uuid), nil)
	if err != nil {
		return err
	}

	deletePipelineResponse := pipelineResponse{}
	err = json.Unmarshal(respBody, &deletePipelineResponse)
	if err != nil {
		return fmt.Errorf("error unmarshalling JSON: %w", err)
	}

	if deletePipelineResponse.Pipeline.UUID == "" {
		errRes := errorResponse{}
		err = json.Unmarshal(respBody, &errRes)
		if err != nil {
			return fmt.Errorf("error unmarshalling JSON: %w", err)
		}
		return fmt.Errorf("error deleting pipeline: %s, Status code: %d", errRes.Error.Exception, errRes.Error.Code)
	}
	return nil
}

func (c *client) ReadPipeline(ctx context.Context, uuid *string) (*pipelineResponse, error) {
	readPipelineResponse := pipelineResponse{}
	body, err := c.makeAPICall(http.MethodGet, path.Join(PipelinesAPIPath, *uuid), nil)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &readPipelineResponse)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling JSON: %w", err)
	}

	if readPipelineResponse.Pipeline.UUID == "" {
		errRes := errorResponse{}
		err = json.Unmarshal(body, &errRes)
		if err != nil {
			return nil, fmt.Errorf("error unmarshalling JSON: %w", err)
		}
		return nil, fmt.Errorf("error getting pipeline: %s, Status code: %d", errRes.Error.Exception, errRes.Error.Code)
	}
	return &readPipelineResponse, nil
}

func (c *client) ReadPipelines(ctx context.Context) (*pipelinesResponse, error) {

	readPipelinesResponse := pipelinesResponse{}
	body, err := c.makeAPICall(http.MethodGet, PipelinesAPIPath, nil)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &readPipelinesResponse)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling JSON: %w", err)
	}
	return &readPipelinesResponse, nil
}

func (c *client) UpdatePipeline(ctx context.Context, uuid *string, pipelineRequest *UpdatePipelineRequest) (*pipelineResponse, error) {
	if !pipelineRequest.Pipeline.Type.IsValid() {
		return nil, fmt.Errorf("invalid pipeline type: %s", pipelineRequest.Pipeline.Type)
	}

	reqBody, err := json.Marshal(pipelineRequest)
	if err != nil {
		return nil, err
	}

	fmt.Println(string(reqBody))

	respBody, err := c.makeAPICall(http.MethodPut, path.Join(PipelinesAPIPath, *uuid), bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, err
	}

	fmt.Printf("respBody: %s\n", string(respBody))

	updatePipelineResponse := pipelineResponse{}
	err = json.Unmarshal(respBody, &updatePipelineResponse)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling JSON: %w", err)
	}

	if updatePipelineResponse.Pipeline.UUID == "" {
		errRes := errorResponse{}
		err = json.Unmarshal(respBody, &errRes)
		if err != nil {
			return &updatePipelineResponse, fmt.Errorf("error unmarshalling JSON: %w", err)
		}
		return &updatePipelineResponse, fmt.Errorf("error updating pipeline: %s, Status code: %d", errRes.Error.Exception, errRes.Error.Code)
	}
	return &updatePipelineResponse, nil
}
