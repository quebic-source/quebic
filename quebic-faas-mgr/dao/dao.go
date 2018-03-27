package dao

import (
	"encoding/json"
	"fmt"
	"quebic-faas/types"

	bolt "github.com/coreos/bbolt"
)

//Add entity.
// Check before save.
// If already exists a object under id. throw error
func Add(db *bolt.DB, entity types.Entity) error {

	//check allready exists
	err := getByID(db, entity, func(savedObj []byte) error {

		if savedObj != nil {
			return fmt.Errorf("object already exists")
		}

		return nil
	})

	if err != nil {
		return err
	}

	return Save(db, entity)
}

//Update entity
// Check before save.
// If unable to found a object under id. Throw object not found error
func Update(db *bolt.DB, entity types.Entity) error {

	//check for id
	err := getByID(db, entity, func(savedObj []byte) error {

		if savedObj == nil {
			return fmt.Errorf("unable to found object")
		}

		return nil
	})

	if err != nil {
		return err
	}

	return Save(db, entity)
}

//Save entity
// If there is no any entity under id add new.
// Otherwise save new entity under previous entity
func Save(db *bolt.DB, entity types.Entity) error {

	objVal := entity.GetReflectObject().Elem()
	typeName := objVal.Type().Name()
	id := entity.GetID()

	typeNameInBytes := []byte(typeName)
	idInBytes := []byte(id)

	entityJSON, err := json.Marshal(entity)
	if err != nil {
		return fmt.Errorf("failed json parse, error : %v", err)
	}

	//log.Printf("saving %s ,id %s, data %s\n", typeName, id, entityJSON)

	return db.Update(func(tx *bolt.Tx) error {

		bucket, err := tx.CreateBucketIfNotExists(typeNameInBytes)
		if err != nil {
			return fmt.Errorf("unable to create bucket for %s, error : %v", typeName, err)
		}

		err = bucket.Put(idInBytes, entityJSON)
		if err != nil {
			return fmt.Errorf("unable to put data for %s, error : %v", typeName, err)
		}

		//log.Printf("saved %s ,id %s\n", typeName, id)

		return nil
	})

}

//GetAll entity
func GetAll(db *bolt.DB, entity types.Entity, fn func(k, v []byte) error) error {
	return getAll(db, entity, fn)
}

//GetByID get by ID
func GetByID(db *bolt.DB, entity types.Entity, fn func(v []byte) error) error {
	return getByID(db, entity, fn)
}

func getAll(db *bolt.DB, entity types.Entity, fn func(k, v []byte) error) error {

	objVal := entity.GetReflectObject().Elem()
	typeName := objVal.Type().Name()

	typeNameInBytes := []byte(typeName)

	return db.View(func(tx *bolt.Tx) error {

		bucket := tx.Bucket([]byte(typeNameInBytes))

		if bucket != nil {
			return bucket.ForEach(fn)
		}

		return nil

	})

}

func getByID(db *bolt.DB, entity types.Entity, fn func(v []byte) error) error {

	objVal := entity.GetReflectObject().Elem()
	typeName := objVal.Type().Name()
	id := entity.GetID()

	typeNameInBytes := []byte(typeName)
	idInBytes := []byte(id)

	return db.View(func(tx *bolt.Tx) error {

		bucket := tx.Bucket([]byte(typeNameInBytes))

		if bucket != nil {
			if fn != nil {
				return fn(bucket.Get(idInBytes))
			}
		}

		return nil

	})

}
