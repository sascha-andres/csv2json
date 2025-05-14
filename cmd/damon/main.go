package main

import (
	"github.com/sascha-andres/reuse/flag"

	"github.com/sascha-andres/csv2json/cmd/damon/rpc"
)

var (
	storageDsn, iface string
	port              uint
)

func init() {
	flag.SetEnvPrefix("CSV2JSON_DAEMON")
	flag.StringVar(&storageDsn, "storage-dsn", "file:file://./file-storage", "storage dsn")
	flag.StringVar(&iface, "interface", "", "interface to listen on")
	flag.UintVar(&port, "port", 50501, "port to listen on")
}

func main() {
	flag.Parse()

	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	r, err := rpc.NewRpc(rpc.WithPort(port), rpc.WithStorageDsn(storageDsn), rpc.WithInterface(iface))
	if err != nil {
		return err
	}
	return r.Run()
}
