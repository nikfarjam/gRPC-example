package logevnt

import "time"

type StartEventRequest struct {
	Guid        string
	Ip          string
	Time        time.Time
	Method      string
	Path        string
	Status      int
	ProcessTime int64
	FullLog     string
}

type StartEventResponse struct {
	Guid string
}

type EndEventRequest struct {
	Guid string
}

type EndEventResponse struct {
	Guid string
}
