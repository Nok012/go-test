package entity

import (
	"gorm.io/gorm"
)

type HighestPollutant struct {
	gorm.Model
	Hour      int
	Pollutant string
}
