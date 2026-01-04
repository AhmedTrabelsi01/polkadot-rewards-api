package services

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/sirupsen/logrus"

	"polkadot-rewards-api/config"
	"polkadot-rewards-api/types"
)

type SubscanClient struct {
	httpClient *http.Client
	baseURL    string
	apiKey     string
	maxPages   int
	rowsPerPage int
	logger     *logrus.Entry
}

func NewSubscanClient(cfg *config.Config) *SubscanClient {
	return &SubscanClient{
		httpClient: &http.Client{
			Timeout: time.Duration(cfg.HTTPTimeoutSeconds) * time.Second,
		},
		baseURL:     cfg.SubscanBaseURL,
		apiKey:      cfg.SubscanAPIKey,
		maxPages:    cfg.MaxPages,
		rowsPerPage: cfg.RowsPerPage,
		logger:      logrus.WithField("service", "subscan"),
	}
}

// GetRewards fetches rewards for a validator from Subscan (single page)
func (c *SubscanClient) GetRewards(ctx context.Context, address string, page, row int) (*types.SubscanRewardSlashResponse, error) {
	url := fmt.Sprintf("%s/api/scan/account/reward_slash", c.baseURL)

	reqBody := types.SubscanRewardSlashRequest{
		Address:  address,
		Page:     page,
		Row:      row,
		Category: "Reward",
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, strings.NewReader(string(jsonBody)))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	if c.apiKey != "" {
		req.Header.Set("X-API-Key", c.apiKey)
	}

	c.logger.WithFields(logrus.Fields{
		"address": address,
		"page":    page,
		"row":     row,
	}).Debug("Making request to Subscan API")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var result types.SubscanRewardSlashResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &result, nil
}

// GetAllRewards fetches all rewards for a validator (handles pagination)
func (c *SubscanClient) GetAllRewards(ctx context.Context, address string) (*types.SubscanRewardSlashResponse, error) {
	firstResp, err := c.GetRewards(ctx, address, 0, c.rowsPerPage)
	if err != nil {
		return nil, err
	}

	if firstResp.Code != 0 {
		return nil, fmt.Errorf("subscan API error: %s", firstResp.Message)
	}

	allRewards := firstResp

	totalCount := firstResp.Data.Count
	if totalCount > c.rowsPerPage {
		pages := (totalCount + c.rowsPerPage - 1) / c.rowsPerPage
		if pages > c.maxPages {
			pages = c.maxPages
		}

		c.logger.WithFields(logrus.Fields{
			"address":    address,
			"totalCount": totalCount,
			"pages":      pages,
		}).Debug("Fetching additional pages")

		for page := 1; page < pages; page++ {
			resp, err := c.GetRewards(ctx, address, page, c.rowsPerPage)
			if err != nil {
				c.logger.WithError(err).WithField("page", page).Warn("Failed to fetch page")
				continue
			}
			allRewards.Data.List = append(allRewards.Data.List, resp.Data.List...)
		}
	}

	c.logger.WithFields(logrus.Fields{
		"address":      address,
		"totalFetched": len(allRewards.Data.List),
	}).Debug("Finished fetching rewards")

	return allRewards, nil
}
