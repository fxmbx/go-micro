## installing grpc here

# go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.27

# go install google.golang.org/grpc/cmd/protoc-gen-go@v1.27

# brew install protobuf

## grpc makes use of proto file that auto generate the code required

### STEPS

    - specify the message they will be sent in request
    -secify the message that will be returned as response
    -specify the service and functions you want to access

message Log{
string name = 1;
string data = 2;
}

message LogRequest{
Log LogEntry = 1;
}

message LogResponse {
string result = 1;
}

service LogService {
rpc WriteLog(LogRequest) returns LogResponse;
}
