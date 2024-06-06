package services

import "smart_electricity_tracker_backend/internal/config"

type BroadcastService struct {
	cfg *config.Config
}

func NewBroadcastService(cfg *config.Config) (*BroadcastService, error) {
	return &BroadcastService{cfg: cfg}, nil
}
