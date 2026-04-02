FROM golang:1.26-alpine AS base
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

FROM base AS builder
# strip to basics for scratch build
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /app/main .

FROM scratch AS load-test
COPY --from=builder /app/main /main
ENTRYPOINT ["/main"]