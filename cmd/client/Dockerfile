FROM golang:1.25-alpine AS builder

WORKDIR /src

# Restore dependencies
COPY . ./

#RUN go mod tidy

# Build executable
RUN go build -o /src/client ./cmd/client/main.go

FROM scratch
WORKDIR /src
COPY --from=builder /src/client ./
EXPOSE 8083
CMD ["/src/client"]
