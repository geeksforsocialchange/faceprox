# syntax=docker/dockerfile:1
FROM golang:1.18-alpine
WORKDIR /go/src/github.com/wheresalice/faceprox/
COPY . .
#RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o faceprox .

FROM scratch
COPY --from=0 /go/src/github.com/wheresalice/faceprox/faceprox /
COPY Procfile /
CMD ["/faceprox"]
