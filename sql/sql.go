package sql

import (
	"database/sql"
	"log"
)

type MyDB struct {
	db *sql.DB
}

//https://qiita.com/fiemon/items/eb38c8d681ed1ae05925
//　Transaction
func (m *MyDB) Transaction(txFunc func(*sql.Tx) error) error {
	tx, err := m.db.Begin()
	//tx, err := db.MustBegin()
	if err != nil {
		return err
	}

	defer func() {
		if p := recover(); p != nil {
			log.Println("recover")
			tx.Rollback()
			panic(p)
		} else if err != nil {
			log.Println("rollback")
			tx.Rollback()
		} else {
			log.Println("commit")
			err = tx.Commit()
		}
	}()
	err = txFunc(tx)
	return err
}
