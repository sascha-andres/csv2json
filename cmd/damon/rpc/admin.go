package rpc

import (
	"context"

	"github.com/rs/xid"

	"github.com/sascha-andres/csv2json/pb"
)

func (s *adminServer) CreateProject(_ context.Context, req *pb.CreateProjectRequest) (*pb.CreateProjectResponse, error) {
	newProjectId := xid.New().String()

	return &pb.CreateProjectResponse{Id: &newProjectId}, nil
}
