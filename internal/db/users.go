package db

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/joseph0x45/tessera/internal/models"
	"github.com/joseph0x45/tessera/internal/shared"
)

func (c *Conn) InsertUser(user *models.User) error {
	const query = `
    insert into users (
      id, app_id, name, password
    )
    values (
      :id, :app_id, :name, :password
    );
  `
	if _, err := c.db.NamedExec(query, user); err != nil {
		return fmt.Errorf("Error while inserting user: %w", err)
	}
	return nil
}

func (c *Conn) GetUsersByAppID(appID string) ([]models.User, error) {
	const query = "select * from users where app_id=?"
	users := []models.User{}
	if err := c.db.Select(&users, query); err != nil {
		return nil, fmt.Errorf("Error while getting users by id: %w", err)
	}
	return users, nil
}

func (c *Conn) GetUser(userID, appID string) (*models.User, error) {
	const query = "select * from users where id=? and app_id=?"
	user := &models.User{}
	err := c.db.Get(user, query, userID, appID)
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
