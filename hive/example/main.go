// SPDX-License-Identifier: Apache-2.0
// Copyright Authors of Cilium

package main

import (
	"github.com/luomu/gofrw/logging"
	"github.com/luomu/gofrw/logging/hooks"
	"log"

	"github.com/spf13/cobra"

	"github.com/luomu/gofrw/hive"
	"github.com/luomu/gofrw/hive/cell"
	"github.com/luomu/gofrw/metrics"
)

var (
	// Create a hive from a set of cells.
	Hive = hive.New(
		metrics.Cell,
		serverCell,         // An HTTP server, depends on HTTPHandler's
		eventsCell,         // Example event source (ExampleEvents)
		helloHandlerCell,   // Handler for /hello
		eventsHandlerCell,  // Handler for /events
		exampleMetricsCell, // Register the metrics
		jobsCell,           // Examples of jobs

		// Constructors are lazy and only invoked if they are a dependency
		// to an "invoke" function or an indirect dependency of a constructor
		// referenced in an invoke. This allows composing "bundles" of modules
		// and then only paying for what's actually used from the bundle.
		//
		// Think of invoke functions as the driver that decides what things
		// should be constructed and how they should integrate with each other.
		//
		// Modules that provide a service to others should usually not have any invoke
		// functions that force object construction whether or not it is needed.
		//
		// In this example we have the server at the top of the dependency tree,
		// so we'll just depend on it here to make sure it gets instantiated.
		cell.Invoke(func(Server) {}),
	)

	// Define a cobra command that runs the hive.
	cmd = &cobra.Command{
		Use: "example",
		Run: func(_ *cobra.Command, args []string) {
			// When we get here, cobra has parsed all the command-line flags and hive
			// can be started.
			// This first populates all configurations from Viper (and via pflag)
			// and then constructs all objects, followed by executing the start
			// hooks in dependency order. It will then block waiting for signals
			// after which it will run the stop hooks in reverse order.
			if err := Hive.Run(); err != nil {
				// Run() can fail if:
				// - There are missing types in the object graph
				// - Executing the lifecycle start or stop hooks fails
				// - Shutdowner.Shutdown() is called with an error
				log.Fatal(err)
			}
		},
	}
)

func main() {
	// Register all configuration flags in the hive to the command
	Hive.RegisterFlags(cmd.Flags())

	initLogging("hive-info.log", "custom", false)
	// Add the "hive" sub-command for inspecting the hive
	cmd.AddCommand(Hive.Command())

	// And finally execute the command to parse the command-line flags and
	// run the hive
	cmd.Execute()
}

func initLogging(logFile, logFormat string, debug bool) error {
	logging.DefaultLogger.Hooks.Add(metrics.NewLoggingHook())
	f := logFormat
	if f == "" {
		f = string(logging.DefaultLogFormat)
	}
	logOptions := logging.LogOptions{
		logging.FormatOpt: f,
	}
	err := logging.SetupLogging([]string{}, logOptions, "hive", debug)
	if err != nil {
		return err
	}

	if len(logFile) != 0 {
		logging.AddHooks(
			hooks.NewFileRotationLogHook(logFile,
				hooks.EnableCompression(),
				hooks.EnableLocalTime(),
				hooks.WithMaxSize(100),
				hooks.WithMaxAge(5),
				hooks.WithMaxBackups(7),
			))
	}

	return nil
}
