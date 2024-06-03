package models

type UserCountDeviceRes struct {
	UserID      string `json:"user_id" gorm:"unique;type:varchar(255);index:;not null;" `
	Username    string `json:"username" gorm:"type:varchar(255);index:;not null;"`
	Name        string `json:"name" gorm:"type:varchar(255);index:;not null;"`
	Role        Role   `json:"role" gorm:"type:varchar(255);index:;not null;"`
	DeviceCount uint   `json:"device_count" gorm:"index;not null;"`
}

type SearchUserCountDeviceListReq struct {
	Username string    `json:"username"`
	Name     string    `json:"name"`
	Role     string    `json:"role"`
	DeviceId string    `json:"device_id"`
	Pageable *Pageable `json:"pageable"`
}

type SearchUserCountDeviceListTest struct {
	Username *string `json:"username"`
	Name     *string `json:"name"`
	Role     *string `json:"role"`
	DeviceId *string `json:"device_id"`
	// Pageable *Pageable `json:"pageable"`
}
