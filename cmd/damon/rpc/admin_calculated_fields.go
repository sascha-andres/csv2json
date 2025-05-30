package rpc

import (
	"context"

	"github.com/sascha-andres/csv2json"
	"github.com/sascha-andres/csv2json/pb"
)

// AddOrUpdateCalculatedFields add one or more calculated fields to a project
func (s *adminServer) AddOrUpdateCalculatedFields(ctx context.Context, req *pb.AddOrUpdateCalculatedFieldsRequest) (*pb.AddOrUpdateCalculatedFieldsResponse, error) {
	r := make(map[string]csv2json.CalculatedField)
	for k, field := range req.CalculatedFields {
		r[k] = csv2json.CalculatedField{
			Kind:     pb.Kind_name[int32(field.Kind)],
			Format:   field.Format,
			Type:     pb.FieldType_name[int32(field.Type)],
			Location: pb.Location_name[int32(field.Location)],
		}
	}
	err := s.storage.CreateCalculatedFields(req.Project, r)
	if err != nil {
		return &pb.AddOrUpdateCalculatedFieldsResponse{Errors: []*pb.Error{{Message: err.Error(), Severity: pb.Severity_CRITICAL}}}, nil
	}
	return &pb.AddOrUpdateCalculatedFieldsResponse{}, nil
}

// RemoveCalculatedFields remove on or more calculated fields from a project
func (s *adminServer) RemoveCalculatedFields(ctx context.Context, req *pb.RemoveCalculatedFieldsRequest) (*pb.RemoveCalculatedFieldsResponse, error) {
	err := s.storage.RemoveCalculatedFields(req.Project, req.CalculatedFields)
	if err != nil {
		return &pb.RemoveCalculatedFieldsResponse{Errors: []*pb.Error{{Message: err.Error(), Severity: pb.Severity_CRITICAL}}}, nil
	}
	return &pb.RemoveCalculatedFieldsResponse{}, nil
}

// ListCalculatedFields for a project
func (s *adminServer) ListCalculatedFields(ctx context.Context, req *pb.ListCalculatedFieldsRequest) (*pb.ListCalculatedFieldsResponse, error) {
	fields, err := s.storage.GetCalculatedFields(req.Project)
	if err != nil {
		return &pb.ListCalculatedFieldsResponse{Errors: []*pb.Error{{Message: err.Error(), Severity: pb.Severity_CRITICAL}}}, nil
	}
	result := make(map[string]*pb.CalculatedField)
	for k, field := range fields {
		result[k] = &pb.CalculatedField{
			Kind:     pb.Kind(pb.Kind_value[field.Kind]),
			Format:   field.Format,
			Type:     pb.FieldType(pb.FieldType_value[field.Type]),
			Location: pb.Location(pb.Location_value[field.Location]),
		}
	}
	return &pb.ListCalculatedFieldsResponse{CalculatedFields: result}, nil
}
