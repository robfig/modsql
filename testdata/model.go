// MACHINE GENERATED BY ModSQL (github.com/kless/modsql)

package testdata

import (
	"database/sql"
	"time"

	"github.com/kless/modsql"
)

// == EDIT
var ENGINE = modsql.Postgres

//var ENGINE = modsql.MySQL
//var ENGINE = modsql.SQLite

// * * *

// Init prepares all statements in "listStatements".
// It hast to be called before of insert data.
func Init(db *sql.DB) {
	for _, v := range listStatements {
		v.Prepare(db, ENGINE)
	}
}

// Close closes all statements in "listStatements".
// Returns the first error, if any.
func Close() error {
	var err, errExit error

	for _, v := range listStatements {
		if err = v.Close(); err != nil && errExit == nil {
			errExit = err
		}
	}
	return errExit
}

var listStatements = []*modsql.Statements{insert}

// * * *

var insert = modsql.NewStatements(map[int]string{
	0:  "INSERT INTO types (t_int, t_int8, t_int16, t_int32, t_int64, t_float32, t_float64, t_string, t_binary, t_byte, t_rune, t_bool) VALUES({P}, {P}, {P}, {P}, {P}, {P}, {P}, {P}, {P}, {P}, {P}, {P})",
	1:  "INSERT INTO default_value (id, d_int8, d_float32, d_string, d_binary, d_byte, d_rune, d_bool) VALUES({P}, {P}, {P}, {P}, {P}, {P}, {P}, {P})",
	2:  "INSERT INTO times (typeId, t_duration, t_datetime) VALUES({P}, {P}, {P})",
	3:  "INSERT INTO account (acc_num, acc_type, acc_descr) VALUES({P}, {P}, {P})",
	4:  "INSERT INTO sub_account (sub_acc, ref_num, ref_type, sub_descr) VALUES({P}, {P}, {P}, {P})",
	5:  "INSERT INTO catalog (catalog_id, name, description, price) VALUES({P}, {P}, {P}, {P})",
	6:  "INSERT INTO magazine (catalog_id, page_count) VALUES({P}, {P})",
	7:  "INSERT INTO mp3 (catalog_id, size, length, filename) VALUES({P}, {P}, {P}, {P})",
	8:  "INSERT INTO book (book_id, title, author) VALUES({P}, {P}, {P})",
	9:  "INSERT INTO chapter (chapter_id, title, book_fk) VALUES({P}, {P}, {P})",
	10: "INSERT INTO {Q}user{Q} (user_id, first_name, last_name) VALUES({P}, {P}, {P})",
	11: "INSERT INTO address (address_id, street, city, state, post_code) VALUES({P}, {P}, {P}, {P}, {P})",
	12: "INSERT INTO user_address (user_id, address_id) VALUES({P}, {P})",
})

// sex
const (
	SEX_FEMALE = iota
	SEX_MALE
)

type Types struct {
	T_int     int
	T_int8    int8
	T_int16   int16
	T_int32   int32
	T_int64   int64
	T_float32 float32
	T_float64 float64
	T_string  string
	T_binary  []byte
	T_byte    byte
	T_rune    rune
	T_bool    bool
}

func (t *Types) Args() ([]interface{}, error) {
	return []interface{}{
		t.T_int, t.T_int8, t.T_int16, t.T_int32, t.T_int64, t.T_float32, t.T_float64, t.T_string, t.T_binary, t.T_byte, t.T_rune, modsql.BoolToSQL(ENGINE, t.T_bool),
	}, nil
}

func (t *Types) StmtInsert() *sql.Stmt { return insert.Stmt[0] }

type Default_value struct {
	Id        int
	D_int8    int8
	D_float32 float32
	D_string  string
	D_binary  []byte
	D_byte    byte
	D_rune    rune
	D_bool    bool
}

func (t *Default_value) Args() ([]interface{}, error) {
	return []interface{}{
		t.Id, t.D_int8, t.D_float32, t.D_string, t.D_binary, t.D_byte, t.D_rune, modsql.BoolToSQL(ENGINE, t.D_bool),
	}, nil
}

