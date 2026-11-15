# Build
FROM golang:1.26-alpine AS build

WORKDIR /src

# Layer dependensi dipisah agar cache build tetap efektif saat kode berubah.
COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/server ./cmd/server

# Runtime
FROM alpine:3.20

RUN apk add --no-cache ca-certificates tzdata \
    && adduser -D -u 10001 appuser

ENV TZ=Asia/Jakarta
ENV PORT=8080

WORKDIR /app
COPY --from=build /out/server /app/server

USER appuser
EXPOSE 8080

ENTRYPOINT ["/app/server"]
