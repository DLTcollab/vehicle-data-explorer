FROM golang:1.15.2
WORKDIR /mam-data-explorer
ADD . /mam-data-explorer
RUN cd /mam-data-explorer && go build
EXPOSE 8080
ENTRYPOINT ./mam-data-explorer