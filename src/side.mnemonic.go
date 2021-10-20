package main

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/crypto/hd"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/sonntuet1997/medical-chain-utils/common"
	"github.com/tendermint/spm/cosmoscmd"
	"github.com/urfave/cli/v2"
)

func genMnemonic(appCtx *cli.Context) error {
	logger = common.InitLogger(appCtx)

	cosmoscmd.SetPrefixes("medichain")
	unsafeKeyring := keyring.NewUnsafe(keyring.NewInMemory())

	_, mnemonic, err := unsafeKeyring.NewMnemonic("admin", keyring.English, "m/44'/118'/0'/0", hd.Secp256k1)
	if err != nil {
		logger.Errorf("error while generating new mnemonic: %v", err)
		return err
	}

	fmt.Printf("New mnemonic generated: %s\n", mnemonic)
	return nil
}