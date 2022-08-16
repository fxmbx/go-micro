## add rabbitmq to docker yml

# go get github.com/rabbitmq/amqp091-go

## RabbitMQ is a message queueing service. Service A you just send the message to a queue and send a response back to the user and Service B(notification) picks the message from the queue and does its processing. This is basically decoupling the services because without it you would have to make like http / rpm call to the notification service to send the email and that adds extra overhead to the performance of the service and it tightly coupled to two services together, since service A is a dependency of the notification service, if the notification service goes down, then service A goes down with it too. Hence the need to decouple these services, That’s why we need to use rabbitMQ. Kafka pretty much does the same, but it can do a lot more things like events streaming, logs aggregation etc… - Adetunji Samuel
