# --- build stage ---------------------------------------------------------
# Node 18 (LTS) — Create React App 5.0.1 is best-tested on 16/18.
FROM node:18-alpine AS build

WORKDIR /app
COPY package*.json ./
RUN npm ci --no-audit --no-fund

COPY . .

# REACT_APP_API_URL is baked into the JS bundle at build time, so it has to
# be passed as a build arg (docker-compose.yml passes it through).
ARG REACT_APP_API_URL=http://localhost:8080
ENV REACT_APP_API_URL=$REACT_APP_API_URL
# More headroom for the webpack build on small VMs.
ENV NODE_OPTIONS=--max-old-space-size=4096
RUN npm run build

# --- runtime stage -------------------------------------------------------
FROM nginx:alpine

COPY --from=build /app/build /usr/share/nginx/html
COPY nginx.conf /etc/nginx/conf.d/default.conf

EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
