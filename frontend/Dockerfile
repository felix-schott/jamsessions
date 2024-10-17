# PROD

# build stage
FROM node:20.17.0-alpine3.20 AS builder

ARG API_ADDRESS
COPY . /workspace
WORKDIR /workspace

RUN API_ADDRESS=${API_ADDRESS} npm ci --omit-dev && npm run build

# final image
FROM node:20.17.0-alpine3.20 AS production
COPY --from=builder /workspace/build/ /workspace/build
COPY --from=builder /workspace/package.json /workspace/package.json

EXPOSE 3000

CMD ["node", "/workspace/build"]