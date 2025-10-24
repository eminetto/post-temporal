FROM golang:1.25-alpine AS builder

WORKDIR /src

# Restore dependencies
COPY . .

RUN go mod tidy

# Build executable
RUN go build -o /src/refund ./

FROM scratch
WORKDIR /src
COPY --from=builder /src/refund ./
EXPOSE 8082
CMD ["/src/refund"]
