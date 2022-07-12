package gophermart

import (
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"github.com/nastradamus39/gophermart/internal/db"
)

var DB *sqlx.DB

func InitDB() (err error) {
	if DB, err = sqlx.Open("postgres", Cfg.DatabaseDsn); err == nil {
		err = migrate()
		if err != nil {
			log.Println(err)
			return
		}
	} else {
		log.Println(err)
	}
	return
}

func migrate() (err error) {
	// миграции
	_, err = DB.Exec(db.Schema)
	if err != nil {
		log.Println(err)
	}

	return
}
