FROM golang:1.18-alpine as builder
WORKDIR /app
COPY . ./
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o go-demo ./cmd/go-demo

FROM alpine:latest  
WORKDIR /root/
COPY conf ./conf
RUN sed -i 's/Host = 127.0.0.1:6379/Host = redis:6379/g' ./conf/app.ini
COPY --from=builder /app/go-demo .
EXPOSE 8000
CMD ["./go-demo"]
