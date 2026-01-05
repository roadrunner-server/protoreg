package protoreg

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/roadrunner-server/errors"
)

type Config struct {
	ProtoPath []string `mapstructure:"proto_path"`
	Files     []string `mapstructure:"files"`
}

func (c *Config) InitDefaults() error { //nolint:gocyclo,gocognit
	const op = errors.Op("protoreg_plugin_config")

	// validate input first
	// both proto_path and files are required and should not be empty
	if len(c.ProtoPath) == 0 || len(c.Files) == 0 {
		return errors.E(op, errors.Errorf("proto_path and files must be specified"))
	}

	// map to track found files
	// struct{} to save memory
	found := make(map[string]struct{}, len(c.Files))

	// O(n*m) where n - proto_paths, m - files
	// but n and m are expected to be small
	// so this is acceptable
	for _, ppath := range c.ProtoPath {
		// do not believe user input - validate each proto_path entry
		if strings.TrimSpace(ppath) == "" {
			return errors.E(op, errors.Errorf("proto_path entry cannot be empty"))
		}
		// take absolute path of proto_path entry
		abs, err := filepath.Abs(ppath)
		if err != nil {
			return errors.E(op, err)
		}

		// validate proto paths with abs path
		if _, err := os.Stat(abs); err != nil {
			return errors.E(op, err)
		}

		// for that proto path, check all user provided files
		// the case is, that this file can be located in any of the proto paths
		// and if we don't find it in any, we return an error (check below)
		// in this loop we check if we already found the file - if yes, skip searching
		// else, try to find it in the current proto_path
		for _, file := range c.Files {
			// check for the edge case when the filename is empty string
			if strings.TrimSpace(file) == "" {
				return errors.E(op, errors.Errorf("proto file path cannot be empty"))
			}
			// if we saw it already, skip searching
			// this avoids redundant os.Stat calls
			// not much optimization, but still...
			if _, ok := found[file]; ok {
				continue
			}

			// check if file exists in the current proto_path
			// if yes, mark it as found
			if _, err := os.Stat(filepath.Join(abs, file)); err == nil {
				found[file] = struct{}{}
			}
		}
	}

	// after searching all proto paths, ensure all files were found
	// this is simple O(n) check
	for _, file := range c.Files {
		if _, ok := found[file]; !ok {
			return errors.E(op, errors.Errorf("proto file '%s' does not exist in any of the proto_paths", file))
		}
	}

	return nil
}
