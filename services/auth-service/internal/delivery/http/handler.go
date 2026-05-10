package http

import (
	"errors"
	"net/http"

	"auth-service/internal/usecase"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	auth *usecase.AuthUsecase
}

func NewAuthHandler(auth *usecase.AuthUsecase) AuthHandler {
	return AuthHandler{auth: auth}
}

func (h AuthHandler) RegisterRoutes(r *gin.Engine, secret string) {
	r.GET("/healthz", h.healthz)

	public := r.Group("/auth")
	public.POST("/register", h.register)
	public.POST("/login", h.login)
	public.POST("/refresh", h.refresh)

	protected := r.Group("/auth")
	protected.Use(AuthMiddleware(secret))
	protected.GET("/me", h.getMe)
	protected.PUT("/me", h.updateMe)
	protected.POST("/logout", h.logout)
}

type registerRequest struct {
	Name     string `json:"name" binding:"required" example:"John Doe"`
	Email    string `json:"email" binding:"required,email" example:"john@example.com"`
	Password string `json:"password" binding:"required,min=8" example:"password123"`
}

// register godoc
// @Summary Register user
// @Tags auth
// @Accept json
// @Produce json
// @Param payload body registerRequest true "Register payload"
// @Success 201 {object} usecase.RegisterOutput
// @Failure 400 {object} map[string]string
// @Failure 409 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /auth/register [post]
func (h AuthHandler) register(c *gin.Context) {
	var req registerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request payload"})
		return
	}

	out, err := h.auth.Register(c.Request.Context(), usecase.RegisterInput{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		switch {
		case errors.Is(err, usecase.ErrEmailAlreadyExists):
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		case errors.Is(err, usecase.ErrInvalidInput):
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}
	c.JSON(http.StatusCreated, out)
}

type loginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// login godoc
// @Summary Login user
// @Tags auth
// @Accept json
// @Produce json
// @Param payload body loginRequest true "Login payload"
// @Success 200 {object} usecase.AuthTokens
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /auth/login [post]
func (h AuthHandler) login(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request payload"})
		return
	}

	out, err := h.auth.Login(c.Request.Context(), usecase.LoginInput{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		switch {
		case errors.Is(err, usecase.ErrInvalidCredentials):
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		case errors.Is(err, usecase.ErrInactiveUser):
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		case errors.Is(err, usecase.ErrInvalidInput):
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	// Set refresh token in HttpOnly cookie
	// Using a long TTL (e.g., 30 days) for the cookie
	c.SetCookie("refresh_token", out.RefreshToken, 30*24*60*60, "/", "", false, true)

	c.JSON(http.StatusOK, out)
}

type refreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// refresh godoc
// @Summary Refresh access token
// @Tags auth
// @Accept json
// @Produce json
// @Param payload body refreshRequest true "Refresh payload"
// @Success 200 {object} usecase.AuthTokens
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /auth/refresh [post]
func (h AuthHandler) refresh(c *gin.Context) {
	var refreshToken string
	var err error

	// Try to get refresh token from cookie first
	refreshToken, err = c.Cookie("refresh_token")
	if err != nil {
		// Fallback to JSON body if cookie is missing
		var req refreshRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request payload: refresh_token is required"})
			return
		}
		refreshToken = req.RefreshToken
	}

	out, err := h.auth.Refresh(c.Request.Context(), usecase.RefreshInput{
		RefreshToken: refreshToken,
	})
	if err != nil {
		switch {
		case errors.Is(err, usecase.ErrInvalidRefreshToken):
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		case errors.Is(err, usecase.ErrInactiveUser):
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		case errors.Is(err, usecase.ErrInvalidInput):
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	// Update refresh token in cookie
	c.SetCookie("refresh_token", out.RefreshToken, 30*24*60*60, "/", "", false, true)

	c.JSON(http.StatusOK, out)
}

type logoutRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// logout godoc
// @Summary Logout user
// @Tags auth
// @Accept json
// @Produce json
// @Param payload body logoutRequest true "Logout payload"
// @Security BearerAuth
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /auth/logout [post]
func (h AuthHandler) logout(c *gin.Context) {
	var refreshToken string
	var err error

	// Try to get refresh token from cookie first
	refreshToken, err = c.Cookie("refresh_token")
	if err != nil {
		// Fallback to JSON body if cookie is missing
		var req logoutRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request payload: refresh_token is required"})
			return
		}
		refreshToken = req.RefreshToken
	}

	if err := h.auth.Logout(c.Request.Context(), usecase.LogoutInput{
		RefreshToken: refreshToken,
	}); err != nil {
		if errors.Is(err, usecase.ErrInvalidInput) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	// Clear refresh token cookie
	c.SetCookie("refresh_token", "", -1, "/", "", false, true)

	c.JSON(http.StatusOK, gin.H{"message": "logout successful"})
}

// getMe godoc
// @Summary Get current user profile
// @Tags auth
// @Produce json
// @Security BearerAuth
// @Success 200 {object} usecase.UserView
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /auth/me [get]
func (h AuthHandler) getMe(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	user, err := h.auth.GetMe(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, user)
}

type updateMeRequest struct {
	Name  *string `json:"name" example:"John Smith"`
	Email *string `json:"email" example:"john.smith@example.com"`
}

// updateMe godoc
// @Summary Update current user profile
// @Tags auth
// @Accept json
// @Produce json
// @Param payload body updateMeRequest true "Update profile payload"
// @Security BearerAuth
// @Success 200 {object} usecase.UserView
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 409 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /auth/me [put]
func (h AuthHandler) updateMe(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req updateMeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request payload"})
		return
	}

	user, err := h.auth.UpdateMe(c.Request.Context(), userID, usecase.UpdateMeInput{
		Name:  req.Name,
		Email: req.Email,
	})
	if err != nil {
		switch {
		case errors.Is(err, usecase.ErrEmailAlreadyExists):
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		case errors.Is(err, usecase.ErrInvalidInput):
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	c.JSON(http.StatusOK, user)
}

// healthz godoc
// @Summary Health check
// @Tags system
// @Produce json
// @Success 200 {object} map[string]string
// @Router /healthz [get]
func (h AuthHandler) healthz(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"service": "auth-service",
		"status":  "ok",
	})
}

