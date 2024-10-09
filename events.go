package daemon

import (
	"context"
	"errors"
	"log"
	"strings"
)

type EventHandler func(ctx context.Context, event string, value interface{}) error

func NewCriticalError(message string) error {
	return CriticalError{Message: message}
}

type CriticalError struct {
	Message string
}

func (t CriticalError) Error() string {
	return t.Message
}

func IsCriticalError(err error) bool {
	return errors.Is(err, CriticalError{})
}

var registeredListeners map[string][]EventHandler = make(map[string][]EventHandler, 0)

// RegisterListener Registering a listener ListenerHandler should not return an error, and if it does at all,
// we should not interrupt the application because of it! this should be done by the plugin itself if necessary
func RegisterListener(name string, listener EventHandler) {

	if strings.TrimSpace(name) == "" {
		log.Printf("warn: [daemon] attempt to register a Listener without a name\n")
		return
	}

	if listener == nil {
		log.Printf("warn: [daemon] attempt to register a Listener pointing to NIL\n")
		return
	}
	if registeredListeners[name] == nil {
		registeredListeners[name] = make([]EventHandler, 0)
	}
	registeredListeners[name] = append(registeredListeners[name], listener)
}

func FireTrigger(name string, value interface{}) error {
	return FireTriggerWithContext(name, nil, value)
}

func FireTriggerWithContext(name string, ctx context.Context, value interface{}) error {

	if registeredListeners == nil {
		return nil
	}

	if listeners, ok := registeredListeners[name]; ok {
		for _, listener := range listeners {
			if listener != nil {
				if err := listener(ctx, name, value); err != nil {
					log.Printf("err: [daemon] trigger %s -> %v\n", name, err)
					if IsCriticalError(err) {
						stopAppGracefully(true)
					}
					// return err
				}
			}
		}
	}
	return nil
}
