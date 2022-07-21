FROM apline:latest

RUN mkdir /app

COPY loggerServiceApp /app

CMD ["/app/loggerServiceApp"]