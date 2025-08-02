package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"

	"portfolio-backend/internal/models"
)

// Success sends a successful API response
func Success(c *gin.Context, data interface{}, message ...string) {
	response := models.APIResponse{
		Data:    data,
		Success: true,
	}

	if len(message) > 0 && message[0] != "" {
		response.Message = &message[0]
	}

	c.JSON(http.StatusOK, response)
}

// Created sends a 201 Created response
func Created(c *gin.Context, data interface{}, message ...string) {
	response := models.APIResponse{
		Data:    data,
		Success: true,
	}

	if len(message) > 0 && message[0] != "" {
		response.Message = &message[0]
	}

	c.JSON(http.StatusCreated, response)
}

// Error sends an error response
func Error(c *gin.Context, statusCode int, err error, message string, details ...map[string]interface{}) {
	errorResponse := models.APIError{
		Error:   err.Error(),
		Message: message,
	}

	if len(details) > 0 {
		errorResponse.Details = details[0]
	}

	// Log the error
	log.Error().
		Err(err).
		Str("path", c.Request.URL.Path).
		Str("method", c.Request.Method).
		Int("status", statusCode).
		Str("message", message).
		Msg("API error response")

	c.JSON(statusCode, errorResponse)
}

// BadRequest sends a 400 Bad Request response
func BadRequest(c *gin.Context, err error, message string, details ...map[string]interface{}) {
	Error(c, http.StatusBadRequest, err, message, details...)
}

// NotFound sends a 404 Not Found response
func NotFound(c *gin.Context, err error, message string, details ...map[string]interface{}) {
	Error(c, http.StatusNotFound, err, message, details...)
}

// InternalServerError sends a 500 Internal Server Error response
func InternalServerError(c *gin.Context, err error, message string, details ...map[string]interface{}) {
	Error(c, http.StatusInternalServerError, err, message, details...)
}

// Unauthorized sends a 401 Unauthorized response
func Unauthorized(c *gin.Context, err error, message string, details ...map[string]interface{}) {
	Error(c, http.StatusUnauthorized, err, message, details...)
}

// Forbidden sends a 403 Forbidden response
func Forbidden(c *gin.Context, err error, message string, details ...map[string]interface{}) {
	Error(c, http.StatusForbidden, err, message, details...)
}

// TooManyRequests sends a 429 Too Many Requests response
func TooManyRequests(c *gin.Context, message string, details ...map[string]interface{}) {
	errorResponse := models.APIError{
		Error:   "rate_limit_exceeded",
		Message: message,
	}

	if len(details) > 0 {
		errorResponse.Details = details[0]
	}

	log.Warn().
		Str("path", c.Request.URL.Path).
		Str("method", c.Request.Method).
		Str("client_ip", c.ClientIP()).
		Msg("Rate limit exceeded")

	c.JSON(http.StatusTooManyRequests, errorResponse)
}

// ValidationError sends a 400 Bad Request response with validation errors
func ValidationError(c *gin.Context, validationErrors map[string]interface{}) {
	errorResponse := models.APIError{
		Error:   "validation_failed",
		Message: "Request validation failed",
		Details: validationErrors,
	}

	log.Warn().
		Str("path", c.Request.URL.Path).
		Str("method", c.Request.Method).
		Interface("validation_errors", validationErrors).
		Msg("Validation error")

	c.JSON(http.StatusBadRequest, errorResponse)
}