package rpc

import (
	"context"
	"strings"

	"github.com/rs/xid"

	"github.com/sascha-andres/csv2json/pb"
	"github.com/sascha-andres/csv2json/storer"
)

// getValueOrEmptyString returns the dereferenced value of the string pointer `s` or an empty string if `s` is nil.
func getValueOrEmptyString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

// CreateProject handles the creation of a new project using the provided request data and stores it in the database.
// Returns a response with the new project's ID or an error if the operation fails.
func (s *adminServer) CreateProject(_ context.Context, req *pb.CreateProjectRequest) (*pb.CreateProjectResponse, error) {
	newProjectId := xid.New().String()
	project := storer.Project{
		Id:   newProjectId,
		Name: req.Name,
	}
	project.Array = req.Array
	project.Named = req.Named
	project.Description = req.Description
	project.NestedPropertyName = req.NestedPropertyName
	ot := storer.OutputType(strings.ToLower(req.GetOutputType().String()))
	project.OutputType = &ot
	project.Separator = req.Separator
	err := s.storage.CreateProject(project)
	if err != nil {
		return &pb.CreateProjectResponse{Errors: []*pb.Error{&pb.Error{Message: err.Error(), Severity: pb.Severity_CRITICAL}}}, nil
	}
	return &pb.CreateProjectResponse{Id: &newProjectId}, nil
}
