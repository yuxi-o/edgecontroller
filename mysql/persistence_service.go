package mysql

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"reflect"

	_ "github.com/go-sql-driver/mysql" // provides the mysql driver
	"github.com/pkg/errors"
	cce "github.com/smartedgemec/controller-ce"
)

// PersistenceService implements cce.PersistenceService.
type PersistenceService struct {
	DB *sql.DB
}

// Create persists a resource.
func (s *PersistenceService) Create(
	ctx context.Context,
	e cce.Entity,
) error {
	bytes, err := json.Marshal(e)
	if err != nil {
		return errors.Wrap(err, "error marshalling")
	}

	_, err = s.DB.ExecContext(
		ctx,
		fmt.Sprintf( //nolint:gosec
			`INSERT INTO %s (entity) VALUES (?)`, e.GetTableName()),
		bytes)
	if err != nil {
		return errors.Wrap(err, "error inserting record")
	}

	return nil
}

// Read retrieves a single resource of the given type by ID.
func (s *PersistenceService) Read(
	ctx context.Context,
	id string,
	zv cce.EntityModel,
) (e cce.Entity, err error) {
	rows, err := s.DB.QueryContext(
		ctx,
		fmt.Sprintf( //nolint:gosec
			`SELECT entity
             FROM %s
             WHERE id = ?`, zv.GetTableName()),
		id)
	if err != nil {
		return nil, errors.Wrap(err, "error running query")
	}

	if !rows.Next() {
		return nil, nil
	}

	if e, err = s.scan(rows, zv); err != nil {
		return nil, err
	}

	return e, nil
}

// Filter retrieves a collection of resources of the given type using a set of
// filters.
func (s *PersistenceService) Filter(
	ctx context.Context,
	zv cce.EntityModel,
	fs []cce.Filter,
) (es []cce.Entity, err error) {
	q := fmt.Sprintf("SELECT entity FROM %s", zv.GetTableName()) //nolint:gosec
	if len(fs) > 0 {
		q += " WHERE "
		for i, f := range fs {
			q += fmt.Sprintf("%s = '%s'", f.Field, f.Value)
			if i < len(fs)-1 {
				q += " AND "
			}
		}
	}

	rows, err := s.DB.QueryContext(
		ctx, q)
	if err != nil {
		return nil, errors.Wrap(err, "error running query")
	}

	for rows.Next() {
		e, err := s.scan(rows, zv)
		if err != nil {
			return nil, err
		}

		es = append(es, e)
	}

	return
}

// ReadAll retrieves all resources of the given type.
func (s *PersistenceService) ReadAll(
	ctx context.Context,
	zv cce.EntityModel,
) (es []cce.Entity, err error) {
	rows, err := s.DB.QueryContext(
		ctx,
		fmt.Sprintf( //nolint:gosec
			"SELECT entity FROM %s", zv.GetTableName()))
	if err != nil {
		return nil, errors.Wrap(err, "error running query")
	}

	for rows.Next() {
		e, err := s.scan(rows, zv)
		if err != nil {
			return nil, err
		}

		es = append(es, e)
	}

	return
}

func (s *PersistenceService) scan(
	rows *sql.Rows,
	zv cce.EntityModel,
) (cce.Entity, error) {
	var bytes []byte
	if err := rows.Scan(&bytes); err != nil {
		return nil, errors.Wrap(err, "error scanning row")
	}

	e := reflect.New(
		reflect.ValueOf(zv).Elem().Type()).Interface().(cce.Entity)
	if err := json.Unmarshal(bytes, e); err != nil {
		return nil, errors.Wrap(err, "error unmarshalling")
	}

	return e, nil
}

// BulkUpdate updates multiple resources.
func (s *PersistenceService) BulkUpdate(
	ctx context.Context,
	es []cce.Entity,
) error {
	for _, e := range es {
		bytes, err := json.Marshal(e)
		if err != nil {
			return errors.Wrap(err, "error marshalling")
		}

		_, err = s.DB.ExecContext(
			ctx,
			fmt.Sprintf( //nolint:gosec
				`UPDATE %s
                 SET entity = ?
                 WHERE id = JSON_EXTRACT(?, "$.id")`,
				e.GetTableName()),
			bytes, bytes)
		if err != nil {
			return errors.Wrap(err, "error updating record")
		}
	}

	return nil
}

// Delete deletes a resource of the given type.
func (s *PersistenceService) Delete(
	ctx context.Context,
	id string,
	zv cce.EntityModel,
) (ok bool, err error) {
	result, err := s.DB.ExecContext(
		ctx,
		fmt.Sprintf( //nolint:gosec
			`DELETE
             FROM %s
             WHERE id = ?`, zv.GetTableName()),
		id)
	if err != nil {
		return false, errors.Wrap(err, "error deleting record")
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return false, errors.Wrap(err, "error getting rows affected")
	}

	if rows != 1 {
		return false, nil
	}

	return true, nil
}
