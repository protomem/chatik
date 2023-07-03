FROM node:alpine AS builder
WORKDIR /app

COPY package.json package-lock.json
COPY yarn.lock yarn.lock
RUN yarn 

COPY ./ ./
RUN yarn build


FROM nginx:alpine

COPY nginx.conf /etc/nginx/nginx.conf
COPY --from=builder /app/build /usr/share/nginx/html
