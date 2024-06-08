package handlers

import (
	"smart_electricity_tracker_backend/internal/config"
	"smart_electricity_tracker_backend/internal/helpers"
	"smart_electricity_tracker_backend/internal/services"

	"github.com/gofiber/fiber/v2"
)

type DeviceHandler struct {
	masterDeviceService *services.MasterDeviceService
	userService         *services.UserService
	cfg                 *config.Config
}

func NewDeviceHandler(masterDeviceService *services.MasterDeviceService, userService *services.UserService, cfg *config.Config) *DeviceHandler {
	return &DeviceHandler{
		masterDeviceService: masterDeviceService,
		userService:         userService,
		cfg:                 cfg,
	}
}

func (h *DeviceHandler) GetDevice(c *fiber.Ctx) error {
	devices, err := h.masterDeviceService.GetAllDevices()
	if err != nil {
		return helpers.ErrorResponse(c, fiber.StatusInternalServerError, "Cannot get devices")
	}

	return helpers.SuccessResponse(c,
		fiber.StatusOK,
		"Get devices successful",
		devices,
	)
}

func (h *DeviceHandler) GetDeviceByUserId(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	userCount, err := h.userService.GetUserDeviceById(&userID)
	if err != nil {
		return helpers.ErrorResponse(c, fiber.StatusInternalServerError, "Cannot get users count device by user id")
	}

	return helpers.SuccessResponse(c,
		fiber.StatusOK,
		"Get users count device by user id successful",
		fiber.Map{
			"data_list": userCount,
		},
	)
}
