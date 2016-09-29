package avatica

import (
	"github.com/Boostport/avatica/message"
	"golang.org/x/net/context"
)

type tx struct {
	conn *conn
}

// Commit commits a transaction
func (t *tx) Commit() error {

	defer t.enableAutoCommit()

	_, err := t.conn.httpClient.post(context.Background(), &message.CommitRequest{
		ConnectionId: t.conn.connectionId,
	})

	return err
}

// Rollback rolls back a transaction
func (t *tx) Rollback() error {

	defer t.enableAutoCommit()

	_, err := t.conn.httpClient.post(context.Background(), &message.RollbackRequest{
		ConnectionId: t.conn.connectionId,
	})

	return err
}

// enableAutoCommit enables auto-commit on the server
func (t *tx) enableAutoCommit() error {

	_, err := t.conn.httpClient.post(context.Background(), &message.ConnectionSyncRequest{
		ConnectionId: t.conn.connectionId,
		ConnProps: &message.ConnectionProperties{
			AutoCommit:           true,
			HasAutoCommit:        true,
			TransactionIsolation: t.conn.config.transactionIsolation,
		},
	})

	return err
}
