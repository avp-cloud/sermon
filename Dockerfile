# Stage-1
FROM golang:1.13-alpine AS builder

RUN apk update && apk add alpine-sdk git && rm -rf /var/cache/apk/*

RUN mkdir -p /app
WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
RUN go build ./

# Stage-2
FROM alpine:latest

RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*

COPY --from=builder /app/sermon .

ENV DB_PATH "sermon.db"
ENV POLL_INTERVAL "15"
ENV PORT "80"
EXPOSE 80

ENTRYPOINT ["./sermon"]