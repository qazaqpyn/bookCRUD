package model

import "time"

type Msg struct {
	Action    string
	Entity    string
	EntityID  string
	Timestamp time.Time
}
