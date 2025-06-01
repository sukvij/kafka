package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Controller handles HTTP requests for user CRUD
type Controller struct {
	DB      *gorm.DB
	Service *Service
}

// NewController creates a new user controller
func NewController(db *gorm.DB) *Controller {
	return &Controller{
		DB:      db,
		Service: NewService(db),
	}
}

// CreateUser handles POST /users
func (ctrl *Controller) CreateUser(c *gin.Context) {
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := ctrl.Service.CreateUser(&user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	c.JSON(http.StatusCreated, user)
}

// GetUser handles GET /users/:id
func (ctrl *Controller) GetUser(c *gin.Context) {
	user, err := ctrl.Service.GetUser(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	c.JSON(http.StatusOK, user)
}

// UpdateUser handles PUT /users/:id
func (ctrl *Controller) UpdateUser(c *gin.Context) {
	var input User
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := ctrl.Service.UpdateUser(c.Param("id"), &input)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

// DeleteUser handles DELETE /users/:id
func (ctrl *Controller) DeleteUser(c *gin.Context) {
	if err := ctrl.Service.DeleteUser(c.Param("id")); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "User deleted"})
}
