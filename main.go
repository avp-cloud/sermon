package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jasonlvhit/gocron"

	"github.com/avp-cloud/sermon/internal/app"
	"github.com/avp-cloud/sermon/internal/pkg/metrics"
	"github.com/avp-cloud/sermon/internal/pkg/models"
)

var pollInterval = uint64(1)

func periodicSweep() {
	metrics.CollectMetrics()
	gocron.Every(pollInterval).Minutes().Do(metrics.CollectMetrics)
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
	r.GET("/overview", app.GetOverview)

	poll := os.Getenv("POLL_INTERVAL")
	if poll != "" {
		pI, err := strconv.ParseUint(poll, 10, 64)
		if err == nil {
			pollInterval = pI
		}
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}
	// Run the server
	r.Run(fmt.Sprintf("0.0.0.0:%s", port))
}
