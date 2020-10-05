FROM golang:alpine as builder
LABEL maintainer="svyatobatko199@gmail.com"
RUN apk update && apk add --no-cache git
WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download 
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o linux_web .



FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=0 /build/linux_web .

EXPOSE 8080

CMD ["./linux_web"]