package detection

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/boobsrate/core/internal/domain"
)

type Client struct {
	baseURL    string
	httpClient *http.Client
}

type detectRequest struct {
	URL string `json:"url"`
}

type detectResponse struct {
	Detections []struct {
		Class string    `json:"class"`
		Score float64   `json:"score"`
		Box   []float64 `json:"box"`
	} `json:"detections"`
}

func NewClient(baseURL string) *Client {
	return &Client{
		baseURL:    baseURL,
		httpClient: &http.Client{},
	}
}

func (c *Client) Detect(ctx context.Context, url string) (domain.DetectionResult, error) {
	reqBody := detectRequest{
		URL: url,
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return domain.DetectionResult{}, fmt.Errorf("marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		c.baseURL+"/detect/url",
		bytes.NewReader(jsonBody),
	)
	if err != nil {
		return domain.DetectionResult{}, fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return domain.DetectionResult{}, fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return domain.DetectionResult{}, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var response detectResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return domain.DetectionResult{}, fmt.Errorf("decode response: %w", err)
	}

	result := domain.DetectionResult{
		Detections: make([]domain.Detection, 0, len(response.Detections)),
	}

	for _, detection := range response.Detections {
		// Convert float64 box coordinates to int
		box := make([]int, len(detection.Box))
		for i, v := range detection.Box {
			box[i] = int(v)
		}

		result.Detections = append(result.Detections, domain.Detection{
			Class: domain.DetectionClass(detection.Class),
			Score: detection.Score,
			Box:   box,
		})
	}

	return result, nil
}
