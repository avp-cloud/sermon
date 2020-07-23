package metrics

import (
	"crypto/tls"
	"fmt"
	"log"
	"math"
	"net/http"
	"net/http/httptrace"
	"strings"
	"sync"
	"time"

	"github.com/avp-cloud/sermon/internal/pkg/models"
)

func GetMetrics(url string) (models.Metrics, error) {
	metricsSet := models.Metrics{}
	req, _ := http.NewRequest("GET", url, nil)

	var start, connect, dns, tlsHandshake time.Time

	trace := &httptrace.ClientTrace{
		DNSStart: func(dsi httptrace.DNSStartInfo) { dns = time.Now() },
		DNSDone: func(ddi httptrace.DNSDoneInfo) {
			metricsSet.DNSTime = time.Since(dns).Seconds()
		},

		TLSHandshakeStart: func() { tlsHandshake = time.Now() },
		TLSHandshakeDone: func(cs tls.ConnectionState, err error) {
			metricsSet.TLSTime = time.Since(tlsHandshake).Seconds()
		},

		ConnectStart: func(network, addr string) { connect = time.Now() },
		ConnectDone: func(network, addr string, err error) {
			metricsSet.ConnectTime = time.Since(connect).Seconds()
		},

		GotFirstResponseByte: func() {
			metricsSet.DNSTime = time.Since(start).Seconds()
		},
	}

	req = req.WithContext(httptrace.WithClientTrace(req.Context(), trace))
	start = time.Now()
	if _, err := http.DefaultTransport.RoundTrip(req); err != nil {
		return metricsSet, err
	}
	metricsSet.TotalTime = time.Since(start).Seconds()
	metricsSet.Timestamp = start.String()
	return metricsSet, nil
}

func GetStatus(url, upCodes string) models.Status {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	resp, err := client.Get(url)
	if err != nil {
		log.Println(err)
		return models.StatusDown
	}
	if !strings.Contains(upCodes, fmt.Sprintf("%d", resp.StatusCode)) {
		return models.StatusDown
	}
	return models.StatusUp
}

func computeOverview() {
	var dbSvcs []models.DBService
	models.DB.Find(&dbSvcs)
	log.Println(fmt.Sprintf("Computing service overview for %d services", len(dbSvcs)))
	models.LiveOverview = &models.Overview{}
	for _, dbSvc := range dbSvcs {
		service := models.DBSvcToSvc(dbSvc)

		// Status overview
		if service.Status == models.StatusUp {
			models.LiveOverview.StatusOverview.Up += 1
		} else {
			models.LiveOverview.StatusOverview.Down += 1
		}

		// Latency
		min := 99.0
		max := 0.0
		avgTotalTime := 0.0
		for _, t := range service.TimeSeriesMetrics {
			avgTotalTime += t.TotalTime
			if t.TotalTime > max {
				max = t.TotalTime
			}
			if t.TotalTime < min {
				min = t.TotalTime
			}
		}
		avgTotalTime = avgTotalTime/float64(len(service.TimeSeriesMetrics))
		if avgTotalTime > 3.0 {
			models.LiveOverview.LatencyOverview.High += 1
		} else {
			models.LiveOverview.LatencyOverview.Low += 1
		}

		// Consistency
		if math.Abs(avgTotalTime - min) > .1 || math.Abs(max - avgTotalTime) > .1 {
			models.LiveOverview.ConsistencyOverview.Inconsistent += 1
		} else {
			models.LiveOverview.ConsistencyOverview.Consistent += 1
		}
		log.Println(fmt.Sprintf("Computed service overview for %d services", len(dbSvcs)))
	}
}

func CollectMetrics() {
	var wg sync.WaitGroup
	var dbSvcs []models.DBService
	models.DB.Find(&dbSvcs)
	log.Println(fmt.Sprintf("Commencing periodic sweep for %d services", len(dbSvcs)))
	for _, dbSvc := range dbSvcs {
		dbSvc := dbSvc
		go func() {
			wg.Add(1)
			met, err := GetMetrics(dbSvc.Endpoint)
			if err != nil {
				log.Println(err)
				wg.Done()
				return
			}
			dbSvc.Status = GetStatus(dbSvc.Endpoint, dbSvc.UpCodes)
			dbSvc = models.FormatDBSvcMetrics(dbSvc, met)
			models.DB.Model(dbSvc).Updates(dbSvc)
			wg.Done()
		}()
	}
	wg.Wait()
	log.Println(fmt.Sprintf("Periodic sweep complete for %d services", len(dbSvcs)))
	computeOverview()
}
