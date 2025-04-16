package service

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/mabzd/snorlax/api"
)

var ErrConflict = errors.New("sql: conflict")

func getSleepDiaryEntryById(db *sql.DB, id int64) (SleepDiaryEntry, error) {
	query := `
		SELECT *
		FROM sleep_diary_entries
		WHERE id = $1
	`
	row := db.QueryRow(query, id)

	var entry SleepDiaryEntry
	err := row.Scan(
		&entry.Id,
		&entry.AccountUuid,
		&entry.Timezone,
		&entry.InBedAt,
		&entry.TriedToSleepAt,
		&entry.SleepDelayInMin,
		&entry.AwakeningsCount,
		&entry.AwakeningsTotalDurationInMin,
		&entry.FinalWakeUpAt,
		&entry.OutOfBedAt,
		&entry.SleepQuality,
		&entry.Comments,
		&entry.CreatedAt,
		&entry.UpdatedAt,
		&entry.Version,
	)
	return entry, err
}

func getSleepDiaryEntriesByFilter(db *sql.DB, filter api.SleepDiaryFilterDto) ([]SleepDiaryEntry, error) {
	whereClause, args := buildWhereClause(filter)
	limitClause := buildLimitClause(filter)

	query := fmt.Sprintf(
		"SELECT * FROM sleep_diary_entries %s ORDER BY tried_to_sleep_at %s",
		whereClause,
		limitClause)

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	capacity := filter.PageSize
	entries := make([]SleepDiaryEntry, 0, capacity)

	for rows.Next() {
		var entry SleepDiaryEntry
		err := rows.Scan(
			&entry.Id,
			&entry.AccountUuid,
			&entry.Timezone,
			&entry.InBedAt,
			&entry.TriedToSleepAt,
			&entry.SleepDelayInMin,
			&entry.AwakeningsCount,
			&entry.AwakeningsTotalDurationInMin,
			&entry.FinalWakeUpAt,
			&entry.OutOfBedAt,
			&entry.SleepQuality,
			&entry.Comments,
			&entry.CreatedAt,
			&entry.UpdatedAt,
			&entry.Version,
		)
		if err != nil {
			return nil, err
		}
		entries = append(entries, entry)
	}

	return entries, rows.Err()
}

func countSleepDiaryEntriesByFilter(db *sql.DB, filter api.SleepDiaryFilterDto) (int64, error) {
	whereClause, args := buildWhereClause(filter)
	query := fmt.Sprintf("SELECT count(*) FROM sleep_diary_entries %s", whereClause)
	row := db.QueryRow(query, args...)
	var count int64
	err := row.Scan(&count)
	return count, err
}

func buildWhereClause(filter api.SleepDiaryFilterDto) (string, []interface{}) {
	var whereClauses []string
	var args []interface{}
	argPos := 1

	if len(filter.AccountUuid) > 0 {
		placeholders := make([]string, len(filter.AccountUuid))
		for i, uuid := range filter.AccountUuid {
			placeholders[i] = fmt.Sprintf("$%d", argPos)
			args = append(args, uuid)
			argPos++
		}
		whereClauses = append(whereClauses, fmt.Sprintf("account_uuid IN (%s)", strings.Join(placeholders, ",")))
	}

	if filter.FromDate != nil {
		whereClauses = append(whereClauses, fmt.Sprintf("tried_to_sleep_at >= $%d", argPos))
		args = append(args, *filter.FromDate)
		argPos++
	}

	if filter.ToDate != nil {
		whereClauses = append(whereClauses, fmt.Sprintf("tried_to_sleep_at < $%d", argPos))
		args = append(args, *filter.ToDate)
		argPos++
	}

	sql := fmt.Sprintf(
		"WHERE %s",
		strings.Join(whereClauses, " AND "))

	return sql, args
}

func buildLimitClause(filter api.SleepDiaryFilterDto) string {
	return fmt.Sprintf(
		"LIMIT %d OFFSET %d",
		filter.PageSize,
		(filter.PageNumber-1)*filter.PageSize)
}

func insertSleepDiaryEntry(db *sql.DB, entry SleepDiaryEntry) (SleepDiaryEntry, error) {
	query := `
		INSERT INTO sleep_diary_entries (
			account_uuid,
			timezone, 
			in_bed_at, 
			tried_to_sleep_at, 
			sleep_delay_in_min, 
			awakenings_count, 
			awakenings_total_duration_in_min, 
			final_wake_up_at, 
			out_of_bed_at, 
			sleep_quality, 
			comments, 
			created_at, 
			updated_at, 
			version
		)
		VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14
		)
		RETURNING id
	`

	var id int64
	err := db.QueryRow(
		query,
		entry.AccountUuid,
		entry.Timezone,
		entry.InBedAt,
		entry.TriedToSleepAt,
		entry.SleepDelayInMin,
		entry.AwakeningsCount,
		entry.AwakeningsTotalDurationInMin,
		entry.FinalWakeUpAt,
		entry.OutOfBedAt,
		entry.SleepQuality,
		entry.Comments,
		entry.CreatedAt,
		entry.UpdatedAt,
		entry.Version,
	).Scan(&id)
	if err != nil {
		return SleepDiaryEntry{}, err
	}

	entry.Id = id
	return entry, nil
}

func updateSleepDiaryEntry(db *sql.DB, entry SleepDiaryEntry) (SleepDiaryEntry, error) {
	tx, err := db.Begin()
	if err != nil {
		return SleepDiaryEntry{}, err
	}
	defer func() {
		tx.Rollback()
	}()

	query := `
		UPDATE sleep_diary_entries
		SET 
			timezone = $1,
			in_bed_at = $2,
			tried_to_sleep_at = $3,
			sleep_delay_in_min = $4,
			awakenings_count = $5,
			awakenings_total_duration_in_min = $6,
			final_wake_up_at = $7,
			out_of_bed_at = $8,
			sleep_quality = $9,
			comments = $10,
			updated_at = $11,
			version = version + 1
		WHERE id = $12
		RETURNING version, account_uuid
	`
	var newVersion int64
	var accountUuid string
	err = tx.QueryRow(
		query,
		entry.Timezone,
		entry.InBedAt,
		entry.TriedToSleepAt,
		entry.SleepDelayInMin,
		entry.AwakeningsCount,
		entry.AwakeningsTotalDurationInMin,
		entry.FinalWakeUpAt,
		entry.OutOfBedAt,
		entry.SleepQuality,
		entry.Comments,
		entry.UpdatedAt,
		entry.Id,
	).Scan(&newVersion, &accountUuid)
	if err != nil {
		return SleepDiaryEntry{}, err
	}

	if entry.Version.Valid && entry.Version.Int64+1 != newVersion {
		return SleepDiaryEntry{}, ErrConflict
	}

	entry.Version = sql.NullInt64{Int64: newVersion, Valid: true}
	entry.AccountUuid = accountUuid
	return entry, tx.Commit()
}
