### Multi-stage Dockerfile for the Vue.js frontend
##
# Stage 1: Build the application. Use an official Node.js image to install
# dependencies and compile the production build via Vite.
FROM node:24-alpine AS builder
WORKDIR /app
# Copy package definitions and install dependencies
COPY package.json package-lock.json ./
RUN npm install --legacy-peer-deps --quiet
# Copy the rest of the frontend source code into the build context
COPY . .
# Build the production files. This produces a `dist` folder with the static assets.
RUN npm run build

# Stage 2: Serve the compiled assets with Nginx. Use a minimal Nginx
# image and copy the build output to the web root. We also provide a
# custom configuration to enable SPA routing so that refreshes do not
# return 404 responses.
FROM nginx:alpine
WORKDIR /usr/share/nginx/html
COPY --from=builder /app/dist .
# Copy custom Nginx configuration
COPY nginx.conf /etc/nginx/conf.d/default.conf
# Expose the Nginx port. The docker-compose file will map this to 3000.
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]