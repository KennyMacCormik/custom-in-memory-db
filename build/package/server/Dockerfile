FROM golang:1.22.5 as test
RUN git clone https://github.com/KennyMacCormik/custom-in-memory-db.git -b task-2
WORKDIR /go/custom-in-memory-db/
RUN go test ./...

FROM golang:1.22.5 as build
RUN git clone https://github.com/KennyMacCormik/custom-in-memory-db.git -b task-2
WORKDIR /go/custom-in-memory-db/
RUN ls -la
RUN go build ./main.go
RUN mv main /bin/server
RUN chmod 550 /bin/server

FROM golang:1.22.5 as run
COPY --from=build /bin/server /bin/server
ENTRYPOINT ["/bin/server"]