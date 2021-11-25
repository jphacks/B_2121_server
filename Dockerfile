FROM golang:1.17.3 AS build-env

ENV GOOS=linux
ENV GOARCH=amd64
ENV CGO_ENABLED=0
ENV GO111MODULE=on

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download -x

COPY . .
RUN go build .

FROM alpine:3.15.0

RUN apk add --no-cache bash ca-certificates curl

COPY --from=build-env /app/B_2121_server /B_2121_server
RUN chmod a+x /B_2121_server
RUN mkdir profileImages
COPY ./migrations /migrations

EXPOSE 8080
CMD ["/B_2121_server"]
