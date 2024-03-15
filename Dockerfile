FROM golang:1.22.1-bookworm

WORKDIR /smoeji

COPY . .

RUN go get

RUN go build -o smoeji

EXPOSE 3000

HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 CMD [ "curl -f http://localhost:3000/_healthy || exit 1" ]

CMD [ "./smoeji" ]