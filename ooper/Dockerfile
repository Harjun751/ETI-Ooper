FROM node:12.18.1
RUN npm install -g http-server
WORKDIR /app

# Add package json to  workdir
COPY ["package.json", "package-lock.json*", "./"]
# Install all required  packages
RUN npm install

# Add source code to workdir
COPY . .
RUN npm run build
EXPOSE 8080
CMD [ "http-server", "dist" ]