package mcapiloader

import "time"

type Store interface{}

type Params interface {
	SetApikey(string)
	SetHash(string)
	SetTs(string)
	SetTimeout(time.Duration)
}
