/*
Copyright © 2023 The Devspace Authors
This file is a part of CLI application for Devspace.
*/

package cmd

import (
	"github.com/spf13/cobra"
)

// tagsCmd represents the tags command
var tagsCmd = &cobra.Command{
	Use:   "tags",
	Short: "Interact with devspace tags.",
	Long: `Interact with devspace tags.

Given a particular devspace, perform actions on its tags like
adding a new tag, list all the tags or see a particular tag`,
	TraverseChildren: true,
}

func init() {
	rootCmd.AddCommand(tagsCmd)

	tagsCmd.PersistentFlags().StringP("devspace", "d", "", "The devspace whose tags are acted upon")
	tagsCmd.PersistentFlags().StringP("token", "k", "", "login token")
	_ = tagsCmd.MarkPersistentFlagRequired("devspace")
	_ = tagsCmd.MarkPersistentFlagRequired("token")
}
