# DEV

FROM node:20.17.0-alpine3.20 AS development

COPY . /workspace
# RUN cd workspace && npm ci && npm run build
WORKDIR /workspace
RUN echo ${API_ADDRESS}
RUN npm i && API_ADDRESS=${API_ADDRESS} npm run build

ENTRYPOINT ["npm", "run", "dev"]

# PROD

# build stage
FROM node:20.17.0-alpine3.20 AS builder

ARG API_ADDRESS
COPY . /workspace
# RUN cd workspace && npm ci && npm run build
RUN cd workspace && npm ci && npm run build

# final image
FROM nginx:1.26-alpine3.20 AS production
COPY --from=builder /workspace/public/ /usr/share/nginx/html