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

func (h *ReportHandler) GetReportVolt(c *fiber.Ctx) error {
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

	var volt struct {
		DeviceID  string    `json:"device_id"`
		Volt      []float32 `json:"volt"`
		CreatedAt []string  `json:"created_at"`
	}

	volt.DeviceID = reports[0].DeviceID
	for _, report := range reports {
		volt.Volt = append(volt.Volt, report.Volt)
		volt.CreatedAt = append(volt.CreatedAt, report.CreatedAt.Format("02/01/2006 15:04:05"))
	}

	return helpers.SuccessResponse(c,
		fiber.StatusOK,
		"Get report successful",
		volt,
	)
}

func (h *ReportHandler) GetReportAmpere(c *fiber.Ctx) error {
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

	var ampere struct {
		DeviceID  string    `json:"device_id"`
		Ampere    []float32 `json:"ampere"`
		CreatedAt []string  `json:"created_at"`
	}

	ampere.DeviceID = reports[0].DeviceID
	for _, report := range reports {
		ampere.Ampere = append(ampere.Ampere, report.Ampere)
		ampere.CreatedAt = append(ampere.CreatedAt, report.CreatedAt.Format("02/01/2006 15:04:05"))
	}

	return helpers.SuccessResponse(c,
		fiber.StatusOK,
		"Get report successful",
		ampere,
	)
}

func (h *ReportHandler) GetReportWatt(c *fiber.Ctx) error {
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

	var watt struct {
		DeviceID  string    `json:"device_id"`
		Watt      []float32 `json:"watt"`
		CreatedAt []string  `json:"created_at"`
	}

	watt.DeviceID = reports[0].DeviceID
	for _, report := range reports {
		watt.Watt = append(watt.Watt, report.Watt)
		watt.CreatedAt = append(watt.CreatedAt, report.CreatedAt.Format("02/01/2006 15:04:05"))
	}

	return helpers.SuccessResponse(c,
		fiber.StatusOK,
		"Get report successful",
		watt,
	)
}
