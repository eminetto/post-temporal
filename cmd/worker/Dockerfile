FROM golang:1.25-alpine AS builder

WORKDIR /src

# Restore dependencies
COPY . ./

RUN go mod tidy

# Build executable
RUN go build -o /src/worker ./cmd/worker/main.go

FROM scratch
WORKDIR /src
COPY --from=builder /src/worker ./
CMD ["/src/worker"]
