package file

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
)

const extraVariablesPathName = "extra_variables"

// CreateExtraVariables to add a static variable to the project
func (s Storer) CreateExtraVariables(projectID string, variables map[string]string) error {
	extraVariables, err := s.loadExtraVariables(projectID)
	if err != nil {
		return err
	}
	if extraVariables == nil {
		extraVariables = make(map[string]string)
	}
	for k, v := range variables {
		extraVariables[k] = v
	}
	return s.saveExtraVariables(projectID, extraVariables)
}

// GetExtraVariables to load all extra variables
func (s Storer) GetExtraVariables(projectID string) (map[string]string, error) {
	return s.loadExtraVariables(projectID)
}

// RemoveExtraVariables to remove one or more extra variables
func (s Storer) RemoveExtraVariables(projectID string, extraVariables []string) error {
	ev, err := s.loadExtraVariables(projectID)
	if err != nil {
		return err
	}
	if ev == nil {
		return nil
	}
	for _, k := range extraVariables {
		delete(ev, k)
	}
	return s.saveExtraVariables(projectID, ev)
}

// ClearExtraVariables to clear the project from extra variables
func (s Storer) ClearExtraVariables(projectID string) error {
	return s.bucket.Delete(context.Background(), getExtraVariablesPathForProject(projectID))
}

// loadExtraVariables deserializes mappings from file storage
func (s Storer) loadExtraVariables(projectID string) (map[string]string, error) {
	path := getExtraVariablesPathForProject(projectID)
	ctx := context.Background()
	known, err := s.bucket.Exists(ctx, path)
	if err != nil {
		return nil, err
	}
	if !known {
		return nil, nil
	}
	reader, err := s.bucket.NewReader(ctx, path, nil)
	if err != nil {
		return nil, err
	}
	var buf bytes.Buffer
	_, err = reader.WriteTo(&buf)
	if err != nil {
		return nil, err
	}
	var mapping map[string]string
	err = json.Unmarshal(buf.Bytes(), &mapping)
	if err != nil {
		return nil, err
	}
	return mapping, nil
}

// saveExtraVariables persists to storage
func (s Storer) saveExtraVariables(id string, variables map[string]string) error {
	path := getExtraVariablesPathForProject(id)
	ctx := context.Background()
	w, err := s.bucket.NewWriter(ctx, path, nil)
	if err != nil {
		return err
	}
	data, err := json.Marshal(variables)
	if err != nil {
		return err
	}
	_, err = w.Write(data)
	if err != nil {
		return err
	}
	return w.Close()
}

// getExtraVariablesPathForProject returns the file storage path for mappings
func getExtraVariablesPathForProject(projectId string) string {
	return fmt.Sprintf("%s/%s.json", projectId, extraVariablesPathName)
}
