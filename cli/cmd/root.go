/*
Copyright © 2023 The Devspace Authors
This file is a part of CLI application for Devspace.
*/

package cmd

import (
	"errors"
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "dev",
	Short: "Searchable encryption based on PEKS",
	Long: `Searchable encryption based on PEKS

Dev is a CLI tool for PEKS utilities for interaction with
Devspace. Dev provides useful methods for key generation,
encryption of data, create shared keys, ciphertext for
given keyword and trapdoor generation.`,

	RunE: func(cmd *cobra.Command, args []string) error {
		return errors.New("'dev' has no default action")
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
