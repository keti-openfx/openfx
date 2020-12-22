package main

import (
	"fmt"
	"log"

	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

// func saveJSON(db *buntdb.DB) []byte {
// 	keys := make(map[string]string)
// 	db.View(func(tx *buntdb.Tx) error {
// 		tx.Ascend("", func(key, val string) bool {
// 			keys[key] = val
// 			//fmt.Printf("key is %s\t", key)
// 			//fmt.Printf("value is %s\n", val)
// 			return true
// 		})
// 		return nil
// 	})
// 	data, _ := json.Marshal(keys)
// 	return data
// }

func main() {
	// Open the data.db file. It will be created if it doesn't exist.
	// db, err := buntdb.Open("oauth2.db")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer db.Close()

	// err = db.View(func(tx *buntdb.Tx) error {
	// 	val, err := tx.Get("clientID")
	// 	if err != nil {
	// 		return err
	// 	}
	// 	fmt.Printf("value is %s\n", val)
	// 	return nil

	// })

	// err = db.Update(func(tx *buntdb.Tx) error {
	// 	_, _, err := tx.Set("test", "echo", nil)
	// 	return err
	// })

	// b := saveJSON(db)
	// fmt.Println(string(b))

	db, err := sql.Open("mysql", "root:root@tcp(10.0.0.91:3308)/OAuth2")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	rows, err := db.Query("SELECT scope FROM role where id ='test'")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	functions := []string{}

	for rows.Next() {
		var function string
		err := rows.Scan(&function)
		if err != nil {
			log.Fatal(err)
		}
		functions = append(functions, function)
	}

	fmt.Println(functions)

	// []string -> value 로
}
