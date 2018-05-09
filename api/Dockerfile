FROM golang:alpine

RUN apk add --no-cache git

ARG app_env
ENV APP_ENV $app_env

COPY ./ /go/src/github.com/gsouzab/tcc/api
WORKDIR /go/src/github.com/gsouzab/tcc/api

RUN go get ./
RUN go build

CMD if [ ${APP_ENV} = production ]; \
	then \
	app; \
	else \
	go get github.com/pilu/fresh && \
	fresh; \
	fi
	
EXPOSE 8000