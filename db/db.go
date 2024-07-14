package db

import (
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Sql struct {
	Db       *sqlx.DB
	Host     string
	Port     int
	Username string
	Password string
	Dbname   string
}

func (s *Sql) Connect() {
	datasource := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		s.Host, s.Port, s.Username, s.Password, s.Dbname)
	//postgres là username d
	s.Db = sqlx.MustConnect("postgres", datasource)

	if err := s.Db.Ping(); err != nil {
		log.Println(err.Error())
	}
	fmt.Println("database ok")
}

func (s *Sql) Close() {
	// Gọi phương thức Close() của cấu trúc sqlx.DB được nhúng vào
	if s.Db != nil {
		s.Db.Close()
	}
}

func (s *Sql) InsertData(TxId string) error {
	query := `INSERT INTO public."TxTransaction" ("TxId") VALUES ($1)`
	_, err := s.Db.Exec(query, TxId)
	if err != nil {
		return err
	}
	fmt.Println("Data inserted successfully database!")
	return nil
}
