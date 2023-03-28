/*
Copyright © 2023 The Devspace Authors
This file is a part of CLI application for Devspace.
*/

package cmd

import (
	"encoding/hex"
	"github.com/bingxueshuang/devspaces/cli/keyio"
	"github.com/bingxueshuang/devspaces/core"
	"github.com/spf13/cobra"
)

// sharedkeyCmd represents the sharedkey command
var sharedkeyCmd = &cobra.Command{
	Use:   "sharedkey",
	Short: "Generate shared key for two users",
	Long: `Generate shared key for two users.

Take public key of one user and private key of
other user and output shared key for the pair.
This operation is commutable, hence both users
can perform independently.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// flags
		oFlag, err := cmd.Flags().GetString("output")
		if err != nil {
			return err
		}
		skFlag, err := cmd.Flags().GetString("skey")
		if err != nil {
			return err
		}
		pkFlag, err := cmd.Flags().GetString("pkey")
		if err != nil {
			return err
		}
		skHex, err := cmd.Flags().GetString("skey-hex")
		if err != nil {
			return err
		}
		pkHex, err := cmd.Flags().GetString("pkey-hex")
		if err != nil {
			return err
		}

		// input
		sk := new(core.SKey)
		err = keyio.ReadKey(sk, skFlag, skHex, false)
		if err != nil {
			return err
		}
		pk := new(core.PKey)
		err = keyio.ReadKey(pk, pkFlag, pkHex, false)
		if err != nil {
			return err
		}

		// core logic
		key := core.SharedKey(pk, sk)

		// output
		err = keyio.WriteString(hex.EncodeToString(key.Bytes()), oFlag, true)
		if err != nil {
			return err
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(sharedkeyCmd)

	sharedkeyCmd.Flags().StringP("output", "o", "", "file to output shared key")
	sharedkeyCmd.Flags().StringP("skey", "s", "", "private key file")
	_ = sharedkeyCmd.MarkFlagFilename("skey")
	sharedkeyCmd.Flags().StringP("pkey", "p", "", "public key file")
	_ = sharedkeyCmd.MarkFlagFilename("pkey")
	sharedkeyCmd.Flags().String("skey-hex", "", "hexadecimal private key")
	sharedkeyCmd.Flags().String("pkey-hex", "", "hexadecimal public key")
	sharedkeyCmd.MarkFlagsMutuallyExclusive("skey", "skey-hex")
	sharedkeyCmd.MarkFlagsMutuallyExclusive("pkey", "pkey-hex")
}
