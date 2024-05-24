FROM golang:alpine

WORKDIR /usr/src/app
COPY . .

RUN go mod download
RUN go mod tidy

RUN go build -o /bin/app

EXPOSE 9019
ENV EVENTSOURCEHOST 0.0.0.0
ENV EVENTSOURCEPORT 09019

ENV GIN_MODE release

CMD ["/bin/app"]