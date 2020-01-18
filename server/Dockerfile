FROM golang:alpine

WORKDIR /usr/app

COPY . .

#RUN make
RUN go build -v ./cmd/apiserver

EXPOSE 5000

CMD ["./apiserver"]
