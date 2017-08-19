package randstr

import (
	"log"
	"math/rand"
	"time"

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
			b := tx.Bucket([]byte("rand"))
			err = b.Delete([]byte(str))
			return err
		})

		if err != nil {
			log.Fatalf("Failed to delete string.\n%v", err)
		}
	}
}

func generateRandomString(n int) string {
	var letterRunes = []rune("0123456789abcdefghijklmnopqrstuvwxyz01234567899876543210ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

	randSource := rand.NewSource(time.Now().UnixNano())
	r := rand.New(randSource)

	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[r.Intn(len(letterRunes))]
	}

	return string(b)
}

func isRandStringUnique(randStr string) bool {
	db, err := bolt.Open("rand.db", 0600, nil)
	if err != nil {
		log.Fatalf("Failed to create or open rand.db file,\n%v", err)
	}
	defer db.Close()

	var v []byte

	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("rand"))
		v = b.Get([]byte(randStr))
		if v == nil {
			err = b.Put([]byte(randStr), []byte(randStr))
		}

		return err
	})

	if v != nil {
		return false
	}

	return true
}
