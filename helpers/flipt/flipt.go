package helpers

import (
	"context"
	"fmt"

	"github.com/Improwised/golang-api/config"
	flipt "go.flipt.io/flipt-grpc"
	"google.golang.org/grpc"
)

type BooleanFlagResponse struct {
	Key         string `json:"key"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Enabled     bool   `json:"enabled,omitempty"` // Use a pointer to handle optional field
}

type Context struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type VarientFlagResponse struct {
	RequestContext Context `json:"request_context"`
	Match          bool    `json:"match,omitempty"` // Use pointer to handle optional field
	FlagKey        string  `json:"flag_key"`
	SegmentKey     string  `json:"segment_key,omitempty"` // Use pointer to handle optional field
	Value          string  `json:"value,omitempty"`       // Use pointer to handle optional field

}

func FliptConnection(cfg config.FliptConfig) (*flipt.FliptClient, error) {

	if !cfg.Enabled {
		return nil, fmt.Errorf("flipt is not enabled")
	}

	fliptServer := cfg.Host + ":" + cfg.Port
	conn, err := grpc.Dial(fliptServer, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	client := flipt.NewFliptClient(conn)
	return &client, nil
}

func GetBooleanFlag(fc flipt.FliptClient, bool_flag string) (BooleanFlagResponse, error) {
	var response BooleanFlagResponse
	flagResp, err := fc.GetFlag(context.Background(), &flipt.GetFlagRequest{
		Key: bool_flag,
	})
	if err != nil {
		return BooleanFlagResponse{}, err
	}

	if flagResp == nil {
		return BooleanFlagResponse{}, fmt.Errorf("flag response is nil")
	}

	response.Name = flagResp.Name
	response.Key = flagResp.Key
	response.Description = flagResp.Description
	if flagResp.Enabled {
		response.Enabled = flagResp.Enabled
	}
	return response, nil
}

func GetVarientFlag(fc flipt.FliptClient, flagKey string, entityId string, contextMap map[string]string) (VarientFlagResponse, error) {
	var response VarientFlagResponse
	resp, err := fc.Evaluate(context.Background(), &flipt.EvaluationRequest{
		FlagKey:  flagKey,
		EntityId: entityId,
		Context:  contextMap,
	})
	if err != nil {
		return VarientFlagResponse{}, err
	}

	if resp == nil {
		return VarientFlagResponse{}, fmt.Errorf("flag is not found")
	}

	response.FlagKey = resp.FlagKey
	response.RequestContext = Context{Key: flagKey, Value: flagKey}

	if resp.Match {
		response.Match = resp.Match
		response.SegmentKey = resp.SegmentKey
		response.Value = resp.Value
	}

	return response, nil
}
