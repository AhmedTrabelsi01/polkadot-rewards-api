package helpers

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"polkadot-rewards-api/config"
	"polkadot-rewards-api/types"
)

func testConfig() *config.Config {
	return &config.Config{
		BlocksPerEra: 14400,
		PlanckToDOT:  1e10,
	}
}

func TestCalculateEra(t *testing.T) {
	cfg := testConfig()

	tests := []struct {
		name     string
		blockNum int64
		expected int
	}{
		{"era 0", 0, 0},
		{"era 0 boundary", 14399, 0},
		{"era 1 start", 14400, 1},
		{"era 1 middle", 20000, 1},
		{"era 2", 28800, 2},
		{"high era", 1440000, 100},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CalculateEra(tt.blockNum, cfg)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestParseAmount(t *testing.T) {
	cfg := testConfig()

	tests := []struct {
		name      string
		amountStr string
		expected  float64
	}{
		{"1 DOT", "10000000000", 1.0},
		{"0.5 DOT", "5000000000", 0.5},
		{"10 DOT", "100000000000", 10.0},
		{"zero", "0", 0.0},
		{"invalid string", "invalid", 0.0},
		{"empty string", "", 0.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ParseAmount(tt.amountStr, cfg)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestFormatAmount(t *testing.T) {
	tests := []struct {
		name     string
		amount   float64
		expected string
	}{
		{"whole number", 1.0, "1.0000"},
		{"decimal", 1.5678, "1.5678"},
		{"rounds down", 1.56784, "1.5678"},
		{"rounds up", 1.56786, "1.5679"},
		{"zero", 0.0, "0.0000"},
		{"large number", 1234.5678, "1234.5678"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FormatAmount(tt.amount)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestTimestampToDate(t *testing.T) {
	tests := []struct {
		name      string
		timestamp int64
		expected  string
	}{
		{"2021-01-01", 1609459200, "2021-01-01"},
		{"2021-02-01", 1612137600, "2021-02-01"},
		{"unix epoch", 0, "1970-01-01"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := TimestampToDate(tt.timestamp)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestTimestampToMonth(t *testing.T) {
	tests := []struct {
		name      string
		timestamp int64
		expected  string
	}{
		{"2021-01", 1609459200, "2021-01"},
		{"2021-02", 1612137600, "2021-02"},
		{"unix epoch", 0, "1970-01"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := TimestampToMonth(tt.timestamp)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestAggregateByEra(t *testing.T) {
	cfg := testConfig()

	tests := []struct {
		name     string
		rewards  []types.SubscanRewardItem
		expected []types.EraReward
	}{
		{
			name:     "empty list",
			rewards:  []types.SubscanRewardItem{},
			expected: []types.EraReward{},
		},
		{
			name: "single reward",
			rewards: []types.SubscanRewardItem{
				{BlockNum: 14400, Amount: "10000000000"},
			},
			expected: []types.EraReward{
				{Era: 1, Amount: "1.0000"},
			},
		},
		{
			name: "multiple rewards same era",
			rewards: []types.SubscanRewardItem{
				{BlockNum: 14400, Amount: "10000000000"},
				{BlockNum: 14500, Amount: "20000000000"},
			},
			expected: []types.EraReward{
				{Era: 1, Amount: "3.0000"},
			},
		},
		{
			name: "multiple eras sorted descending",
			rewards: []types.SubscanRewardItem{
				{BlockNum: 14400, Amount: "10000000000"},
				{BlockNum: 28800, Amount: "20000000000"},
				{BlockNum: 43200, Amount: "30000000000"},
			},
			expected: []types.EraReward{
				{Era: 3, Amount: "3.0000"},
				{Era: 2, Amount: "2.0000"},
				{Era: 1, Amount: "1.0000"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := AggregateByEra(tt.rewards, cfg)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestAggregateByDay(t *testing.T) {
	cfg := testConfig()

	tests := []struct {
		name     string
		rewards  []types.SubscanRewardItem
		expected []types.DailyReward
	}{
		{
			name:     "empty list",
			rewards:  []types.SubscanRewardItem{},
			expected: []types.DailyReward{},
		},
		{
			name: "single reward",
			rewards: []types.SubscanRewardItem{
				{BlockTimestamp: 1609459200, Amount: "10000000000"},
			},
			expected: []types.DailyReward{
				{Date: "2021-01-01", Amount: "1.0000"},
			},
		},
		{
			name: "multiple rewards same day",
			rewards: []types.SubscanRewardItem{
				{BlockTimestamp: 1609459200, Amount: "10000000000"},
				{BlockTimestamp: 1609459200, Amount: "20000000000"},
			},
			expected: []types.DailyReward{
				{Date: "2021-01-01", Amount: "3.0000"},
			},
		},
		{
			name: "multiple days sorted descending",
			rewards: []types.SubscanRewardItem{
				{BlockTimestamp: 1609459200, Amount: "10000000000"},
				{BlockTimestamp: 1609545600, Amount: "20000000000"},
			},
			expected: []types.DailyReward{
				{Date: "2021-01-02", Amount: "2.0000"},
				{Date: "2021-01-01", Amount: "1.0000"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := AggregateByDay(tt.rewards, cfg)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestAggregateByMonth(t *testing.T) {
	cfg := testConfig()

	tests := []struct {
		name     string
		rewards  []types.SubscanRewardItem
		expected []types.MonthlyReward
	}{
		{
			name:     "empty list",
			rewards:  []types.SubscanRewardItem{},
			expected: []types.MonthlyReward{},
		},
		{
			name: "single reward",
			rewards: []types.SubscanRewardItem{
				{BlockTimestamp: 1609459200, Amount: "10000000000"},
			},
			expected: []types.MonthlyReward{
				{Month: "2021-01", Amount: "1.0000"},
			},
		},
		{
			name: "multiple rewards same month",
			rewards: []types.SubscanRewardItem{
				{BlockTimestamp: 1609459200, Amount: "10000000000"},
				{BlockTimestamp: 1609545600, Amount: "20000000000"},
			},
			expected: []types.MonthlyReward{
				{Month: "2021-01", Amount: "3.0000"},
			},
		},
		{
			name: "multiple months sorted descending",
			rewards: []types.SubscanRewardItem{
				{BlockTimestamp: 1609459200, Amount: "10000000000"},
				{BlockTimestamp: 1612137600, Amount: "20000000000"},
			},
			expected: []types.MonthlyReward{
				{Month: "2021-02", Amount: "2.0000"},
				{Month: "2021-01", Amount: "1.0000"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := AggregateByMonth(tt.rewards, cfg)
			assert.Equal(t, tt.expected, result)
		})
	}
}
