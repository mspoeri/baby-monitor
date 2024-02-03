# Use an official Node.js runtime as the base image
FROM node:20-alpine

RUN apk update && apk add ffmpeg

# Set the working directory in the container to /app
WORKDIR /app

# Copy the rest of the application's code to the working directory
COPY server.js server.js

# Make port 8081 available outside the container
EXPOSE 8081

# Start the application when the container is run
CMD [ "node", "server.js" ]
