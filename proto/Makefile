proto: slave.pb.go

%.pb.go: %.proto
	protoc --go_out=. $<

clean:
	rm -f *.pb.go
