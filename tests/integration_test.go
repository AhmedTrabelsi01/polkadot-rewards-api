package tests

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"polkadot-rewards-api/config"
	"polkadot-rewards-api/routes"
	"polkadot-rewards-api/types"
)

const testValidatorAddress = "14ShUZUYUR35RBZW6uVVt1zXDxmSQddkeDdXf1JkMA6P721N"

func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)

	cfg := &config.Config{
		Port:               "8080",
		Environment:        "test",
		LogLevel:           "debug",
		SubscanBaseURL:     "https://polkadot.api.subscan.io",
		SubscanAPIKey:      "",
		MaxPages:           2,
		RowsPerPage:        10,
		HTTPTimeoutSeconds: 30,
		BlocksPerEra:       14400,
		PlanckToDOT:        1e10,
	}

	router := gin.New()
	router.Use(gin.Recovery())

	routes.SetupRoutes(router, cfg)

	return router
}

// TestGetRewardsPerEra tests the era rewards endpoint
func TestGetRewardsPerEra(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	router := setupTestRouter()

	req, _ := http.NewRequest("GET", "/validators/"+testValidatorAddress+"/rewards/eras", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response types.EraRewardsResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, testValidatorAddress, response.Validator)
	assert.NotEmpty(t, response.Rewards, "Expected non-empty rewards array")

	if len(response.Rewards) > 0 {
		reward := response.Rewards[0]
		assert.Greater(t, reward.Era, 0, "Era should be greater than 0")
		assert.NotEmpty(t, reward.Amount, "Amount should not be empty")
	}
}

// TestGetDailyRewards tests the daily rewards endpoint
func TestGetDailyRewards(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	router := setupTestRouter()

	req, _ := http.NewRequest("GET", "/validators/"+testValidatorAddress+"/rewards/daily", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response types.DailyRewardsResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, testValidatorAddress, response.Validator)
	assert.NotEmpty(t, response.Rewards, "Expected non-empty rewards array")

	// Verify daily reward structure
	if len(response.Rewards) > 0 {
		reward := response.Rewards[0]
		assert.NotEmpty(t, reward.Date, "Date should not be empty")
		assert.Regexp(t, `^\d{4}-\d{2}-\d{2}$`, reward.Date, "Date should be in YYYY-MM-DD format")
		assert.NotEmpty(t, reward.Amount, "Amount should not be empty")
	}
}

// TestGetMonthlyRewards tests the monthly rewards endpoint
func TestGetMonthlyRewards(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	router := setupTestRouter()

	req, _ := http.NewRequest("GET", "/validators/"+testValidatorAddress+"/rewards/monthly", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response types.MonthlyRewardsResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, testValidatorAddress, response.Validator)
	assert.NotEmpty(t, response.Rewards, "Expected non-empty rewards array")

	// Verify monthly reward structure
	if len(response.Rewards) > 0 {
		reward := response.Rewards[0]
		assert.NotEmpty(t, reward.Month, "Month should not be empty")
		assert.Regexp(t, `^\d{4}-\d{2}$`, reward.Month, "Month should be in YYYY-MM format")
		assert.NotEmpty(t, reward.Amount, "Amount should not be empty")
	}
}