package rpc

import (
	"context"
	"strings"

	"github.com/rs/xid"
	"google.golang.org/protobuf/proto"

	"github.com/sascha-andres/csv2json/pb"
	"github.com/sascha-andres/csv2json/storer"
)

// ListProjects returns all known projects
func (s *adminServer) ListProjects(_ context.Context, _ *pb.ListProjectsRequest) (*pb.ListProjectsResponse, error) {
	projects, err := s.storage.ListProjects()
	if err != nil {
		return &pb.ListProjectsResponse{Errors: []*pb.Error{{Message: err.Error(), Severity: pb.Severity_CRITICAL}}}, nil
	}
	resp := &pb.ListProjectsResponse{}
	resp.Projects = make([]*pb.ListProjectData, 0)
	for _, project := range projects {
		resp.Projects = append(resp.Projects, &pb.ListProjectData{
			Id:          project.Id,
			Name:        project.Name,
			Description: project.Description,
		})
	}
	return resp, nil
}

// RemoveProject and all related data
func (s *adminServer) RemoveProject(_ context.Context, req *pb.RemoveProjectRequest) (*pb.RemoveProjectResponse, error) {
	err := s.storage.RemoveProject(req.Project)
	if err != nil {
		resp := &pb.RemoveProjectResponse{Errors: []*pb.Error{}}
		for _, e := range err {
			resp.Errors = append(resp.Errors, &pb.Error{Message: e.Error(), Severity: pb.Severity_WARN})
		}
		return resp, nil
	}
	return &pb.RemoveProjectResponse{}, nil
}

// CreateProject handles the creation of a new project using the provided request data and stores it in the database.
// Returns a response with the new project's ID or an error if the operation fails.
func (s *adminServer) CreateProject(_ context.Context, req *pb.CreateProjectRequest) (*pb.CreateProjectResponse, error) {
	newProjectId := xid.New().String()
	project := storer.Project{
		Id:   newProjectId,
		Name: req.Name,
	}
	if req.Array == nil || !*req.Array {
		project.Array = proto.Bool(false)
	} else {
		project.Array = proto.Bool(true)
	}
	if req.Named == nil || !*req.Named {
		project.Named = proto.Bool(false)
	} else {
		project.Named = proto.Bool(true)
	}
	if req.Description == nil {
		project.Description = proto.String("")
	} else {
		project.Description = req.Description
	}
	if req.NestedPropertyName == nil {
		project.NestedPropertyName = proto.String("")
	} else {
		project.NestedPropertyName = req.NestedPropertyName
	}
	ot := storer.OutputType(strings.ToLower(req.GetOutputType().String()))
	project.OutputType = &ot
	if req.Separator == nil {
		project.Separator = proto.String("")
	} else {
		project.Separator = req.Separator
	}
	err := s.storage.CreateProject(project)
	if err != nil {
		return &pb.CreateProjectResponse{Errors: []*pb.Error{&pb.Error{Message: err.Error(), Severity: pb.Severity_CRITICAL}}}, nil
	}
	return &pb.CreateProjectResponse{Id: &newProjectId}, nil
}
