package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"sync"

	"github.com/avp-cloud/sermon/internal/pkg/metrics"
	"github.com/gin-gonic/gin"
	"github.com/jasonlvhit/gocron"

	"github.com/avp-cloud/sermon/internal/app"
	"github.com/avp-cloud/sermon/internal/pkg/models"
)

var pollInterval = uint64(15)

func collectMetrics() {
	var wg sync.WaitGroup
	var dbSvcs []models.DBService
	models.DB.Find(&dbSvcs)
	log.Println(fmt.Sprintf("Commencing periodic sweep for %d services", len(dbSvcs)))
	for _, dbSvc := range dbSvcs {
		go func() {
			wg.Add(1)
			met, err := metrics.GetMetrics(dbSvc.Endpoint)
			if err != nil {
				log.Println(err)
				wg.Done()
				return
			}
			dbSvc.Status = metrics.GetStatus(dbSvc.Endpoint, dbSvc.UpCodes)
			models.FormatDBSvcMetrics(&dbSvc, met)
			models.DB.Model(&dbSvc).Updates(dbSvc)
			wg.Done()
		}()
	}
	wg.Wait()
	log.Println(fmt.Sprintf("Periodic sweep complete for %d services", len(dbSvcs)))
}

func periodicSweep() {
	gocron.Every(pollInterval).Minutes().Do(collectMetrics)
	<-gocron.Start()
}

func main() {
	r := gin.Default()

	// Connect to database
	models.ConnectDatabase()
	go periodicSweep()

	// Routes
	r.GET("/services", app.FindServices)
	r.GET("/services/:id", app.FindService)
	r.POST("/services", app.CreateService)
	r.PATCH("/services/:id", app.UpdateService)
	r.DELETE("/services/:id", app.DeleteService)

	poll := os.Getenv("POLL_INTERVAL")
	if poll != "" {
		pI, err := strconv.ParseUint(poll, 10, 64)
		if err == nil {
			pollInterval = pI
		}
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "80"
	}
	// Run the server
	r.Run(fmt.Sprintf("0.0.0.0:%s", port))
}
