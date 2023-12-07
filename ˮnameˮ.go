package ˮnameˮ

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"go.k6.io/k6/output"
)

func init() {
	output.RegisterExtension("ˮnameˮ",
		func(p output.Params) (output.Output, error) { return newExtension(p) },
	)
}

type extension struct {
	output.SampleBuffer
	logger  logrus.FieldLogger
	flusher *output.PeriodicFlusher

	config      *config
	description string
	name        string
}

var _ output.Output = (*extension)(nil)

// newExtension creates an instance of the extension.
func newExtension(params output.Params) (*extension, error) {
	cfg, err := newConfig(params)
	if err != nil {
		return nil, err
	}

	return &extension{
		logger: params.Logger,
		config: cfg,
		name:   params.OutputType,
		description: fmt.Sprintf(
			"%s: interval=%s, address=%s", params.OutputType, cfg.interval.String(), cfg.address,
		),
	}, nil
}

func (e *extension) log() *logrus.Entry {
	return e.logger.WithField("extension", e.name)
}

// Description returns a human-readable description of the output that will be shown in `k6 run`.
func (e *extension) Description() string {
	return e.description
}

// Stop flushes all remaining metrics and finalizes the test run.
func (e *extension) Stop() error {
	e.log().Debug("Stopping...")

	e.flusher.Stop()

	e.log().Debug("Stopped!")

	return nil
}

// Start performs initialization tasks prior to Engine using the output.
func (e *extension) Start() error {
	e.log().Debug("Starting...")

	var err error

	e.flusher, err = output.NewPeriodicFlusher(e.config.interval, e.flush)
	if err != nil {
		return err
	}

	e.log().Debug("Started!")

	return nil
}

func (e *extension) flush() {
	containers := e.GetBufferedSamples()

	var count int

	for _, sc := range containers {
		samples := sc.GetSamples()
		count += len(samples)

		// Here we actually write or accumulate to then write in batches
	}

	if count > 0 {
		e.log().WithField("count", count).Debug("Metrics processed")
	}
}
