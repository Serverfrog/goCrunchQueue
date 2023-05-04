
FROM debian:bookworm-slim as baseimage
LABEL authors="serverfrog"

RUN apt-get update
RUN apt-get full-upgrade -y
RUN apt-get install -y ffmpeg
RUN rm -rf /var/lib/apt/lists/*

FROM node:lts as builderfrontend
WORKDIR /app

COPY gocrunchqueue-ui .

RUN yarn install \
  --prefer-offline \
  --frozen-lockfile \
  --non-interactive \
  --production=false

RUN yarn build

RUN rm -rf node_modules && \
  NODE_ENV=production yarn install \
  --prefer-offline \
  --pure-lockfile \
  --non-interactive \
  --production=true

FROM rust:bookworm as buildercrunchy
RUN git clone https://github.com/crunchy-labs/crunchy-cli.git
WORKDIR crunchy-cli
RUN cargo build --release

FROM golang:bullseye as buildergo
WORKDIR /usr/src/app
COPY ./cmd ./cmd
COPY ./internal ./internal
COPY ./go.mod ./
COPY ./go.sum ./


RUN go mod download && go mod verify

RUN go build -race -ldflags "-extldflags '-static'" -o goCrunchQueue cmd/Lets.go

FROM baseimage

RUN mkdir /goCrunchQueue

WORKDIR /goCrunchQueue

COPY --from=buildercrunchy /crunchy-cli/target/release/crunchy-cli /usr/bin/crunchy-cli
COPY --from=buildergo /usr/src/app/goCrunchQueue /goCrunchQueue/goCrunchQueue
COPY --from=builderfrontend /app/dist/ /goCrunchQueue/ui/

COPY config/ /goCrunchQueue/config
COPY docker-entrypoint.sh /docker-entrypoint.sh

ENV ETP_RT=""
ENV CREDENTIALS=""

RUN chmod +x /goCrunchQueue/goCrunchQueue
RUN chmod +x /usr/bin/crunchy-cli

VOLUME /goCrunchQueue/queue
VOLUME /goCrunchQueue/media-destination
VOLUME /root/.config/crunchy-cli
VOLUME /goCrunchQueue/config

ENTRYPOINT ["/docker-entrypoint.sh"]
CMD ["/goCrunchQueue/goCrunchQueue", "-config=/goCrunchQueue/config/config.yaml"]