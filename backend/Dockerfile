# Use a multi-stage build to keep the final image small
FROM golang:1.20-bullseye AS builder

# Set the working directory
WORKDIR /app

# Copy the go.mod and go.sum files
COPY go.mod go.sum ./

# Download the dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the Go application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Use a minimal base image for the final stage
FROM debian:bullseye-slim


RUN apt-get update && apt-get install -y \
netcat \
wkhtmltopdf \
postgresql-client \
&& apt-get clean \
&& rm -rf /var/lib/apt/lists/*
# Set the working directory

WORKDIR /app

# Copy the built Go binary from the builder stage
COPY --from=builder /app/main .

# Copy the wait-for-it.sh script to the container
COPY wait-for-it.sh /usr/local/bin/wait-for-it.sh
RUN chmod +x /usr/local/bin/wait-for-it.sh

# Copy environment variables file
COPY .env ./

# Set the entrypoint
CMD [ "./main" ]
# ENTRYPOINT ["/usr/local/bin/wait-for-it.sh", "db:3306", "--", "/app/main"]
# Set the entrypoint to run migrations and then start the main application
# ENTRYPOINT ["/usr/local/bin/wait-for-it.sh", "db:3306", "go run /app/migrate.go && ./main"]

