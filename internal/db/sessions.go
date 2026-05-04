package db

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/joseph0x45/tessera/internal/models"
	"github.com/joseph0x45/tessera/internal/shared"
)

func (c *Conn) InsertSession(session *models.Session, tx *sqlx.Tx) error {
	const query = `
    insert into sessions (
      id, session_user_id
    )
    values (
      :id, :session_user_id
    );
  `
	var err error
	if tx != nil {
		_, err = tx.NamedExec(query, session)
	} else {
		_, err = c.db.NamedExec(query, session)
	}
	if err != nil {
		return fmt.Errorf("Error while inserting session: %w", err)
	}
	return nil
}

func (c *Conn) DeleteSession(sessionID string) error {
	const query = "delete from sessions where id=?"
	if _, err := c.db.Exec(query, sessionID); err != nil {
		return fmt.Errorf("Error while deleting session: %w", err)
	}
	return nil
}

func (c *Conn) GetSessionByID(sessionID string) (*models.Session, error) {
	const query = "select * from sessions where id=?"
	session := &models.Session{}
	err := c.db.Get(session, query, sessionID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, shared.ErrSessionNotFound
		}
		return nil, fmt.Errorf("Error while getting session by id: %w", err)
	}
	return session, nil
}
