/*
Copyright © 2023 The Devspace Authors
This file is a part of CLI application for Devspace.
*/

package cmd

import (
	"github.com/bingxueshuang/devspaces/cli/keyio"
	"github.com/bingxueshuang/devspaces/core"
	"github.com/spf13/cobra"
)

// keygenCmd represents the keygen command
var keygenCmd = &cobra.Command{
	Use:   "keygen",
	Short: "Generate public-secret key pair",
	Long: `Generate public-secret key pair.

Create a new user public key and private key pair
and output the private key to the user.
`,
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		// flags
		skFlag, err := cmd.Flags().GetString("skey")
		if err != nil {
			return err
		}
		pkFlag, err := cmd.Flags().GetString("pkey")
		if err != nil {
			return err
		}

		// core logic
		sk, pk, err := core.KeyGen()
		if err != nil {
			return err
		}

		// output
		err = keyio.WriteKey(sk, skFlag, true)
		if err != nil {
			return err
		}
		err = keyio.WriteKey(pk, pkFlag, false)
		if err != nil {
			return err
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(keygenCmd)

	keygenCmd.Flags().StringP("skey", "s", "", "file to output secret key")
	keygenCmd.Flags().StringP("pkey", "p", "", "file to output public key")
}
