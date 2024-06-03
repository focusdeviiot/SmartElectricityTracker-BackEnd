package models

type UserDeviceFromDB struct {
	UserID   string `json:"user_id"`
	DeviceID string `json:"device_id"`
}

type UpdateUserDeviceReq struct {
	UserID   string   `json:"user_id"`
	DeviceID []string `json:"device_id"`
}
