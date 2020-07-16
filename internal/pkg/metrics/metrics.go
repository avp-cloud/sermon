package metrics

import (
	"crypto/tls"
	"log"
	"net/http"
	"net/http/httptrace"
	"strings"
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
	if !strings.Contains(upCodes, string(resp.StatusCode)) {
		return models.StatusDown
	}
	return models.StatusUp
}
