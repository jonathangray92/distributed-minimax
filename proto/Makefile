all: proto

proto: slave.pb.go user.pb.go

%.pb.go: %.proto
	protoc --go_out=. $<

clean:
	rm -f *.pb.go
