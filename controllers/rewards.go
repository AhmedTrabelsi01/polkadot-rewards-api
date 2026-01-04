package controllers

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"polkadot-rewards-api/config"
	"polkadot-rewards-api/helpers"
	"polkadot-rewards-api/types"
)

type RewardsClient interface {
	GetAllRewards(ctx context.Context, address string) (*types.SubscanRewardSlashResponse, error)
}

type RewardsController struct {
	cfg           *config.Config
	rewardsClient RewardsClient
	logger        *logrus.Entry
}

func NewRewardsController(cfg *config.Config, client RewardsClient) *RewardsController {
	return &RewardsController{
		cfg:           cfg,
		rewardsClient: client,
		logger:        logrus.WithField("controller", "rewards"),
	}
}

// GetRewardsPerEra handles GET /validators/:validatorAddress/rewards/eras
func (rc *RewardsController) GetRewardsPerEra(c *gin.Context) {
	validatorAddress := c.Param("validatorAddress")

	if validatorAddress == "" {
		rc.logger.Warn("Missing validator address in request")
		c.JSON(http.StatusBadRequest, types.ErrorResponse{
			Error:   "Bad Request",
			Message: "validator address is required",
		})
		return
	}

	rc.logger.WithField("validator", validatorAddress).Debug("Fetching rewards per era")

	rewards, err := rc.rewardsClient.GetAllRewards(c.Request.Context(), validatorAddress)
	if err != nil {
		rc.logger.WithError(err).WithField("validator", validatorAddress).Error("Failed to fetch rewards")
		c.JSON(http.StatusInternalServerError, types.ErrorResponse{
			Error:   "Internal Server Error",
			Message: "failed to fetch rewards from Subscan",
		})
		return
	}

	eraRewards := helpers.AggregateByEra(rewards.Data.List, rc.cfg)

	response := types.EraRewardsResponse{
		Validator: validatorAddress,
		Rewards:   eraRewards,
	}

	rc.logger.WithFields(logrus.Fields{
		"validator": validatorAddress,
		"count":     len(eraRewards),
	}).Debug("Successfully fetched rewards per era")

	c.JSON(http.StatusOK, response)
}

// GetDailyRewards handles GET /validators/:validatorAddress/rewards/daily
func (rc *RewardsController) GetDailyRewards(c *gin.Context) {
	validatorAddress := c.Param("validatorAddress")

	if validatorAddress == "" {
		rc.logger.Warn("Missing validator address in request")
		c.JSON(http.StatusBadRequest, types.ErrorResponse{
			Error:   "Bad Request",
			Message: "validator address is required",
		})
		return
	}

	rc.logger.WithField("validator", validatorAddress).Debug("Fetching daily rewards")

	rewards, err := rc.rewardsClient.GetAllRewards(c.Request.Context(), validatorAddress)
	if err != nil {
		rc.logger.WithError(err).WithField("validator", validatorAddress).Error("Failed to fetch rewards")
		c.JSON(http.StatusInternalServerError, types.ErrorResponse{
			Error:   "Internal Server Error",
			Message: "failed to fetch rewards from Subscan",
		})
		return
	}

	dailyRewards := helpers.AggregateByDay(rewards.Data.List, rc.cfg)

	response := types.DailyRewardsResponse{
		Validator: validatorAddress,
		Rewards:   dailyRewards,
	}

	rc.logger.WithFields(logrus.Fields{
		"validator": validatorAddress,
		"count":     len(dailyRewards),
	}).Debug("Successfully fetched daily rewards")

	c.JSON(http.StatusOK, response)
}

// GetMonthlyRewards handles GET /validators/:validatorAddress/rewards/monthly
func (rc *RewardsController) GetMonthlyRewards(c *gin.Context) {
	validatorAddress := c.Param("validatorAddress")

	if validatorAddress == "" {
		rc.logger.Warn("Missing validator address in request")
		c.JSON(http.StatusBadRequest, types.ErrorResponse{
			Error:   "Bad Request",
			Message: "validator address is required",
		})
		return
	}

	rc.logger.WithField("validator", validatorAddress).Debug("Fetching monthly rewards")

	rewards, err := rc.rewardsClient.GetAllRewards(c.Request.Context(), validatorAddress)
	if err != nil {
		rc.logger.WithError(err).WithField("validator", validatorAddress).Error("Failed to fetch rewards")
		c.JSON(http.StatusInternalServerError, types.ErrorResponse{
			Error:   "Internal Server Error",
			Message: "failed to fetch rewards from Subscan",
		})
		return
	}

	monthlyRewards := helpers.AggregateByMonth(rewards.Data.List, rc.cfg)

	response := types.MonthlyRewardsResponse{
		Validator: validatorAddress,
		Rewards:   monthlyRewards,
	}

	rc.logger.WithFields(logrus.Fields{
		"validator": validatorAddress,
		"count":     len(monthlyRewards),
	}).Debug("Successfully fetched monthly rewards")

	c.JSON(http.StatusOK, response)
}
