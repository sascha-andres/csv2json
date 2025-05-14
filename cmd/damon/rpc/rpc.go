package rpc

import (
	"fmt"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/sascha-andres/csv2json/internal/persistence"
	"github.com/sascha-andres/csv2json/internal/types"
	"github.com/sascha-andres/csv2json/pb"
)

// OptionFunc defines a function signature for configuring an Rpc instance with specific options or parameters.
type OptionFunc func(*Rpc) error

type (
	// adminServer implements the AdminService interface and provides RPC methods for managing projects.
	adminServer struct {
		pb.UnimplementedAdminServiceServer

		// storage provides access to a Storer interface for managing projects, mappings, and related data storage operations.
		storage types.Storer
	}

	// Rpc represents a gRPC server configuration and dependencies for RPC service execution.
	// It holds server details and storage interface for project management operations.
	Rpc struct {

		// iface specifies the network interface to bind the RPC server to, e.g., "localhost" or "0.0.0.0".
		iface string

		// port specifies the port number to bind the RPC server to.
		port uint

		// storageDsn specifies the DSN for the storage backend to use for project management operations.
		storageDsn string

		// s represents the gRPC server instance used to handle gRPC-based RPC communication in the application.
		s *grpc.Server

		// p provides access to a Storer interface for managing projects, mappings, and related data storage operations.
		p types.Storer
	}
)

// WithPort sets the port number for the Rpc server configuration. It returns an OptionFunc to apply this setting.
func WithPort(port uint) OptionFunc {
	return func(r *Rpc) error {
		r.port = port
		return nil
	}
}

// WithStorageDsn sets the storage DSN for the Rpc instance and initializes its storage layer using the provided DSN.
func WithStorageDsn(dsn string) OptionFunc {
	return func(r *Rpc) error {
		r.storageDsn = dsn
		s, err := persistence.GetStorer(r.storageDsn)
		if err != nil {
			return err
		}
		r.p = s
		return nil
	}
}

// WithInterface sets the network interface for the Rpc instance. It returns an OptionFunc to apply this configuration.
func WithInterface(iface string) OptionFunc {
	return func(r *Rpc) error {
		r.iface = iface
		return nil
	}
}

// NewRpc initializes and returns a new Rpc instance configured with the provided options or an error if configuration fails.
func NewRpc(opts ...OptionFunc) (*Rpc, error) {
	r := &Rpc{}
	for _, opt := range opts {
		if err := opt(r); err != nil {
			return nil, err
		}
	}
	r.s = grpc.NewServer()
	return r, nil
}

// Run starts the RPC server and blocks until it is stopped.
func (r *Rpc) Run() error {
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", r.iface, r.port))
	if err != nil {
		return err
	}
	pb.RegisterAdminServiceServer(r.s, &adminServer{storage: r.p})
	reflection.Register(r.s)
	if err := r.s.Serve(lis); err != nil {
		return err
	}
	return nil
}
