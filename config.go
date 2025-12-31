package protoreg

import (
	"os"
	"path/filepath"

	"github.com/roadrunner-server/errors"
)

type Config struct {
	ProtoPath []string `mapstructure:"proto_path"`
	Files     []string `mapstructure:"files"`
}

func (c *Config) InitDefaults() error { //nolint:gocyclo,gocognit
	const op = errors.Op("protoreg_plugin_config")

	protoPaths := make([]string, 0, len(c.ProtoPath))
	for _, path := range c.ProtoPath {
		if path == "" {
			continue
		}

		importPath, err := filepath.Abs(path)
		if err != nil {
			return errors.E(op, err)
		}

		if _, err := os.Stat(path); err != nil {
			if os.IsNotExist(err) {
				return errors.E(op, errors.Errorf("proto_path '%s' does not exist", importPath))
			}

			return errors.E(op, err)
		}

		protoPaths = append(protoPaths, path)
	}

	if len(protoPaths) == 0 {
		return errors.E(op, errors.Errorf("no proto_path specified"))
	}

	files := make([]string, 0)
	for _, path := range c.Files {
		if path == "" {
			continue
		}

		exists := false

		for _, protoPath := range protoPaths {
			if _, err := os.Stat(filepath.Join(protoPath, path)); err == nil {
				exists = true
			}
		}

		if !exists {
			return errors.E(op, errors.Errorf("proto file '%s' does not exist", path))
		}

		files = append(files, path)
	}

	if len(files) == 0 {
		return errors.E(op, errors.Errorf("no proto files specified"))
	}

	c.Files = files

	return nil
}
