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
		// results, err := p.client.ReadInputRegisters(0, 16)
		// if err != nil {
		// 	log.Println(err)
		// 	time.Sleep(5 * time.Second)
		// 	continue
		// }

		readRegisters(p.client, 30001, "Voltage")
		readRegisters(p.client, 30007, "Current")
		readRegisters(p.client, 30013, "Active power")
		readRegisters(p.client, 30019, "Apparent power")
		readRegisters(p.client, 30025, "Reactive power")
		readRegisters(p.client, 30031, "Power factor")
		readRegisters(p.client, 30071, "Frequency")
		readRegisters(p.client, 30073, "Import active energy")
		readRegisters(p.client, 30075, "Export active energy")
		readRegisters(p.client, 30077, "Import reactive energy")
		readRegisters(p.client, 30079, "Export reactive energy")

		// voltage := Float32FromBytes(results[0:4])
		// current := Float32FromBytes(results[6:10])
		// power := Float32FromBytes(results[12:16])
		// kwh := Float32FromBytes(results[18:22])
		// kwg := Float32FromBytes(results[25:29])

		// err = p.usageRepository.StoreElectricityUsage(voltage, current, power, 1) // 1: สมมติว่าเป็น user ID
		// if err != nil {
		// 	log.Println(err)
		// }

		// data := map[string]float32{
		// 	"voltage": voltage,
		// 	"current": current,
		// 	"power":   power,
		// }

		// Call the broadcast function
		// broadcastFunc(data)

		// logf.Info("Voltage: ", voltage)
		// logf.Info("Current: ", current)
		// logf.Info("Power: ", power)
		// logf.Info("kWh: ", kwh)
		// logf.Info("kWg: ", kwg)

		time.Sleep(1 * time.Second)
	}
}

func readRegisters(client modbus.Client, address uint16, paramName string) {
	quantity := uint16(2) // Each float is 4 bytes, which is 2 registers
	results, err := client.ReadInputRegisters(address-30001, quantity)
	if err != nil {
		logf.Infof("Error reading %s: %v\n", paramName, err)
		return
	}

	if len(results) != 4 {
		logf.Infof("Unexpected length of results for %s: %v\n", paramName, results)
		return
	}

	// Convert bytes to float32
	value := Float32FromBytes(results)
	logf.Infof("%s: %f\n", paramName, value)
}

func float32frombytes(bytes []byte) float32 {
	bits := uint32(bytes[3]) | uint32(bytes[2])<<8 | uint32(bytes[1])<<16 | uint32(bytes[0])<<24
	return float32(bits)
}

func Float32FromBytes(bytes []byte) float32 {
	bits := binary.BigEndian.Uint32(bytes)
	float := math.Float32frombits(bits)
	return float
}
