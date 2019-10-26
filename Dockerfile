# syntax=docker/dockerfile:experimental

FROM golang:1-alpine
RUN apk add --no-cache openssh-client git make
COPY . /tmp/build
RUN --mount=type=ssh,required \
    mkdir -p -m 0600 ~/.ssh && ssh-keyscan bitbucket.org >> ~/.ssh/known_hosts && \
    git config --global url."git@bitbucket.org:".insteadOf "https://bitbucket.org/" && \
    cd /tmp/build && make wxfetcher

FROM alpine
RUN apk add --no-cache ca-certificates && update-ca-certificates && \
    mkdir /etc/wxfetcher
COPY --from=0 /tmp/build/bin/ /usr/local/bin/
EXPOSE 9967 9968
ENTRYPOINT [ "wxfetcher" ]
CMD [ "-config", "/etc/wxfetcher/config.json" ]
