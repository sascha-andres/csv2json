package file

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"

	"gocloud.dev/blob"
	"gocloud.dev/gcerrors"

	"github.com/sascha-andres/csv2json/storer"
)

// ListProjects returns all known projects
func (s Storer) ListProjects() ([]storer.Project, error) {
	iter := s.bucket.List(&blob.ListOptions{
		Delimiter: "/",
		Prefix:    "projects/",
	})
	result := make([]storer.Project, 0)
	ctx := context.Background()
	for {
		obj, err := iter.Next(ctx)
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		reader, err := s.bucket.NewReader(ctx, obj.Key, nil)
		if err != nil {
			return nil, err
		}
		var buf bytes.Buffer
		_, err = reader.WriteTo(&buf)
		if err != nil {
			return nil, err
		}
		var p storer.Project
		err = json.Unmarshal(buf.Bytes(), &p)
		if err != nil {
			return nil, err
		}
		result = append(result, p)
	}
	return result, nil
}

// RemoveProject removes project data (incl all run data)
func (s Storer) RemoveProject(id string) []error {
	result := make([]error, 0)
	if err := s.ClearMappings(id); err != nil {
		errCode := gcerrors.Code(err)
		if errCode != gcerrors.NotFound {
			result = append(result, err)
		}
	}
	if err := s.ClearExtraVariables(id); err != nil {
		errCode := gcerrors.Code(err)
		if errCode != gcerrors.NotFound {
			result = append(result, err)
		}
	}
	if err := s.ClearCalculatedFields(id); err != nil {
		errCode := gcerrors.Code(err)
		if errCode != gcerrors.NotFound {
			result = append(result, err)
		}
	}
	if err := s.bucket.Delete(context.Background(), projectPathForId(storer.Project{Id: id})); err != nil {
		result = append(result, err)
	}
	return result
}

// CreateProject is used to create a project
func (s Storer) CreateProject(p storer.Project) error {
	w, err := s.bucket.NewWriter(context.Background(), projectPathForId(p), nil)
	if err != nil {
		return err
	}
	data, err := json.Marshal(p)
	if err != nil {
		return err
	}
	_, err = w.Write(data)
	if err != nil {
		return err
	}
	return w.Close()
}

// projectPathForId returns the project path in storage
func projectPathForId(p storer.Project) string {
	return fmt.Sprintf("%s/project.json", p.Id)
}
