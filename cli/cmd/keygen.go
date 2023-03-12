/*
Copyright © 2023 The Devspace Authors
This file is a part of CLI application for Devspace.
*/

package cmd

import (
	"encoding/hex"
	"fmt"
	"github.com/bingxueshuang/devspaces/core"
	"github.com/spf13/cobra"
	"os"
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
		sk, pk, err := core.KeyGen()
		skHex := hex.EncodeToString(sk.Bytes())
		pkHex := hex.EncodeToString(pk.Bytes())
		if skFlag == "" {
			fmt.Print(skHex)
		} else {
			err := os.WriteFile(skFlag, []byte(skHex), 0666)
			if err != nil {
				return err
			}
		}
		if pkFlag != "" {
			err := os.WriteFile(pkFlag, []byte(pkHex), 0666)
			if err != nil {
				return err
			}
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(keygenCmd)

	keygenCmd.Flags().StringP("skey", "s", "", "file to output secret key")
	keygenCmd.Flags().StringP("pkey", "p", "", "file to output public key")
}
