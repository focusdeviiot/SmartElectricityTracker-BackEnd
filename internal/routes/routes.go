package routes

import (
	"smart_electricity_tracker_backend/internal/config"
	"smart_electricity_tracker_backend/internal/handlers"
	"smart_electricity_tracker_backend/internal/middleware"
	"smart_electricity_tracker_backend/internal/models"
	"smart_electricity_tracker_backend/internal/repositories"
	"smart_electricity_tracker_backend/internal/services"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	socketio "github.com/googollee/go-socket.io"
	"gorm.io/gorm"
)

func Setup(app *fiber.App, cfg *config.Config, db *gorm.DB) {
	server := socketio.NewServer(nil)
	authMiddleware := middleware.NewAuthMiddleware(cfg)

	log.Info("Starting power meter service")
	powerMeterService, err := services.NewPowerMeterService(cfg) //server,usageRepo)
	if err != nil {
		log.Fatal(err)
	}

	log.Info("Reading and storing power data")
	go powerMeterService.ReadAndStorePowerData()

	// dependencies
	log.Info("Setting up routes")
	userRepo := repositories.NewUserRepository(db)
	refreshTokenRepo := repositories.NewRefreshTokenRepository(db)

	userService := services.NewUserService(userRepo, refreshTokenRepo, cfg.JWTSecret, cfg.JWTExpiration, cfg.RefreshTokenExpiration, cfg)

	userHandler := handlers.NewUserHandler(userService, cfg)

	server.OnConnect("/", func(s socketio.Conn) error {
		s.SetContext("")
		log.Infof("connected:", s.ID())
		return nil
	})

	server.OnError("/", func(s socketio.Conn, e error) {
		log.Infof("meet error:", e)
	})

	server.OnDisconnect("/", func(s socketio.Conn, reason string) {
		log.Infof("closed", reason)
	})
	go server.Serve()
	defer server.Close()

	api := app.Group("/api")
	// Authentication
	api.Post("/login", userHandler.Login)
	api.Post("/logout", userHandler.Logout)
	api.Post("/refresh-Token", userHandler.RefreshToken)
	api.Get("/check-token", authMiddleware.Authenticate(), userHandler.CheckToken)
	// api.Post("/register", userHandler.Register)

	// // Electricity Bill
	// data := api.Group("/data", authMiddleware.Authenticate(), authMiddleware.Permission([]models.Role{models.USER, models.ADMIN}))
	// data.Get("/power-meter", userHandler.GetPowerMeter)
	// data.Get("/electricity-bill", userHandler.GetElectricityBill)

	// Admin
	admin := api.Group("/admin", authMiddleware.Authenticate(), authMiddleware.Permission([]models.Role{models.ADMIN}))
	admin.Get("/users", userHandler.GetUsers)

	admin.Get("/user", userHandler.GetUser)
	admin.Post("/user", userHandler.Register)
	admin.Put("/user", userHandler.UpdateUser)
	admin.Delete("/user", userHandler.DeleteUser)
	// admin.Get("/user/:username", userHandler.GetUserByUsername)

	admin.Post("/users-count-device", userHandler.GetAllUsersCountDevice)
	admin.Get("/users-device", userHandler.GetUserDeviceById)
	admin.Put("/users-device", userHandler.UpdateUserDevice)

	// admin.Get("/user_device", userHandler.GetUserDevices)
	// admin.Get("/user_device/:id", userHandler.GetUserDevice)
	// admin.Post("/user_device", userHandler.CreateUserDevice)
	// admin.Put("/user_device/:id", userHandler.UpdateUserDevice)
	// admin.Delete("/user_device/:id", userHandler.DeleteUserDevice)

	// admin.Get("/electricity-cost", userHandler.GetElectricityCost)
	// admin.Get("/electricity-cost/:id", userHandler.GetElectricityCost)
	// admin.Put("/electricity-cost/:id", userHandler.UpdateElectricityCost)
}
