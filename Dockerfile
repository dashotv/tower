############################
# STEP 1 build executable binary
############################
FROM golang:1.14-alpine AS builder

RUN apk add --no-cache --update curl \
    bash \
    grep \
    sed \
    jq \
    ca-certificates \
    openssl \
    git \
	make \
	gcc \
	musl-dev

WORKDIR /go/src/app
COPY . .

RUN make deps
RUN make production
RUN make install

############################
# STEP 2 build a small image
############################
FROM alpine
# Copy our static executable.
WORKDIR /root/
COPY --from=builder /go/bin/tower .
COPY --from=builder /go/src/app/.tower.production.yaml ./.tower.yaml
CMD ["./tower", "server"]
