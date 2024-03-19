FROM golang:1.22-bookworm as builder

WORKDIR /usr/local/src

COPY . .

RUN go build -o token-service .

RUN chmod +x token-service

FROM gcr.io/distroless/base-nossl-debian12

WORKDIR /usr/local/src

COPY --from=builder /usr/local/src/configuration ./configuration
COPY --from=builder /usr/local/src/token-service .

ENV ENVIRONMENT prod
ENTRYPOINT [ "/usr/local/src/token-service" ]