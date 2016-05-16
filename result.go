package avatica

import "errors"

type result struct {
	affectedRows int64
	insertId     int64
}

// LastInsertId returns the database's auto-generated ID
// after, for example, an INSERT into a table with primary
// key.
func (r *result) LastInsertId() (int64, error) {
	return 0, errors.New("Use 'SELECT CURRENT VALUE FOR your.sequence' to get the last inserted id. For more information, see: https://phoenix.apache.org/sequences.html.")
}

// RowsAffected returns the number of rows affected by the
// query.
func (r *result) RowsAffected() (int64, error) {
	return r.affectedRows, nil
}
