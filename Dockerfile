# syntax=docker/dockerfile:experimental

FROM golang:1-alpine
RUN apk add --no-cache openssh-client git make
COPY . /tmp/build
RUN --mount=type=ssh,required \
    cd /tmp/build && make wxfetcher

FROM alpine
RUN apk add --no-cache ca-certificates && update-ca-certificates && \
    mkdir /etc/wxfetcher
COPY --from=0 /tmp/build/bin/ /usr/local/bin/
EXPOSE 9967 9968
ENTRYPOINT [ "wxfetcher" ]
CMD [ "-config", "/etc/wxfetcher/config.json" ]
