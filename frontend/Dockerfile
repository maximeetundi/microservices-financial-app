# Build stage
FROM node:18-alpine AS builder

WORKDIR /app

# Build arguments
ARG API_URL=https://api.app.maximeetundi.store

# Set environment for build
ENV API_URL=$API_URL
ENV API_BASE_URL=$API_URL
ENV NITRO_PRESET=node-server

# Copy package files
COPY package*.json ./

# Install dependencies
RUN npm install

# Copy source code
COPY . .

# Build the application
RUN npm run build

# Debug: show output structure
RUN ls -la .output/ && ls -la .output/public/ || true

# Production stage
FROM node:18-alpine AS runner

WORKDIR /app

# Install dumb-init for proper signal handling
RUN apk add --no-cache dumb-init

# Create app user
RUN addgroup --system --gid 1001 nodejs
RUN adduser --system --uid 1001 nuxtjs

# Copy entire .output directory
COPY --from=builder --chown=nuxtjs:nodejs /app/.output ./.output

# Also copy node_modules for any runtime deps
COPY --from=builder --chown=nuxtjs:nodejs /app/node_modules ./node_modules

USER nuxtjs

# Expose port
EXPOSE 3000

# Set environment
ENV NODE_ENV=production
ENV NITRO_PORT=3000
ENV NITRO_HOST=0.0.0.0

# Start the application
ENTRYPOINT ["dumb-init", "--"]
CMD ["node", ".output/server/index.mjs"]
