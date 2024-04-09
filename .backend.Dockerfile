FROM golang:latest
WORKDIR /backend/
COPY ./backend/go.mod /backend/
COPY ./backend/go.sum /backend/
RUN go install github.com/cosmtrek/air@latest
COPY ./backend/ /backend/
COPY .env /backend/
RUN go mod tidy
EXPOSE 1323
ENTRYPOINT ["air"]