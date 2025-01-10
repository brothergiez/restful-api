package middlewares

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

var sensitiveFields = []string{"password", "token", "secret", "id_token"}

func anonymizeSensitiveData(data map[string]interface{}) map[string]interface{} {
	for _, field := range sensitiveFields {
		if _, exists := data[field]; exists {
			data[field] = "******"
		}
	}
	return data
}

type CustomResponseWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w *CustomResponseWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func LoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()

		var requestBody map[string]interface{}
		if c.Request.Body != nil {
			bodyBytes, err := ioutil.ReadAll(c.Request.Body)
			if err == nil {
				_ = json.Unmarshal(bodyBytes, &requestBody)
				c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
			}
		}

		anonymizedBody := anonymizeSensitiveData(requestBody)

		customWriter := &CustomResponseWriter{
			ResponseWriter: c.Writer,
			body:           &bytes.Buffer{},
		}
		c.Writer = customWriter

		c.Next()

		var responseData map[string]interface{}
		_ = json.Unmarshal(customWriter.body.Bytes(), &responseData)

		anonymizedResponse := anonymizeSensitiveData(responseData)

		requestHeadersJSON, _ := json.Marshal(c.Request.Header)
		responseHeadersJSON, _ := json.Marshal(c.Writer.Header())

		logEntry := map[string]interface{}{
			"time":            startTime.Format(time.RFC3339),
			"level":           "info",
			"method":          c.Request.Method,
			"path":            c.Request.URL.Path,
			"requestHeaders":  json.RawMessage(requestHeadersJSON),
			"responseHeaders": json.RawMessage(responseHeadersJSON),
			"RequestBody":     anonymizedBody,
			"responseBody":    anonymizedResponse,
			"status":          c.Writer.Status(),
			"duration_ms":     time.Since(startTime).Milliseconds(),
		}

		logEntryJSON, _ := json.Marshal(logEntry)

		currentFlags := log.Flags()
		log.SetFlags(0)
		log.Println()
		log.SetFlags(currentFlags)

		log.Println(string(logEntryJSON))
	}
}
