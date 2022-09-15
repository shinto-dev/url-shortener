package shortid

import (
	"math/rand"
	"strconv"
	"time"
	"url-shortner/platform/data"

	"github.com/itchyny/base58-go"
	"gorm.io/gorm"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type IDGenerator struct {
	db       *gorm.DB
	encoding *base58.Encoding
}

func NewShortIDGenerator(db *gorm.DB) *IDGenerator {
	return &IDGenerator{db: db, encoding: base58.FlickrEncoding}
}

func (i *IDGenerator) NewBase58ID() (string, error) {
	id, err := i.generateNextAutoIncrementID()
	if err != nil {
		return "", err
	}

	//todo randomize the id: https://www.enjoyalgorithms.com/blog/design-a-url-shortening-service-like-tiny-url
	id = id*10 + rand.Int63n(10)

	// TODO custom implement this without string conversion in the middle
	encoded, err := i.encoding.Encode([]byte(strconv.FormatInt(id, 10)))
	if err != nil {
		return "", err
	}

	return string(encoded), nil
}

func (i *IDGenerator) generateNextAutoIncrementID() (int64, error) {
	var id int64

	err := data.WithTransaction(i.db, func(tx *gorm.DB) error {
		err := tx.Exec("REPLACE INTO unique_ids_64 (stub) VALUES (?)", "s").Error
		if err != nil {
			return err
		}

		err = tx.Raw("SELECT LAST_INSERT_ID()").Scan(&id).Error
		if err != nil {
			return err
		}

		return nil
	})

	return id, err
}
