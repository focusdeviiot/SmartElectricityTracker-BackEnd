package models

import "time"

type ReportRes struct {
	ID        int64     `json:"id"`
	DeviceID  string    `json:"device_id"`
	Volt      float32   `json:"volt"`
	Ampere    float32   `json:"ampere"`
	Watt      float32   `json:"watt"`
	CreatedAt time.Time `json:"created_at"`
}
