# Build stage
FROM node:18-alpine AS builder

# Set working directory
WORKDIR /app

# Configure Chinese npm registry for faster downloads
RUN npm config set registry https://registry.npmmirror.com

# Copy package files
COPY web/package.json web/package-lock.json web/pnpm-lock.yaml web/yarn.lock ./web/

# Install dependencies
WORKDIR /app/web
RUN npm install

# Copy source code
WORKDIR /app
COPY web ./web

# Build the application
WORKDIR /app/web
RUN npm run build

# Production stage
FROM nginx:alpine

# Remove default nginx website
RUN rm -rf /usr/share/nginx/html/*

# Copy built files from builder stage
COPY --from=builder /app/web/dist /usr/share/nginx/html

# Copy nginx configuration
COPY deploy/frontend/nginx.conf /etc/nginx/nginx.conf

# Expose port
EXPOSE 80

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --quiet --tries=1 --spider http://localhost/ || exit 1

# Start nginx
CMD ["nginx", "-g", "daemon off;"]