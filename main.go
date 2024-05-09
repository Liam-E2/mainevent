package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)


func main() {
	router := gin.Default()

	// Initialize new streaming server
	stream := NewServer()

	subscribe := router.Group("/events/subscribe/:name")

	subscribe.GET("/", HeadersMiddleware(), stream.serveHTTP(), func(c *gin.Context) {
		name := c.Param("name") // Get event name from path param
		log.Printf("Subscribe request for name %s", name)
		v, ok := c.Get("clientChan")
		if !ok {
			return
		}
		clientChan, ok := v.(ClientChan)
		if !ok {
			return
		}
		c.Stream(func(w io.Writer) bool {
			// Stream message to client from message channel
			if msg, ok := <-clientChan.Chan; ok {
				c.SSEvent(name, msg)
				return true
			}
			return false
		})
	})

	publish := router.Group("/events/publish")
	publish.POST("/", func(c *gin.Context){
		// Pass body data directly as an SSE
		json_bytes, err := io.ReadAll(c.Request.Body)
		if err != nil {
			errResp := map[string]string{"error": "Could not read Body"}
			c.AbortWithStatusJSON(http.StatusBadRequest, errResp)
			return 
		}
		var json_data SSEPub
		err = json.Unmarshal(json_bytes, &json_data)
		if err != nil {
			errResp := map[string]string{"error": "Misformatted Question"}
			c.AbortWithStatusJSON(http.StatusBadRequest, errResp)
			return 
		}
		
		json_data.Data = string(json_bytes)
		stream.Message <- json_data
		c.String(http.StatusOK, json_data.Name)
	})

	// Parse Static files
	router.StaticFile("/", "./index.html")

	router.Run(":8001")
}



func HeadersMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// SSE
		c.Writer.Header().Set("Content-Type", "text/event-stream")
		c.Writer.Header().Set("Cache-Control", "no-cache")
		c.Writer.Header().Set("Connection", "keep-alive")
		c.Writer.Header().Set("Transfer-Encoding", "chunked")

		// CORS
        c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
        c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
        c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
        c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

        if c.Request.Method == "OPTIONS" {
            c.AbortWithStatus(204)
            return
        }
		c.Next()
	}
}
