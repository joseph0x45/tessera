package db

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/joseph0x45/tessera/internal/models"
	"github.com/joseph0x45/tessera/internal/shared"
)

func (c *Conn) InsertApp(app *models.App) error {
	const query = `
    insert into apps (
      id, name
    )
    values (
      :id, :name
    );
  `
	if _, err := c.db.NamedExec(query, app); err != nil {
		return fmt.Errorf("Error while inserting app: %w", err)
	}
	return nil
}

func (c *Conn) GetAllApps() ([]models.App, error) {
	apps := []models.App{}
	const query = "select * from apps"
	if err := c.db.Select(&apps, query); err != nil {
		return nil, fmt.Errorf("Error while getting apps: %w", err)
	}
	return apps, nil
}

func (c *Conn) DeleteApp(appID string) error {
	const query = "delete from apps where id=?"
	if _, err := c.db.Exec(query, appID); err != nil {
		return fmt.Errorf("Error while deleting app: %w", err)
	}
	return nil
}

func (c *Conn) GetAppByName(appName string) (*models.App, error) {
	const query = "select * from apps where name=?"
	app := &models.App{}
	err := c.db.Get(app, query, appName)
	if err == nil {
		return app, nil
	}
	if errors.Is(err, sql.ErrNoRows) {
		return nil, shared.ErrAppNotFound
	}
	return nil, fmt.Errorf("Error while getting app by name: %w", err)
}

func (c *Conn) AppNameIsTaken(appName string) bool {
  return false
}
