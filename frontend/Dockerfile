FROM node:20.3-alpine3.17 as builder

WORKDIR /swimlogs

COPY package.json package-lock.json ./

RUN npm install

COPY . .

RUN npm run build

FROM nginx:1.25.1

WORKDIR /

RUN rm -rf /usr/share/nginx/html/*

COPY --from=builder /swimlogs/dist /usr/share/nginx/html
