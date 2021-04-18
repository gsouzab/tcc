FROM node:9.11.1-alpine

WORKDIR /app
COPY package.json .

# Install dependencies
RUN yarn install

ADD . /app

CMD npm start
