package store

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"sync"

	"github.com/tidwall/buntdb"
	"gopkg.in/oauth2.v4"
	"gopkg.in/oauth2.v4/models"
)

// NewClientStore create client store
func NewClientStore() *ClientStore {
	stor := &ClientStore{
		data: make(map[string]oauth2.ClientInfo),
	}
	db, err := buntdb.Open("./db/client.db")
	if err != nil {
		log.Println("can't load buntdb.")
		return nil
	}
	// 가져오기 View 구현 필요
	err = db.View(func(tx *buntdb.Tx) error {
		err := tx.Ascend("", func(key, value string) bool {
			var tm models.Client
			err = json.Unmarshal([]byte(value), &tm)
			if err != nil {
				return false
			}

			stor.data[key] = &tm
			return true
		})
		return err
	})
	return stor
}

// ClientStore client information store
type ClientStore struct {
	sync.RWMutex
	data map[string]oauth2.ClientInfo
}

// GetByID according to the ID for the client information
func (cs *ClientStore) GetByID(ctx context.Context, id string) (oauth2.ClientInfo, error) {
	cs.RLock()
	defer cs.RUnlock()

	if c, ok := cs.data[id]; ok {
		return c, nil
	}
	return nil, errors.New("not found")
}

// Set set client information
func (cs *ClientStore) Set(id string, cli oauth2.ClientInfo) (err error) {
	cs.Lock()
	defer cs.Unlock()

	cs.data[id] = cli

	db, err := buntdb.Open("./db/client.db")
	if err != nil {
		log.Println("can't load buntdb.")
		return err
	}
	jv, err := json.Marshal(cli)

	db.Update(func(tx *buntdb.Tx) error {
		tx.Set(id, string(jv), nil)
		return nil
	})

	return
}
