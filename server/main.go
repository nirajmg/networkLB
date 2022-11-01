package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/caarlos0/env"
	_ "github.com/lib/pq"
)

type User struct {
	UserId      string
	Description string
}

type DB struct {
	conn     *sql.DB
	Password string `env:"PQ_PASS" envDefault:"root"`
	User     string `env:"PQ_USER" envDefault:"root"`
	Host     string `env:"PQ_HOST" envDefault:"postgresql-hl"`
	Port     int64  `env:"PQ_PORT" envDefault:"5432"`
	Database string `env:"PQ_DB" envDefault:"mydb"`
}

func addUser(w http.ResponseWriter, r *http.Request) {

	var u User
	err := json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		fmt.Println(err)
		return
	}
	if err := db.insert(u.UserId, u.Description); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		fmt.Println(err)
		return
	}
	fmt.Fprintf(w, "Success")
}
func getUser(w http.ResponseWriter, r *http.Request) {
	var u User
	err := json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		fmt.Println(err)
		return
	}

	u.Description, err = db.getByid(u.UserId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		fmt.Println(err)
		return
	}
	res, err := json.Marshal(u)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		fmt.Println(err)
		return

	}
	fmt.Fprintf(w, string(res))
}

func health(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Healthy")
}

var db = DB{}

func main() {

	env.Parse(&db)

	if err := db.NewConn(); err != nil {
		panic(err)
	}

	http.HandleFunc("/add", addUser)
	http.HandleFunc("/get", getUser)
	http.HandleFunc("/", health)

	if err := http.ListenAndServe(":80", nil); err != nil {
		panic(err)
	}

}

func (db *DB) NewConn() error {
	connString := "host=%s port=%d user=%s password=%s dbname=%s sslmode=disable"
	var err error
	print(fmt.Sprintf(connString, db.Host, db.Port, db.User, db.Password, db.Database))
	db.conn, err = sql.Open("postgres", fmt.Sprintf(connString, db.Host, db.Port, db.User, db.Password, db.Database))
	if err != nil {
		return err
	}

	query := `create table if not exists nlb(user_id VARCHAR(255) primary key, description VARCHAR(255) NOT NULL)`
	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()

	if _, err := db.conn.ExecContext(ctx, query); err != nil {
		return err
	}
	return nil
}

func (db *DB) insert(userid, description string) error {
	insertDynStmt := `insert into "nlb"("user_id", "description") values($1, $2)`
	if _, err := db.conn.Exec(insertDynStmt, userid, description); err != nil {
		return err
	}
	return nil

}

func (db *DB) getByid(user_id string) (string, error) {
	sqlStatement := `SELECT  description FROM nlb WHERE user_id=$1;`
	var description string

	row := db.conn.QueryRow(sqlStatement, user_id)
	if err := row.Scan(&description); err != nil {
		return "", err
	}
	return description, nil

}
