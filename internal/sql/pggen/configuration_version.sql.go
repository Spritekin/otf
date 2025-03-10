// Code generated by pggen. DO NOT EDIT.

package pggen

import (
	"context"
	"fmt"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v4"
)

const insertConfigurationVersionSQL = `INSERT INTO configuration_versions (
    configuration_version_id,
    created_at,
    auto_queue_runs,
    source,
    speculative,
    status,
    workspace_id
) VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6,
    $7
);`

type InsertConfigurationVersionParams struct {
	ID            pgtype.Text
	CreatedAt     pgtype.Timestamptz
	AutoQueueRuns bool
	Source        pgtype.Text
	Speculative   bool
	Status        pgtype.Text
	WorkspaceID   pgtype.Text
}

// InsertConfigurationVersion implements Querier.InsertConfigurationVersion.
func (q *DBQuerier) InsertConfigurationVersion(ctx context.Context, params InsertConfigurationVersionParams) (pgconn.CommandTag, error) {
	ctx = context.WithValue(ctx, "pggen_query_name", "InsertConfigurationVersion")
	cmdTag, err := q.conn.Exec(ctx, insertConfigurationVersionSQL, params.ID, params.CreatedAt, params.AutoQueueRuns, params.Source, params.Speculative, params.Status, params.WorkspaceID)
	if err != nil {
		return cmdTag, fmt.Errorf("exec query InsertConfigurationVersion: %w", err)
	}
	return cmdTag, err
}

// InsertConfigurationVersionBatch implements Querier.InsertConfigurationVersionBatch.
func (q *DBQuerier) InsertConfigurationVersionBatch(batch genericBatch, params InsertConfigurationVersionParams) {
	batch.Queue(insertConfigurationVersionSQL, params.ID, params.CreatedAt, params.AutoQueueRuns, params.Source, params.Speculative, params.Status, params.WorkspaceID)
}

// InsertConfigurationVersionScan implements Querier.InsertConfigurationVersionScan.
func (q *DBQuerier) InsertConfigurationVersionScan(results pgx.BatchResults) (pgconn.CommandTag, error) {
	cmdTag, err := results.Exec()
	if err != nil {
		return cmdTag, fmt.Errorf("exec InsertConfigurationVersionBatch: %w", err)
	}
	return cmdTag, err
}

const insertConfigurationVersionStatusTimestampSQL = `INSERT INTO configuration_version_status_timestamps (
    configuration_version_id,
    status,
    timestamp
) VALUES (
    $1,
    $2,
    $3
)
RETURNING *;`

type InsertConfigurationVersionStatusTimestampParams struct {
	ID        pgtype.Text
	Status    pgtype.Text
	Timestamp pgtype.Timestamptz
}

type InsertConfigurationVersionStatusTimestampRow struct {
	ConfigurationVersionID pgtype.Text        `json:"configuration_version_id"`
	Status                 pgtype.Text        `json:"status"`
	Timestamp              pgtype.Timestamptz `json:"timestamp"`
}

// InsertConfigurationVersionStatusTimestamp implements Querier.InsertConfigurationVersionStatusTimestamp.
func (q *DBQuerier) InsertConfigurationVersionStatusTimestamp(ctx context.Context, params InsertConfigurationVersionStatusTimestampParams) (InsertConfigurationVersionStatusTimestampRow, error) {
	ctx = context.WithValue(ctx, "pggen_query_name", "InsertConfigurationVersionStatusTimestamp")
	row := q.conn.QueryRow(ctx, insertConfigurationVersionStatusTimestampSQL, params.ID, params.Status, params.Timestamp)
	var item InsertConfigurationVersionStatusTimestampRow
	if err := row.Scan(&item.ConfigurationVersionID, &item.Status, &item.Timestamp); err != nil {
		return item, fmt.Errorf("query InsertConfigurationVersionStatusTimestamp: %w", err)
	}
	return item, nil
}

// InsertConfigurationVersionStatusTimestampBatch implements Querier.InsertConfigurationVersionStatusTimestampBatch.
func (q *DBQuerier) InsertConfigurationVersionStatusTimestampBatch(batch genericBatch, params InsertConfigurationVersionStatusTimestampParams) {
	batch.Queue(insertConfigurationVersionStatusTimestampSQL, params.ID, params.Status, params.Timestamp)
}

