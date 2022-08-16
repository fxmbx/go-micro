 
FROM alpine:latest

RUN mkdir /app

# COMMENTING THIS OUT BECAUSE OF THE MAKEFILE   
# COPY --from=build /app/brokerApp /app
COPY listenerApp /app

CMD [ "/app/listenerApp" ]