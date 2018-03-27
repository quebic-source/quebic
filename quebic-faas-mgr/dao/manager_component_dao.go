package dao

import (
	"encoding/json"
	"fmt"
	"quebic-faas/types"
	"time"

	bolt "github.com/coreos/bbolt"
)

//AddManagerComponentDockerImageID set DockerImageID
func AddManagerComponentDockerImageID(db *bolt.DB, component *types.ManagerComponent, dockerImageID string) error {

	err := getByID(db, component, func(savedObj []byte) error {

		if savedObj == nil {
			return fmt.Errorf("unable to found component")
		}

		json.Unmarshal(savedObj, component)

		return nil
	})

	if err != nil {
		return err
	}

	component.DockerImageID = dockerImageID
	return Save(db, component)

}

//AddManagerComponentDockerContainerID set ContainerID
func AddManagerComponentDockerContainerID(db *bolt.DB, component *types.ManagerComponent, containerID string) error {

	err := getByID(db, component, func(savedObj []byte) error {

		if savedObj == nil {
			return fmt.Errorf("unable to found component")
		}

		json.Unmarshal(savedObj, component)

		return nil
	})

	if err != nil {
		return err
	}

	component.DockerContainerID = containerID
	return Save(db, component)

}

//AddManagerComponentLog add log
func AddManagerComponentLog(db *bolt.DB, component *types.ManagerComponent, log types.EntityLog) error {

	log.Time = time.Now().String()

	err := getByID(db, component, func(savedObj []byte) error {

		if savedObj == nil {
			return fmt.Errorf("unable to found component")
		}

		json.Unmarshal(savedObj, component)

		return nil
	})

	if err != nil {
		return err
	}

	component.Log = log
	return Save(db, component)
}
