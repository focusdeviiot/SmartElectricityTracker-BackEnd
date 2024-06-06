package handlers

import (
	"smart_electricity_tracker_backend/internal/config"
	"smart_electricity_tracker_backend/internal/helpers"
	"smart_electricity_tracker_backend/internal/services"

	"github.com/gofiber/fiber/v2"
)

type ReportHandler struct {
	reportService *services.ReportService
	cfg           *config.Config
}

func NewReportHandler(reportService *services.ReportService, cfg *config.Config) *ReportHandler {
	return &ReportHandler{reportService: reportService, cfg: cfg}
}

func (h *ReportHandler) GetReport(c *fiber.Ctx) error {
	var body struct {
		Device_id string `json:"device_id"`
		DateFrom  string `json:"date_from"`
		DateTo    string `json:"date_to"`
	}

	if err := c.BodyParser(&body); err != nil {
		return helpers.ErrorResponse(c, fiber.StatusBadRequest, "Cannot parse JSON")
	}

	reports, err := h.reportService.GetReportByDeviceAndDate(&body.Device_id, &body.DateFrom, &body.DateTo)
	if err != nil {
		return helpers.ErrorResponse(c, fiber.StatusInternalServerError, "Cannot get report")
	}

	return helpers.SuccessResponse(c,
		fiber.StatusOK,
		"Get report successful",
		reports,
	)
}
