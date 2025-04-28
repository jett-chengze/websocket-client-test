FROM golang:latest
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o /websocket-client-test
EXPOSE 8081
CMD ["/websocket-client-test"]