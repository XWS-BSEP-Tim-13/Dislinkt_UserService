protoc -I ./infrastructure/grpc/proto --go_out ./infrastructure/grpc/proto --go_opt paths=source_relative --go-grpc_out ./infrastructure/grpc/proto --go-grpc_opt paths=source_relative --grpc-gateway_out ./infrastructure/grpc/proto --grpc-gateway_opt paths=source_relative ./infrastructure/grpc/proto/user_service.proto