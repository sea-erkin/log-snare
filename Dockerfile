FROM golang:alpine AS builder

ENV USER=logsnare
ENV UID=10001

RUN adduser \
    --disabled-password \
    --gecos "" \
    --home "/nonexistent" \
    --shell "/sbin/nologin" \
    --no-create-home \
    --uid "${UID}" \
    "${USER}"

RUN apk update && apk add --no-cache git gcc musl-dev

COPY . /app

WORKDIR /app

RUN go mod tidy

RUN cd /app/web/cmd && GOOS=linux go build -o log-snare

FROM alpine:latest

COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group
COPY --from=builder /app /app
RUN chown -R logsnare:logsnare /app
USER logsnare:logsnare

EXPOSE 8080

WORKDIR /app/web/cmd

ENTRYPOINT ["/app/web/cmd/log-snare"]
CMD ["-r"]
