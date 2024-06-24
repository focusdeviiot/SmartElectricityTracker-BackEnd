package routes

import (
	"smart_electricity_tracker_backend/internal/config"
	"smart_electricity_tracker_backend/internal/external"
	"smart_electricity_tracker_backend/internal/handlers"
	"smart_electricity_tracker_backend/internal/middleware"
	"smart_electricity_tracker_backend/internal/models"
	"smart_electricity_tracker_backend/internal/repositories"
	"smart_electricity_tracker_backend/internal/services"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/websocket/v2"

	"gorm.io/gorm"
)

func Setup(app *fiber.App, cfg *config.Config, db *gorm.DB) {
	authMiddleware := middleware.NewAuthMiddleware(cfg)

	// dependencies
	log.Info("Setting up dependencie")
	userRepo := repositories.NewUserRepository(db)
	refreshTokenRepo := repositories.NewRefreshTokenRepository(db)
	reportRepo := repositories.NewReportRepository(db)
	deviceRepo := repositories.NewDeviceRepository(db)

	userService := services.NewUserService(userRepo, refreshTokenRepo, cfg.JWTSecret, cfg.JWTExpiration, cfg.RefreshTokenExpiration, cfg)
	reportService := services.NewReportService(reportRepo, cfg)
	deivceService := services.NewMasterDeviceService(deviceRepo)

	userHandler := handlers.NewUserHandler(userService, cfg)
	reportHandler := handlers.NewReportHandler(reportService, cfg)
	deviceHandler := handlers.NewDeviceHandler(deivceService, userService, cfg)
	syncTimeHandler := handlers.NewSyncTimeHandler(cfg)

	wsHandler := external.NewWebSocketHandler(userRepo, cfg)
	wsTimeHandler := handlers.NewWSTimeHandler()

	log.Info("Starting power meter service")
	powerMeterService, err := services.NewPowerMeterService(cfg, reportRepo, wsHandler)
	if err != nil {
		log.Fatal(err)
	}
	go wsHandler.Start()

	log.Info("Reading and storing power data")
	go powerMeterService.ReadAndStorePowerData()
	go powerMeterService.Broadcast()
	go powerMeterService.RecordData()
	go wsTimeHandler.StartBroadcast()

	log.Info("Setting up routes")
	api := app.Group("/api")
	// Authentication
	api.Post("/login", userHandler.Login)
	api.Post("/logout", userHandler.Logout)
	api.Post("/refresh-Token", userHandler.RefreshToken)
	api.Get("/check-token", authMiddleware.Authenticate(), userHandler.CheckToken)

	// Optional
	api.Get("/devices", authMiddleware.Authenticate(), deviceHandler.GetDevice)
	api.Get("/devices-byuser", authMiddleware.Authenticate(), deviceHandler.GetDeviceByUserId)

	// Report
	api.Post("/reports_volt", authMiddleware.Authenticate(), reportHandler.GetReportVolt)
	api.Post("/reports_ampere", authMiddleware.Authenticate(), reportHandler.GetReportAmpere)
	api.Post("/reports_watt", authMiddleware.Authenticate(), reportHandler.GetReportWatt)

	// Admin
	admin := api.Group("/admin", authMiddleware.Authenticate(), authMiddleware.Permission([]models.Role{models.ADMIN}))
	admin.Get("/users", userHandler.GetUsers)
	admin.Get("/user", userHandler.GetUser)
	admin.Post("/user", userHandler.Register)
	admin.Put("/user", userHandler.UpdateUser)
	admin.Delete("/user", userHandler.DeleteUser)
	admin.Post("/users-count-device", userHandler.GetAllUsersCountDevice)
	admin.Get("/users-device", userHandler.GetUserDeviceById)
	admin.Put("/users-device", userHandler.UpdateUserDevice)

	admin.Post("/sync-time", syncTimeHandler.SyncTime)

	// WebSocket endpoint
	ws := app.Group("/ws")
	ws.Get("/real-time", websocket.New(wsHandler.HandleWebSocket))
	ws.Get("/time", websocket.New(wsTimeHandler.HandleWSTime))
}
