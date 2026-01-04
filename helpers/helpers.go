package helpers

import (
	"fmt"
	"sort"
	"strconv"
	"time"

	"polkadot-rewards-api/config"
	"polkadot-rewards-api/types"
)

func CalculateEra(blockNum int64, cfg *config.Config) int {
	return int(blockNum / cfg.BlocksPerEra)
}

func ParseAmount(amountStr string, cfg *config.Config) float64 {
	amount, err := strconv.ParseFloat(amountStr, 64)
	if err != nil {
		return 0
	}
	return amount / cfg.PlanckToDOT
}

func FormatAmount(amount float64) string {
	return fmt.Sprintf("%.4f", amount)
}

func TimestampToDate(timestamp int64) string {
	t := time.Unix(timestamp, 0).UTC()
	return t.Format("2006-01-02")
}

func TimestampToMonth(timestamp int64) string {
	t := time.Unix(timestamp, 0).UTC()
	return t.Format("2006-01")
}

func AggregateByEra(rewards []types.SubscanRewardItem, cfg *config.Config) []types.EraReward {
	eraRewards := make(map[int]float64)

	for _, reward := range rewards {
		era := CalculateEra(reward.BlockNum, cfg)
		amountDOT := ParseAmount(reward.Amount, cfg)
		eraRewards[era] += amountDOT
	}

	result := make([]types.EraReward, 0, len(eraRewards))
	for era, amount := range eraRewards {
		result = append(result, types.EraReward{
			Era:    era,
			Amount: FormatAmount(amount),
		})
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].Era > result[j].Era
	})

	return result
}

func AggregateByDay(rewards []types.SubscanRewardItem, cfg *config.Config) []types.DailyReward {
	dailyRewards := make(map[string]float64)

	for _, reward := range rewards {
		date := TimestampToDate(reward.BlockTimestamp)
		amountDOT := ParseAmount(reward.Amount, cfg)
		dailyRewards[date] += amountDOT
	}

	result := make([]types.DailyReward, 0, len(dailyRewards))
	for date, amount := range dailyRewards {
		result = append(result, types.DailyReward{
			Date:   date,
			Amount: FormatAmount(amount),
		})
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].Date > result[j].Date
	})

	return result
}

func AggregateByMonth(rewards []types.SubscanRewardItem, cfg *config.Config) []types.MonthlyReward {
	monthlyRewards := make(map[string]float64)

	for _, reward := range rewards {
		month := TimestampToMonth(reward.BlockTimestamp)
		amountDOT := ParseAmount(reward.Amount, cfg)
		monthlyRewards[month] += amountDOT
	}

	result := make([]types.MonthlyReward, 0, len(monthlyRewards))
	for month, amount := range monthlyRewards {
		result = append(result, types.MonthlyReward{
			Month:  month,
			Amount: FormatAmount(amount),
		})
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].Month > result[j].Month
	})

	return result
}
