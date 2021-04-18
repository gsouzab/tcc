FROM golang:alpine

RUN apk add --no-cache git
ARG app_env=development
ENV APP_ENV $app_env

WORKDIR /go/src/app
COPY . .

RUN go get -d -v ./...
RUN go install -v ./...

CMD if [ ${APP_ENV} = production ]; \
	then \
	app; \
	else \
	go get github.com/pilu/fresh && \
	fresh; \
	fi
	
EXPOSE 8000