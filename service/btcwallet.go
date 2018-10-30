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
	log.Println("latest private key depth is " + (string)(mPriv3.Depth()))

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
	log.Println(err)

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
