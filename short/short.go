package short

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/shweta20039/shortme/base"
	"github.com/shweta20039/shortme/conf"
	"github.com/shweta20039/shortme/sequence"
	_ "github.com/shweta20039/shortme/sequence/db"
)

type shorter struct {
	readDB   *sql.DB
	writeDB  *sql.DB
	sequence sequence.Sequence
}

func (shorter *shorter) Analytics(shortURL string) (int, error) {
	analyticsQuery := fmt.Sprintf(`SELECT fetch_count FROM short WHERE short_url=?`)
	rows, err := shorter.readDB.Query(analyticsQuery, shortURL)
	if err != nil {
		log.Printf("short read db query error. %v", err)
		return 0, errors.New("short read db query error")
	}

	defer rows.Close()
	var fetch_count int
	err = rows.Scan(&fetch_count)
	if err != nil {
		log.Printf("error while fetching count. %v", err)
		return 0, errors.New("fetch count read db query error")
	}
	return fetch_count, nil

}
func (shorter *shorter) Expand(shortURL string) (string, error) {
	selectSQL := fmt.Sprintf(`SELECT long_url, fetch_count FROM short WHERE short_url=?`)
	rows, err := shorter.readDB.Query(selectSQL, shortURL)
	if err != nil {
		log.Printf("short read db query error. %v", err)
		return "", errors.New("short read db query error")
	}
	var longURL string
	var fetch_count int
	defer rows.Close()

	err = rows.Scan(&longURL, &fetch_count)
	if err != nil {
		log.Printf("short read db query rows scan error. %v", err)
		return "", errors.New("short read db query rows scan error")
	}

	updateQuery := fmt.Sprintf("UPDATE short SET fetch_count = ? WHERE short_url = ?")
	var stmt *sql.Stmt
	defer stmt.Close()
	stmt, err = shorter.writeDB.Prepare(updateQuery)
	if err != nil {
		log.Printf("update db query error. %v", err)
		return "", errors.New("update db query error")
	}
	fetch_count++
	_, err = stmt.Exec(fetch_count, shortURL)
	if err != nil {
		log.Printf("short write db insert error. %v", err)
		return "", errors.New("short write db insert error")
	}
	return longURL, nil
}

func (shorter *shorter) Short(longURL string) (shortURL string, err error) {
	for {
		var seq uint64
		seq, err = shorter.sequence.NextSequence()
		if err != nil {
			log.Printf("get next sequence error. %v", err)
			return "", errors.New("get next sequence error")
		}

		shortURL = base.Int2String(seq)
		if _, exists := conf.Conf.Common.BlackShortURLsMap[shortURL]; exists {
			continue
		} else {
			break
		}
	}

	insertSQL := fmt.Sprintf(`INSERT INTO short(long_url, short_url, fetch_count) VALUES(?, ?, 0)`)

	var stmt *sql.Stmt
	stmt, err = shorter.writeDB.Prepare(insertSQL)
	if err != nil {
		log.Printf("short write db prepares error. %v", err)
		return "", errors.New("short write db prepares error")
	}
	defer stmt.Close()

	_, err = stmt.Exec(longURL, shortURL)
	if err != nil {
		log.Printf("short write db insert error. %v", err)
		return "", errors.New("short write db insert error")
	}

	return shortURL, nil
}

var Shorter shorter

func Start() {
	//DB connection and intialisation here
	log.Println("shorter starts")
}
