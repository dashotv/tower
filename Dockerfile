############################
# STEP 1 build executable binary
############################
FROM golang:1.21-alpine AS builder

WORKDIR /go/src/app
COPY . .

RUN go install

############################
# STEP 2 build a small image
############################
FROM alpine
# Copy our static executable.
WORKDIR /root/
RUN apk add --no-cache ffmpeg
COPY --from=builder /go/bin/tower .
COPY --from=builder /go/src/app/.env.vault .
CMD ["./tower", "server"]
