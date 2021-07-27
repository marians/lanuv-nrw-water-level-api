FROM golang:1.16-alpine3.14 AS build
WORKDIR /project
COPY . /project/
RUN go build .

FROM alpine:3.14
RUN apk add --no-cache tzdata
COPY --from=build /project/lanuv-nrw-water-level-api /
ENV GIN_MODE release
EXPOSE 8080
ENTRYPOINT ["/lanuv-nrw-water-level-api"]
