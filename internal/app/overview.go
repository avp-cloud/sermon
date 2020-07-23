package app

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/avp-cloud/sermon/internal/pkg/models"
)

// GET /overview
// Get an overview
func GetOverview(c *gin.Context) {
	c.JSON(http.StatusOK, *models.LiveOverview)
}
