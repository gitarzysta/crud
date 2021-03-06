package crud_test

import (
	"database/sql"
	"github.com/azer/crud"
	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
	"time"
)

var DB *crud.DB

type UserProfile struct {
	Id       int    `json:"id" sql:"auto-increment primary-key required"`
	Name     string `json:"name" sql:"required"`
	Bio      string `json:"bio" sql:"type=text"`
	Email    string `json:"e-mail" sql:"name=email"`
	Modified int64  `json:"modified" sql:"name=modified"`
}

type UserProfileNull struct {
	Id       sql.NullInt64  `json:"id" sql:"auto-increment primary-key required"`
	Name     sql.NullString `json:"name" sql:"required"`
	Bio      sql.NullString `json:"bio" sql:"type=text"`
	Email    sql.NullString `json:"e-mail" sql:"name=email"`
	Modified sql.NullInt64  `json:"modified" sql:"name=modified"`
}

type Post struct {
	Id        int       `json:"id" sql:"auto-increment primary-key required table-name=renamed_post"`
	Title     string    `json:"title"`
	Text      string    `json:"text"`
	CreatedAt time.Time `json:"created_at"`
}

type Foo struct {
	Id     int
	APIKey string
	YOLO   bool
	Beast  string
}

type EmbeddedFoo struct {
	Foo
	Span int
	Eggs string
}

type FooSlice []Foo
type FooPTRSlice []*Foo

func init() {
	var err error
	DB, err = crud.Connect("mysql", os.Getenv("DATABASE_URL"))
	if err != nil {
		panic(err)
	}
}

func TestPing(t *testing.T) {
	assert.Nil(t, DB.Ping())
}

func TestExecuteSQL(t *testing.T) {
	result, err := DB.Client.Exec("SHOW TABLES LIKE 'shouldnotexist'")
	assert.Nil(t, err)

	l, err := result.LastInsertId()
	assert.Equal(t, err, nil)
	assert.Equal(t, l, int64(0))

	a, err := result.RowsAffected()
	assert.Equal(t, err, nil)
	assert.Equal(t, a, int64(0))
}

func TestCreateTables(t *testing.T) {
	err := DB.CreateTables(UserProfile{}, Post{})
	assert.Nil(t, err)
	assert.True(t, DB.CheckIfTableExists("user_profile"))
	assert.True(t, DB.CheckIfTableExists("renamed-post"))
}

func TestDropTables(t *testing.T) {
	err := DB.DropTables(UserProfile{}, Post{})
	assert.Nil(t, err)
	assert.False(t, DB.CheckIfTableExists("user_profile"))
	assert.False(t, DB.CheckIfTableExists("post"))
}
