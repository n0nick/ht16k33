package models

import (
	"context"
	"strconv"
	"time"

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
)

const (
	wait = 250 * time.Millisecond // default wait between animation frames
)

func init() {
	resource.RegisterComponent(generic.API, Seg14X4,
		resource.Registration[resource.Resource, *Config]{
			Constructor: newHt16k33DisplaySeg14X4,
		},
	)
}

type Config struct {
	// I2C address to use for connecting to HT16K33 display.
	Address string `json:"address,omitempty"`
	// Whether to skip intro animation
	SkipIntro bool `json:"skip_intro,omitempty"`
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

	err = s.initDisplay()
	if err != nil {
		return nil, err
	}

	// a little intro animation
	if !s.cfg.SkipIntro {
		s.animateIntro()
	}

	return s, nil
}

func (s *ht16k33DisplaySeg14X4) Name() resource.Name {
	return s.name
}

func (s *ht16k33DisplaySeg14X4) Reconfigure(ctx context.Context, deps resource.Dependencies, conf resource.Config) error {
	cfg, err := resource.NativeConfig[*Config](conf)
	if err != nil {
		return err
	}

	if cfg.Address != s.cfg.Address {
		if s.display != nil {
			s.display.Halt()
		}

		err = s.initDisplay()
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *ht16k33DisplaySeg14X4) initDisplay() error {
	s.logger.Info("Initializing HT16K33 display.")

	if _, err := host.Init(); err != nil {
		return err
	}

	bus, err := i2creg.Open("")
	if err != nil {
		return err
	}
	s.bus = bus

	address := uint16(0x70)
	if s.cfg.Address != "" {
		addressParsed, err := strconv.ParseUint(s.cfg.Address, 0, 16)
		if err != nil {
			return err
		}
		address = uint16(addressParsed)
		s.logger.Debugf("Will use custom address %s", address)
	}

	display, err := ht16k33.NewAlphaNumericDisplay(bus, address)
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
	s.logger.Debugf("Received command: %v", cmd)

	if s.display == nil {
		return map[string]interface{}{"error": "uninitialized"}, nil
	}

	for key, value := range cmd {
		switch key {
		case "print":
			callOnSubstrings(value.(string), 4, func(st string) {
				s.display.WriteString(st)
				time.Sleep(wait)
			})
		case "clear":
			s.display.Halt()
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

func (s *ht16k33DisplaySeg14X4) animateIntro() {
	for _, f := range []string{"   *", "  * ", " *  ", "*   "} {
		s.display.WriteString(f)
		time.Sleep(wait)
	}
	s.display.Halt()
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
