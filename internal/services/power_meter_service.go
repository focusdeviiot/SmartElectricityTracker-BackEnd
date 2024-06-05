package services

import (
	"encoding/binary"
	"math"
	"smart_electricity_tracker_backend/internal/config"
	"time"

	"github.com/goburrow/modbus"
	logf "github.com/gofiber/fiber/v2/log"
)

type PowerMeterService struct {
	client modbus.Client
	// usageRepository  *repositories.ElectricityUsageRepository
}

func NewPowerMeterService(cfg *config.Config) (*PowerMeterService, error) {
	handler := modbus.NewRTUClientHandler(cfg.Devices.USB)
	handler.BaudRate = cfg.Devices.BaudRate
	handler.DataBits = cfg.Devices.DataBits
	handler.Parity = cfg.Devices.Parity
	handler.StopBits = cfg.Devices.StopBits
	handler.SlaveId = byte(cfg.Devices.DEVICE01.SlaveId) // Convert int to byte
	handler.Timeout = cfg.Devices.TimeOut * time.Second

	if err := handler.Connect(); err != nil {
		logf.Info("Error connecting:", err)
		return nil, err
	}
	defer handler.Close()

	client := modbus.NewClient(handler)
	return &PowerMeterService{
		client: client,
		// usageRepository: usageRepo,
	}, nil
}

func (p *PowerMeterService) ReadAndStorePowerData() { //(broadcastFunc func(data interface{})) {
	for {

		address := uint16(30001)
		quantity := uint16(18) // Read all registers from 30001 to 30080 (40 registers)

		results, err := p.client.ReadInputRegisters(address-30001, quantity)
		if err != nil {
			logf.Infof("Error reading registers: %v\n", err)
			return
		}

		values := parseRegisters(results)
		logf.Infof("Voltage: %f\n", values[0])
		logf.Infof("Current: %f\n", values[1])
		logf.Infof("Active power: %f\n", values[2])

		time.Sleep(1 * time.Second)
	}
}

func parseRegisters(results []byte) []float32 {
	values := make([]float32, 3)
	for i := 0; i < 3; i++ {
		start := i * 12
		end := start + 4
		values[i] = Float32FromBytes(results[start:end])
	}
	return values
}

func Float32FromBytes(bytes []byte) float32 {
	bits := binary.BigEndian.Uint32(bytes)
	float := math.Float32frombits(bits)
	return float
}
