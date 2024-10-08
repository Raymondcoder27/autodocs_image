# Stage 1 - Build the frontend
FROM node:18 AS frontend_build

WORKDIR /app/frontend

COPY frontend/package*.json ./
RUN npm install
COPY frontend/ ./
RUN npm run build

# Install serve globally
RUN npm install -g serve

# Stage 2 - Build the backend
FROM golang:1.20 AS backend_build

WORKDIR /app/backend

COPY backend/go.mod backend/go.sum ./
RUN go mod download
COPY backend/ ./
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# Stage 3 - Combine frontend and backend into one image
# Use a minimal base image for the final stage
FROM debian:bullseye-slim

# Install required packages
RUN apt-get update && apt-get install -y \
   netcat \
   wkhtmltopdf \
   postgresql-client \
   curl \
   && apt-get clean \
   && rm -rf /var/lib/apt/lists/*

# Install Node.js and npm to use `serve`
RUN curl -fsSL https://deb.nodesource.com/setup_14.x | bash -
RUN apt-get install -y nodejs
RUN npm install -g serve

# Copy frontend build to the image
COPY --from=frontend_build /app/frontend/dist /usr/share/nginx/html

# Copy backend binary and .env file to the image
COPY --from=backend_build /app/backend/main /usr/local/bin/main
# COPY --from=backend_build /app/backend/.env /usr/local/bin/.env

# Copy start.sh script to the image and make it executable
COPY start.sh /usr/local/bin/start.sh
RUN chmod +x /usr/local/bin/start.sh

# Set environment variables
# ENV PORT=8080
# ENV MINIO_ROOT_USER=minioadmin
# ENV MINIO_ROOT_PASSWORD=minioadmin

# Expose ports
EXPOSE 3550
EXPOSE 4100

# Start backend and frontend
ENTRYPOINT ["/usr/local/bin/start.sh"]

