package services

import (
	"encoding/binary"
	"math"
	"smart_electricity_tracker_backend/internal/config"
	"smart_electricity_tracker_backend/internal/external"
	"smart_electricity_tracker_backend/internal/models"
	"smart_electricity_tracker_backend/internal/repositories"
	"sync"
	"time"

	"github.com/goburrow/modbus"
	"github.com/gofiber/fiber/v2/log"
)

type PowerMeterService struct {
	client     modbus.Client
	handler    *modbus.RTUClientHandler
	mu         sync.Mutex
	sharedData map[string]map[string]float32
	reportRepo *repositories.ReportRepository
	ws         *external.WebSocketHandler
	cfg        *config.Config
}

func NewPowerMeterService(cfg *config.Config, reportRepo *repositories.ReportRepository, ws *external.WebSocketHandler) (*PowerMeterService, error) {
	handler := modbus.NewRTUClientHandler(cfg.PowerMeter.Device)
	handler.BaudRate = cfg.Devices.BaudRate
	handler.DataBits = cfg.Devices.DataBits
	handler.Parity = cfg.Devices.Parity
	handler.StopBits = cfg.Devices.StopBits
	handler.SlaveId = byte(cfg.Devices.DEVICE01.SlaveId) // Convert int to byte
	handler.Timeout = cfg.Devices.TimeOut * time.Second

	if err := handler.Connect(); err != nil {
		log.Info("Error connecting:", err)
		return nil, err
	}
	defer handler.Close()

	client := modbus.NewClient(handler)

	return &PowerMeterService{
		client:     client,
		handler:    handler,
		sharedData: make(map[string]map[string]float32),
		reportRepo: reportRepo,
		ws:         ws,
		cfg:        cfg,
	}, nil
}

func (p *PowerMeterService) ReadAndStorePowerData() { //(broadcastFunc func(data interface{})) {
	for {
		address := uint16(30001)
		quantity := uint16(18) // Read all registers from 30001 to 30080 (40 registers)

		p.handler.SlaveId = byte(p.cfg.Devices.DEVICE01.SlaveId)
		results1, err := p.client.ReadInputRegisters(address-30001, quantity)
		if err != nil {
			log.Infof("Error reading registers 1: %v\n", err)
			time.Sleep(5 * time.Second)
			continue
		}
		values1 := parseRegisters(results1)
		// log.Infof("Values 1: %v\n", values1)

		time.Sleep(100 * time.Millisecond)
		p.handler.SlaveId = byte(p.cfg.Devices.DEVICE02.SlaveId)
		results2, err := p.client.ReadInputRegisters(address-30001, quantity)
		if err != nil {
			log.Infof("Error reading registers 2: %v\n", err)
			time.Sleep(5 * time.Second)
			continue
		}
		values2 := parseRegisters(results2)
		// log.Infof("Values 2: %v\n", values2)

		time.Sleep(100 * time.Millisecond)
		p.handler.SlaveId = byte(p.cfg.Devices.DEVICE03.SlaveId)
		results3, err := p.client.ReadInputRegisters(address-30001, quantity)
		if err != nil {
			log.Infof("Error reading registers 3: %v\n", err)
			time.Sleep(5 * time.Second)
			continue
		}
		values3 := parseRegisters(results3)
		// log.Infof("Values 3: %v\n", values3)

		p.mu.Lock()
		p.sharedData[p.cfg.Devices.DEVICE01.DeviceId] = map[string]float32{
			"voltage":      float32(math.Abs(float64(values1[0]))),
			"current":      float32(math.Abs(float64(values1[1]))),
			"active_power": float32(math.Abs(float64(values1[2]))),
		}
		p.sharedData[p.cfg.Devices.DEVICE02.DeviceId] = map[string]float32{
			"voltage":      float32(math.Abs(float64(values2[0]))),
			"current":      float32(math.Abs(float64(values2[1]))),
			"active_power": float32(math.Abs(float64(values2[2]))),
		}
		p.sharedData[p.cfg.Devices.DEVICE03.DeviceId] = map[string]float32{
			"voltage":      float32(math.Abs(float64(values3[0]))),
			"current":      float32(math.Abs(float64(values3[1]))),
			"active_power": float32(math.Abs(float64(values3[2]))),
		}
		p.mu.Unlock()

		time.Sleep(p.cfg.Devices.LoopReadTime * time.Second)
	}
}

func (p *PowerMeterService) Broadcast() {
	for {
		nextTick := time.Now().Truncate(p.cfg.Devices.LoopbroadcastTime * time.Second).Add(p.cfg.Devices.LoopbroadcastTime * time.Second)
		time.Sleep(time.Until(nextTick))

		p.mu.Lock()
		data := p.sharedData
		p.mu.Unlock()
		// log.Infof("Broadcasting data: %v\n", data)

		p.ws.Broadcast(data)
	}
}

func (p *PowerMeterService) RecordData() {
	for {
		nextTick := time.Now().Truncate(p.cfg.Devices.LoopRecordTime * time.Second).Add(p.cfg.Devices.LoopRecordTime * time.Second)
		time.Sleep(time.Until(nextTick))

		p.mu.Lock()
		data := p.sharedData
		p.mu.Unlock()

		device01, err := data[p.cfg.Devices.DEVICE01.DeviceId]
		if !err {
			log.Infof("Device %s not found in shared data\n", p.cfg.Devices.DEVICE01.DeviceId)
			continue
		}

		device02, err := data[p.cfg.Devices.DEVICE02.DeviceId]
		if !err {
			log.Infof("Device %s not found in shared data\n", p.cfg.Devices.DEVICE02.DeviceId)
			continue
		}

		device03, err := data[p.cfg.Devices.DEVICE03.DeviceId]
		if !err {
			log.Infof("Device %s not found in shared data\n", p.cfg.Devices.DEVICE03.DeviceId)
			continue
		}

		record := &models.RecodePowermeter{
			DeviceID: p.cfg.Devices.DEVICE01.DeviceId,
			Volt:     device01["voltage"],
			Ampere:   device01["current"],
			Watt:     device01["active_power"],
		}
		if err := p.reportRepo.RecordPowermeter(record); err != nil {
			log.Infof("Error recording power meter data: %v\n", err)
		}

		record = &models.RecodePowermeter{
			DeviceID: p.cfg.Devices.DEVICE02.DeviceId,
			Volt:     device02["voltage"],
			Ampere:   device02["current"],
			Watt:     device02["active_power"],
		}
		if err := p.reportRepo.RecordPowermeter(record); err != nil {
			log.Infof("Error recording power meter data: %v\n", err)
		}

		record = &models.RecodePowermeter{
			DeviceID: p.cfg.Devices.DEVICE03.DeviceId,
			Volt:     device03["voltage"],
			Ampere:   device03["current"],
			Watt:     device03["active_power"],
		}
		if err := p.reportRepo.RecordPowermeter(record); err != nil {
			log.Infof("Error recording power meter data: %v\n", err)
		}
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
