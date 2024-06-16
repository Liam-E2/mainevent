package eventsource

import (
	"io"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func RunServer(addr string, middlewares ...gin.HandlerFunc) error {
	router := gin.Default()

	// Middleware
	router.Use(middlewares...)

	// Initialize new Event Server
	stream := NewServer()

	// Handle Event Subscriptions
	// /events/subscribe/name; name of SSE, created in events.Server
	subscribe := router.Group("/events/subscribe")
	subscribe.GET("/:name", stream.serveHTTP(), func(c *gin.Context) {
		c.Writer.Header().Set("Content-Type", "text/event-stream")
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
				log.Printf("%v", msg)
				c.SSEvent(name, msg)
				return true
			}
			return false
		})
	})

	publish := router.Group("/events/publish")

	publish.POST("/", func(c *gin.Context) {
		c.Writer.Header().Set("Content-Type", "application/json")

		// Pass body json directly as SSE
		json_bytes, err := io.ReadAll(c.Request.Body)
		if err != nil {
			errResp := map[string]string{"error": "Could not read Body"}
			c.AbortWithStatusJSON(http.StatusBadRequest, errResp)
			return
		}
		defer c.Request.Body.Close()

		json_data := SSEPub{c.Request.Header.Get("X-Event-Name"), string(json_bytes)}
		stream.Message <- json_data
		c.String(http.StatusOK, json_data.Name)
	})

	// Serve Docs as Markdown
	docs := router.Group("/docs")
	docs.GET("/", func(c *gin.Context) {
		c.File("./Readme.md")
		c.AbortWithStatus(http.StatusOK)
	})

	// Serve simple demo frontend
	router.StaticFile("/", "./index.html")

	return router.Run(addr)
}
