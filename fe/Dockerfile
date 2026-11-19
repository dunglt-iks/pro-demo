FROM node:20-bookworm-slim
WORKDIR /app
ENV NEXT_TELEMETRY_DISABLED=1
COPY package.json package-lock.json ./
RUN npm ci --no-audit --no-fund
COPY . .
ARG BE_URL=http://be:8180
ENV BE_URL=$BE_URL
RUN npm run build
ENV NODE_ENV=production
EXPOSE 3100
CMD ["node_modules/.bin/next", "start", "-p", "3100"]
