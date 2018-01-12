package randstr

import (
	"log"
	"math/rand"
	"time"
	"encoding/base64"

	"github.com/boltdb/bolt"
)

// GenerateRandomString func
// 랜덤 문자열 생성 함수
func GenerateRandomString(n int) string {
	randStr := generateRandomString(n)

	for !isRandStringUnique(randStr) {
		randStr = generateRandomString(n)
	}

	return randStr
}

func generateRandomBytes(n int) ([]byte, error) {
    b := make([]byte, n)
    _, err := rand.Read(b)
    // Note that err == nil only if we read len(b) bytes.
    if err != nil {
        return nil, err
    }

    return b, nil
}

// GenerateRandomString returns a URL-safe, base64 encoded
// securely generated random string.
func generateRandomString(s int) string {
    b, _ := generateRandomBytes(s)
    return base64.URLEncoding.EncodeToString(b)
}

// UseRandString func
// Delete used string from db
func UseRandString(str string) {
	if !isRandStringUnique(str) {
		db, err := bolt.Open("rand.db", 0600, nil)
		if err != nil {
			log.Fatalf("Failed to create or open rand.db file,\n%v", err)
		}
		defer db.Close()

		err = db.Update(func(tx *bolt.Tx) error {
			b, err := tx.CreateBucketIfNotExists([]byte("rand"))
			err = b.Delete([]byte(str))
			return err
		})

		if err != nil {
			log.Fatalf("Failed to delete string.\n%v", err)
		}
	}
}

func isRandStringUnique(randStr string) bool {
	db, err := bolt.Open("rand.db", 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		log.Fatalf("Failed to create or open rand.db file,\n%v", err)
	}
	defer db.Close()

	var v []byte

	err = db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte("rand"))
		v = b.Get([]byte(randStr))
		if v != nil {
			return nil
		}

		err = b.Put([]byte(randStr), []byte(randStr))

		return err
	})

	if err != nil {
		log.Fatal(err)
	}

	if v != nil {
		return false
	}

	return true
}
