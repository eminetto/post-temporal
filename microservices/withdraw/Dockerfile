FROM golang:1.25-alpine AS builder

WORKDIR /src

# Restore dependencies
COPY . .

RUN go mod tidy

# Build executable
RUN go build -o /src/withdraw ./

FROM scratch
WORKDIR /src
COPY --from=builder /src/withdraw ./
EXPOSE 8081
CMD ["/src/withdraw"]
