# Build stage (SvelteKit 3 requires Node 22.17+)
FROM node:22-alpine AS builder

WORKDIR /app

# Install dependencies first (better caching)
COPY package*.json ./
RUN npm ci

# Copy source and build
COPY . .
RUN npm run build

# Runtime stage: adapter-node output, started with `node build`
FROM node:22-alpine

WORKDIR /app

# Copy built application
COPY --from=builder /app/build ./build
COPY --from=builder /app/package*.json ./

# Install production dependencies only (adapter-node keeps `dependencies` external)
RUN npm ci --omit=dev

EXPOSE 3000

ENV PORT=3000

CMD ["node", "build"]
