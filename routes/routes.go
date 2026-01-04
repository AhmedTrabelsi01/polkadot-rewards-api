package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"polkadot-rewards-api/config"
	"polkadot-rewards-api/controllers"
	"polkadot-rewards-api/services"
)

func SetupRoutes(router *gin.Engine, cfg *config.Config) {
	subscanClient := services.NewSubscanClient(cfg)
	rewardsController := controllers.NewRewardsController(cfg, subscanClient)

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "healthy"})
	})

	validators := router.Group("/validators/:validatorAddress/rewards")
	{
		validators.GET("/eras", rewardsController.GetRewardsPerEra)
		validators.GET("/daily", rewardsController.GetDailyRewards)
		validators.GET("/monthly", rewardsController.GetMonthlyRewards)
	}
}
