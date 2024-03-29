FROM ghcr.io/hazmi35/node:18-dev-alpine as build-stage

LABEL name "Nezu Fuzzier (Docker Build)"
LABEL maintainer "KagChi"

COPY . .

RUN npm ci

RUN npm run build

RUN npm prune --production

FROM oven/bun

LABEL name "Nezu Fuzzier Production"
LABEL maintainer "KagChi"

COPY --from=build-stage /tmp/build/package.json .
COPY --from=build-stage /tmp/build/package-lock.json .
COPY --from=build-stage /tmp/build/bun.lockb .
COPY --from=build-stage /tmp/build/node_modules ./node_modules
COPY --from=build-stage /tmp/build/src ./src

CMD ["bun", "start"]