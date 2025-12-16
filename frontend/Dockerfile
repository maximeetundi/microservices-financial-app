# Build stage
FROM node:18-alpine AS builder

WORKDIR /app

# Build arguments
ARG API_URL=https://api.app.maximeetundi.store

# Set environment for build
ENV API_URL=$API_URL
ENV API_BASE_URL=$API_URL

# Copy package files
COPY package*.json ./

# Install dependencies
RUN npm install

# Copy source code
COPY . .

# Build the application
RUN npm run build

# Production stage
FROM node:18-alpine AS runner

WORKDIR /app

# Install dumb-init for proper signal handling
RUN apk add --no-cache dumb-init

# Create app user
RUN addgroup --system --gid 1001 nodejs
RUN adduser --system --uid 1001 nuxtjs

# Copy built application
COPY --from=builder --chown=nuxtjs:nodejs /app/.output ./.output

USER nuxtjs

# Expose port
EXPOSE 3000

# Set environment
ENV NODE_ENV=production

# Start the application
ENTRYPOINT ["dumb-init", "--"]
CMD ["node", ".output/server/index.mjs"]
