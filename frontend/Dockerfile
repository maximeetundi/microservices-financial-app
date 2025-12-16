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

# Production stage - use simple static file server
FROM node:18-alpine AS runner

WORKDIR /app

# Install serve globally
RUN npm install -g serve

# Create app user
RUN addgroup --system --gid 1001 nodejs
RUN adduser --system --uid 1001 nuxtjs

# Copy generated static files
COPY --from=builder --chown=nuxtjs:nodejs /app/.output/public ./public

USER nuxtjs

# Expose port
EXPOSE 3000

# Start with serve
CMD ["serve", "-s", "public", "-l", "3000"]

