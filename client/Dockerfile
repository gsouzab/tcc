FROM node:16.10.0-alpine3.14

WORKDIR /app
COPY package.json .

# Install dependencies
RUN yarn install

ADD . /app

CMD npm start
