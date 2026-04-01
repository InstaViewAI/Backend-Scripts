package clevertap

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// Event represents a CleverTap event or profile update
type Event struct {
	Identity    string                 `json:"identity"`
	Type        string                 `json:"type"` // "profile" or "event"
	Timestamp   int64                  `json:"ts"`
	ProfileData map[string]interface{} `json:"profileData,omitempty"`
	EventData   map[string]interface{} `json:"evtData,omitempty"`
}

// UploadEventsReq is the request payload for /1/upload
type UploadEventsReq struct {
	Events []Event `json:"d"`
}

// UploadEventsRes represents the CleverTap response
type UploadEventsRes struct {
	Status      string                   `json:"status"`
	Processed   int                      `json:"processed"`
	Unprocessed []map[string]interface{} `json:"unprocessed,omitempty"`
	Error       string                   `json:"error,omitempty"`
}

// Client holds CleverTap credentials and HTTP client
type Client struct {
	baseURL    string
	accountID  string
	passcode   string
	httpClient *http.Client
}

func NewClient(baseURL, accountID, passcode string) *Client {
	return &Client{
		baseURL:    baseURL,
		accountID:  accountID,
		passcode:   passcode,
		httpClient: &http.Client{Timeout: 10 * time.Second},
	}
}

// UploadEvents sends events or profile updates to CleverTap
func (c *Client) UploadEvents(ctx context.Context, reqData *UploadEventsReq) (*UploadEventsRes, error) {
	payload, err := json.Marshal(reqData)
	if err != nil {
		return nil, fmt.Errorf("dev=marshal-failed: %w", err)
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		c.baseURL+"/1/upload",
		bytes.NewBuffer(payload),
	)
	if err != nil {
		return nil, fmt.Errorf("new-request-failed: %w", err)
	}

	req.Header.Set("X-CleverTap-Account-Id", c.accountID)
	req.Header.Set("X-CleverTap-Passcode", c.passcode)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("http-request-failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("response-read-failed: %w", err)
	}

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return nil, fmt.Errorf("clevertap-status=%s body=%s", resp.Status, strings.TrimSpace(string(body)))
	}

	var res UploadEventsRes
	if err := json.Unmarshal(body, &res); err != nil {
		return nil, fmt.Errorf("response-decode-failed: %w", err)
	}

	return &res, nil
}

// CreateUser creates a new CleverTap user (profile push)
func (c *Client) CreateUser(ctx context.Context, identity string, attributes map[string]interface{}) (*UploadEventsRes, error) {
	event := Event{
		Identity:    identity,
		Type:        "profile",
		Timestamp:   time.Now().Unix(),
		ProfileData: attributes, // must be profileData for profile creation
	}

	req := &UploadEventsReq{
		Events: []Event{event},
	}

	return c.UploadEvents(ctx, req)
}

// UpdateUser updates ONLY the provided attributes (partial profile update)
func (c *Client) UpdateUser(ctx context.Context, identity string, updates map[string]interface{}) (*UploadEventsRes, error) {
	event := Event{
		Identity:    identity,
		Type:        "profile",
		Timestamp:   time.Now().Unix(),
		ProfileData: updates, // partial update merges
	}

	req := &UploadEventsReq{
		Events: []Event{event},
	}

	return c.UploadEvents(ctx, req)
}
