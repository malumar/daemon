package stuff

import (
	"context"
	"fmt"
	"github.com/malumar/daemon"
	"github.com/malumar/daemon/plugins"
)

func DoTheReallyImportantThings() {
	fmt.Println("[stuff] Hello World")
}

func init() {
	daemon.RegisterListener(plugins.TRIGGER_BEFORE_APP_INIT, func(ctx context.Context, event string, value interface{}) error {
		fmt.Println("[stuff] if you see this, it means that the TRIGGER_BEFORE_APP_INIT trigger was called")
		return nil
	})

	daemon.RegisterListener(plugins.TRIGGER_APPLICATION_EXIT, func(ctx context.Context, event string, value interface{}) error {
		fmt.Println("[stuff] i received plugins.TRIGGER_APPLICATION_EXIT")
		return nil
	})
}
