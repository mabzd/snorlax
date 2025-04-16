package utils

import (
	"database/sql"
	"time"
)

func ToNullTime(t *time.Time) sql.NullTime {
	if t != nil {
		return sql.NullTime{Time: t.UTC(), Valid: true}
	}
	return sql.NullTime{}
}

func FromNullTime(t sql.NullTime) *time.Time {
	if t.Valid {
		return &t.Time
	}
	return nil
}

func ToNullInt32(i *int) sql.NullInt32 {
	if i != nil {
		return sql.NullInt32{Int32: int32(*i), Valid: true}
	}
	return sql.NullInt32{}
}

func FromNullInt32(i sql.NullInt32) *int {
	if i.Valid {
		val := int(i.Int32)
		return &val
	}
	return nil
}

func ToNullInt64(i *int64) sql.NullInt64 {
	if i != nil {
		return sql.NullInt64{Int64: int64(*i), Valid: true}
	}
	return sql.NullInt64{}
}

func FromNullInt64(i sql.NullInt64) *int64 {
	if i.Valid {
		val := int64(i.Int64)
		return &val
	}
	return nil
}

func ToNullString(s *string) sql.NullString {
	if s != nil {
		return sql.NullString{String: *s, Valid: true}
	}
	return sql.NullString{}
}

func FromNullString(s sql.NullString) *string {
	if s.Valid {
		return &s.String
	}
	return nil
}
