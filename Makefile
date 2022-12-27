.PHONY: proto
## proto: compiles .proto files
proto:
	docker run --rm -v $(PWD):/defs namely/protoc-all -f grpc/interface.proto -l go -o . --go-source-relative

