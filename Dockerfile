FROM golang:1.8

RUN mkdir -p ./src/github.com/tldr-cell/transmission-telegram

WORKDIR ./src/github.com/tldr-cell/transmission-telegram

ADD . ./

RUN go get

RUN go build -o transmission-telegram ./

#ENTRYPOINT ["./transmission-telegram"]

CMD ["sh","-c","./transmission-telegram -token=$BOT_TOKEN -masters=$MASTERS -url=$TRANSMISSION_URL -username=$TRANSMISSION_USERNAME -password=$TRANSMISSION_PASSWORD"]