func (t *Default_value) StmtInsert() *sql.Stmt { return insert.Stmt[1] }

type Times struct {
	TypeId     int
	T_duration time.Duration
	T_datetime time.Time
}

func (t *Times) Args() ([]interface{}, error) {
	t0, err := time.Parse(time.RFC3339, t.T_datetime.String())
	if err != nil {
		return nil, err
	}
	return []interface{}{
		t.TypeId, modsql.ReplTime.Replace(t.T_duration.String()), t0.String(),
	}, nil
}

func (t *Times) StmtInsert() *sql.Stmt { return insert.Stmt[2] }

type Account struct {
	Acc_num   int
	Acc_type  int
	Acc_descr string
}

func (t *Account) Args() ([]interface{}, error) {
	return []interface{}{
		t.Acc_num, t.Acc_type, t.Acc_descr,
	}, nil
}

func (t *Account) StmtInsert() *sql.Stmt { return insert.Stmt[3] }

type Sub_account struct {
	Sub_acc   int
	Ref_num   int
	Ref_type  int
	Sub_descr string
}

func (t *Sub_account) Args() ([]interface{}, error) {
	return []interface{}{
		t.Sub_acc, t.Ref_num, t.Ref_type, t.Sub_descr,
	}, nil
}

func (t *Sub_account) StmtInsert() *sql.Stmt { return insert.Stmt[4] }

type Catalog struct {
	Catalog_id  int
	Name        string
	Description string
	Price       float32
}

func (t *Catalog) Args() ([]interface{}, error) {
	return []interface{}{
		t.Catalog_id, t.Name, t.Description, t.Price,
	}, nil
}

func (t *Catalog) StmtInsert() *sql.Stmt { return insert.Stmt[5] }

type Magazine struct {
	Catalog_id int
	Page_count string
}

func (t *Magazine) Args() ([]interface{}, error) {
	return []interface{}{
		t.Catalog_id, t.Page_count,
	}, nil
}

func (t *Magazine) StmtInsert() *sql.Stmt { return insert.Stmt[6] }

type Mp3 struct {
	Catalog_id int
	Size       int
	Length     float32
	Filename   string
}

func (t *Mp3) Args() ([]interface{}, error) {
	return []interface{}{
		t.Catalog_id, t.Size, t.Length, t.Filename,
	}, nil
}

func (t *Mp3) StmtInsert() *sql.Stmt { return insert.Stmt[7] }

type Book struct {
	Book_id int
	Title   string
	Author  string
}

func (t *Book) Args() ([]interface{}, error) {
	return []interface{}{
		t.Book_id, t.Title, t.Author,
	}, nil
}

func (t *Book) StmtInsert() *sql.Stmt { return insert.Stmt[8] }

type Chapter struct {
	Chapter_id int
	Title      string
	Book_fk    int
}

func (t *Chapter) Args() ([]interface{}, error) {
	return []interface{}{
		t.Chapter_id, t.Title, t.Book_fk,
	}, nil
}

func (t *Chapter) StmtInsert() *sql.Stmt { return insert.Stmt[9] }

type User struct {
	User_id    int
	First_name string
	Last_name  string
}

func (t *User) Args() ([]interface{}, error) {
	return []interface{}{
		t.User_id, t.First_name, t.Last_name,
	}, nil
}

func (t *User) StmtInsert() *sql.Stmt { return insert.Stmt[10] }

type Address struct {
	Address_id int
	Street     string
	City       string
	State      string
	Post_code  string
}

func (t *Address) Args() ([]interface{}, error) {
	return []interface{}{
		t.Address_id, t.Street, t.City, t.State, t.Post_code,
	}, nil
}

func (t *Address) StmtInsert() *sql.Stmt { return insert.Stmt[11] }

type User_address struct {
	User_id    int
	Address_id int
}

func (t *User_address) Args() ([]interface{}, error) {
	return []interface{}{
		t.User_id, t.Address_id,
	}, nil
}

func (t *User_address) StmtInsert() *sql.Stmt { return insert.Stmt[12] }
