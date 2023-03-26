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

// pubkeyCmd represents the pubkey command
var pubkeyCmd = &cobra.Command{
	Use:   "pubkey",
	Short: "Generate Public Key from a given Private key",
	Long: `Generate Public Key from a given Private key

Take user private key and output the corresponding public key`,
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		// flags
		skFile, err := cmd.Flags().GetString("skey")
		if err != nil {
			return err
		}
		skHex, err := cmd.Flags().GetString("skey-hex")
		if err != nil {
			return err
		}
		pkFile, err := cmd.Flags().GetString("pkey")
		if err != nil {
			return err
		}

		// input
		sk := new(core.SKey)
		err = keyio.ReadKey(sk, skFile, skHex, true)
		if err != nil {
			return err
		}

		// core logic
		pk := new(core.PKey)
		err = pk.FromSKey(sk)
		if err != nil {
			return err
		}

		// output
		err = keyio.WriteFile(pk.Bytes(), pkFile, true)
		if err != nil {
			return err
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(pubkeyCmd)

	pubkeyCmd.Flags().StringP("pkey", "p", "", "file to output public key")
	pubkeyCmd.Flags().StringP("skey", "s", "", "private key file")
	_ = pubkeyCmd.MarkFlagFilename("skey") // error happens only when flag does not exist
	pubkeyCmd.Flags().String("skey-hex", "", "hexadecimal user private key")
	pubkeyCmd.MarkFlagsMutuallyExclusive("skey", "skey-hex")
}
