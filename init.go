package daemon

import (
	"github.com/malumar/daemon/plugins"
)

func init() {
	FireTrigger(plugins.TRIGGER_BEFORE_APP_RUN, nil)
}
