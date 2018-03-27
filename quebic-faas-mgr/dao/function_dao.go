package dao

import (
	"encoding/json"
	"fmt"
	"quebic-faas/types"
	"time"

	bolt "github.com/coreos/bbolt"
)

//AddFunctionDockerImageID set DockerImageID
func AddFunctionDockerImageID(db *bolt.DB, function *types.Function, dockerImageID string) error {

	err := getByID(db, function, func(savedObj []byte) error {

		if savedObj == nil {
			return fmt.Errorf("unable to found function")
		}

		json.Unmarshal(savedObj, function)

		return nil
	})

	if err != nil {
		return err
	}

	function.DockerImageID = dockerImageID
	return Save(db, function)
}

//AddFunctionLog add function log
func AddFunctionLog(db *bolt.DB, function *types.Function, log types.EntityLog) error {

	log.Time = time.Now().String()

	err := getByID(db, function, func(savedObj []byte) error {

		if savedObj == nil {
			return fmt.Errorf("unable to found function")
		}

		json.Unmarshal(savedObj, function)

		return nil
	})

	if err != nil {
		return err
	}

	function.Log = log
	return Save(db, function)
}
