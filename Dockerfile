FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 go build -o /out/task_tracker ./cmd/task_tracker

FROM alpine:3.20
WORKDIR /app
COPY --from=builder /out/task_tracker /app/task_tracker
COPY config/config.yaml /app/config/config.yaml
ENTRYPOINT ["/app/task_tracker"]


