# Build stage
FROM node:18-alpine AS builder

WORKDIR /app

# Copy package files
COPY package*.json ./

# Install dependencies
RUN npm install

# Copy source code
COPY . .

# Generate static site
RUN npm run generate

# Debug: show what was generated
RUN echo "=== .output structure ===" && ls -la .output/ || true
RUN echo "=== .output/public structure ===" && ls -la .output/public/ || true
RUN echo "=== Looking for index.html ===" && find .output -name "index.html" 2>/dev/null || true

# Production stage - use nginx for proper static file serving
FROM nginx:alpine AS runner

# Copy generated static files to nginx html folder
COPY --from=builder /app/.output/public /usr/share/nginx/html

# Configure nginx for SPA routing
RUN echo 'server { \
    listen 3000; \
    server_name localhost; \
    root /usr/share/nginx/html; \
    index index.html; \
    location / { \
        try_files $uri $uri/ /index.html; \
    } \
    location /_nuxt/ { \
        expires 1y; \
        add_header Cache-Control "public, immutable"; \
    } \
}' > /etc/nginx/conf.d/default.conf

# Expose port
EXPOSE 3000

CMD ["nginx", "-g", "daemon off;"]
