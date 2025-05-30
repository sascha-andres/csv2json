package rpc

import (
	"context"
	"fmt"
	"strings"

	"github.com/sascha-andres/csv2json"
	"github.com/sascha-andres/csv2json/pb"
)

// AddOrUpdateMapping to a project (or multiple)
func (s *adminServer) AddOrUpdateMapping(ctx context.Context, req *pb.AddOrUpdateMappingRequest) (*pb.AddOrUpdateMappingResponse, error) {
	mappings := make(map[string]csv2json.ColumnConfiguration)
	for key, configuration := range req.Mappings {
		mappings[key] = csv2json.ColumnConfiguration{Property: configuration.Property, Type: strings.ToLower(configuration.GetType().String())}
	}
	at, err := s.storage.CreateMappings(req.Project, mappings)
	if err != nil {
		return &pb.AddOrUpdateMappingResponse{Errors: []*pb.Error{{Message: err.Error(), Severity: pb.Severity_CRITICAL}}}, nil
	}
	res := &pb.AddOrUpdateMappingResponse{}
	if res.ActionsTaken == nil {
		res.ActionsTaken = make(map[string]pb.ActionTaken)
	}
	for key, taken := range at {
		res.ActionsTaken[key] = taken
	}
	return res, nil
}

// RemoveMapping (s)
func (s *adminServer) RemoveMapping(ctx context.Context, req *pb.RemoveMappingRequest) (*pb.RemoveMappingResponse, error) {
	err := s.storage.RemoveMappings(req.Project, req.Mappings)
	if err != nil {
		return &pb.RemoveMappingResponse{Errors: []*pb.Error{{Message: err.Error(), Severity: pb.Severity_CRITICAL}}}, nil
	}
	return &pb.RemoveMappingResponse{}, nil
}

// ListMappings for project
func (s *adminServer) ListMappings(ctx context.Context, req *pb.ListMappingsRequest) (*pb.ListMappingsResponse, error) {
	mappings, err := s.storage.GetMappings(req.Project)
	if err != nil {
		return &pb.ListMappingsResponse{Errors: []*pb.Error{{Message: err.Error(), Severity: pb.Severity_CRITICAL}}}, nil
	}
	result := &pb.ListMappingsResponse{}
	result.Mappings = make(map[string]*pb.ColumnConfiguration)
	for key, configuration := range mappings {
		result.Mappings[key] = &pb.ColumnConfiguration{Property: configuration.Property}
		switch configuration.Type {
		case "string":
			result.Mappings[key].Type = pb.FieldType_STRING
		case "integer":
			result.Mappings[key].Type = pb.FieldType_INT
		case "float":
			result.Mappings[key].Type = pb.FieldType_FLOAT
		case "bool":
			result.Mappings[key].Type = pb.FieldType_BOOL
		default:
			return &pb.ListMappingsResponse{Errors: []*pb.Error{{Message: fmt.Sprintf("unknown type %s", configuration.Type), Severity: pb.Severity_CRITICAL}}}, nil
		}
	}
	return result, nil
}
