package main

import (
	"fmt"
)

// CheckErrorListener looks for any errors and if there are some it will pass them to the result channel.
func CheckErrorListener(err error, result chan ListenerResult) {
	if CheckErrorInternal(err) {
		result <- ListenerResult{Result: Result{err: err}}
	}
}

// CheckError looks for any errors and if there are some it will pass them to the result channel.
func CheckError(err error, result chan Result) {
	if CheckErrorInternal(err) {
		result <- Result{err: err}
	}
}

// CheckErrorInternal looks for any errors and if there are some return true.
func CheckErrorInternal(err error) bool {
	return err != nil
}

// OutputUserMessage shows the result of the channel to the user console.
func OutputUserMessage(result Result) {
	if result.err == nil {
		fmt.Println(result.res)
	}
}

// CheckListenerChannelIsClose return true if the given channel is closed.
func CheckListenerChannelIsClose(channel chan ListenerResult) bool {
	select {
	case <-channel:
		// channel was closed
		return true
	default:
		return false
	}
}
