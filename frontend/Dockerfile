# syntax=docker/dockerfile:1.6

FROM oven/bun:1.2 AS deps
WORKDIR /app
COPY package.json bun.lock ./
RUN bun install --frozen-lockfile

FROM deps AS build
ARG FRONTEND_API_URL=http://api:8080
ENV FRONTEND_API_URL=$FRONTEND_API_URL
COPY . .
RUN bun run build

FROM oven/bun:1.2 AS runner
WORKDIR /app
ENV NODE_ENV=production \
    FRONTEND_PORT=4173 \
    FRONTEND_API_URL=http://api:8080
COPY --from=build /app/dist ./dist
COPY package.json bun.lock ./
COPY scripts ./scripts
RUN bun install --frozen-lockfile --production
EXPOSE 4173
CMD ["bun", "run", "serve"]
