FROM golang:alpine

WORKDIR /usr/app

ENV CGO_ENABLED=0

COPY . .

#RUN make
RUN go build -v ./cmd/apiserver


EXPOSE 5000

CMD ["./apiserver"]
