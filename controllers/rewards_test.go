package controllers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"polkadot-rewards-api/config"
	"polkadot-rewards-api/types"
)

type mockRewardsClient struct {
	response *types.SubscanRewardSlashResponse
	err      error
}

func (m *mockRewardsClient) GetAllRewards(ctx context.Context, address string) (*types.SubscanRewardSlashResponse, error) {
	return m.response, m.err
}

func testConfig() *config.Config {
	return &config.Config{
		BlocksPerEra: 14400,
		PlanckToDOT:  1e10,
	}
}

func setupRouter(controller *RewardsController) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/validators/:validatorAddress/rewards/eras", controller.GetRewardsPerEra)
	router.GET("/validators/:validatorAddress/rewards/daily", controller.GetDailyRewards)
	router.GET("/validators/:validatorAddress/rewards/monthly", controller.GetMonthlyRewards)
	return router
}

func mockResponse(items []types.SubscanRewardItem) *types.SubscanRewardSlashResponse {
	return &types.SubscanRewardSlashResponse{
		Code:    0,
		Message: "Success",
		Data: types.SubscanRewardSlashData{
			Count: len(items),
			List:  items,
		},
	}
}

func TestGetRewardsPerEra(t *testing.T) {
	tests := []struct {
		name           string
		mockResponse   *types.SubscanRewardSlashResponse
		mockErr        error
		expectedStatus int
		checkBody      func(t *testing.T, body []byte)
	}{
		{
			name: "success with rewards",
			mockResponse: mockResponse([]types.SubscanRewardItem{
				{BlockNum: 14400, Amount: "10000000000", BlockTimestamp: 1609459200},
				{BlockNum: 28800, Amount: "20000000000", BlockTimestamp: 1609545600},
			}),
			expectedStatus: http.StatusOK,
			checkBody: func(t *testing.T, body []byte) {
				var resp types.EraRewardsResponse
				require.NoError(t, json.Unmarshal(body, &resp))
				assert.Equal(t, "testaddress", resp.Validator)
				assert.NotEmpty(t, resp.Rewards)
			},
		},
		{
			name:           "client error",
			mockErr:        errors.New("subscan API error"),
			expectedStatus: http.StatusInternalServerError,
			checkBody: func(t *testing.T, body []byte) {
				var resp types.ErrorResponse
				require.NoError(t, json.Unmarshal(body, &resp))
				assert.Equal(t, "Internal Server Error", resp.Error)
			},
		},
		{
			name:           "empty response",
			mockResponse:   mockResponse([]types.SubscanRewardItem{}),
			expectedStatus: http.StatusOK,
			checkBody: func(t *testing.T, body []byte) {
				var resp types.EraRewardsResponse
				require.NoError(t, json.Unmarshal(body, &resp))
				assert.Empty(t, resp.Rewards)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := &mockRewardsClient{response: tt.mockResponse, err: tt.mockErr}
			controller := NewRewardsController(testConfig(), mockClient)
			router := setupRouter(controller)

			req, _ := http.NewRequest("GET", "/validators/testaddress/rewards/eras", nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			if tt.checkBody != nil {
				tt.checkBody(t, w.Body.Bytes())
			}
		})
	}
}

func TestGetDailyRewards(t *testing.T) {
	tests := []struct {
		name           string
		mockResponse   *types.SubscanRewardSlashResponse
		mockErr        error
		expectedStatus int
		checkBody      func(t *testing.T, body []byte)
	}{
		{
			name: "success with rewards",
			mockResponse: mockResponse([]types.SubscanRewardItem{
				{Amount: "10000000000", BlockTimestamp: 1609459200},
				{Amount: "20000000000", BlockTimestamp: 1609459200},
			}),
			expectedStatus: http.StatusOK,
			checkBody: func(t *testing.T, body []byte) {
				var resp types.DailyRewardsResponse
				require.NoError(t, json.Unmarshal(body, &resp))
				assert.Equal(t, "testaddress", resp.Validator)
				assert.NotEmpty(t, resp.Rewards)
			},
		},
		{
			name:           "client error",
			mockErr:        errors.New("subscan API error"),
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name:           "empty response",
			mockResponse:   mockResponse([]types.SubscanRewardItem{}),
			expectedStatus: http.StatusOK,
			checkBody: func(t *testing.T, body []byte) {
				var resp types.DailyRewardsResponse
				require.NoError(t, json.Unmarshal(body, &resp))
				assert.Empty(t, resp.Rewards)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := &mockRewardsClient{response: tt.mockResponse, err: tt.mockErr}
			controller := NewRewardsController(testConfig(), mockClient)
			router := setupRouter(controller)

			req, _ := http.NewRequest("GET", "/validators/testaddress/rewards/daily", nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			if tt.checkBody != nil {
				tt.checkBody(t, w.Body.Bytes())
			}
		})
	}
}

func TestGetMonthlyRewards(t *testing.T) {
	tests := []struct {
		name           string
		mockResponse   *types.SubscanRewardSlashResponse
		mockErr        error
		expectedStatus int
		checkBody      func(t *testing.T, body []byte)
	}{
		{
			name: "success with rewards",
			mockResponse: mockResponse([]types.SubscanRewardItem{
				{Amount: "10000000000", BlockTimestamp: 1609459200},
				{Amount: "20000000000", BlockTimestamp: 1612137600},
			}),
			expectedStatus: http.StatusOK,
			checkBody: func(t *testing.T, body []byte) {
				var resp types.MonthlyRewardsResponse
				require.NoError(t, json.Unmarshal(body, &resp))
				assert.Equal(t, "testaddress", resp.Validator)
				assert.NotEmpty(t, resp.Rewards)
			},
		},
		{
			name:           "client error",
			mockErr:        errors.New("subscan API error"),
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name:           "empty response",
			mockResponse:   mockResponse([]types.SubscanRewardItem{}),
			expectedStatus: http.StatusOK,
			checkBody: func(t *testing.T, body []byte) {
				var resp types.MonthlyRewardsResponse
				require.NoError(t, json.Unmarshal(body, &resp))
				assert.Empty(t, resp.Rewards)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := &mockRewardsClient{response: tt.mockResponse, err: tt.mockErr}
			controller := NewRewardsController(testConfig(), mockClient)
			router := setupRouter(controller)

			req, _ := http.NewRequest("GET", "/validators/testaddress/rewards/monthly", nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			if tt.checkBody != nil {
				tt.checkBody(t, w.Body.Bytes())
			}
		})
	}
}
