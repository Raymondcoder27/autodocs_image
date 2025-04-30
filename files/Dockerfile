FROM golang:1.22-bullseye AS builder
RUN apt-get update && apt-get install -y \
    wkhtmltopdf \
    netcat \
    curl \
    postgresql-client

# Copy your Go app and build
WORKDIR /app
COPY . .
RUN go build -o pdf-generator .

# Final lightweight image
FROM debian:bullseye-slim as runtime
RUN apt-get update && apt-get install -y wkhtmltopdf curl
COPY --from=builder /app/pdf-generator /usr/local/bin/pdf-generator
CMD ["pdf-generator"]
