package user

import (
	"context"
	"fmt"
	"log"

	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgx/v4"
	"github.com/mittz/roleplay-webapp-portal/database"
)

type User struct {
	Userkey   string `json:"userkey"`
	LDAP      string `json:"ldap"`
	Team      string `json:"team"`
	Region    string `json:"region"`
	SubRegion string `json:"sub_region"`
	Role      string `json:"role"`
}

type BulkUsers struct {
	Users []User `json:"users"`
}

func NewUser() User {
	return User{}
}

func NewBulkUsers() BulkUsers {
	return BulkUsers{}
}

func GetUser(userkey string) User {
	dbPool := database.GetDatabaseConnection()

	user := User{}
	if err := dbPool.QueryRow(context.Background(), "select * from users where userkey=$1", userkey).Scan(
		&user.Userkey,
		&user.LDAP,
		&user.Team,
		&user.Region,
		&user.SubRegion,
		&user.Role,
	); err != nil && err != pgx.ErrNoRows {
		log.Printf("QueryRow failed: %v\n", err)
	}

	return user
}

func GetUsers() []*User {
	dbPool := database.GetDatabaseConnection()

	var users []*User
	if err := pgxscan.Select(context.Background(), dbPool, &users, `SELECT userkey, ldap, team, region, sub_region, role from users`); err != nil {
		log.Println(err)
	}

	return users
}

func (b BulkUsers) OverrideDatabase() error {
	dbPool := database.GetDatabaseConnection()

	queryTableCreation := `
	DROP TABLE IF EXISTS users;
	CREATE TABLE users (
		userkey character varying(50) NOT NULL,
		ldap character varying(20) NOT NULL,
		team character varying(100) NOT NULL,
		region character varying(20) NOT NULL,
		sub_region character varying(40) NOT NULL,
		role character varying(100) NOT NULL,
		PRIMARY KEY(userkey)
	);
	GRANT ALL ON users TO PUBLIC;
	`

	if _, err := dbPool.Exec(context.Background(), queryTableCreation); err != nil {
		return fmt.Errorf("Table recreation failed: %v\n", err)
	}

	if _, err := dbPool.CopyFrom(
		context.Background(),
		pgx.Identifier{"users"},
		[]string{"userkey", "ldap", "team", "region", "sub_region", "role"},
		pgx.CopyFromSlice(len(b.Users), func(i int) ([]interface{}, error) {
			return []interface{}{
				b.Users[i].Userkey,
				b.Users[i].LDAP,
				b.Users[i].Team,
				b.Users[i].Region,
				b.Users[i].SubRegion,
				b.Users[i].Role,
			}, nil
		}),
	); err != nil {
		return fmt.Errorf("Bulk insertion failed: %v\n", err)
	}

	return nil
}
