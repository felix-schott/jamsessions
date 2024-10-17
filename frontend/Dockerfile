# DEV

FROM node:20.17.0-alpine3.20 AS development

COPY . /workspace
# RUN cd workspace && npm ci && npm run build
WORKDIR /workspace
RUN npm i && API_ADDRESS=${API_ADDRESS} npm run build

ENTRYPOINT ["npm", "run", "dev"]

# PROD

# build stage
FROM node:20.17.0-alpine3.20 AS builder

ARG API_ADDRESS
COPY . /workspace
WORKDIR /workspace

RUN API_ADDRESS=${API_ADDRESS} npm ci && npm run build

EXPOSE 80

# final image
FROM nginx:1.26-alpine3.20 AS production
COPY --from=builder /workspace/public/ /usr/share/nginx/html