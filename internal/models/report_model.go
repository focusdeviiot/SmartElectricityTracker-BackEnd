package models

type ReportRes struct {
	ID        int64   `json:"id"`
	DeviceID  string  `json:"device_id"`
	Volt      float32 `json:"volt"`
	Ampere    float32 `json:"ampere"`
	Watt      float32 `json:"watt"`
	CreatedAt string  `json:"created_at"`
}