// InsertConfigurationVersionStatusTimestampScan implements Querier.InsertConfigurationVersionStatusTimestampScan.
func (q *DBQuerier) InsertConfigurationVersionStatusTimestampScan(results pgx.BatchResults) (InsertConfigurationVersionStatusTimestampRow, error) {
	row := results.QueryRow()
	var item InsertConfigurationVersionStatusTimestampRow
	if err := row.Scan(&item.ConfigurationVersionID, &item.Status, &item.Timestamp); err != nil {
		return item, fmt.Errorf("scan InsertConfigurationVersionStatusTimestampBatch row: %w", err)
	}
	return item, nil
}

const findConfigurationVersionsByWorkspaceIDSQL = `SELECT
    configuration_versions.configuration_version_id,
    configuration_versions.created_at,
    configuration_versions.auto_queue_runs,
    configuration_versions.source,
    configuration_versions.speculative,
    configuration_versions.status,
    configuration_versions.workspace_id,
    (
        SELECT array_agg(t.*) AS configuration_version_status_timestamps
        FROM configuration_version_status_timestamps t
        WHERE t.configuration_version_id = configuration_versions.configuration_version_id
        GROUP BY configuration_version_id
    ) AS configuration_version_status_timestamps,
    (ingress_attributes.*)::"ingress_attributes"
FROM configuration_versions
JOIN workspaces USING (workspace_id)
LEFT JOIN ingress_attributes USING (configuration_version_id)
WHERE workspaces.workspace_id = $1
LIMIT $2
OFFSET $3;`

type FindConfigurationVersionsByWorkspaceIDParams struct {
	WorkspaceID pgtype.Text
	Limit       int
	Offset      int
}

type FindConfigurationVersionsByWorkspaceIDRow struct {
	ConfigurationVersionID               pgtype.Text                            `json:"configuration_version_id"`
	CreatedAt                            pgtype.Timestamptz                     `json:"created_at"`
	AutoQueueRuns                        bool                                   `json:"auto_queue_runs"`
	Source                               pgtype.Text                            `json:"source"`
	Speculative                          bool                                   `json:"speculative"`
	Status                               pgtype.Text                            `json:"status"`
	WorkspaceID                          pgtype.Text                            `json:"workspace_id"`
	ConfigurationVersionStatusTimestamps []ConfigurationVersionStatusTimestamps `json:"configuration_version_status_timestamps"`
	IngressAttributes                    *IngressAttributes                     `json:"ingress_attributes"`
}

