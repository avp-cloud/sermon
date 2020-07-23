package app

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/avp-cloud/sermon/internal/pkg/metrics"
	"github.com/avp-cloud/sermon/internal/pkg/models"
)

// GET /services
// Find all services
func FindServices(c *gin.Context) {
	var dbSvcs []models.DBService
	models.DB.Select("id,name,status,metadata,endpoint,up_codes").Find(&dbSvcs)
	c.JSON(http.StatusOK, dbSvcs)
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
	var input models.CreateServiceInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// Create service
	dbSvc := models.DBService{}
	dbSvc.Name = input.Name
	dbSvc.Endpoint = input.Endpoint
	dbSvc.UpCodes = input.UpCodes
	dbSvc.Metadata = input.Metadata
	dbSvc.Status = metrics.GetStatus(input.Endpoint, input.UpCodes)
	met, err := metrics.GetMetrics(input.Endpoint)
	if err != nil {
		log.Println(err)
	}
	dbSvc = models.FormatDBSvcMetrics(dbSvc, met)
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
	var input models.UpdateServiceInput
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
