.PHONY: build

build: zxcrpc.pb.go

%.pb.go: proto/%.proto
	protoc -I proto/ --go_out=plugins=grpc:. $<
