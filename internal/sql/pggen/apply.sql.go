// Code generated by pggen. DO NOT EDIT.

package pggen

import (
	"context"
	"fmt"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v4"
)

const insertApplySQL = `INSERT INTO applies (
    run_id,
    status
) VALUES (
    $1,
    $2
);`

// InsertApply implements Querier.InsertApply.
func (q *DBQuerier) InsertApply(ctx context.Context, runID pgtype.Text, status pgtype.Text) (pgconn.CommandTag, error) {
	ctx = context.WithValue(ctx, "pggen_query_name", "InsertApply")
	cmdTag, err := q.conn.Exec(ctx, insertApplySQL, runID, status)
	if err != nil {
		return cmdTag, fmt.Errorf("exec query InsertApply: %w", err)
	}
	return cmdTag, err
}

// InsertApplyBatch implements Querier.InsertApplyBatch.
func (q *DBQuerier) InsertApplyBatch(batch genericBatch, runID pgtype.Text, status pgtype.Text) {
	batch.Queue(insertApplySQL, runID, status)
}

// InsertApplyScan implements Querier.InsertApplyScan.
func (q *DBQuerier) InsertApplyScan(results pgx.BatchResults) (pgconn.CommandTag, error) {
	cmdTag, err := results.Exec()
	if err != nil {
		return cmdTag, fmt.Errorf("exec InsertApplyBatch: %w", err)
	}
	return cmdTag, err
}

const updateAppliedChangesByIDSQL = `UPDATE applies
SET report = (
    $1,
    $2,
    $3
)
WHERE run_id = $4
RETURNING run_id
;`

type UpdateAppliedChangesByIDParams struct {
	Additions    int
	Changes      int
	Destructions int
	RunID        pgtype.Text
}

// UpdateAppliedChangesByID implements Querier.UpdateAppliedChangesByID.
func (q *DBQuerier) UpdateAppliedChangesByID(ctx context.Context, params UpdateAppliedChangesByIDParams) (pgtype.Text, error) {
	ctx = context.WithValue(ctx, "pggen_query_name", "UpdateAppliedChangesByID")
	row := q.conn.QueryRow(ctx, updateAppliedChangesByIDSQL, params.Additions, params.Changes, params.Destructions, params.RunID)
	var item pgtype.Text
	if err := row.Scan(&item); err != nil {
		return item, fmt.Errorf("query UpdateAppliedChangesByID: %w", err)
	}
	return item, nil
}

// UpdateAppliedChangesByIDBatch implements Querier.UpdateAppliedChangesByIDBatch.
func (q *DBQuerier) UpdateAppliedChangesByIDBatch(batch genericBatch, params UpdateAppliedChangesByIDParams) {
	batch.Queue(updateAppliedChangesByIDSQL, params.Additions, params.Changes, params.Destructions, params.RunID)
}

// UpdateAppliedChangesByIDScan implements Querier.UpdateAppliedChangesByIDScan.
func (q *DBQuerier) UpdateAppliedChangesByIDScan(results pgx.BatchResults) (pgtype.Text, error) {
	row := results.QueryRow()
	var item pgtype.Text
	if err := row.Scan(&item); err != nil {
		return item, fmt.Errorf("scan UpdateAppliedChangesByIDBatch row: %w", err)
	}
	return item, nil
}

const updateApplyStatusByIDSQL = `UPDATE applies
SET status = $1
WHERE run_id = $2
RETURNING run_id
;`

// UpdateApplyStatusByID implements Querier.UpdateApplyStatusByID.
func (q *DBQuerier) UpdateApplyStatusByID(ctx context.Context, status pgtype.Text, runID pgtype.Text) (pgtype.Text, error) {
	ctx = context.WithValue(ctx, "pggen_query_name", "UpdateApplyStatusByID")
	row := q.conn.QueryRow(ctx, updateApplyStatusByIDSQL, status, runID)
	var item pgtype.Text
	if err := row.Scan(&item); err != nil {
		return item, fmt.Errorf("query UpdateApplyStatusByID: %w", err)
	}
	return item, nil
}

// UpdateApplyStatusByIDBatch implements Querier.UpdateApplyStatusByIDBatch.
func (q *DBQuerier) UpdateApplyStatusByIDBatch(batch genericBatch, status pgtype.Text, runID pgtype.Text) {
	batch.Queue(updateApplyStatusByIDSQL, status, runID)
}

// UpdateApplyStatusByIDScan implements Querier.UpdateApplyStatusByIDScan.
func (q *DBQuerier) UpdateApplyStatusByIDScan(results pgx.BatchResults) (pgtype.Text, error) {
	row := results.QueryRow()
	var item pgtype.Text
	if err := row.Scan(&item); err != nil {
		return item, fmt.Errorf("scan UpdateApplyStatusByIDBatch row: %w", err)
	}
	return item, nil
}
