# Simple Golang SSE Pub/Sub

Runs a gin server at http://localhost:8080 which runs minimal event bus middleware that uses Server-Sent Events to demonstrate:
  - Accepting subscriptions to named event topics
  - Publishing events to named event topics
  - A demo frontend which displays events published to the "stream" topic
  - Simple all-origin CORS configuration from scratch

## Endpoints
### POST /events/publish

Publishes the event data to the topic specified by X-Event-Name.

    Request Headers:
        X-Event-Name: Name of topic to publish event to
    Request Body:
        json-encoded event data
    Response Headers:
        See HeadersMiddleware(), as well as Content-Type: text/event-stream
    Response Body:
        "topic name"

### GET /events/subscribe/{name}
Subscribes to the topic specified by name. Opens long-lived HTTP connection which recieves published SSE e.g. by EventSource in browser.

    Response Headers:
        see HeadersMiddleware()

### GET /docs
Returns this markdown document.
