FROM golang as build
WORKDIR src
COPY . .
ENV CGO_ENABLED=0
RUN make build

FROM alpine/git as app
EXPOSE 1737

WORKDIR /
RUN mkdir wiki
COPY --from=build /go/src/mycorrhiza /usr/bin
RUN mkdir config

VOLUME /config

WORKDIR /wiki
VOLUME /wiki
ENTRYPOINT ["mycorrhiza"]
CMD ["/wiki"]
