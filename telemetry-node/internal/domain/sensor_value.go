package domain

import "time"

// SensorValue represent simple version of sensor measurement
type SensorValue struct {
	SensorName string
	Value      int64
	Timestamp  time.Time
}
