FROM golang:1.17-alpine3.15 AS build
WORKDIR /project
COPY . /project/
RUN go build .

FROM alpine:3.15
RUN apk add --no-cache tzdata
COPY --from=build /project/lanuv-nrw-water-level-api /

ENV GIN_MODE release
ENV PORT 8080
EXPOSE 8080

ENTRYPOINT ["/lanuv-nrw-water-level-api"]
