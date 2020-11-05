FROM golang:1.15.2
WORKDIR /vehicle-data-explorer
ADD . /vehicle-data-explorer
RUN cd /vehicle-data-explorer && go build -o app
EXPOSE 8080
ENTRYPOINT ./app
