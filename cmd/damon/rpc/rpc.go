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
	adminServer struct {
		pb.UnimplementedAdminServiceServer
	}

	Rpc struct {
		iface      string
		port       uint
		storageDsn string

		s *grpc.Server
		p types.Storer
	}
)

func WithPort(port uint) OptionFunc {
	return func(r *Rpc) error {
		r.port = port
		return nil
	}
}

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

func WithInterface(iface string) OptionFunc {
	return func(r *Rpc) error {
		r.iface = iface
		return nil
	}
}

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

func (r *Rpc) Run() error {
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", r.iface, r.port))
	if err != nil {
		return err
	}
	pb.RegisterAdminServiceServer(r.s, &adminServer{})
	reflection.Register(r.s)
	if err := r.s.Serve(lis); err != nil {
		return err
	}
	return nil
}
