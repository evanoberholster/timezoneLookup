package pb

//go:generate go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
//go:generate protoc --proto_path=. --go_out=. --go_opt=paths=source_relative ./timezone.proto
