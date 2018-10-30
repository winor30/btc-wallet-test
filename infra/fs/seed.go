package fs

import (
	"github.com/winor30/btc-wallet-test/util"
	"io/ioutil"
)

// SaveSeed save ramdom seed for HD wallet
func SaveSeed(seed []byte) error {
	err := util.CreateSecretDir()
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(util.SeedFilePath, seed, 0666)
	if err != nil {
		return err
	}

	return nil
}
