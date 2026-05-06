// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2025 kdsmith18542

package main

import (
	"fmt"
	"io"
	"log"
	"strings"

	"github.com/chzyer/readline"
	"github.com/kdsmith18542/pwny/internal/cli"
	"github.com/kdsmith18542/pwny/internal/payload"
	"github.com/spf13/cobra"
	"os"
)

type CLIState struct {
	CurrentModule string
	Options       map[string]interface{}
}

var state = &CLIState{
	Options: make(map[string]interface{}),
}

func main() {
	rl, err := readline.NewEx(&readline.Config{
		Prompt:          "\033[31mpwny\033[0m > ",
		HistoryFile:     "/tmp/pwny_history",
		AutoComplete:    cli.GetDefaultCompleter(),
		InterruptPrompt: "^C",
		EOFPrompt:       "exit",
	})
	if err != nil {
		log.Fatalf("failed to initialize readline: %v", err)
	}
	defer rl.Close()

	rootCmd := &cobra.Command{Use: "pwny"}
	
	// Add commands
	rootCmd.AddCommand(&cobra.Command{
		Use:   "use [module]",
		Short: "Select a module to use",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			state.CurrentModule = args[0]
			rl.SetPrompt(fmt.Sprintf("\033[31mpwny\033[0m %s > ", state.CurrentModule))
			fmt.Printf("Using module: %s\n", state.CurrentModule)
		},
	})

	rootCmd.AddCommand(&cobra.Command{
		Use:   "set [option] [value]",
		Short: "Set an option for the current module",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			state.Options[args[0]] = args[1]
			fmt.Printf("%s => %v\n", args[0], args[1])
		},
	})

	rootCmd.AddCommand(&cobra.Command{
		Use:   "show options",
		Short: "Show options for the current module",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Current Options:")
			for k, v := range state.Options {
				fmt.Printf("  %s: %v\n", k, v)
			}
		},
	})

	rootCmd.AddCommand(&cobra.Command{
		Use:   "run",
		Short: "Run the current module",
		Run: func(cmd *cobra.Command, args []string) {
			if state.CurrentModule == "" {
				fmt.Println("Error: No module selected. Use 'use [module]' first.")
				return
			}
			fmt.Printf("Running module: %s...\n", state.CurrentModule)
			// TODO: Call API /api/v1/modules/{name}/run
		},
	})

	sessionsCmd := &cobra.Command{
		Use:   "sessions",
		Short: "Manage active sessions",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Active Sessions:")
			// TODO: Call API /api/v1/sessions
		},
	}

	sessionsCmd.AddCommand(&cobra.Command{
		Use:   "interact [id]",
		Short: "Interact with a session",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			cli.InteractWithSession(args[0], rl)
		},
	})

	rootCmd.AddCommand(sessionsCmd)

	payloadCmd := &cobra.Command{
		Use:   "payload",
		Short: "Generate a payload (MSFVenom compatible flags)",
		Run: func(cmd *cobra.Command, args []string) {
			p, _ := cmd.Flags().GetString("payload")
			format, _ := cmd.Flags().GetString("format")
			output, _ := cmd.Flags().GetString("out")
			encoder, _ := cmd.Flags().GetString("encoder")
			iterations, _ := cmd.Flags().GetInt("iterations")

			opts := make(map[string]interface{})
			// Parse key=value arguments
			for _, arg := range args {
				parts := strings.SplitN(arg, "=", 2)
				if len(parts) == 2 {
					opts[strings.ToUpper(parts[0])] = parts[1]
				}
			}

			fmt.Printf("Generating payload %s in %s format...\n", p, format)
			
			// Call internal generator directly for CLI
			gen := payload.NewDefaultGenerator()
			payloadOpts := payload.PayloadOptions{
				Name:       p,
				Format:     format,
				Encoder:    encoder,
				Iterations: iterations,
				Options:    opts,
			}
			
			res, err := gen.Generate(payloadOpts)
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				return
			}

			if output != "" {
				err := os.WriteFile(output, []byte(res.Payload), 0644)
				if err != nil {
					fmt.Printf("Error writing to file: %v\n", err)
				} else {
					fmt.Printf("Payload saved to %s (%d bytes)\n", output, res.Size)
				}
			} else {
				fmt.Println(res.Payload)
			}
		},
	}

	payloadCmd.Flags().StringP("payload", "p", "", "Payload to use")
	payloadCmd.Flags().StringP("format", "f", "raw", "Output format")
	payloadCmd.Flags().StringP("out", "o", "", "Output file")
	payloadCmd.Flags().StringP("encoder", "e", "", "Encoder to use")
	payloadCmd.Flags().IntP("iterations", "i", 1, "Number of encoding iterations")

	rootCmd.AddCommand(payloadCmd)

	// If arguments provided, run command and exit (MSFVenom mode)
	if len(os.Args) > 1 {
		if err := rootCmd.Execute(); err != nil {
			os.Exit(1)
		}
		return
	}

	fmt.Println("Pwny CLI v0.1.0")

	for {
		line, err := rl.Readline()
		if err != nil {
			if err == readline.ErrInterrupt {
				if len(line) == 0 {
					break
				} else {
					continue
				}
			} else if err == io.EOF {
				break
			}
			break
		}

		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		args := strings.Fields(line)
		rootCmd.SetArgs(args)
		rootCmd.Execute()
	}
}
