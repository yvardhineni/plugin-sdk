package specs

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
)

// Source is the spec for a source plugin
type Source struct {
	// Name of the source plugin to use
	Name string `json:"name,omitempty"`
	// Version of the source plugin to use
	Version string `json:"version,omitempty"`
	// Path is the canonical path to the source plugin in a given registry
	// For example:
	// in github the path will be: org/repo
	// For the local registry the path will be the path to the binary: ./path/to/binary
	// For the gRPC registry the path will be the address of the gRPC server: host:port
	Path string `json:"path,omitempty"`
	// Registry can be github,local,grpc.
	Registry    Registry `json:"registry,omitempty"`
	Concurrency uint64   `json:"concurrency,omitempty"`
	// Tables to sync from the source plugin
	Tables []string `json:"tables,omitempty"`
	// SkipTables defines tables to skip when syncing data. Useful if a glob pattern is used in Tables
	SkipTables []string `json:"skip_tables,omitempty"`
	// Destinations are the names of destination plugins to send sync data to
	Destinations []string `json:"destinations,omitempty"`
	// Spec defines plugin specific configuration
	// This is different in every source plugin.
	Spec interface{} `json:"spec,omitempty"`
}

func (s *Source) SetDefaults() {
	if s.Registry.String() == "" {
		s.Registry = RegistryGithub
	}
	if s.Path == "" {
		s.Path = s.Name
	}
	if s.Registry == RegistryGithub && !strings.Contains(s.Path, "/") {
		s.Path = "cloudquery/" + s.Path
	}
	if s.Tables == nil {
		s.Tables = []string{"*"}
	}
}

// UnmarshalSpec unmarshals the internal spec into the given interface
func (s *Source) UnmarshalSpec(out interface{}) error {
	b, err := json.Marshal(s.Spec)
	if err != nil {
		return err
	}
	dec := json.NewDecoder(bytes.NewReader(b))
	dec.UseNumber()
	dec.DisallowUnknownFields()
	return dec.Decode(out)
}

func (s *Source) Validate() error {
	if s.Name == "" {
		return fmt.Errorf("name is required")
	}
	if s.Version == "" {
		return fmt.Errorf("version is required")
	}
	if !strings.HasPrefix(s.Version, "v") {
		return fmt.Errorf("version must start with v")
	}
	if len(s.Destinations) == 0 {
		return fmt.Errorf("at least one destination is required")
	}
	return nil
}
