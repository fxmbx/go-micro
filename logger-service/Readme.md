# go get go.mongodb/org/mongo-driver/mongo

# go get go.mongodb/org/mongo-driver/mongo/options

## rpc remote procedure call is basically allowing service A access finctionalities of service B like they are codded in service B

# - on the server side you lregister a new rpcServer then listen

# - on the server side you listen over tcp via a specified address

# - then accept all incoming dials and serve each request. on the server

# - on the client side you dial via tcp and server port

# - once you get a client you then call the rpcserver and rpcfunction you need then pass in the rpc requeired payload
