package db

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/joseph0x45/tessera/internal/models"
	"github.com/joseph0x45/tessera/internal/shared"
)

func (c *Conn) SetMetadata(metaData *models.MetaData) error {
	const query = `
    insert into app_metadata (
      key, value
    )
    values (
      :key, :value
    )
    on conflict(key)
    do update set value=excluded.value;
  `
	if _, err := c.db.NamedExec(query, metaData); err != nil {
		return fmt.Errorf("Error while inserting metadata: %w", err)
	}
	return nil
}

func (c *Conn) GetMetadata(key string) (*string, error) {
	const query = "select value from app_metadata where key=?"
	value := ""
	err := c.db.Get(&value, query, key)
	if err == nil {
		return &value, nil
	}
	if errors.Is(err, sql.ErrNoRows) {
		return nil, shared.ErrValueNotFound
	}
	return nil, fmt.Errorf("Error while getting metadata: %w", err)
}
