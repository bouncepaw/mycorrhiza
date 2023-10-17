FROM golang:alpine as build
WORKDIR src
COPY . .
ENV CGO_ENABLED=0
RUN go build -o /out/mycorrhiza .

FROM alpine/git as app
EXPOSE 1737

RUN apk add --no-cache curl
HEALTHCHECK CMD curl -Ns localhost:1737 || exit 1

WORKDIR /
RUN mkdir wiki
COPY --from=build /out/mycorrhiza /usr/bin

WORKDIR /wiki
VOLUME /wiki
ENTRYPOINT ["mycorrhiza"]
CMD ["/wiki"]
