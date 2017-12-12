package ltrymgr

import (
	"time"
)

type Ltryif interface {
	GetGameName() string

	GetCurrentExpect() int

	GetOpenCode() []int

	GetOpenCodeStr() string

	GetCurrentOpenTime() time.Time

	GetNextExpect() int

	GetNextOpenTime() time.Time
}
