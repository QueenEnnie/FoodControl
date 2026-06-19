FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /food-control ./cmd/food-control

FROM scratch

COPY --from=builder /food-control /usr/local/bin/food-control

EXPOSE 8080

ENTRYPOINT ["food-control"]
