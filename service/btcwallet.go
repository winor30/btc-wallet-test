package service

import (
	"github.com/btcsuite/btcutil/hdkeychain"
	"github.com/winor30/btc-wallet-test/entity"
	"github.com/winor30/btc-wallet-test/infra/db"
	"github.com/winor30/btc-wallet-test/infra/fs"
	"github.com/winor30/btc-wallet-test/util"
	"math"

	"log"
)

// Setup set up params
func Setup() {
	util.InitialzeParams()
}

// Initialize create 3 extend private key in 1, 2 3 layer
// Also we can choise main/test net
func Initialize(net string) {
	netParams, err := util.GetNetParams(net)

	if err != nil {
		log.Fatalln(err)
		return
	}

	seed, err := hdkeychain.GenerateSeed(hdkeychain.RecommendedSeedLen)
	if err != nil {
		log.Fatalln(err)
		return
	}

	mPriv1, err := hdkeychain.NewMaster(seed, netParams)
	if err != nil {
		log.Fatalln(err)
		return
	}

	// layer 2: 44'
	mPriv2, err := mPriv1.Child(hdkeychain.HardenedKeyStart + 44)
	if err != nil {
		log.Fatalln(err)
		return
	}

	// layer 3: main(0') or test(1')
	mPriv3, err := mPriv2.Child(hdkeychain.HardenedKeyStart + util.IndexFromNet[net])
	if err != nil {
		log.Fatalln(err)
		return
	}

	// save seed
	err = fs.SaveSeed(seed)
	if err != nil {
		log.Fatalln(err)
		return
	}

	key1 := createKey(net, 0, []uint32{0}, mPriv1)
	key2 := createKey(net, 0, []uint32{0, hdkeychain.HardenedKeyStart + 44}, mPriv2)
	key3 := createKey(net, 0, []uint32{0, hdkeychain.HardenedKeyStart + 44, hdkeychain.HardenedKeyStart + util.IndexFromNet[net]}, mPriv3)

	err = db.CreateDb()
	if err != nil {
		log.Fatalln(err)
		return
	}

	keys := []*entity.Key{&key1, &key2, &key3}
	err = db.SaveKeys(keys)
	if err != nil {
		log.Fatalln(err)
		return
	}
}

// AddAccount is account for btcwallet
func AddAccount(name string, net string) {
	index4 := db.GetNewestIndexForAccount()

	err := db.SaveAccount(entity.Account{
		Name:   name,
		Index4: index4,
	})

	key, err := db.GetKey(
		0,
		0,
		[]uint32{
			0,
			hdkeychain.HardenedKeyStart + 44,
			hdkeychain.HardenedKeyStart + util.IndexFromNet[net],
			math.MaxUint32,
			math.MaxUint32,
			math.MaxUint32,
		},
	)

	exKey, err := hdkeychain.NewKeyFromString(key.Value)

	accountExKey, err := exKey.Child(hdkeychain.HardenedKeyStart + index4)
	accountKey := createKey(net, 0, []uint32{0, hdkeychain.HardenedKeyStart + 44, hdkeychain.HardenedKeyStart + util.IndexFromNet[net], hdkeychain.HardenedKeyStart + index4}, accountExKey)
	err = db.SaveKey(&accountKey)

	pubAccountExKey, err := accountExKey.Neuter()
	if err != nil {
		log.Fatalln(err)
		return
	}
	pubAccountKey := createKey(net, 1, []uint32{0, hdkeychain.HardenedKeyStart + 44, hdkeychain.HardenedKeyStart + util.IndexFromNet[net], hdkeychain.HardenedKeyStart + index4}, pubAccountExKey)
	err = db.SaveKey(&pubAccountKey)
	if err != nil {
		log.Fatalln(err)
		return
	}

	externalExKey, err := pubAccountExKey.Child(0)
	if err != nil {
		log.Fatalln(err)
		return
	}
	externalKey := createKey(
		net,
		1,
		[]uint32{
			0,
			hdkeychain.HardenedKeyStart + 44,
			hdkeychain.HardenedKeyStart + util.IndexFromNet[net],
			hdkeychain.HardenedKeyStart + index4,
			0},
		externalExKey)
	err = db.SaveKey(&externalKey)
	if err != nil {
		log.Fatalln(err)
		return
	}

	internalExKey, err := pubAccountExKey.Child(1)
	if err != nil {
		log.Fatalln(err)
		return
	}
	internalKey := createKey(
		net,
		1,
		[]uint32{
			0,
			hdkeychain.HardenedKeyStart + 44,
			hdkeychain.HardenedKeyStart + util.IndexFromNet[net],
			hdkeychain.HardenedKeyStart + index4,
			1},
		internalExKey)
	err = db.SaveKey(&internalKey)
	if err != nil {
		log.Fatalln(err)
		return
	}
}

// Receive create receive address
func Receive(name string, net string, index6 uint32) (string, error) {
	return getAddress(name, net, 0, index6)
}

// Change create Change address
func Change(name string, net string, index6 uint32) (string, error) {
	return getAddress(name, net, 1, index6)
}

func getAddress(name string, net string, index5 uint32, index6 uint32) (string, error) {
	account, err := db.GetAccount(name)
	if err != nil {
		log.Fatalln(err)
		return "", err
	}

	key, err := db.GetKey(
		0,
		1,
		[]uint32{
			0,
			hdkeychain.HardenedKeyStart + 44,
			hdkeychain.HardenedKeyStart + util.IndexFromNet[net],
			hdkeychain.HardenedKeyStart + account.Index4,
			index5,
			math.MaxUint32,
		},
	)

	pubAccountExKey, err := hdkeychain.NewKeyFromString(key.Value)
	if err != nil {
		log.Fatalln(err)
		return "", err
	}

	receiveKey, err := pubAccountExKey.Child(index6)
	if err != nil {
		log.Fatalln(err)
		return "", err
	}

	address, err := receiveKey.Address(util.ParamsFromNet[net])
	if err != nil {
		log.Fatalln(err)
		return "", err
	}

	return address.String(), nil
}

func createKey(
	net string,
	keyType uint8,
	indexes []uint32,
	exKey *hdkeychain.ExtendedKey) entity.Key {
	return entity.Key{
		Net:     uint8(util.IndexFromNet[net]),
		Type:    keyType,
		Indexes: indexes,
		Value:   exKey.String(),
	}
}
