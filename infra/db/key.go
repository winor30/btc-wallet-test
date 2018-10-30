package db

import (
	"database/sql"
	"fmt"
	"log"
	// to use sqlite3
	_ "github.com/mattn/go-sqlite3"
	"github.com/winor30/btc-wallet-test/entity"
	"github.com/winor30/btc-wallet-test/util"
	"math"
)

func exec(query string) (sql.Result, error) {
	db, err := sql.Open("sqlite3", util.DbFilePath)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	return db.Exec(query)
}

func queryAccount(query string) ([]entity.Account, error) {
	db, err := sql.Open("sqlite3", util.DbFilePath)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var accounts []entity.Account
	for rows.Next() {
		var name sql.NullString
		var index4 sql.NullInt64
		if err := rows.Scan(&name, &index4); err != nil {
			log.Fatalln(err)
		}
		accounts = append(
			accounts,
			entity.Account{
				Name:   name.String,
				Index4: uint32(index4.Int64),
			},
		)
	}

	return accounts, nil
}

func queryKey(query string) ([]entity.Key, error) {
	db, err := sql.Open("sqlite3", util.DbFilePath)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var keys []entity.Key
	for rows.Next() {
		var net sql.NullInt64
		var keyType sql.NullInt64
		var indexes [6]sql.NullInt64
		var value sql.NullString
		if err := rows.Scan(
			&net,
			&keyType,
			&indexes[0],
			&indexes[1],
			&indexes[2],
			&indexes[3],
			&indexes[4],
			&indexes[5],
			&value,
		); err != nil {
			log.Fatalln(err)
		}

		keys = append(
			keys,
			entity.Key{
				Net:  uint8(net.Int64),
				Type: uint8(keyType.Int64),
				Indexes: []uint32{
					(uint32)(indexes[0].Int64),
					(uint32)(indexes[1].Int64),
					(uint32)(indexes[2].Int64),
					(uint32)(indexes[3].Int64),
					(uint32)(indexes[4].Int64),
					(uint32)(indexes[5].Int64),
				},
				Value: value.String,
			},
		)
	}

	return keys, nil
}

// SaveKeys save some private or public keys for HD wallet
func SaveKeys(keys []*entity.Key) error {
	for _, key := range keys {
		err := SaveKey(key)
		if err != nil {
			return err
		}
	}

	return nil
}

// SaveKey save private or public key for HD wallet
func SaveKey(key *entity.Key) error {
	query := `
		INSERT INTO
		keys(net, type, index1, index2, index3, index4, index5, index6, value)
		VALUES(%d, %d, %d, %d, %d, %d, %d, %d, "%s")
	`
	fillMaxUint32(&key.Indexes)

	query = fmt.Sprintf(
		query,
		key.Net,
		key.Type,
		key.Indexes[0],
		key.Indexes[1],
		key.Indexes[2],
		key.Indexes[3],
		key.Indexes[4],
		key.Indexes[5],
		key.Value,
	)
	_, err := exec(query)
	return err
}

// GetKey get public or private key from db
func GetKey(net uint8, keyType uint8, indexes []uint32) (entity.Key, error) {

	query := fmt.Sprintf(`
		SELECT net, type, index1, index2, index3, index4, index5, index6, value from keys
		WHERE net=%d and type=%d and index1=%d and index2=%d and index3=%d and index4=%d and index5=%d and index6=%d;
		`,
		net,
		keyType,
		indexes[0],
		indexes[1],
		indexes[2],
		indexes[3],
		indexes[4],
		indexes[5],
	)
	keys, err := queryKey(query)

	if len(keys) != 1 {
		return entity.NullKey, err
	}

	return keys[0], err
}

// CreateDb is create db
func CreateDb() error {
	_, err := exec(`
		CREATE TABLE IF NOT EXISTS "keys" (net INTEGER, type INTEGER, index1 INTEGER, index2 INTEGER, index3 INTEGER, index4 INTEGER, index5 INTEGER, index6 INTEGER, value TEXT,
		PRIMARY KEY (net, type, index1, index2, index3, index4, index5, index6))
	`)
	if err != nil {
		return err
	}

	_, err = exec(`
		CREATE TABLE IF NOT EXISTS "accounts" (name TEXT PRIMARY KEY, index4 INTEGER)
	`)
	if err != nil {
		return err
	}

	return nil
}

func fillMaxUint32(indexes *([]uint32)) {
	for i := len(*indexes); i < 6; i++ {
		*indexes = append(*indexes, math.MaxUint32)
	}
}

// GetNewestIndexForAccount get index for account in order to create new account
func GetNewestIndexForAccount() uint32 {
	accounts, err := queryAccount(`
		SELECT name, max(index4) from accounts;
	`)

	log.Println(accounts, err)

	if len(accounts) != 0 && accounts[0] != entity.NullAccount {
		return accounts[0].Index4 + 1
	}

	return 0
}

// SaveAccount save account
func SaveAccount(account entity.Account) error {
	query := fmt.Sprintf(`
		INSERT INTO
		accounts(name, index4)
		VALUES("%s", %d)
		`,
		account.Name,
		account.Index4,
	)

	_, err := exec(query)
	return err
}
