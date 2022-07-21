#build a tiny docker image 
FROM alpine:latest

RUN mkdir /app

# COMMENTING THIS OUT BECAUSE OF THE MAKEFILE   
# COPY --from=build /app/authApp /app
COPY mailerApp /app
COPY templates /templates

CMD [ "/app/mailerApp" ]