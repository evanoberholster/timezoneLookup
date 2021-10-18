package cp

//go:generate go install capnproto.org/go/capnp/v3/capnpc-go@latest
//go:generate sh -c "capnp compile -I$(go env GOPATH)/pkg/mod/capnproto.org/go/capnp/v3@v3.0.0-alpha.1/std -ogo *.capnp"
