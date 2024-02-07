# Use the Nginx image from Docker Hub
FROM nginx:alpine

# Copy the HTML file to the Nginx html directory
COPY index.html /usr/share/nginx/html

# Copy the Nginx configuration file
COPY nginx.conf /etc/nginx/conf.d/default.conf
