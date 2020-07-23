package models

import (
	"encoding/json"
)

var LiveOverview = &Overview{}

type DBService struct {
	// ID specifies the id of the service
	ID uint `json:"id" gorm:"primary_key"`
	// Name specifies the name of service
	Name string `json:"name"`
	// Endpoint specifies the endpoint of service
	Endpoint string `json:"endpoint"`
	// UpCodes specifies the comma separated http codes that signify UP status
	UpCodes string `json:"upCodes"`
	// Tags specifies metadata for the service
	Metadata string `json:"metadata"`
	// Status specifies the status of the endpoint (up/down)
	Status Status `json:"status"`
	// Metrics specifies json string of the last collected metrics
	Metrics string `json:"metrics"`
	// TimeSeriesMetrics json string of last few iteration metrics
	TimeSeriesMetrics string `json:"timeSeriesMetrics"`
}

type Service struct {
	// ID specifies the id of the service
	ID uint `json:"id" gorm:"primary_key"`
	// Name specifies the name of service
	Name string `json:"name"`
	// Endpoint specifies the endpoint of service
	Endpoint string `json:"endpoint"`
	// UpCodes specifies the comma separated http codes that signify UP status
	UpCodes string `json:"upCodes"`
	// Tags specifies metadata for the service
	Metadata string `json:"metadata"`
	// Status specifies the status of the endpoint (up/down)
	Status Status `json:"status"`
	// Metrics specifies the last collected metrics
	Metrics Metrics `json:"metrics"`
	// TimeSeriesMetrics for last few iteration metrics
	TimeSeriesMetrics []Metrics `json:"timeSeriesMetrics"`
}

// Metrics represents a set of metrics
type Metrics struct {
	//Timestamp
	Timestamp string `json:"timeStamp"`
	// DNSTime ...
	DNSTime float64 `json:"dnsTime"`
	// ConnectTime ...
	ConnectTime float64 `json:"connectTime"`
	// TLSTime ...
	TLSTime float64 `json:"tlsTime"`
	// TotalTime ...
	TotalTime float64 `json:"totalTime"`
}

type CreateServiceInput struct {
	// Name specifies the name of service
	Name string `json:"name" binding:"required"`

	// Endpoint specifies the endpoint of service
	Endpoint string `json:"endpoint" binding:"required"`

	// UpCodes specifies the comma separated http codes that signify UP status
	UpCodes string `json:"upCodes" binding:"required"`

	// Tags specifies metadata for the service
	Metadata string `json:"metadata"`
}

type UpdateServiceInput struct {
	// Name specifies the name of service
	Name string `json:"name"`

	// Endpoint specifies the endpoint of service
	Endpoint string `json:"endpoint"`

	// UpCodes specifies the comma separated http codes that signify UP status
	UpCodes string `json:"upCodes"`

	// Tags specifies metadata for the service
	Metadata string `json:"metadata"`
}

type Status string

const (
	StatusDown Status = "DOWN"
	StatusUp   Status = "UP"
)

type Overview struct {
	StatusOverview      StatusOverview      `json:"status"`
	LatencyOverview     LatencyOverview     `json:"latency"`
	ConsistencyOverview ConsistencyOverview `json:"consistency"`
}

type StatusOverview struct {
	Up   int `json:"up"`
	Down int `json:"down"`
}

type LatencyOverview struct {
	High int `json:"high"`
	Low  int `json:"low"`
}

type ConsistencyOverview struct {
	Consistent   int `json:"consistent"`
	Inconsistent int `json:"inconsistent"`
}

func DBSvcToSvc(dbSvc DBService) Service {
	var met Metrics
	var tmet []Metrics
	_ = json.Unmarshal([]byte(dbSvc.Metrics), &met)
	_ = json.Unmarshal([]byte(dbSvc.TimeSeriesMetrics), &tmet)
	return Service{
		ID:                dbSvc.ID,
		Name:              dbSvc.Name,
		Endpoint:          dbSvc.Endpoint,
		UpCodes:           dbSvc.UpCodes,
		Metadata:          dbSvc.Metadata,
		Status:            dbSvc.Status,
		Metrics:           met,
		TimeSeriesMetrics: tmet,
	}
}

func FormatDBSvcMetrics(dbSvc DBService, metrics Metrics) DBService {
	var ms []Metrics
	_ = json.Unmarshal([]byte(dbSvc.TimeSeriesMetrics), &ms)
	if len(ms) == 30 {
		// maintain 30 records
		ms = ms[1:]
	}
	ms = append(ms, metrics)
	metB, _ := json.Marshal(metrics)
	dbSvc.Metrics = string(metB)
	msB, _ := json.Marshal(ms)
	dbSvc.TimeSeriesMetrics = string(msB)
	return dbSvc
}
