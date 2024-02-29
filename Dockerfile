FROM golang:1.22 as builder

WORKDIR /app

COPY . .

RUN go mod download

RUN go mod tidy

RUN cd /app/web/cmd && GOOS=linux go build -o log-snare

EXPOSE 8080

WORKDIR /app/web/cmd

ENTRYPOINT ["/app/web/cmd/log-snare"]
CMD ["-r"]
