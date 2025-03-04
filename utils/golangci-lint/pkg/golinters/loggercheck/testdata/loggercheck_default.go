//golangcitest:args -Eloggercheck
package loggercheck

import (
	"fmt"
	"log/slog"

	"github.com/go-logr/logr"
	"go.uber.org/zap"
	"k8s.io/klog/v2"
)

func ExampleDefaultLogr() {
	log := logr.Discard()
	log = log.WithValues("key")                                         // want `odd number of arguments passed as key-value pairs for logging`
	log.Info("message", "key1", "value1", "key2", "value2", "key3")     // want `odd number of arguments passed as key-value pairs for logging`
	log.Error(fmt.Errorf("error"), "message", "key1", "value1", "key2") // want `odd number of arguments passed as key-value pairs for logging`
	log.Error(fmt.Errorf("error"), "message", "key1", "value1", "key2", "value2")
}

func ExampleDefaultKlog() {
	klog.InfoS("message", "key1") // want `odd number of arguments passed as key-value pairs for logging`
}

func ExampleZapSugarNotChecked() {
	sugar := zap.NewExample().Sugar()
	defer sugar.Sync()
	sugar.Infow("message", "key1", "value1", "key2") // want `odd number of arguments passed as key-value pairs for logging`
}

func ExampleSlog() {
	logger := slog.With("key1", "value1")

	logger.Info("msg", "key1") // want `odd number of arguments passed as key-value pairs for logging`
	slog.Info("msg", "key1")   // want `odd number of arguments passed as key-value pairs for logging`
}
