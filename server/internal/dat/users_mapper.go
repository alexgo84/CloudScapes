package dat

import (
	"CloudScapes/pkg/wire"
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
)

type UsersMapper struct {
	txn *sqlx.Tx
	ctx context.Context
}

type User struct {
	Id           int64     `json:"id" db:"id"`
	Created      time.Time `json:"created" db:"created_at"`
	PasswordHash string    `json:"-" db:"password_hash"`
	wire.NewUser
}

func NewUsersMapper(ctx context.Context, txn *sqlx.Tx) UsersMapper {
	return UsersMapper{
		txn: txn,
		ctx: ctx,
	}
}

func (am *UsersMapper) CreateUser(newUser *wire.NewUser) (*User, error) {

	u := User{NewUser: *newUser}

	passHash, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	u.PasswordHash = string(passHash)

	err = namedGet(am.txn, "INSERT INTO users (name, email, accountid, password_hash) VALUES (:name, :email, :accountid, :password_hash) RETURNING id, created_at", &u)
	if err != nil {
		return nil, err
	}

	return &u, nil
}

func (am *UsersMapper) GetUsers(accountID int64) ([]User, error) {
	users := []User{} // assign to empty array so that no result case does not return null
	err := am.txn.SelectContext(am.ctx, &users, "select * from users WHERE accountid = $1 ORDER BY id desc", accountID)
	if errors.Is(err, sql.ErrNoRows) {
		return []User{}, nil
	}
	if err != nil {
		return nil, err
	}
	return users, nil
}
