package protoreg_test

import (
	"log/slog"
	"testing"
	mocklogger "tests/mock"

	"github.com/roadrunner-server/config/v5"
	"github.com/roadrunner-server/endure/v2"
	"github.com/roadrunner-server/protoreg/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"google.golang.org/genproto/googleapis/api/annotations"
	"google.golang.org/protobuf/proto"
)

func TestProtoregInit(t *testing.T) {
	cont := endure.New(slog.LevelDebug)

	cfg := &config.Plugin{
		Version: "2023.3.0",
		Path:    "configs/.rr-protoreg-init.yaml",
	}

	l, oLogger := mocklogger.ZapTestLogger(zap.DebugLevel)

	plugin := &protoreg.Plugin{}

	err := cont.RegisterAll(
		cfg,
		l,
		plugin,
	)
	assert.NoError(t, err)

	err = cont.Init()
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, 1, oLogger.FilterMessageSnippet("protoreg initialized").Len())

	registry := plugin.ProtoRegistry()
	assert.NotNil(t, registry)

	_, err = registry.FindMethodByFullPath("service.v1.Test/Echo")
	assert.NoError(t, err)

	unknown, err := registry.FindMethodByFullPath("service.v1.Test/Unknown")
	assert.Nil(t, unknown)
}

func TestProtoregInitExtensions(t *testing.T) {
	cont := endure.New(slog.LevelDebug)

	cfg := &config.Plugin{
		Version: "2023.3.0",
		Path:    "configs/.rr-protoreg-init-extensions.yaml",
	}

	l, oLogger := mocklogger.ZapTestLogger(zap.DebugLevel)

	plugin := &protoreg.Plugin{}

	err := cont.RegisterAll(
		cfg,
		l,
		plugin,
	)
	assert.NoError(t, err)

	err = cont.Init()
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, 1, oLogger.FilterMessageSnippet("protoreg initialized").Len())

	registry := plugin.ProtoRegistry()
	assert.NotNil(t, registry)

	method, err := registry.FindMethodByFullPath("endpoint.v1.DocumentService/GetDocument")
	assert.NoError(t, err)
	assert.NotNil(t, method)

	opts := method.UnwrapMethod().Options()
	rule, ok := proto.GetExtension(opts, annotations.E_Http).(*annotations.HttpRule)
	assert.True(t, ok)
	assert.NotNil(t, rule)
}

func TestProtoregInitDuplicate(t *testing.T) {
	cont := endure.New(slog.LevelDebug)

	cfg := &config.Plugin{
		Version: "2023.3.0",
		Path:    "configs/.rr-protoreg-init-duplicate.yaml",
	}

	l, oLogger := mocklogger.ZapTestLogger(zap.DebugLevel)

	plugin := &protoreg.Plugin{}

	err := cont.RegisterAll(
		cfg,
		l,
		plugin,
	)
	assert.NoError(t, err)

	err = cont.Init()
	assert.Error(t, err)

	require.Equal(t, 0, len(oLogger.All()))

	registry := plugin.ProtoRegistry()
	assert.Nil(t, registry)
}

func TestProtoregInitGrpcDisabled(t *testing.T) {
	cont := endure.New(slog.LevelDebug)

	cfg := &config.Plugin{
		Version: "2023.3.0",
		Path:    "configs/.rr-protoreg-init-grpc-disabled.yaml",
	}

	l, oLogger := mocklogger.ZapTestLogger(zap.DebugLevel)

	plugin := &protoreg.Plugin{}

	err := cont.RegisterAll(
		cfg,
		l,
		plugin,
	)
	assert.NoError(t, err)

	err = cont.Init()
	if err != nil {
		t.Fatal(err)
	}

	// Init was skipped
	require.Equal(t, 0, len(oLogger.All()))

	registry := plugin.ProtoRegistry()
	assert.Nil(t, registry)
}
