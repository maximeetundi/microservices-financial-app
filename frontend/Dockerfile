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

# Production stage - use nginx for proper static file serving
FROM nginx:alpine AS runner

# Remove default nginx config and html
RUN rm -rf /etc/nginx/conf.d/default.conf /usr/share/nginx/html/*

# Copy generated static files to nginx html folder
COPY --from=builder /app/.output/public /usr/share/nginx/html

# Create custom nginx config
RUN printf 'server {\n\
    listen 3000;\n\
    server_name localhost;\n\
    root /usr/share/nginx/html;\n\
    index index.html;\n\
\n\
    location / {\n\
        try_files $uri $uri/ /index.html;\n\
    }\n\
\n\
    location /_nuxt/ {\n\
        expires 1y;\n\
        add_header Cache-Control "public, immutable";\n\
    }\n\
}\n' > /etc/nginx/conf.d/app.conf

# Expose port
EXPOSE 3000

CMD ["nginx", "-g", "daemon off;"]

