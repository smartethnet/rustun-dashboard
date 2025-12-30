package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/smartethnet/rustun-dashboard/internal/model"
	"github.com/smartethnet/rustun-dashboard/internal/service"
)

type ClusterHandler struct {
	routeService *service.RouteService
}

func NewClusterHandler(routeService *service.RouteService) *ClusterHandler {
	return &ClusterHandler{
		routeService: routeService,
	}
}

// ListClusters godoc
// @Summary List all clusters
// @Description Get all clusters with client counts
// @Tags clusters
// @Accept json
// @Produce json
// @Success 200 {object} model.Response{data=[]model.Cluster}
// @Failure 500 {object} model.ErrorResponse
// @Router /api/clusters [get]
func (h *ClusterHandler) ListClusters(c *gin.Context) {
	clusters, err := h.routeService.GetAllClusters()
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponseWithCode(
			http.StatusInternalServerError,
			"Failed to get clusters",
			err.Error(),
		))
		return
	}

	c.JSON(http.StatusOK, model.SuccessResponse(clusters))
}

// GetCluster godoc
// @Summary Get a cluster
// @Description Get a specific cluster with all its clients
// @Tags clusters
// @Accept json
// @Produce json
// @Param name path string true "Cluster name"
// @Success 200 {object} model.Response{data=object}
// @Failure 404 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /api/clusters/{name} [get]
func (h *ClusterHandler) GetCluster(c *gin.Context) {
	clusterName := c.Param("name")

	cluster, clients, err := h.routeService.GetCluster(clusterName)
	if err != nil {
		c.JSON(http.StatusNotFound, model.ErrorResponseWithCode(
			http.StatusNotFound,
			"Cluster not found",
			err.Error(),
		))
		return
	}

	c.JSON(http.StatusOK, model.SuccessResponse(gin.H{
		"cluster": cluster,
		"clients": clients,
	}))
}

// DeleteCluster godoc
// @Summary Delete a cluster
// @Description Delete a cluster and all its clients
// @Tags clusters
// @Accept json
// @Produce json
// @Param name path string true "Cluster name"
// @Success 200 {object} model.Response
// @Failure 404 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /api/clusters/{name} [delete]
func (h *ClusterHandler) DeleteCluster(c *gin.Context) {
	clusterName := c.Param("name")

	err := h.routeService.DeleteCluster(clusterName)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "cluster not found" {
			statusCode = http.StatusNotFound
		}
		c.JSON(statusCode, model.ErrorResponseWithCode(
			statusCode,
			"Failed to delete cluster",
			err.Error(),
		))
		return
	}

	c.JSON(http.StatusOK, model.SuccessResponse(gin.H{
		"message": "Cluster deleted successfully",
	}))
}
