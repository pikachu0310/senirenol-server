package core

import (
	"github.com/pikachu0310/senirenol-server/core/internal/handler"
	docs "github.com/pikachu0310/senirenol-server/docs" // Swagger docs
	"github.com/pikachu0310/senirenol-server/frontend"

	echoSwagger "github.com/swaggo/echo-swagger"

	"github.com/labstack/echo/v4"
)

func SetupRoutes(h *handler.Handler, e *echo.Echo) {
	e.StaticFS("/", frontend.UI)

	// Swagger UI
	docs.SwaggerInfo.Host = "senirenol.trap.games"
	docs.SwaggerInfo.BasePath = "/api/v1"
	docs.SwaggerInfo.Schemes = []string{"https"}
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	v1API := e.Group("/api/v1")

	// ping API
	pingAPI := v1API.Group("/ping")
	{
		pingAPI.GET("", h.Ping)
	}

	// user API
	userAPI := v1API.Group("/users")
	{
		userAPI.POST("", h.RegisterUser)
		userAPI.POST("/update", h.UpdateUserName)
		userAPI.GET("/:userID", h.GetUser)
		userAPI.GET("/:userID/stats", h.GetUserStats)
	}

	// chart API
	chartAPI := v1API.Group("/charts")
	{
		chartAPI.POST("", h.UpsertChart)
		chartAPI.GET("/ranking", h.GetChartRanking)
	}

	// song API
	v1API.GET("/songs/playcount", h.GetSongPlaycountRanking)

	// score API
	v1API.POST("/scores", h.SubmitScore)
}
