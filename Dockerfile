# builder image
FROM golang:1.14-alpine AS builder

RUN set -ex && apk --update --no-cache add \
    bash \
    make \
    cmake  \
    git \
    gcc \
    musl-dev
    
WORKDIR /app
COPY .netrc /root/.netrc
COPY . .
ENV GO111MODULE=on
RUN make dep
# RUN make test
RUN make clean
RUN make bin
# final image
FROM alpine:3.10.1

COPY --from=builder  /app/bin/linux/${APP} /bin/${APP}

RUN set -ex && apk --update --no-cache add \
    curl \
    openssl \ 
    ca-certificates \
    && update-ca-certificates

ENTRYPOINT ["/bin/simulator"]
