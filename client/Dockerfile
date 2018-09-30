FROM node:alpine

RUN npm install -g yarn

ADD . /usr/src/app

# Install dependencies
WORKDIR /usr/src/app

CMD npm start
