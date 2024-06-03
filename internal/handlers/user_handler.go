package handlers

import (
	"smart_electricity_tracker_backend/internal/config"
	"smart_electricity_tracker_backend/internal/helpers"
	"smart_electricity_tracker_backend/internal/models"
	"smart_electricity_tracker_backend/internal/services"

	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	userService *services.UserService
	cfg         *config.Config
}

func NewUserHandler(userService *services.UserService, cfg *config.Config) *UserHandler {
	return &UserHandler{userService: userService, cfg: cfg}
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

func (h *UserHandler) GetUsers(c *fiber.Ctx) error {
	users, err := h.userService.GetUsers()
	if err != nil {
		return helpers.ErrorResponse(c, fiber.StatusInternalServerError, "Cannot get users")
	}

	return helpers.SuccessResponse(c,
		fiber.StatusOK,
		"Get users successful",
		fiber.Map{
			"users": users,
		},
	)
}

func (h *UserHandler) GetUser(c *fiber.Ctx) error {
	id := c.Params("id")
	user, err := h.userService.GetUserById(id)
	if err != nil {
		return helpers.ErrorResponse(c, fiber.StatusInternalServerError, "Cannot get user")
	}

	return helpers.SuccessResponse(c,
		fiber.StatusOK,
		"Get user successful",
		fiber.Map{
			"user": user,
		},
	)
}

func (h *UserHandler) UpdateUser(c *fiber.Ctx) error {
	id := c.Params("id")

	usernameAdmin := h.cfg.AdminUser.Username
	userAdmin, err := h.userService.GetUserByUsername(usernameAdmin)
	if err != nil {
		return helpers.ErrorResponse(c, fiber.StatusInternalServerError, "Cannot get main userAdmin")
	}

	if id == userAdmin.ID.String() {
		return helpers.ErrorResponse(c, fiber.StatusBadRequest, "Cannot update main userAdmin")
	}

	var body struct {
		Password string `json:"password"`
		Role     string `json:"role"`
		Name     string `json:"name"`
	}

	if err := c.BodyParser(&body); err != nil {
		return helpers.ErrorResponse(c, fiber.StatusBadRequest, "Cannot parse JSON")
	}

	if err := h.userService.UpdateUser(id, body.Password, body.Role, body.Name); err != nil {
		return helpers.ErrorResponse(c, fiber.StatusInternalServerError, "Cannot update user")
	}

	return helpers.SuccessResponse(c,
		fiber.StatusOK,
		"Update user successful",
		fiber.Map{},
	)
}

func (h *UserHandler) DeleteUser(c *fiber.Ctx) error {
	id := c.Params("id")
	user_id := c.Locals("user_id").(string)
	if id == user_id {
		return helpers.ErrorResponse(c, fiber.StatusBadRequest, "Cannot delete own user")
	}

	usernameAdmin := h.cfg.AdminUser.Username
	userAdmin, err := h.userService.GetUserByUsername(usernameAdmin)
	if err != nil {
		return helpers.ErrorResponse(c, fiber.StatusInternalServerError, "Cannot get main userAdmin")
	}

	if id == userAdmin.ID.String() {
		return helpers.ErrorResponse(c, fiber.StatusBadRequest, "Cannot delete main userAdmin")
	}

	if err := h.userService.DeleteUser(id); err != nil {
		return helpers.ErrorResponse(c, fiber.StatusInternalServerError, "Cannot delete user")
	}

	return helpers.SuccessResponse(c,
		fiber.StatusOK,
		"Delete user successful",
		fiber.Map{},
	)
}

func (h *UserHandler) GetUserByUsername(c *fiber.Ctx) error {
	username := c.Params("username")
	user, err := h.userService.GetUserByUsername(username)
	if err != nil {
		return helpers.ErrorResponse(c, fiber.StatusInternalServerError, "Cannot get user")
	}

	return helpers.SuccessResponse(c,
		fiber.StatusOK,
		"Get user successful",
		fiber.Map{
			"user": user,
		},
	)
}

func (h *UserHandler) GetAllUsersCountDevice(c *fiber.Ctx) error {
	body := models.SearchUserCountDeviceListReq{}

	if err := c.BodyParser(&body); err != nil {
		return helpers.ErrorResponse(c, fiber.StatusBadRequest, "Cannot parse JSON")
	}

	userCount, pageable, err := h.userService.GetAllUsersCountDevice(&body)
	if err != nil {
		return helpers.ErrorResponse(c, fiber.StatusInternalServerError, "Cannot get users count device")
	}

	return helpers.SuccessResponse(c,
		fiber.StatusOK,
		"Get users count device successful",
		fiber.Map{
			"data_list": userCount,
			"pageable":  pageable,
		},
	)
}

func (h *UserHandler) GetUserDeviceById(c *fiber.Ctx) error {
	userID := c.Query("user_id")
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

func (h *UserHandler) UpdateUserDevice(c *fiber.Ctx) error {
	body := models.UpdateUserDeviceReq{}

	if err := c.BodyParser(&body); err != nil {
		return helpers.ErrorResponse(c, fiber.StatusBadRequest, "Cannot parse JSON")
	}

	if body.UserID == "" {
		return helpers.ErrorResponse(c, fiber.StatusBadRequest, "User ID is required")
	}
	if err := h.userService.UpdateUserDevice(&body); err != nil {
		return helpers.ErrorResponse(c, fiber.StatusInternalServerError, "Cannot update user device")
	}

	return helpers.SuccessResponse(c,
		fiber.StatusOK,
		"Update user device successful",
		fiber.Map{},
	)
}
