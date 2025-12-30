package protoreg

import (
	"os"
	"path/filepath"

	"github.com/roadrunner-server/errors"
)

type Config struct {
	ImportPaths []string `mapstructure:"import_paths"`
	Proto       []string `mapstructure:"proto"`
}

func (c *Config) InitDefaults() error { //nolint:gocyclo,gocognit
	const op = errors.Op("protoreg_plugin_config")

	importPaths := make([]string, 0, len(c.ImportPaths))
	for _, path := range c.ImportPaths {
		if path == "" {
			continue
		}

		importPath, err := filepath.Abs(path)
		if err != nil {
			return errors.E(op, err)
		}

		if _, err := os.Stat(path); err != nil {
			if os.IsNotExist(err) {
				return errors.E(op, errors.Errorf("import path '%s' does not exist", importPath))
			}

			return errors.E(op, err)
		}

		importPaths = append(importPaths, path)
	}

	if len(importPaths) == 0 {
		return errors.E(op, errors.Errorf("no import path specified"))
	}

	protos := make([]string, 0)
	for _, path := range c.Proto {
		if path == "" {
			continue
		}

		exists := false

		for _, importPath := range importPaths {
			if _, err := os.Stat(filepath.Join(importPath, path)); err == nil {
				exists = true
			}
		}

		if !exists {
			return errors.E(op, errors.Errorf("proto file '%s' does not exist", path))
		}

		protos = append(protos, path)
	}

	if len(protos) == 0 {
		return errors.E(op, errors.Errorf("no proto specified"))
	}

	c.Proto = protos

	return nil
}
