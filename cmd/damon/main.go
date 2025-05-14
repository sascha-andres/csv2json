package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/sascha-andres/reuse/flag"

	"github.com/sascha-andres/csv2json/cmd/damon/rpc"
)

var (
	storageDsn, iface string
	port              uint
)

// init initializes the command-line flags and environment variables.
func init() {
	flag.SetEnvPrefix("CSV2JSON_DAEMON")
	flag.StringVar(&storageDsn, "storage-dsn", "file:file://./file-storage", "storage dsn")
	flag.StringVar(&iface, "interface", "", "interface to listen on")
	flag.UintVar(&port, "port", 50501, "port to listen on")
}

// main parses flags, executes the application logic via the run function, and handles any errors by panicking.
func main() {
	flag.Parse()

	if err := run(); err != nil {
		panic(err)
	}
}

// run initializes an RPC server with provided configurations and handles graceful shutdown.
// It returns an error if any step in the process fails.
func run() error {
	r, err := rpc.NewRpc(rpc.WithPort(port), rpc.WithStorageDsn(storageDsn), rpc.WithInterface(iface))
	if err != nil {
		return err
	}

	// Create a channel to listen for OS signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Run the RPC server in a goroutine
	errChan := make(chan error, 1)
	go func() {
		fmt.Printf("Starting RPC server on %s:%d\n", iface, port)
		if err := r.Run(); err != nil {
			errChan <- err
		}
	}()

	// Wait for a signal or an error from the server
	select {
	case sig := <-sigChan:
		fmt.Printf("Received signal: %v, initiating graceful shutdown\n", sig)
		r.GracefulStop()
		fmt.Println("Server gracefully stopped")
	case err := <-errChan:
		fmt.Printf("Server error: %v\n", err)
		return err
	}

	return nil
}
