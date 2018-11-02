package util

import (
	"errors"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/joho/godotenv"
	"log"
	"os"
)

const (
	main          = "main"
	test          = "test"
	secretDirName = ".data"
	seedFileName  = "seed"
	dbFileName    = "wallet.db"
)

var (
	// SecretDirPath is secredi dir path
	SecretDirPath string
	// SeedFilePath is seed file for HD wallet
	SeedFilePath string
	// DbFilePath is db for extend private & public key
	DbFilePath string
	// IndexFromNet is index from btc network. It is declared in BIP0044
	IndexFromNet = map[string]uint32{
		main: 0,
		test: 1,
	}

	// ParamsFromNet is index from btc network. It is declared in BIP0044
	ParamsFromNet = map[string]*chaincfg.Params{
		main: &chaincfg.MainNetParams,
		test: &chaincfg.TestNet3Params,
	}
)

// GetNetParams get net params by string
func GetNetParams(net string) (*chaincfg.Params, error) {
	switch net {
	case main:
		return &chaincfg.MainNetParams, nil
	case test:
		return &chaincfg.TestNet3Params, nil
	}
	return &chaincfg.Params{}, errors.New("invalid network name")
}

// InitialzeParams initilaze parameters
func InitialzeParams() error {
	err := envLoad()
	if err != nil {
		return err
	}
	base := os.Getenv("BASE_DIR")
	SecretDirPath = base + "/" + secretDirName
	SeedFilePath = SecretDirPath + "/" + seedFileName
	DbFilePath = SecretDirPath + "/" + dbFileName
	return nil
}

// CreateSecretDir is create directory
func CreateSecretDir() error {
	msg := os.RemoveAll(SecretDirPath)
	if msg != nil {
		log.Println(msg)
	}
	msg = os.Mkdir(SecretDirPath, 0755)
	if msg != nil {
		log.Println(msg)
	}

	return nil
}

func envLoad() error {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
		return err
	}
	return nil
}
