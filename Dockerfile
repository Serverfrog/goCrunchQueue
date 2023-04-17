FROM rust:bookworm as buildercrunchy
RUN git clone https://github.com/crunchy-labs/crunchy-cli.git
WORKDIR crunchy-cli
RUN cargo build --release

FROM golang:bullseye as buildergo
WORKDIR /usr/src/app
COPY ./ ./
RUN go mod download && go mod verify

RUN go build -race -ldflags "-extldflags '-static'" -o goCrunchQueue cmd/Lets.go

FROM gcr.io/distroless/base-debian11
LABEL authors="serverfrog"

COPY --from=buildercrunchy /crunchy-cli/target/release/crunchy-cli /usr/bin/crunchy-cli
COPY --from=buildergo /usr/src/app/goCrunchQueue /goCrunchQueue

VOLUME /queue
VOLUME /media-destination

CMD ["/goCrunchQueue", "-config=/config/config.yml"]