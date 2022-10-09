install protoc

https://grpc.io/docs/languages/go/quickstart/

go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28

go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2

protoc -Igreet/proto --go_out=. --go_opt=module=github.com/golang-grpc --go-grpc_out=. --go-grpc_opt=module=github.com/golang-grpc greet/proto/greet.proto


protoc -Icalculator/proto --go_out=. --go_opt=module=github.com/golang-grpc --go-grpc_out=. --go-grpc_opt=module=github.com/golang-grpc calculator/proto/calculator.proto


protoc -Iblog/proto --go_out=. --go_opt=module=github.com/golang-grpc --go-grpc_out=. --go-grpc_opt=module=github.com/golang-grpc blog/proto/blog.proto