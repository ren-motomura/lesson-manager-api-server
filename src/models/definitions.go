package models

import (
	"time"
)

type Gender int

const (
	GenderUndefined Gender = iota
	GenderMale
	GenderFemale
)

var TimeNull time.Time = time.Date(1900, 1, 1, 0, 0, 0, 0, time.UTC)
