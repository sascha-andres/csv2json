package file

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"github.com/sascha-andres/csv2json"
	"github.com/sascha-andres/csv2json/pb"
)

const mappingsPathName = "mappings"

// CreateMappings persists mappings for a given project
func (s Storer) CreateMappings(projectID string, columns map[string]csv2json.ColumnConfiguration) (map[string]pb.ActionTaken, error) {
	mappings, err := s.loadMappings(projectID)
	if err != nil {
		return nil, err
	}
	result := make(map[string]pb.ActionTaken)
	if mappings == nil {
		mappings = make(map[string]csv2json.ColumnConfiguration)
	}
	for k, v := range columns {
		if _, ok := mappings[k]; !ok {
			result[k] = pb.ActionTaken_ADDED
		} else {
			result[k] = pb.ActionTaken_UPDATED
		}
		mappings[k] = v
	}
	return result, s.saveMappings(projectID, mappings)
}

// RemoveMappings removes all mappings provided from storage
func (s Storer) RemoveMappings(projectID string, columns []string) error {
	mappings, err := s.loadMappings(projectID)
	if err != nil {
		return err
	}
	if mappings == nil {
		return nil
	}
	for _, k := range columns {
		delete(mappings, k)
	}
	return s.saveMappings(projectID, mappings)
}

// GetMappings reads all mappings and returns those
func (s Storer) GetMappings(projectID string) (map[string]csv2json.ColumnConfiguration, error) {
	return s.loadMappings(projectID)
}

// ClearMappings is called when project ist removed
func (s Storer) ClearMappings(projectID string) error {
	return s.bucket.Delete(context.Background(), getMappingPathForProject(projectID))
}

// saveMappings writes all mappings to storage
func (s Storer) saveMappings(projectID string, columns map[string]csv2json.ColumnConfiguration) error {
	path := getMappingPathForProject(projectID)
	ctx := context.Background()
	w, err := s.bucket.NewWriter(ctx, path, nil)
	if err != nil {
		return err
	}
	data, err := json.Marshal(columns)
	if err != nil {
		return err
	}
	_, err = w.Write(data)
	if err != nil {
		return err
	}
	return w.Close()
}

// loadMappings deserializes mappings from file storage
func (s Storer) loadMappings(projectID string) (map[string]csv2json.ColumnConfiguration, error) {
	path := getMappingPathForProject(projectID)
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
	var mapping map[string]csv2json.ColumnConfiguration
	err = json.Unmarshal(buf.Bytes(), &mapping)
	if err != nil {
		return nil, err
	}
	return mapping, nil
}

// getMappingPathForProject returns the file storage path for mappings
func getMappingPathForProject(projectId string) string {
	return fmt.Sprintf("admin/projects/%s/%s.json", projectId, mappingsPathName)
}
