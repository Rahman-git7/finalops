package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/finalops/analytics-service/internal/model"
	"github.com/google/uuid"
)

type WorkoutClient struct {
	baseURL    string
	httpClient *http.Client
}

func NewWorkoutClient(baseURL string) *WorkoutClient {
	return &WorkoutClient{
		baseURL: baseURL,
		httpClient: &http.Client{Timeout: 10 * time.Second},
	}
}

func (c *WorkoutClient) GetSessionsInRange(ctx context.Context, userID uuid.UUID, from, to string) ([]model.SessionRangeItem, error) {
	url := fmt.Sprintf("%s/internal/sessions/range?user_id=%s&from=%s&to=%s",
		c.baseURL, userID.String(), from, to)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("workout-service returned %d", resp.StatusCode)
	}
	var items []model.SessionRangeItem
	if err := json.NewDecoder(resp.Body).Decode(&items); err != nil {
		return nil, err
	}
	return items, nil
}

func (c *WorkoutClient) GetExerciseHistory(ctx context.Context, userID, exerciseID uuid.UUID) ([]model.ExerciseSetHistory, error) {
	url := fmt.Sprintf("%s/internal/sessions/exercise/%s?user_id=%s",
		c.baseURL, exerciseID.String(), userID.String())
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("workout-service returned %d", resp.StatusCode)
	}
	var history []model.ExerciseSetHistory
	if err := json.NewDecoder(resp.Body).Decode(&history); err != nil {
		return nil, err
	}
	return history, nil
}
