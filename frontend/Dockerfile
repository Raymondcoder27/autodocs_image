# Build the vuejs application
FROM node:18 AS build

# Setting working directory
WORKDIR /app

# Install dependencies
COPY package*.json ./
RUN npm install

# Copy the application code
COPY . .

# Build the application 
RUN npm run build

# Install serve
RUN npm install -g serve

# Expose port executable
EXPOSE 3550

# Start the server
CMD [ "serve", "-s", "dist", "-l", "3550" ]


