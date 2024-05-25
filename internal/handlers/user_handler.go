package handlers

import (
	"smart_electricity_tracker_backend/internal/helpers"
	"smart_electricity_tracker_backend/internal/services"

	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	userService *services.UserService
}

func NewUserHandler(userService *services.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

func (h *UserHandler) Login(c *fiber.Ctx) error {
	var body struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.BodyParser(&body); err != nil {
		return helpers.ErrorResponse(c, fiber.StatusBadRequest, "Cannot parse JSON")
	}

	accessToken, refreshToken, err := h.userService.Authenticate(body.Username, body.Password)
	if err != nil {
		return helpers.ErrorResponse(c, fiber.StatusBadRequest, "Username or password is incorrect")
	}

	return helpers.SuccessResponse(c,
		fiber.StatusOK,
		"Login successful",
		fiber.Map{
			"access_token":  accessToken,
			"refresh_token": refreshToken,
		},
	)
}

func (h *UserHandler) Register(c *fiber.Ctx) error {
	var body struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Role     string `json:"role"`
		Name     string `json:"name"`
	}

	if err := c.BodyParser(&body); err != nil {
		return helpers.ErrorResponse(c, fiber.StatusBadRequest, "Cannot parse JSON")
	}

	if err := h.userService.CreateUser(body.Username, body.Password, body.Role, body.Name); err != nil {
		return helpers.ErrorResponse(c, fiber.StatusInternalServerError, "Cannot create user")
	}

	return helpers.SuccessResponse(c,
		fiber.StatusCreated,
		"Register successful",
		fiber.Map{},
	)
}

func (h *UserHandler) RefreshToken(c *fiber.Ctx) error {
	var body struct {
		RefreshToken string `json:"refresh_token"`
	}

	if err := c.BodyParser(&body); err != nil {
		return helpers.ErrorResponse(c, fiber.StatusBadRequest, "Cannot parse JSON")
	}

	accessToken, newRefreshToken, err := h.userService.RefreshToken(body.RefreshToken)
	if err != nil {
		return helpers.ErrorResponse(c, fiber.StatusUnauthorized, "Invalid refresh token")
	}

	return helpers.SuccessResponse(c,
		fiber.StatusOK,
		"Refresh token successful",
		fiber.Map{
			"access_token":  accessToken,
			"refresh_token": newRefreshToken,
		},
	)
}

func (h *UserHandler) Logout(c *fiber.Ctx) error {
	var body struct {
		RefreshToken string `json:"refresh_token"`
	}

	if err := c.BodyParser(&body); err != nil {
		return helpers.ErrorResponse(c, fiber.StatusBadRequest, "Cannot parse JSON")
	}

	if err := h.userService.Logout(body.RefreshToken); err != nil {
		return helpers.ErrorResponse(c, fiber.StatusInternalServerError, "Logout failed")
	}

	return helpers.SuccessResponse(c,
		fiber.StatusOK,
		"Logout successful",
		fiber.Map{},
	)
}

func (h *UserHandler) CheckToken(c *fiber.Ctx) error {
	return helpers.SuccessResponse(c,
		fiber.StatusOK,
		"Token valid",
		fiber.Map{},
	)
}
