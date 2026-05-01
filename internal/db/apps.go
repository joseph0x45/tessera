package db

import (
	"fmt"

	"github.com/joseph0x45/tessera/internal/models"
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
