package db

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	"github.com/joseph0x45/tessera/internal/models"
	"github.com/joseph0x45/tessera/internal/shared"
)

func (c *Conn) UserExistsInApp(userName, appID string) bool {
	const query = "select exists(select 1 from users where name=? and app_id=?)"
	exists := false
	err := c.db.QueryRow(query, userName, appID).Scan(&exists)
	if err != nil {
		log.Println("Error while checking if user exists in app:", err)
		return true
	}
	return exists
}

func (c *Conn) InsertUser(user *models.User, tx *sqlx.Tx) error {
	const query = `
    insert into users (
      id, app_id, name, password
    )
    values (
      :id, :app_id, :name, :password
    );
  `
	if c.UserExistsInApp(user.Name, user.AppID) {
		return shared.ErrUserExistsInApp
	}
	var err error
	if tx != nil {
		_, err = tx.NamedExec(query, user)
	} else {
		_, err = c.db.NamedExec(query, user)
	}
	if err != nil {
		return fmt.Errorf("Error while inserting user: %w", err)
	}
	return nil
}

func (c *Conn) InsertUserAndSession(user *models.User, session *models.Session) error {
	tx, err := c.db.Beginx()
	if err != nil {
		return fmt.Errorf("Error while starting transaction: %w", err)
	}
	defer tx.Rollback()

	if err := c.InsertUser(user, tx); err != nil {
		return err
	}

	if err := c.InsertSession(session, tx); err != nil {
		return err
	}
	return tx.Commit()
}

func (c *Conn) GetUsersByAppID(appID string) ([]models.User, error) {
	const query = "select * from users where app_id=?"
	users := []models.User{}
	if err := c.db.Select(&users, query, appID); err != nil {
		return nil, fmt.Errorf("Error while getting users by id: %w", err)
	}
	return users, nil
}

func (c *Conn) GetUserByNameAndAppID(userName, appID string) (*models.User, error) {
	const query = "select * from users where name=? and app_id=?"
	user := &models.User{}
	err := c.db.Get(user, query, userName, appID)
	if err == nil {
		return user, nil
	}
	if errors.Is(err, sql.ErrNoRows) {
		return nil, shared.ErrUserNotFound
	}
	return nil, fmt.Errorf("Error while getting user: %w", err)
}

func (c *Conn) DeleteUser(userID string) error {
	const query = "delete from users where id=?"
	if _, err := c.db.Exec(query, userID); err != nil {
		return fmt.Errorf("Error while deleting user: %w", err)
	}
	return nil
}
