FROM golang:1.18-alpine as builder
WORKDIR /app
COPY . ./
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o gogolook2024 ./cmd/gogolook2024

FROM alpine:latest  
WORKDIR /root/
COPY conf ./conf
COPY --from=builder /app/gogolook2024 .
EXPOSE 8000
CMD ["./gogolook2024"]
