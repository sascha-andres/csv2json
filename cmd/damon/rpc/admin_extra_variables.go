package rpc

import (
	"context"

	"github.com/sascha-andres/csv2json/pb"
)

// AddExtraVariable to reference in calculated fields
func (s *adminServer) AddExtraVariable(ctx context.Context, req *pb.AddExtraVariableRequest) (*pb.AddExtraVariableResponse, error) {
	err := s.storage.CreateExtraVariables(req.Project, req.ExtraVariables)
	if err != nil {
		return &pb.AddExtraVariableResponse{Errors: []*pb.Error{{Message: err.Error(), Severity: pb.Severity_CRITICAL}}}, nil
	}
	return &pb.AddExtraVariableResponse{}, nil
}

// RemoveExtraVariable from project
func (s *adminServer) RemoveExtraVariable(ctx context.Context, req *pb.RemoveExtraVariableRequest) (*pb.RemoveExtraVariableResponse, error) {
	err := s.storage.RemoveExtraVariables(req.Project, req.ExtraVariables)
	if err != nil {
		return &pb.RemoveExtraVariableResponse{Errors: []*pb.Error{{Message: err.Error(), Severity: pb.Severity_CRITICAL}}}, nil
	}
	return &pb.RemoveExtraVariableResponse{}, nil
}

// ListExtraVariables for project
func (s *adminServer) ListExtraVariables(ctx context.Context, req *pb.ListExtraVariablesRequest) (*pb.ListExtraVariablesResponse, error) {
	r, err := s.storage.GetExtraVariables(req.Project)
	if err != nil {
		return &pb.ListExtraVariablesResponse{Errors: []*pb.Error{{Message: err.Error(), Severity: pb.Severity_CRITICAL}}}, nil
	}
	return &pb.ListExtraVariablesResponse{ExtraVariables: r}, nil
}