// FindConfigurationVersionsByWorkspaceID implements Querier.FindConfigurationVersionsByWorkspaceID.
func (q *DBQuerier) FindConfigurationVersionsByWorkspaceID(ctx context.Context, params FindConfigurationVersionsByWorkspaceIDParams) ([]FindConfigurationVersionsByWorkspaceIDRow, error) {
	ctx = context.WithValue(ctx, "pggen_query_name", "FindConfigurationVersionsByWorkspaceID")
	rows, err := q.conn.Query(ctx, findConfigurationVersionsByWorkspaceIDSQL, params.WorkspaceID, params.Limit, params.Offset)
	if err != nil {
		return nil, fmt.Errorf("query FindConfigurationVersionsByWorkspaceID: %w", err)
	}
	defer rows.Close()
	items := []FindConfigurationVersionsByWorkspaceIDRow{}
	configurationVersionStatusTimestampsArray := q.types.newConfigurationVersionStatusTimestampsArray()
	ingressAttributesRow := q.types.newIngressAttributes()
	for rows.Next() {
		var item FindConfigurationVersionsByWorkspaceIDRow
		if err := rows.Scan(&item.ConfigurationVersionID, &item.CreatedAt, &item.AutoQueueRuns, &item.Source, &item.Speculative, &item.Status, &item.WorkspaceID, configurationVersionStatusTimestampsArray, ingressAttributesRow); err != nil {
			return nil, fmt.Errorf("scan FindConfigurationVersionsByWorkspaceID row: %w", err)
		}
		if err := configurationVersionStatusTimestampsArray.AssignTo(&item.ConfigurationVersionStatusTimestamps); err != nil {
			return nil, fmt.Errorf("assign FindConfigurationVersionsByWorkspaceID row: %w", err)
		}
		if err := ingressAttributesRow.AssignTo(&item.IngressAttributes); err != nil {
			return nil, fmt.Errorf("assign FindConfigurationVersionsByWorkspaceID row: %w", err)
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("close FindConfigurationVersionsByWorkspaceID rows: %w", err)
	}
	return items, err
}

// FindConfigurationVersionsByWorkspaceIDBatch implements Querier.FindConfigurationVersionsByWorkspaceIDBatch.
func (q *DBQuerier) FindConfigurationVersionsByWorkspaceIDBatch(batch genericBatch, params FindConfigurationVersionsByWorkspaceIDParams) {
	batch.Queue(findConfigurationVersionsByWorkspaceIDSQL, params.WorkspaceID, params.Limit, params.Offset)
}

// FindConfigurationVersionsByWorkspaceIDScan implements Querier.FindConfigurationVersionsByWorkspaceIDScan.
func (q *DBQuerier) FindConfigurationVersionsByWorkspaceIDScan(results pgx.BatchResults) ([]FindConfigurationVersionsByWorkspaceIDRow, error) {
	rows, err := results.Query()
	if err != nil {
		return nil, fmt.Errorf("query FindConfigurationVersionsByWorkspaceIDBatch: %w", err)
	}
	defer rows.Close()
	items := []FindConfigurationVersionsByWorkspaceIDRow{}
	configurationVersionStatusTimestampsArray := q.types.newConfigurationVersionStatusTimestampsArray()
	ingressAttributesRow := q.types.newIngressAttributes()
	for rows.Next() {
		var item FindConfigurationVersionsByWorkspaceIDRow
		if err := rows.Scan(&item.ConfigurationVersionID, &item.CreatedAt, &item.AutoQueueRuns, &item.Source, &item.Speculative, &item.Status, &item.WorkspaceID, configurationVersionStatusTimestampsArray, ingressAttributesRow); err != nil {
			return nil, fmt.Errorf("scan FindConfigurationVersionsByWorkspaceIDBatch row: %w", err)
		}
		if err := configurationVersionStatusTimestampsArray.AssignTo(&item.ConfigurationVersionStatusTimestamps); err != nil {
			return nil, fmt.Errorf("assign FindConfigurationVersionsByWorkspaceID row: %w", err)
		}
		if err := ingressAttributesRow.AssignTo(&item.IngressAttributes); err != nil {
			return nil, fmt.Errorf("assign FindConfigurationVersionsByWorkspaceID row: %w", err)
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("close FindConfigurationVersionsByWorkspaceIDBatch rows: %w", err)
	}
	return items, err
}

const countConfigurationVersionsByWorkspaceIDSQL = `SELECT count(*)
FROM configuration_versions
WHERE configuration_versions.workspace_id = $1
;`

// CountConfigurationVersionsByWorkspaceID implements Querier.CountConfigurationVersionsByWorkspaceID.
func (q *DBQuerier) CountConfigurationVersionsByWorkspaceID(ctx context.Context, workspaceID pgtype.Text) (int, error) {
	ctx = context.WithValue(ctx, "pggen_query_name", "CountConfigurationVersionsByWorkspaceID")
	row := q.conn.QueryRow(ctx, countConfigurationVersionsByWorkspaceIDSQL, workspaceID)
	var item int
	if err := row.Scan(&item); err != nil {
		return item, fmt.Errorf("query CountConfigurationVersionsByWorkspaceID: %w", err)
	}
	return item, nil
}

// CountConfigurationVersionsByWorkspaceIDBatch implements Querier.CountConfigurationVersionsByWorkspaceIDBatch.
func (q *DBQuerier) CountConfigurationVersionsByWorkspaceIDBatch(batch genericBatch, workspaceID pgtype.Text) {
	batch.Queue(countConfigurationVersionsByWorkspaceIDSQL, workspaceID)
}

// CountConfigurationVersionsByWorkspaceIDScan implements Querier.CountConfigurationVersionsByWorkspaceIDScan.
func (q *DBQuerier) CountConfigurationVersionsByWorkspaceIDScan(results pgx.BatchResults) (int, error) {
	row := results.QueryRow()
	var item int
	if err := row.Scan(&item); err != nil {
		return item, fmt.Errorf("scan CountConfigurationVersionsByWorkspaceIDBatch row: %w", err)
	}
	return item, nil
}

const findConfigurationVersionByIDSQL = `SELECT
    configuration_versions.configuration_version_id,
    configuration_versions.created_at,
    configuration_versions.auto_queue_runs,
    configuration_versions.source,
    configuration_versions.speculative,
    configuration_versions.status,
    configuration_versions.workspace_id,
    (
        SELECT array_agg(t.*) AS configuration_version_status_timestamps
        FROM configuration_version_status_timestamps t
        WHERE t.configuration_version_id = configuration_versions.configuration_version_id
        GROUP BY configuration_version_id
    ) AS configuration_version_status_timestamps,
    (ingress_attributes.*)::"ingress_attributes"
FROM configuration_versions
JOIN workspaces USING (workspace_id)
LEFT JOIN ingress_attributes USING (configuration_version_id)
WHERE configuration_version_id = $1;`

type FindConfigurationVersionByIDRow struct {
	ConfigurationVersionID               pgtype.Text                            `json:"configuration_version_id"`
	CreatedAt                            pgtype.Timestamptz                     `json:"created_at"`
	AutoQueueRuns                        bool                                   `json:"auto_queue_runs"`
	Source                               pgtype.Text                            `json:"source"`
	Speculative                          bool                                   `json:"speculative"`
	Status                               pgtype.Text                            `json:"status"`
	WorkspaceID                          pgtype.Text                            `json:"workspace_id"`
	ConfigurationVersionStatusTimestamps []ConfigurationVersionStatusTimestamps `json:"configuration_version_status_timestamps"`
	IngressAttributes                    *IngressAttributes                     `json:"ingress_attributes"`
}

// FindConfigurationVersionByID implements Querier.FindConfigurationVersionByID.
func (q *DBQuerier) FindConfigurationVersionByID(ctx context.Context, configurationVersionID pgtype.Text) (FindConfigurationVersionByIDRow, error) {
	ctx = context.WithValue(ctx, "pggen_query_name", "FindConfigurationVersionByID")
	row := q.conn.QueryRow(ctx, findConfigurationVersionByIDSQL, configurationVersionID)
	var item FindConfigurationVersionByIDRow
	configurationVersionStatusTimestampsArray := q.types.newConfigurationVersionStatusTimestampsArray()
	ingressAttributesRow := q.types.newIngressAttributes()
	if err := row.Scan(&item.ConfigurationVersionID, &item.CreatedAt, &item.AutoQueueRuns, &item.Source, &item.Speculative, &item.Status, &item.WorkspaceID, configurationVersionStatusTimestampsArray, ingressAttributesRow); err != nil {
		return item, fmt.Errorf("query FindConfigurationVersionByID: %w", err)
	}
	if err := configurationVersionStatusTimestampsArray.AssignTo(&item.ConfigurationVersionStatusTimestamps); err != nil {
		return item, fmt.Errorf("assign FindConfigurationVersionByID row: %w", err)
	}
	if err := ingressAttributesRow.AssignTo(&item.IngressAttributes); err != nil {
		return item, fmt.Errorf("assign FindConfigurationVersionByID row: %w", err)
	}
	return item, nil
}

// FindConfigurationVersionByIDBatch implements Querier.FindConfigurationVersionByIDBatch.
func (q *DBQuerier) FindConfigurationVersionByIDBatch(batch genericBatch, configurationVersionID pgtype.Text) {
	batch.Queue(findConfigurationVersionByIDSQL, configurationVersionID)
}

// FindConfigurationVersionByIDScan implements Querier.FindConfigurationVersionByIDScan.
func (q *DBQuerier) FindConfigurationVersionByIDScan(results pgx.BatchResults) (FindConfigurationVersionByIDRow, error) {
	row := results.QueryRow()
	var item FindConfigurationVersionByIDRow
	configurationVersionStatusTimestampsArray := q.types.newConfigurationVersionStatusTimestampsArray()
	ingressAttributesRow := q.types.newIngressAttributes()
	if err := row.Scan(&item.ConfigurationVersionID, &item.CreatedAt, &item.AutoQueueRuns, &item.Source, &item.Speculative, &item.Status, &item.WorkspaceID, configurationVersionStatusTimestampsArray, ingressAttributesRow); err != nil {
		return item, fmt.Errorf("scan FindConfigurationVersionByIDBatch row: %w", err)
	}
	if err := configurationVersionStatusTimestampsArray.AssignTo(&item.ConfigurationVersionStatusTimestamps); err != nil {
		return item, fmt.Errorf("assign FindConfigurationVersionByID row: %w", err)
	}
	if err := ingressAttributesRow.AssignTo(&item.IngressAttributes); err != nil {
		return item, fmt.Errorf("assign FindConfigurationVersionByID row: %w", err)
	}
	return item, nil
}

const findConfigurationVersionLatestByWorkspaceIDSQL = `SELECT
    configuration_versions.configuration_version_id,
    configuration_versions.created_at,
    configuration_versions.auto_queue_runs,
    configuration_versions.source,
    configuration_versions.speculative,
    configuration_versions.status,
    configuration_versions.workspace_id,
    (
        SELECT array_agg(t.*) AS configuration_version_status_timestamps
        FROM configuration_version_status_timestamps t
        WHERE t.configuration_version_id = configuration_versions.configuration_version_id
        GROUP BY configuration_version_id
    ) AS configuration_version_status_timestamps,
    (ingress_attributes.*)::"ingress_attributes"
FROM configuration_versions
JOIN workspaces USING (workspace_id)
LEFT JOIN ingress_attributes USING (configuration_version_id)
WHERE workspace_id = $1
ORDER BY configuration_versions.created_at DESC;`

type FindConfigurationVersionLatestByWorkspaceIDRow struct {
	ConfigurationVersionID               pgtype.Text                            `json:"configuration_version_id"`
	CreatedAt                            pgtype.Timestamptz                     `json:"created_at"`
	AutoQueueRuns                        bool                                   `json:"auto_queue_runs"`
	Source                               pgtype.Text                            `json:"source"`
	Speculative                          bool                                   `json:"speculative"`
	Status                               pgtype.Text                            `json:"status"`
	WorkspaceID                          pgtype.Text                            `json:"workspace_id"`
	ConfigurationVersionStatusTimestamps []ConfigurationVersionStatusTimestamps `json:"configuration_version_status_timestamps"`
	IngressAttributes                    *IngressAttributes                     `json:"ingress_attributes"`
}

// FindConfigurationVersionLatestByWorkspaceID implements Querier.FindConfigurationVersionLatestByWorkspaceID.
func (q *DBQuerier) FindConfigurationVersionLatestByWorkspaceID(ctx context.Context, workspaceID pgtype.Text) (FindConfigurationVersionLatestByWorkspaceIDRow, error) {
	ctx = context.WithValue(ctx, "pggen_query_name", "FindConfigurationVersionLatestByWorkspaceID")
	row := q.conn.QueryRow(ctx, findConfigurationVersionLatestByWorkspaceIDSQL, workspaceID)
	var item FindConfigurationVersionLatestByWorkspaceIDRow
	configurationVersionStatusTimestampsArray := q.types.newConfigurationVersionStatusTimestampsArray()
	ingressAttributesRow := q.types.newIngressAttributes()
	if err := row.Scan(&item.ConfigurationVersionID, &item.CreatedAt, &item.AutoQueueRuns, &item.Source, &item.Speculative, &item.Status, &item.WorkspaceID, configurationVersionStatusTimestampsArray, ingressAttributesRow); err != nil {
		return item, fmt.Errorf("query FindConfigurationVersionLatestByWorkspaceID: %w", err)
	}
	if err := configurationVersionStatusTimestampsArray.AssignTo(&item.ConfigurationVersionStatusTimestamps); err != nil {
		return item, fmt.Errorf("assign FindConfigurationVersionLatestByWorkspaceID row: %w", err)
	}
	if err := ingressAttributesRow.AssignTo(&item.IngressAttributes); err != nil {
		return item, fmt.Errorf("assign FindConfigurationVersionLatestByWorkspaceID row: %w", err)
	}
	return item, nil
}

// FindConfigurationVersionLatestByWorkspaceIDBatch implements Querier.FindConfigurationVersionLatestByWorkspaceIDBatch.
func (q *DBQuerier) FindConfigurationVersionLatestByWorkspaceIDBatch(batch genericBatch, workspaceID pgtype.Text) {
	batch.Queue(findConfigurationVersionLatestByWorkspaceIDSQL, workspaceID)
}

// FindConfigurationVersionLatestByWorkspaceIDScan implements Querier.FindConfigurationVersionLatestByWorkspaceIDScan.
func (q *DBQuerier) FindConfigurationVersionLatestByWorkspaceIDScan(results pgx.BatchResults) (FindConfigurationVersionLatestByWorkspaceIDRow, error) {
	row := results.QueryRow()
	var item FindConfigurationVersionLatestByWorkspaceIDRow
	configurationVersionStatusTimestampsArray := q.types.newConfigurationVersionStatusTimestampsArray()
	ingressAttributesRow := q.types.newIngressAttributes()
	if err := row.Scan(&item.ConfigurationVersionID, &item.CreatedAt, &item.AutoQueueRuns, &item.Source, &item.Speculative, &item.Status, &item.WorkspaceID, configurationVersionStatusTimestampsArray, ingressAttributesRow); err != nil {
		return item, fmt.Errorf("scan FindConfigurationVersionLatestByWorkspaceIDBatch row: %w", err)
	}
	if err := configurationVersionStatusTimestampsArray.AssignTo(&item.ConfigurationVersionStatusTimestamps); err != nil {
		return item, fmt.Errorf("assign FindConfigurationVersionLatestByWorkspaceID row: %w", err)
	}
	if err := ingressAttributesRow.AssignTo(&item.IngressAttributes); err != nil {
		return item, fmt.Errorf("assign FindConfigurationVersionLatestByWorkspaceID row: %w", err)
	}
	return item, nil
}

const findConfigurationVersionByIDForUpdateSQL = `SELECT
    configuration_versions.configuration_version_id,
    configuration_versions.created_at,
    configuration_versions.auto_queue_runs,
    configuration_versions.source,
    configuration_versions.speculative,
    configuration_versions.status,
    configuration_versions.workspace_id,
    (
        SELECT array_agg(t.*) AS configuration_version_status_timestamps
        FROM configuration_version_status_timestamps t
        WHERE t.configuration_version_id = configuration_versions.configuration_version_id
        GROUP BY configuration_version_id
    ) AS configuration_version_status_timestamps,
    (ingress_attributes.*)::"ingress_attributes"
FROM configuration_versions
JOIN workspaces USING (workspace_id)
LEFT JOIN ingress_attributes USING (configuration_version_id)
WHERE configuration_version_id = $1
FOR UPDATE OF configuration_versions;`

type FindConfigurationVersionByIDForUpdateRow struct {
	ConfigurationVersionID               pgtype.Text                            `json:"configuration_version_id"`
	CreatedAt                            pgtype.Timestamptz                     `json:"created_at"`
	AutoQueueRuns                        bool                                   `json:"auto_queue_runs"`
	Source                               pgtype.Text                            `json:"source"`
	Speculative                          bool                                   `json:"speculative"`
	Status                               pgtype.Text                            `json:"status"`
	WorkspaceID                          pgtype.Text                            `json:"workspace_id"`
	ConfigurationVersionStatusTimestamps []ConfigurationVersionStatusTimestamps `json:"configuration_version_status_timestamps"`
	IngressAttributes                    *IngressAttributes                     `json:"ingress_attributes"`
}

// FindConfigurationVersionByIDForUpdate implements Querier.FindConfigurationVersionByIDForUpdate.
func (q *DBQuerier) FindConfigurationVersionByIDForUpdate(ctx context.Context, configurationVersionID pgtype.Text) (FindConfigurationVersionByIDForUpdateRow, error) {
	ctx = context.WithValue(ctx, "pggen_query_name", "FindConfigurationVersionByIDForUpdate")
	row := q.conn.QueryRow(ctx, findConfigurationVersionByIDForUpdateSQL, configurationVersionID)
	var item FindConfigurationVersionByIDForUpdateRow
	configurationVersionStatusTimestampsArray := q.types.newConfigurationVersionStatusTimestampsArray()
	ingressAttributesRow := q.types.newIngressAttributes()
	if err := row.Scan(&item.ConfigurationVersionID, &item.CreatedAt, &item.AutoQueueRuns, &item.Source, &item.Speculative, &item.Status, &item.WorkspaceID, configurationVersionStatusTimestampsArray, ingressAttributesRow); err != nil {
		return item, fmt.Errorf("query FindConfigurationVersionByIDForUpdate: %w", err)
	}
	if err := configurationVersionStatusTimestampsArray.AssignTo(&item.ConfigurationVersionStatusTimestamps); err != nil {
		return item, fmt.Errorf("assign FindConfigurationVersionByIDForUpdate row: %w", err)
	}
	if err := ingressAttributesRow.AssignTo(&item.IngressAttributes); err != nil {
		return item, fmt.Errorf("assign FindConfigurationVersionByIDForUpdate row: %w", err)
	}
	return item, nil
}

// FindConfigurationVersionByIDForUpdateBatch implements Querier.FindConfigurationVersionByIDForUpdateBatch.
func (q *DBQuerier) FindConfigurationVersionByIDForUpdateBatch(batch genericBatch, configurationVersionID pgtype.Text) {
	batch.Queue(findConfigurationVersionByIDForUpdateSQL, configurationVersionID)
}

// FindConfigurationVersionByIDForUpdateScan implements Querier.FindConfigurationVersionByIDForUpdateScan.
func (q *DBQuerier) FindConfigurationVersionByIDForUpdateScan(results pgx.BatchResults) (FindConfigurationVersionByIDForUpdateRow, error) {
	row := results.QueryRow()
	var item FindConfigurationVersionByIDForUpdateRow
	configurationVersionStatusTimestampsArray := q.types.newConfigurationVersionStatusTimestampsArray()
	ingressAttributesRow := q.types.newIngressAttributes()
	if err := row.Scan(&item.ConfigurationVersionID, &item.CreatedAt, &item.AutoQueueRuns, &item.Source, &item.Speculative, &item.Status, &item.WorkspaceID, configurationVersionStatusTimestampsArray, ingressAttributesRow); err != nil {
		return item, fmt.Errorf("scan FindConfigurationVersionByIDForUpdateBatch row: %w", err)
	}
	if err := configurationVersionStatusTimestampsArray.AssignTo(&item.ConfigurationVersionStatusTimestamps); err != nil {
		return item, fmt.Errorf("assign FindConfigurationVersionByIDForUpdate row: %w", err)
	}
	if err := ingressAttributesRow.AssignTo(&item.IngressAttributes); err != nil {
		return item, fmt.Errorf("assign FindConfigurationVersionByIDForUpdate row: %w", err)
	}
	return item, nil
}

const downloadConfigurationVersionSQL = `SELECT config
FROM configuration_versions
WHERE configuration_version_id = $1;`

// DownloadConfigurationVersion implements Querier.DownloadConfigurationVersion.
func (q *DBQuerier) DownloadConfigurationVersion(ctx context.Context, configurationVersionID pgtype.Text) ([]byte, error) {
	ctx = context.WithValue(ctx, "pggen_query_name", "DownloadConfigurationVersion")
	row := q.conn.QueryRow(ctx, downloadConfigurationVersionSQL, configurationVersionID)
	item := []byte{}
	if err := row.Scan(&item); err != nil {
		return item, fmt.Errorf("query DownloadConfigurationVersion: %w", err)
	}
	return item, nil
}

// DownloadConfigurationVersionBatch implements Querier.DownloadConfigurationVersionBatch.
func (q *DBQuerier) DownloadConfigurationVersionBatch(batch genericBatch, configurationVersionID pgtype.Text) {
	batch.Queue(downloadConfigurationVersionSQL, configurationVersionID)
}

// DownloadConfigurationVersionScan implements Querier.DownloadConfigurationVersionScan.
func (q *DBQuerier) DownloadConfigurationVersionScan(results pgx.BatchResults) ([]byte, error) {
	row := results.QueryRow()
	item := []byte{}
	if err := row.Scan(&item); err != nil {
		return item, fmt.Errorf("scan DownloadConfigurationVersionBatch row: %w", err)
	}
	return item, nil
}

const updateConfigurationVersionErroredByIDSQL = `UPDATE configuration_versions
SET
    status = 'errored'
WHERE configuration_version_id = $1
RETURNING configuration_version_id;`

// UpdateConfigurationVersionErroredByID implements Querier.UpdateConfigurationVersionErroredByID.
func (q *DBQuerier) UpdateConfigurationVersionErroredByID(ctx context.Context, id pgtype.Text) (pgtype.Text, error) {
	ctx = context.WithValue(ctx, "pggen_query_name", "UpdateConfigurationVersionErroredByID")
	row := q.conn.QueryRow(ctx, updateConfigurationVersionErroredByIDSQL, id)
	var item pgtype.Text
	if err := row.Scan(&item); err != nil {
		return item, fmt.Errorf("query UpdateConfigurationVersionErroredByID: %w", err)
	}
	return item, nil
}

// UpdateConfigurationVersionErroredByIDBatch implements Querier.UpdateConfigurationVersionErroredByIDBatch.
func (q *DBQuerier) UpdateConfigurationVersionErroredByIDBatch(batch genericBatch, id pgtype.Text) {
	batch.Queue(updateConfigurationVersionErroredByIDSQL, id)
}

// UpdateConfigurationVersionErroredByIDScan implements Querier.UpdateConfigurationVersionErroredByIDScan.
func (q *DBQuerier) UpdateConfigurationVersionErroredByIDScan(results pgx.BatchResults) (pgtype.Text, error) {
	row := results.QueryRow()
	var item pgtype.Text
	if err := row.Scan(&item); err != nil {
		return item, fmt.Errorf("scan UpdateConfigurationVersionErroredByIDBatch row: %w", err)
	}
	return item, nil
}

const updateConfigurationVersionConfigByIDSQL = `UPDATE configuration_versions
SET
    config = $1,
    status = 'uploaded'
WHERE configuration_version_id = $2
RETURNING configuration_version_id;`

// UpdateConfigurationVersionConfigByID implements Querier.UpdateConfigurationVersionConfigByID.
func (q *DBQuerier) UpdateConfigurationVersionConfigByID(ctx context.Context, config []byte, id pgtype.Text) (pgtype.Text, error) {
	ctx = context.WithValue(ctx, "pggen_query_name", "UpdateConfigurationVersionConfigByID")
	row := q.conn.QueryRow(ctx, updateConfigurationVersionConfigByIDSQL, config, id)
	var item pgtype.Text
	if err := row.Scan(&item); err != nil {
		return item, fmt.Errorf("query UpdateConfigurationVersionConfigByID: %w", err)
	}
	return item, nil
}

// UpdateConfigurationVersionConfigByIDBatch implements Querier.UpdateConfigurationVersionConfigByIDBatch.
func (q *DBQuerier) UpdateConfigurationVersionConfigByIDBatch(batch genericBatch, config []byte, id pgtype.Text) {
	batch.Queue(updateConfigurationVersionConfigByIDSQL, config, id)
}

// UpdateConfigurationVersionConfigByIDScan implements Querier.UpdateConfigurationVersionConfigByIDScan.
func (q *DBQuerier) UpdateConfigurationVersionConfigByIDScan(results pgx.BatchResults) (pgtype.Text, error) {
	row := results.QueryRow()
	var item pgtype.Text
	if err := row.Scan(&item); err != nil {
		return item, fmt.Errorf("scan UpdateConfigurationVersionConfigByIDBatch row: %w", err)
	}
	return item, nil
}

const deleteConfigurationVersionByIDSQL = `DELETE
FROM configuration_versions
WHERE configuration_version_id = $1
RETURNING configuration_version_id;`

// DeleteConfigurationVersionByID implements Querier.DeleteConfigurationVersionByID.
func (q *DBQuerier) DeleteConfigurationVersionByID(ctx context.Context, id pgtype.Text) (pgtype.Text, error) {
	ctx = context.WithValue(ctx, "pggen_query_name", "DeleteConfigurationVersionByID")
	row := q.conn.QueryRow(ctx, deleteConfigurationVersionByIDSQL, id)
	var item pgtype.Text
	if err := row.Scan(&item); err != nil {
		return item, fmt.Errorf("query DeleteConfigurationVersionByID: %w", err)
	}
	return item, nil
}

// DeleteConfigurationVersionByIDBatch implements Querier.DeleteConfigurationVersionByIDBatch.
func (q *DBQuerier) DeleteConfigurationVersionByIDBatch(batch genericBatch, id pgtype.Text) {
	batch.Queue(deleteConfigurationVersionByIDSQL, id)
}

// DeleteConfigurationVersionByIDScan implements Querier.DeleteConfigurationVersionByIDScan.
func (q *DBQuerier) DeleteConfigurationVersionByIDScan(results pgx.BatchResults) (pgtype.Text, error) {
	row := results.QueryRow()
	var item pgtype.Text
	if err := row.Scan(&item); err != nil {
		return item, fmt.Errorf("scan DeleteConfigurationVersionByIDBatch row: %w", err)
	}
	return item, nil
}
