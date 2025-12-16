# Build stage
FROM node:18-alpine AS builder

WORKDIR /app

# Copy package files (only package.json, not lock file)
COPY package.json ./

# Clean and install dependencies
RUN rm -f package-lock.json && npm cache clean --force && npm install

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

# Fix Nuxt bug: server looks in server/chunks/public but files are in public
# Create symlink to fix the path resolution
RUN mkdir -p .output/server/chunks && \
    ln -s /app/.output/public .output/server/chunks/public

USER nuxtjs

# Expose port
EXPOSE 3000

# Set environment
ENV NODE_ENV=production

# Start the application
ENTRYPOINT ["dumb-init", "--"]
CMD ["node", ".output/server/index.mjs"]

