# COMMENTING THIS OUT BECAUSE OF THE MAKEFILE   
# FROM golang:1.18-alpine as build

# RUN mkdir /app

# COPY . /app

# WORKDIR /app

# RUN CGO_ENABLED=0 go build -o brokerApp ./cmd/api

# RUN chmod +x /app/brokerApp

##build a tiny docker image 
FROM alpine:latest

RUN mkdir /app

# COMMENTING THIS OUT BECAUSE OF THE MAKEFILE   
# COPY --from=build /app/brokerApp /app
COPY brokerApp /app

CMD [ "/app/brokerApp" ]