FROM golang:alpine

WORKDIR /usr/src/app
COPY . .

RUN go mod download
RUN go mod tidy

RUN go build -o /bin/app

CMD ["/bin/event-source"]