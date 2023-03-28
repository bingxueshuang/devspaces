/*
Copyright © 2023 The Devspace Authors
This file is a part of CLI application for Devspace.
*/

package cmd

import (
	"github.com/spf13/cobra"
)

// spaceCmd represents the space command
var spaceCmd = &cobra.Command{
	Use:   "space",
	Short: "Perform actions on a devspace",
	Long: `Perform actions on a devspace.

Create a devspace, list the tags in a devspace, send a devspace
collaboration request or send a message on devspace.`,
}

func init() {
	rootCmd.AddCommand(spaceCmd)

	spaceCmd.PersistentFlags().StringP("token", "k", "", "login token")
	_ = spaceCmd.MarkPersistentFlagRequired("token")
}
