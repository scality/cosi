FROM golang:1.22.6 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

# Copy the source code into the container
COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o scality-cosi-driver ./cmd/scality-cosi-driver

FROM gcr.io/distroless/static:latest
COPY --from=builder /app/scality-cosi-driver /scality-cosi-driver
ENTRYPOINT ["/scality-cosi-driver"]
