package middleware

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"vacaytracker-api/internal/dto"
)

func TestErrorMiddleware_NormalRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.Use(ErrorMiddleware())
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var body map[string]interface{}
	err := json.Unmarshal(rec.Body.Bytes(), &body)
	require.NoError(t, err)
	assert.Equal(t, "ok", body["status"])
}

func TestErrorMiddleware_Panic(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.Use(ErrorMiddleware())
	router.GET("/test", func(c *gin.Context) {
		panic("something went horribly wrong")
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)

	var body map[string]interface{}
	err := json.Unmarshal(rec.Body.Bytes(), &body)
	require.NoError(t, err)
	assert.Equal(t, "INTERNAL_ERROR", body["code"])
	assert.Equal(t, "An internal error occurred", body["message"])
}

func TestErrorMiddleware_AppErrorSetViaGinError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.Use(ErrorMiddleware())
	router.GET("/test", func(c *gin.Context) {
		appErr := dto.NewAppError("CUSTOM_ERROR", "Something custom failed", http.StatusBadRequest)
		_ = c.Error(appErr)
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)

	var body map[string]interface{}
	err := json.Unmarshal(rec.Body.Bytes(), &body)
	require.NoError(t, err)
	assert.Equal(t, "CUSTOM_ERROR", body["code"])
	assert.Equal(t, "Something custom failed", body["message"])
}

func TestErrorMiddleware_GenericErrorSetViaGinError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.Use(ErrorMiddleware())
	router.GET("/test", func(c *gin.Context) {
		_ = c.Error(assert.AnError) // testify's generic error
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)

	var body map[string]interface{}
	err := json.Unmarshal(rec.Body.Bytes(), &body)
	require.NoError(t, err)
	assert.Equal(t, "INTERNAL_ERROR", body["code"])
}

func TestErrorMiddleware_PanicWithNilValue(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.Use(ErrorMiddleware())
	router.GET("/test", func(c *gin.Context) {
		panic(nil)
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rec := httptest.NewRecorder()

	// This should not itself panic — the middleware recovers from it.
	// In Go 1.21+ panic(nil) is treated as a real panic with a *runtime.PanicNilError.
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

func TestErrorMiddleware_MultipleErrors(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.Use(ErrorMiddleware())
	router.GET("/test", func(c *gin.Context) {
		// Add multiple errors; the middleware processes the last one
		_ = c.Error(assert.AnError)
		appErr := dto.NewAppError("LAST_ERROR", "The last error wins", http.StatusUnprocessableEntity)
		_ = c.Error(appErr)
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)

	var body map[string]interface{}
	err := json.Unmarshal(rec.Body.Bytes(), &body)
	require.NoError(t, err)
	assert.Equal(t, "LAST_ERROR", body["code"])
}

func TestErrorMiddleware_HandlerAlreadyWroteResponse(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.Use(ErrorMiddleware())
	router.GET("/test", func(c *gin.Context) {
		// Handler writes its own response; no errors set
		c.JSON(http.StatusCreated, gin.H{"id": "123"})
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusCreated, rec.Code)

	var body map[string]interface{}
	err := json.Unmarshal(rec.Body.Bytes(), &body)
	require.NoError(t, err)
	assert.Equal(t, "123", body["id"])
}
