//////////////////////////////////////////////////////////////////////
//
// Given is a mock process which runs indefinitely and blocks the
// program. Right now the only way to stop the program is to send a
// SIGINT (Ctrl-C). Killing a process like that is not graceful, so we
// want to try to gracefully stop the process first.
//
// Change the program to do the following:
//   1. On SIGINT try to gracefully stop the process using
//          `proc.Stop()`
//   2. If SIGINT is called again, just kill the program (last resort)
//

package main

import (
	"fmt"
	"os"
	"os/signal"
)

func main() {

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)

	// Create a process
	proc := MockProcess{}

	// fire goroutine to watch for sigint
	go func() {

		// get the first sigint
		sig := <-sigChan
		fmt.Printf("\nreceived %q", sig)

		// start the clean up process for
		// graceful shutdown
		go proc.Stop()

		// wait for second interrupt
		// to kill the program
		<-sigChan
		os.Exit(0)
	}()

	// Run the process (blocking)
	proc.Run()
}
