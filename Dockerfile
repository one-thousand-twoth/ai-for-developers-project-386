FROM golang:1.27.3-alpine AS build

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY cmd/ ./cmd/
COPY gen/ ./gen/
COPY internal/ ./internal/
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/server ./cmd/server \
    && mkdir /out/data

FROM scratch

COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /out/server /server
COPY --from=build --chown=65532:65532 /out/data /data

USER 65532:65532

ENV LISTEN_ADDR=:8080 \
    DATABASE_PATH=/data/calendar.sqlite \
    CALENDAR_TIMEZONE=Etc/UTC

EXPOSE 8080

ENTRYPOINT ["/server"]
