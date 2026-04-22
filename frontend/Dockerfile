FROM node:20-alpine AS builder
ARG VITE_API_BASE
ENV VITE_API_BASE=${VITE_API_BASE}

WORKDIR /app
COPY package.json package-lock.json* ./
RUN npm install

COPY . .
RUN npm run build

FROM nginx:1.25-alpine
RUN rm /etc/nginx/conf.d/default.conf
COPY nginx.conf /etc/nginx/conf.d/default.conf
COPY --from=builder /app/dist /usr/share/nginx/html

EXPOSE 3000

ARG COMMIT_SHA=unknown
ARG BUILD_DATE=unknown
LABEL org.opencontainers.image.source="https://github.com/Qulive-edu/ArchitectureCoursework" \
      org.opencontainers.image.revision=${COMMIT_SHA} \
      org.opencontainers.image.created=${BUILD_DATE}

EXPOSE 3000
CMD ["nginx", "-g", "daemon off;"]
