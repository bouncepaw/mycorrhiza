# Build the binary
FROM golang:1.14-alpine AS build

WORKDIR /src
COPY . .
ENV CGO_ENABLED=0

RUN go generate && go build -o /out/mycorrhiza .

# Create image for running Mycorrhiza
FROM alpine:3 AS app

EXPOSE 1737
VOLUME ["/wiki-data"]
ENTRYPOINT ["/mycorrhiza", "/wiki-data"]

RUN apk add --no-cache git
COPY --from=build /out/mycorrhiza /

