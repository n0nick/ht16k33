package models

import (
	"context"

	"go.viam.com/rdk/components/generic"
	"go.viam.com/rdk/logging"
	"go.viam.com/rdk/resource"
	"go.viam.com/utils/rpc"

	"periph.io/x/conn/v3/i2c"
	"periph.io/x/conn/v3/i2c/i2creg"
	"periph.io/x/devices/v3/ht16k33"
	"periph.io/x/host/v3"
)

var (
	Seg14X4 = resource.NewModel("n0nick", "ht16k33-display", "seg_14_x_4")
	// errUnimplemented = errors.New("unimplemented")
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

	bus     i2c.BusCloser
	display *ht16k33.Display

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

	s.initDisplay()
	return s, nil
}

func (s *ht16k33DisplaySeg14X4) Name() resource.Name {
	return s.name
}

func (s *ht16k33DisplaySeg14X4) Reconfigure(ctx context.Context, deps resource.Dependencies, conf resource.Config) error {
	if s.display == nil {
		s.initDisplay()
	}

	return nil
}

func (s *ht16k33DisplaySeg14X4) initDisplay() error {
	if _, err := host.Init(); err != nil {
		return err
	}

	bus, err := i2creg.Open("")
	if err != nil {
		return err
	}
	s.bus = bus

	display, err := ht16k33.NewAlphaNumericDisplay(bus, 0x70) // TODO from config
	if err != nil {
		return err
	}

	s.display = display

	return display.Halt()
}

func (s *ht16k33DisplaySeg14X4) NewClientFromConn(ctx context.Context, conn rpc.ClientConn, remoteName string, name resource.Name, logger logging.Logger) (resource.Resource, error) {
	panic("not implemented")
}

func (s *ht16k33DisplaySeg14X4) DoCommand(ctx context.Context, cmd map[string]interface{}) (map[string]interface{}, error) {
	if s.display == nil {
		return map[string]interface{}{"error": "uninitialized"}, nil
	}

	for key, value := range cmd {
		switch key {
		case "print":
			s.display.WriteString(value.(string))
			callOnSubstrings(value.(string), 4, func(st string) { s.display.WriteString(st) })
		default:
			return map[string]interface{}{"error": "unknown command"}, nil
		}

		break // support 1 command at a time
	}

	return cmd, nil
}

func (s *ht16k33DisplaySeg14X4) Close(context.Context) error {
	if s.bus != nil {
		s.bus.Close()
	}

	s.cancelFunc()
	return nil
}

// callOnSubstrings calls the function `f` on each n-length substring of `st`, rotating one character at a time.
func callOnSubstrings(st string, n int, f func(string)) {
	// Iterate through the string one character at a time
	for i := 0; i < len(st); i++ {
		// Get the substring of length `n` starting from index `i`
		end := i + n
		if end > len(st) {
			break // Stop if the substring exceeds the string length
		}
		f(st[i:end])
	}
}
