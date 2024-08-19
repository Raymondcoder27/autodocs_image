#!/bin/sh

# Export environment variables from .env file
if [ -f /usr/local/bin/.env ]; then
  export $(grep -v '^#' /usr/local/bin/.env | xargs)
fi

# Start the backend application
/usr/local/bin/main &
# Start nginx
# nginx -g 'daemon off;'
# Serve frontend using serve
serve -s /usr/share/nginx/html -l 80



#!/bin/sh

# Export environment variables from .env file
# if [ -f /usr/local/bin/.env ]; then
  # export $(grep -v '^#' /usr/local/bin/.env | xargs)
# fi

# Start the backend application
# /usr/local/bin/main &

# Start serving the frontend on port 4000
# serve -s /usr/share/nginx/html -l 4000 &

# Wait indefinitely to keep the container running
# wait
