package models

import (
	"context"
	"errors"

	"go.viam.com/rdk/components/generic"
	"go.viam.com/rdk/logging"
	"go.viam.com/rdk/resource"
	"go.viam.com/utils/rpc"
)

var (
	Seg14X4          = resource.NewModel("", "ht16k33-display", "seg_14_x_4")
	errUnimplemented = errors.New("unimplemented")
)

func init() {
	resource.RegisterComponent(generic.API, Seg14X4,
		resource.Registration[resource.Resource, *Config]{
			Constructor: newHt16k33DisplaySeg14X4,
		},
	)
}

type Config struct {
	// Put config attributes here

	/* if your model  does not need a config,
	   replace *Config in the init function with resource.NoNativeConfig */

	/* Uncomment this if your model does not need to be validated
	   and has no implicit dependecies. */
	// resource.TriviallyValidateConfig

}

func (cfg *Config) Validate(path string) ([]string, error) {
	// Add config validation code here
	return nil, nil
}

type ht16k33DisplaySeg14X4 struct {
	name resource.Name

	logger logging.Logger
	cfg    *Config

	cancelCtx  context.Context
	cancelFunc func()

	/* Uncomment this if your model does not need to reconfigure. */
	// resource.TriviallyReconfigurable

	// Uncomment this if the model does not have any goroutines that
	// need to be shut down while closing.
	// resource.TriviallyCloseable

}

func newHt16k33DisplaySeg14X4(ctx context.Context, deps resource.Dependencies, rawConf resource.Config, logger logging.Logger) (resource.Resource, error) {
	conf, err := resource.NativeConfig[*Config](rawConf)
	if err != nil {
		return nil, err
	}

	cancelCtx, cancelFunc := context.WithCancel(context.Background())

	s := &ht16k33DisplaySeg14X4{
		name:       rawConf.ResourceName(),
		logger:     logger,
		cfg:        conf,
		cancelCtx:  cancelCtx,
		cancelFunc: cancelFunc,
	}
	return s, nil
}

func (s *ht16k33DisplaySeg14X4) Name() resource.Name {
	return s.name
}

func (s *ht16k33DisplaySeg14X4) Reconfigure(ctx context.Context, deps resource.Dependencies, conf resource.Config) error {
	// Put reconfigure code here
	return errUnimplemented
}

func (s *ht16k33DisplaySeg14X4) NewClientFromConn(ctx context.Context, conn rpc.ClientConn, remoteName string, name resource.Name, logger logging.Logger) (resource.Resource, error) {
	panic("not implemented")
}

func (s *ht16k33DisplaySeg14X4) DoCommand(ctx context.Context, cmd map[string]interface{}) (map[string]interface{}, error) {
	panic("not implemented")
}

func (s *ht16k33DisplaySeg14X4) Close(context.Context) error {
	// Put close code here
	s.cancelFunc()
	return nil
}
