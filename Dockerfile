FROM alpine/git
EXPOSE 1737

WORKDIR /
RUN mkdir wiki
RUN mkdir config
RUN wget "https://github.com/bouncepaw/mycorrhiza/releases/download/v1.1.0/mycorrhiza-v1.1.0-linux-386.tar.gz" -O mycorrhiza.tar.gz
RUN tar -xf mycorrhiza.tar.gz -C /usr/bin
RUN rm mycorrhiza.tar.gz

VOLUME /config

WORKDIR /wiki
VOLUME /wiki
ENTRYPOINT ["mycorrhiza"]
CMD ["/wiki"]
