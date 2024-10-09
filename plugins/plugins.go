package plugins

const (
	TRIGGER_BEFORE_TEMPLATES_INIT      = "SYS:BEFORE_TEMPLATES_INIT"
	TRIGGER_TEMPLATE_REFRESH_REQUESTED = "SYS:TEMPLATE_REFRESH_REQUESTED"
	TRIGGER_TEMPLATE_REFRESH_COMPLETED = "SYS:TEMPLATE_REFRESH_COMPLETED"

	// called before the arguments are inserted into the template
	// can be used to set default values
	TRIGGER_BEFORE_TEMPLATE_ARGS_SET = "SYS:BEFORE_TEMPLATE_ARGS_SET"

	// called after inserting arguments into the template just before rendering
	// you can no longer assign values to templates here
	TRIGGER_BEFORE_RENDER_TEMPLATE = "SYS:_BEFORE_RENDER_TEMPLATE"

	// When we force all daemons to stop
	// we need to end all functions in the background, after this command TRIGGER_APPLICATION_EXIT will come
	TRIGGER_FINISH_APP = "SYS:_APPLICATION_FINISH_APP"

	// When we end the application
	// At this point we should not send any more signals and all go routines should be completed
	TRIGGER_APPLICATION_EXIT = "SYS:_APPLICATION_EXIT"

	// It's time to register your flags
	TRIGGER_BEFORE_FLAGS_PREPARED = "SYS:_BEFORE_FLAGS_PREPARED"

	// Call when all flags are already registered
	TRIGGER_AFTER_FLAGS_PREPARED = "SYS:_AFTER_FLAGS_PREPARED"

	// This is the time to register your commands, it is executed after BEFORE_FLAGS
	TRIGGER_BEFORE_COMMANDS_REGISTER = "SYS:_BEFORE_COMMANDS_REGISTER"

	// This is just before the commands are invoked
	TRIGGER_AFTER_COMMANDS_REGISTER = "SYS:_AFTER_COMMANDS_REGISTER"

	// Called after executing CLI commands
	TRIGGER_AFTER_COMMANDS_EXECUTED = "SYS:_AFTER_COMMANDS_EXECUTED"

	// Before starting the application initialization
	TRIGGER_BEFORE_APP_INIT = "SYS:_BEFORE_APP_INIT"

	// Just after basic initialization is completed, but just before it is completed
	// and OnAppInit() has not yet been called to initialize modules
	TRIGGER_BEFORE_MODULES_INIT = "SYS:_BEFORE_MODULES_HANDLER_INIT"

	// Once the application has completed initialization
	TRIGGER_AFTER_APP_INIT = "SYS:_AFTER_APP_INIT"

	// Before starting the application
	TRIGGER_BEFORE_APP_RUN = "SYS:_BEFORE_APP_RUN"

	// After starting the application
	TRIGGER_AFTER_APP_RUN = "SYS:_AFTER_APP_RUN"

	// When the configuration is required to be reloaded by modules
	// (maybe you need to load them from the database) that store them in the global cache system,
	// it means that the cache has been cleared
	TRIGGER_RELOAD_CACHED_CONFIG = "SYS:_RELOAD_CACHED_CONFIG"

	TRIGGER_RELOAD_APP = "SYS:_APPLICATION_RELOAD"
)
