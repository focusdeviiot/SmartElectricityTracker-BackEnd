package helpers

import "fmt"

type ReportType string

const (
	Voltage ReportType = "volt"
	Current ReportType = "ampere"
	Power   ReportType = "watt"
)

func (e *ReportType) Scan(value interface{}) error {
	switch v := value.(type) {
	case []byte:
		*e = ReportType(v)
	case string:
		*e = ReportType(v)
	default:
		return fmt.Errorf("unsupported type: %T", value)
	}
	return nil
}

func StringToReportType(s string) (ReportType, error) {
	switch s {
	case "volt":
		return Voltage, nil
	case "ampere":
		return Current, nil
	case "watt":
		return Power, nil
	default:
		return "", fmt.Errorf("invalid report type: %s", s)
	}
}
