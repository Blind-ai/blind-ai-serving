# build stage
FROM golang:alpine AS build-env
RUN apk --no-cache add build-base git bzr mercurial gcc
ADD . /src
RUN cd /src && go build cmd/AI/lungh/main.go

# final stage
FROM alpine
RUN mkdir /app
COPY --from=build-env /src/main /app/
WORKDIR /app
ENTRYPOINT ./main
