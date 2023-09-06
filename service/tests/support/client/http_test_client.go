// Package client is the HTTP client for the service.
package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	gradingAPI "github.com/mnabbasabadi/grading/api/v1"
)

// GradeAPITestClient is the HTTP client for the service.
type GradeAPITestClient struct {
	client *gradingAPI.ClientWithResponses
}

// NewGradingAPITestClient creates a new GradingAPITestClient.
func NewGradingAPITestClient(url string, opts ...gradingAPI.ClientOption) (*GradeAPITestClient, error) {
	client, err := gradingAPI.NewClientWithResponses(fmt.Sprintf("http://%s", url), opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create client: %w", err)
	}
	return &GradeAPITestClient{
		client: client,
	}, nil
}

// GetLiveness returns the liveness of the service.
func (c *GradeAPITestClient) GetLiveness(ctx context.Context) error {
	resp, err := c.client.GetLivenessWithResponse(ctx)
	if err != nil {
		return fmt.Errorf("failed to get liveness: %w", err)
	}
	if resp.StatusCode() != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode())
	}
	return nil
}

// GetGPA ...
func (c *GradeAPITestClient) GetGPA(ctx context.Context, scaleType gradingAPI.ScaleType, limit, offset int) (gradingAPI.GPAResponse, error) {
	resp, err := c.client.GetGPAWithResponse(ctx, &gradingAPI.GetGPAParams{
		ScaleType: (*gradingAPI.GetGPAParamsScaleType)(&scaleType),
		Limit:     &limit,
		Offset:    &offset,
	})
	if err != nil {
		return gradingAPI.GPAResponse{}, fmt.Errorf("failed to get gpa: %w", err)
	}
	if resp.StatusCode() != http.StatusOK {
		return gradingAPI.GPAResponse{}, fmt.Errorf("unexpected status code: %d", resp.StatusCode())
	}
	var ret gradingAPI.GPAResponse
	if err := json.Unmarshal(resp.Body, &ret); err != nil {
		return gradingAPI.GPAResponse{}, fmt.Errorf("failed to unmarshal response: %w", err)
	}
	return ret, nil

}
