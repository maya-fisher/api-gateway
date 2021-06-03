#build stage
FROM golang:alpine AS builder
RUN apk add --no-cache git
ENV GO111MODULE=on
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o api-gateway -v

#final stage
FROM golang:alpine
# RUN apk --no-cache add curl
COPY app.env .
COPY --from=builder /app/api-gateway /api-gateway
EXPOSE 8081
ENTRYPOINT ["/api-gateway"]