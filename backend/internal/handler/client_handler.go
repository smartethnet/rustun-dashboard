package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/smartethnet/rustun-dashboard/internal/model"
	"github.com/smartethnet/rustun-dashboard/internal/service"
)

type ClientHandler struct {
	routeService *service.RouteService
}

func NewClientHandler(routeService *service.RouteService) *ClientHandler {
	return &ClientHandler{
		routeService: routeService,
	}
}

// ListClients godoc
// @Summary List all clients
// @Description Get all clients across all clusters
// @Tags clients
// @Accept json
// @Produce json
// @Param cluster query string false "Filter by cluster name"
// @Success 200 {object} model.Response{data=[]model.Client}
// @Failure 500 {object} model.ErrorResponse
// @Router /api/clients [get]
func (h *ClientHandler) ListClients(c *gin.Context) {
	clusterFilter := c.Query("cluster")

	var clients []model.Client
	var err error

	if clusterFilter != "" {
		clients, err = h.routeService.GetClientsByCluster(clusterFilter)
	} else {
		clients, err = h.routeService.GetAllClients()
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponseWithCode(
			http.StatusInternalServerError,
			"Failed to get clients",
			err.Error(),
		))
		return
	}

	c.JSON(http.StatusOK, model.SuccessResponse(clients))
}

// GetClient godoc
// @Summary Get a client
// @Description Get a specific client by cluster and identity
// @Tags clients
// @Accept json
// @Produce json
// @Param cluster path string true "Cluster name"
// @Param identity path string true "Client identity"
// @Success 200 {object} model.Response{data=model.Client}
// @Failure 404 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /api/clients/{cluster}/{identity} [get]
func (h *ClientHandler) GetClient(c *gin.Context) {
	cluster := c.Param("cluster")
	identity := c.Param("identity")

	client, err := h.routeService.GetClient(cluster, identity)
	if err != nil {
		c.JSON(http.StatusNotFound, model.ErrorResponseWithCode(
			http.StatusNotFound,
			"Client not found",
			err.Error(),
		))
		return
	}

	c.JSON(http.StatusOK, model.SuccessResponse(client))
}

// CreateClient godoc
// @Summary Create a client
// @Description Add a new client to the configuration
// @Tags clients
// @Accept json
// @Produce json
// @Param client body model.Client true "Client configuration"
// @Success 201 {object} model.Response{data=model.Client}
// @Failure 400 {object} model.ErrorResponse
// @Failure 409 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /api/clients [post]
func (h *ClientHandler) CreateClient(c *gin.Context) {
	var client model.Client
	if err := c.ShouldBindJSON(&client); err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponseWithCode(
			http.StatusBadRequest,
			"Invalid request body",
			err.Error(),
		))
		return
	}

	if err := h.routeService.CreateClient(client); err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "client already exists" {
			statusCode = http.StatusConflict
		}
		c.JSON(statusCode, model.ErrorResponseWithCode(
			statusCode,
			"Failed to create client",
			err.Error(),
		))
		return
	}

	c.JSON(http.StatusCreated, model.SuccessResponse(client))
}

// UpdateClient godoc
// @Summary Update a client
// @Description Update an existing client configuration
// @Tags clients
// @Accept json
// @Produce json
// @Param cluster path string true "Cluster name"
// @Param identity path string true "Client identity"
// @Param client body model.Client true "Updated client configuration"
// @Success 200 {object} model.Response{data=model.Client}
// @Failure 400 {object} model.ErrorResponse
// @Failure 404 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /api/clients/{cluster}/{identity} [put]
func (h *ClientHandler) UpdateClient(c *gin.Context) {
	cluster := c.Param("cluster")
	identity := c.Param("identity")

	var client model.Client
	if err := c.ShouldBindJSON(&client); err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponseWithCode(
			http.StatusBadRequest,
			"Invalid request body",
			err.Error(),
		))
		return
	}

	if err := h.routeService.UpdateClient(cluster, identity, client); err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "client not found" {
			statusCode = http.StatusNotFound
		}
		c.JSON(statusCode, model.ErrorResponseWithCode(
			statusCode,
			"Failed to update client",
			err.Error(),
		))
		return
	}

	c.JSON(http.StatusOK, model.SuccessResponse(client))
}

// DeleteClient godoc
// @Summary Delete a client
// @Description Remove a client from the configuration
// @Tags clients
// @Accept json
// @Produce json
// @Param cluster path string true "Cluster name"
// @Param identity path string true "Client identity"
// @Success 200 {object} model.Response
// @Failure 404 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /api/clients/{cluster}/{identity} [delete]
func (h *ClientHandler) DeleteClient(c *gin.Context) {
	cluster := c.Param("cluster")
	identity := c.Param("identity")

	if err := h.routeService.DeleteClient(cluster, identity); err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "client not found" {
			statusCode = http.StatusNotFound
		}
		c.JSON(statusCode, model.ErrorResponseWithCode(
			statusCode,
			"Failed to delete client",
			err.Error(),
		))
		return
	}

	c.JSON(http.StatusOK, model.SuccessResponse(gin.H{
		"message": "Client deleted successfully",
	}))
}
