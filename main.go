package main

import (
	"context"
	"ht16k33-display/models"

	"go.viam.com/rdk/components/generic"
	"go.viam.com/rdk/logging"
	"go.viam.com/rdk/module"
	"go.viam.com/utils"
)

func main() {
	utils.ContextualMain(mainWithArgs, module.NewLoggerFromArgs("ht16k33-display"))
}

func mainWithArgs(ctx context.Context, args []string, logger logging.Logger) error {
	ht16k33Display, err := module.NewModuleFromArgs(ctx)
	if err != nil {
		return err
	}

	if err = ht16k33Display.AddModelFromRegistry(ctx, generic.API, models.Seg14X4); err != nil {
		return err
	}

	err = ht16k33Display.Start(ctx)
	defer ht16k33Display.Close(ctx)
	if err != nil {
		return err
	}

	<-ctx.Done()
	return nil
}
