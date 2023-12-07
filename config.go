package ˮnameˮ

import (
	"errors"
	"fmt"
	"net/url"
	"time"

	"go.k6.io/k6/output"
)

const (
	defaultInterval = 2 * time.Second
	defaultAddress  = "127.0.0.1:7"
)

type config struct {
	address  string
	interval time.Duration
}

func newConfig(params output.Params) (*config, error) {
	cfg := &config{
		address:  defaultAddress,
		interval: defaultInterval,
	}

	if v, has := params.Environment[envAddress]; has {
		cfg.address = v
	}

	if v, has := params.Environment[envInterval]; has {
		var err error

		cfg.interval, err = time.ParseDuration(v)
		if err != nil {
			return nil, fmt.Errorf("%w: error parsing environment variable %s: %s", errConfig, envInterval, err.Error())
		}
	}

	if params.ConfigArgument == "" {
		return cfg, nil
	}

	values, err := url.ParseQuery(params.ConfigArgument)
	if err != nil {
		return nil, err
	}

	if values.Has(paramAddress) {
		cfg.address = values.Get(paramAddress)
	}

	if values.Has(paramInterval) {
		var err error

		cfg.interval, err = time.ParseDuration(values.Get(paramInterval))
		if err != nil {
			return nil, fmt.Errorf("%w: error parsing parameter %s: %s", errConfig, paramInterval, err.Error())
		}
	}

	return cfg, nil
}

var errConfig = errors.New("config error")

const (
	envPrefix = "ˮenvPrefixˮ_"

	envInterval = envPrefix + "INTERVAL"
	envAddress  = envPrefix + "ADDRESS"

	paramInterval = "interval"
	paramAddress  = "address"
)
