package middlewares

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestLoggingMiddleware(t *testing.T) {
	var logBuffer bytes.Buffer
	log.SetOutput(&logBuffer)
	defer log.SetOutput(ioutil.Discard)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(LoggingMiddleware())

	router.POST("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	requestBody := map[string]interface{}{
		"username": "admin",
		"password": "12345",
	}
	bodyBytes, _ := json.Marshal(requestBody)

	req := httptest.NewRequest(http.MethodPost, "/test", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")

	resp := httptest.NewRecorder()
	startTime := time.Now()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)

	logLines := logBuffer.String()

	logLines = logLines[strings.Index(logLines, "{"):]
	var logEntry map[string]interface{}
	err := json.Unmarshal([]byte(logLines), &logEntry)
	assert.NoError(t, err)

	assert.Equal(t, "POST", logEntry["method"])
	assert.Equal(t, "/test", logEntry["path"])
	assert.Equal(t, float64(200), logEntry["status"])

	expectedRequestBody := map[string]interface{}{
		"username": "admin",
		"password": "******",
	}
	assert.Equal(t, expectedRequestBody, logEntry["RequestBody"])

	expectedResponseBody := map[string]interface{}{
		"message": "success",
	}
	assert.Equal(t, expectedResponseBody, logEntry["responseBody"])

	durationMs, ok := logEntry["duration_ms"].(float64)
	assert.True(t, ok)
	assert.GreaterOrEqual(t, durationMs, float64(0))

	timeStr, ok := logEntry["time"].(string)
	assert.True(t, ok)
	parsedTime, err := time.Parse(time.RFC3339, timeStr)
	assert.NoError(t, err)
	assert.WithinDuration(t, startTime, parsedTime, time.Second)
}
