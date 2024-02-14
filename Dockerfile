FROM golang:1.21.3-bullseye as builder

RUN apt update \
    && apt install -y build-essential \
    && apt clean

WORKDIR /app

COPY . .
RUN go mod download

RUN CGO_ENABLED=1 CGO_LDFLAGS="-ldl" go build -ldflags '-extldflags "-static"' -o main .

FROM golang:1.21.3-bullseye as final
RUN addgroup --group appgroup && adduser appuser && usermod -aG appgroup appuser

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /app/main /apps/backend/

RUN chown -R appuser:appgroup /apps/backend
USER appuser
WORKDIR /apps/backend

CMD ["./main"]
