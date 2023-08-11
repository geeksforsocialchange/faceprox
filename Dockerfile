# syntax=docker/dockerfile:1
FROM golang:1.21-alpine as build
WORKDIR /go/src/github.com/geeksforsocialchange/faceprox/
COPY . .
RUN apk --no-cache add ca-certificates
RUN CGO_ENABLED=0 GOOS=linux go build -o faceprox .

FROM scratch
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /go/src/github.com/geeksforsocialchange/faceprox/faceprox /
COPY Procfile /
CMD ["/faceprox"]
