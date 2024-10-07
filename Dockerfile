FROM golang:1.23

LABEL contributers="Zahra, Dana, Ghadeer & Salman"
LABEL project="Forum"
LABEL version ="1.23"

WORKDIR /app

COPY . .

RUN go build -o forum .

CMD ["./forum"]
