package launcher

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/samuelngs/axis/models"
)

// Start - start program with given scope and entrypoint option
func Start(close chan struct{}, scope *models.Scope, opts *models.ApplicationEntryPoint) {
	defer func() {
		close <- struct{}{}
	}()

	// init commands array
	commands := []string{}

	// build execute command
	fmt.Printf("# %v\n  │\n", opts.EntryPoint)
	l := len(opts.Command)
	for i, v := range opts.Command {
		if arg := scope.Compile(v); strings.TrimSpace(arg) != "" {
			commands = append(commands, arg)
			if i != (l - 1) {
				fmt.Printf("  ├ %v\n", arg)
			} else {
				fmt.Printf("  └ %v\n", arg)
			}
		}
	}
	fmt.Printf("# [%v]\n", strings.Join(commands, " "))

	// start program with entrypoint and commands
	c := make(chan error)
	go func() {
		cmd := exec.Command(opts.EntryPoint, commands...)
		defer func() {
			cmd.Process.Kill()
		}()
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Start(); err != nil {
			log.Fatal(err)
		}
		c <- cmd.Wait()
	}()

	// print out error logs if there is any
	if err := <-c; err != nil {
		log.Fatal(err)
	}
}