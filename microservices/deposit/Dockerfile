FROM golang:1.25-alpine AS builder

WORKDIR /src

# Restore dependencies
COPY . .

RUN go mod tidy

# Build executable
RUN go build -o /src/deposit ./

FROM scratch
WORKDIR /src
COPY --from=builder /src/deposit ./
EXPOSE 8080
CMD ["/src/deposit"]
