package protoreg_test

import (
	"log/slog"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"testing"
	"time"

	"github.com/roadrunner-server/config/v5"
	"github.com/roadrunner-server/endure/v2"
	"github.com/roadrunner-server/logger/v5"
	"github.com/roadrunner-server/protoreg/v5"
	"github.com/stretchr/testify/assert"
)

func TestProtoregInit(t *testing.T) {
	cont := endure.New(slog.LevelDebug)

	cfg := &config.Plugin{
		Version: "2023.3.0",
		Path:    "configs/.rr-protoreg-init.yaml",
	}

	plugin := &protoreg.Plugin{}

	err := cont.RegisterAll(
		cfg,
		&logger.Plugin{},
		plugin,
	)
	assert.NoError(t, err)

	err = cont.Init()
	if err != nil {
		t.Fatal(err)
	}

	ch, err := cont.Serve()
	assert.NoError(t, err)

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	wg := &sync.WaitGroup{}
	wg.Add(1)

	stopCh := make(chan struct{}, 1)

	go func() {
		defer wg.Done()
		for {
			select {
			case e := <-ch:
				assert.Fail(t, "error", e.Error.Error())
				err = cont.Stop()
				if err != nil {
					assert.FailNow(t, "error", err.Error())
				}
			case <-sig:
				err = cont.Stop()
				if err != nil {
					assert.FailNow(t, "error", err.Error())
				}
				return
			case <-stopCh:
				err = cont.Stop()
				if err != nil {
					assert.FailNow(t, "error", err.Error())
				}
				return
			}
		}
	}()

	time.Sleep(time.Second)

	registry := plugin.ProtoRegistry()
	assert.NotNil(t, registry)

	_, err = registry.FindMethodByFullPath("service.v1.Test/Echo")
	assert.NoError(t, err)

	unknown, err := registry.FindMethodByFullPath("service.v1.Test/Unknown")
	assert.Nil(t, unknown)

	stopCh <- struct{}{}

	wg.Wait()
}

func TestProtoregInitDuplicate(t *testing.T) {
	cont := endure.New(slog.LevelDebug)

	cfg := &config.Plugin{
		Version: "2023.3.0",
		Path:    "configs/.rr-protoreg-init-duplicate.yaml",
	}

	plugin := &protoreg.Plugin{}

	err := cont.RegisterAll(
		cfg,
		&logger.Plugin{},
		plugin,
	)
	assert.NoError(t, err)

	err = cont.Init()
	assert.Error(t, err)

	registry := plugin.ProtoRegistry()
	assert.Nil(t, registry)
}
