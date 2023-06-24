FROM golang:1.17

WORKDIR /Backend
COPY . .

RUN go get -u github.com/gorilla/mux
RUN go get github.com/rs/cors
RUN go get -d

EXPOSE 4000 

CMD ["go", "run", "main.go"]
