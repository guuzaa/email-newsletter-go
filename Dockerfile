FROM golang:1.24.2-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o target/email-newsletter ./cmd/.


FROM alpine:latest AS runtime
WORKDIR /app
COPY --from=builder /app/target/email-newsletter .
COPY configuration ./configuration
ENV APP_ENVIRONMENT=production
ENV LOG_LEVEL=warn
ENTRYPOINT ["./email-newsletter"]
