package db

import (
	bolt "github.com/coreos/bbolt"
)

const defaultDBPath string = "quebic-faas-mgr.db"

//GetDb get database
func GetDb() (*bolt.DB, error) {
	return bolt.Open(defaultDBPath, 0600, nil)
}
