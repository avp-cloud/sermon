package app

import (
	"github.com/avp-cloud/sermon/internal/pkg/metrics"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/avp-cloud/sermon/internal/pkg/models"
)

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

// GET /services
// Find all services
func FindServices(c *gin.Context) {
	var dbSvcs []models.DBService
	models.DB.Find(&dbSvcs)
	var services []models.Service
	for _, dbSvc := range dbSvcs {
		services = append(services, models.DBSvcToSvc(dbSvc))
	}
	c.JSON(http.StatusOK, services)
}

// GET /services/:id
// Find a service
func FindService(c *gin.Context) {
	// Get model if exist
	var dbSvc models.DBService
	if err := models.DB.Where("id = ?", c.Param("id")).First(&dbSvc).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}

	c.JSON(http.StatusOK, models.DBSvcToSvc(dbSvc))
}

// POST /services
// Create new service
func CreateService(c *gin.Context) {
	// Validate input
	var input CreateServiceInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// Create service
	dbSvc := models.DBService{
		Name:     input.Name,
		Endpoint: input.Endpoint,
		UpCodes:  input.UpCodes,
		Metadata: input.Metadata,
		Status:   metrics.GetStatus(input.Endpoint, input.UpCodes),
	}
	met, err := metrics.GetMetrics(input.Endpoint)
	if err != nil {
		log.Println(err)
	}
	models.FormatDBSvcMetrics(&dbSvc, met)
	models.DB.Create(&dbSvc)

	c.JSON(http.StatusOK, models.DBSvcToSvc(dbSvc))
}

// PATCH /services/:id
// Update a service
func UpdateService(c *gin.Context) {
	// Get model if exist
	var dbSvc models.DBService
	if err := models.DB.Where("id = ?", c.Param("id")).First(&dbSvc).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}

	// Validate input
	var input UpdateServiceInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	models.DB.Model(&dbSvc).Updates(input)

	c.JSON(http.StatusOK, models.DBSvcToSvc(dbSvc))
}

// DELETE /services/:id
// Delete a service
func DeleteService(c *gin.Context) {
	// Get model if exist
	var dbSvc models.DBService
	if err := models.DB.Where("id = ?", c.Param("id")).First(&dbSvc).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}

	models.DB.Delete(&dbSvc)

	c.JSON(http.StatusOK, "success")
}
