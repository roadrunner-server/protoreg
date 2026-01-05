package protoreg

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInitDefaults(t *testing.T) {
	tests := []struct {
		name      string
		protoPath []string
		files     []string
		assertion assert.ErrorAssertionFunc
	}{
		{
			name:      "empty config",
			protoPath: nil,
			files:     nil,
			assertion: assert.Error,
		},
		{
			name: "empty string in files",
			protoPath: []string{
				"./tests/proto/commonapis",
				"./tests/proto/serviceapis",
			},
			files:     []string{""},
			assertion: assert.Error,
		},
		{
			name: "whitespace strings in files",
			protoPath: []string{
				"./tests/proto/commonapis",
				"./tests/proto/serviceapis",
			},
			files:     []string{"        "},
			assertion: assert.Error,
		},
		{
			name: "whitespace strings in files",
			protoPath: []string{
				"./tests/proto/commonapis",
				"./tests/proto/serviceapis",
			},
			files:     []string{"\n"},
			assertion: assert.Error,
		},
		{
			name: "whitespace strings in files",
			protoPath: []string{
				"./tests/proto/commonapis",
				"./tests/proto/serviceapis",
			},
			files:     []string{"\t"},
			assertion: assert.Error,
		},
		{
			name: "whitespace strings in files",
			protoPath: []string{
				"./tests/proto/commonapis",
				"./tests/proto/serviceapis",
			},
			files:     []string{"        ", "\n", "\t", "  "},
			assertion: assert.Error,
		},
		{
			name: "unknown file",
			protoPath: []string{
				"./tests/proto/commonapis",
				"./tests/proto/serviceapis",
			},
			files:     []string{"unknown/v1/message.proto"},
			assertion: assert.Error,
		},
		{
			name: "valid config",
			protoPath: []string{
				"./tests/proto/commonapis",
				"./tests/proto/serviceapis",
			},
			files:     []string{"service/v1/service.proto"},
			assertion: assert.NoError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := Config{
				ProtoPath: tt.protoPath,
				Files:     tt.files,
			}
			tt.assertion(t, c.InitDefaults())
		})
	}
}

func TestInitDefaults_EmptyInputs(t *testing.T) {
	tests := []struct {
		name      string
		protoPath []string
		files     []string
		assertion assert.ErrorAssertionFunc
	}{
		{
			name:      "both empty",
			protoPath: []string{},
			files:     []string{},
			assertion: assert.Error,
		},
		{
			name:      "empty proto_path",
			protoPath: []string{},
			files:     []string{"service/v1/service.proto"},
			assertion: assert.Error,
		},
		{
			name:      "empty files",
			protoPath: []string{"./tests/proto/serviceapis"},
			files:     []string{},
			assertion: assert.Error,
		},
		{
			name:      "nil slices",
			protoPath: nil,
			files:     nil,
			assertion: assert.Error,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := Config{
				ProtoPath: tt.protoPath,
				Files:     tt.files,
			}
			tt.assertion(t, c.InitDefaults())
		})
	}
}

func TestInitDefaults_ProtoPathValidation(t *testing.T) {
	tests := []struct {
		name      string
		protoPath []string
		files     []string
		assertion assert.ErrorAssertionFunc
	}{
		{
			name:      "non-existent proto_path",
			protoPath: []string{"./non_existent_path"},
			files:     []string{"service/v1/service.proto"},
			assertion: assert.Error,
		},
		{
			name:      "first proto_path invalid second valid",
			protoPath: []string{"./invalid", "./tests/proto/serviceapis"},
			files:     []string{"service/v1/service.proto"},
			assertion: assert.Error, // fails on first invalid path
		},
		{
			name:      "valid relative path",
			protoPath: []string{"./tests/proto/serviceapis"},
			files:     []string{"service/v1/service.proto"},
			assertion: assert.NoError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := Config{
				ProtoPath: tt.protoPath,
				Files:     tt.files,
			}
			tt.assertion(t, c.InitDefaults())
		})
	}
}

func TestInitDefaults_FileResolutionAcrossPaths(t *testing.T) {
	tests := []struct {
		name      string
		protoPath []string
		files     []string
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "file in first proto_path",
			protoPath: []string{
				"./tests/proto/serviceapis",
				"./tests/proto/commonapis",
			},
			files:     []string{"service/v1/service.proto"},
			assertion: assert.NoError,
		},
		{
			name: "file in second proto_path only",
			protoPath: []string{
				"./tests/proto/serviceapis",
				"./tests/proto/commonapis",
			},
			files:     []string{"common/v1/message.proto"}, // only in commonapis
			assertion: assert.NoError,
		},
		{
			name: "files distributed across different proto_paths",
			protoPath: []string{
				"./tests/proto/serviceapis",
				"./tests/proto/commonapis",
			},
			files: []string{
				"service/v1/service.proto", // in serviceapis
				"common/v1/message.proto",  // in commonapis
			},
			assertion: assert.NoError,
		},
		{
			name: "file not found in any proto_path",
			protoPath: []string{
				"./tests/proto/serviceapis",
				"./tests/proto/commonapis",
			},
			files:     []string{"nonexistent/v1/unknown.proto"},
			assertion: assert.Error,
		},
		{
			name: "one file found one not found",
			protoPath: []string{
				"./tests/proto/serviceapis",
				"./tests/proto/commonapis",
			},
			files: []string{
				"service/v1/service.proto",
				"nonexistent/v1/unknown.proto",
			},
			assertion: assert.Error,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := Config{
				ProtoPath: tt.protoPath,
				Files:     tt.files,
			}
			tt.assertion(t, c.InitDefaults())
		})
	}
}

func TestInitDefaults_EarlyTermination(t *testing.T) {
	// File exists in both proto_paths - should be found in first, skipped in second
	// This tests the early termination optimization
	c := Config{
		ProtoPath: []string{
			"./tests/proto/serviceapis",
			"./tests/proto/serviceapis", // duplicate path
		},
		Files: []string{"service/v1/service.proto"},
	}

	assert.NoError(t, c.InitDefaults())
}

func TestInitDefaults_EdgeCases(t *testing.T) {
	tests := []struct {
		name      string
		protoPath []string
		files     []string
		assertion assert.ErrorAssertionFunc
	}{
		{
			name:      "single proto_path single file",
			protoPath: []string{"./tests/proto/serviceapis"},
			files:     []string{"service/v1/service.proto"},
			assertion: assert.NoError,
		},
		{
			name:      "multiple files same proto_path",
			protoPath: []string{"./tests/proto/serviceapis"},
			files: []string{
				"service/v1/service.proto",
				"service/v1/service_dup.proto",
			},
			assertion: assert.NoError,
		},
		{
			name:      "duplicate files in list",
			protoPath: []string{"./tests/proto/serviceapis"},
			files: []string{
				"service/v1/service.proto",
				"service/v1/service.proto",
			},
			assertion: assert.NoError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := Config{
				ProtoPath: tt.protoPath,
				Files:     tt.files,
			}
			tt.assertion(t, c.InitDefaults())
		})
	}
}

func TestInitDefaults_ErrorMessages(t *testing.T) {
	t.Run("missing file error message", func(t *testing.T) {
		c := Config{
			ProtoPath: []string{"./tests/proto/serviceapis"},
			Files:     []string{"missing/file.proto"},
		}

		err := c.InitDefaults()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "does not exist in any of the proto_paths")
	})

	t.Run("empty inputs error message", func(t *testing.T) {
		c := Config{
			ProtoPath: []string{},
			Files:     []string{},
		}

		err := c.InitDefaults()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "proto_path and files must be specified")
	})
}
