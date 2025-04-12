package service

import (
	"database/sql"
	"time"

	"github.com/mabzd/snorlax/api"
)

type SleepDiaryEntry struct {
	Id                           int64
	AccountUuid                  string
	InBedAt                      sql.NullTime
	TriedToSleepAt               time.Time
	SleepDelayInMin              sql.NullInt32
	AwakeningsCount              sql.NullInt32
	AwakeningsTotalDurationInMin sql.NullInt32
	FinalWakeUpAt                time.Time
	OutOfBedAt                   sql.NullTime
	SleepQuality                 api.SleepQuality
	Comments                     sql.NullString
	CreatedAt                    time.Time
	UpdatedAt                    time.Time
	Version                      sql.NullInt64
}

func fromAddSleepDiaryEntryDto(dto api.AddSleepDiaryEntryDto) SleepDiaryEntry {
	entry := SleepDiaryEntry{
		AccountUuid: dto.AccountUuid,
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
		Version:     sql.NullInt64{Int64: 1, Valid: true},
	}
	assignDtoToEntry(dto.SleepDiaryEntryDataDto, &entry)
	return entry
}

func fromUpdateSleepDiaryEntryDto(dto api.UpdateSleepDiaryEntryDto) SleepDiaryEntry {
	entry := SleepDiaryEntry{
		UpdatedAt: time.Now().UTC(),
		Version:   toNullInt64(dto.Version),
	}
	assignDtoToEntry(dto.SleepDiaryEntryDataDto, &entry)
	return entry
}

func toSleepDiaryEntryDto(entry SleepDiaryEntry) api.SleepDiaryEntryDto {
	dto := api.SleepDiaryEntryDto{
		Id:          entry.Id,
		AccountUuid: entry.AccountUuid,
		Version:     entry.Version.Int64,
	}
	assignEntryToDto(entry, &dto.SleepDiaryEntryDataDto)
	return dto
}

func assignEntryToDto(src SleepDiaryEntry, dst *api.SleepDiaryEntryDataDto) {
	dst.InBedAt = fromNullTime(src.InBedAt)
	dst.TriedToSleepAt = src.TriedToSleepAt
	dst.SleepDelayInMin = fromNullInt32(src.SleepDelayInMin)
	dst.AwakeningsCount = fromNullInt32(src.AwakeningsCount)
	dst.AwakeningsTotalDurationInMin = fromNullInt32(src.AwakeningsTotalDurationInMin)
	dst.FinalWakeUpAt = src.FinalWakeUpAt
	dst.OutOfBedAt = fromNullTime(src.OutOfBedAt)
	dst.SleepQuality = src.SleepQuality
	dst.Comments = fromNullString(src.Comments)
}

func assignDtoToEntry(src api.SleepDiaryEntryDataDto, dst *SleepDiaryEntry) {
	dst.InBedAt = toNullTime(src.InBedAt)
	dst.TriedToSleepAt = src.TriedToSleepAt
	dst.SleepDelayInMin = toNullInt32(src.SleepDelayInMin)
	dst.AwakeningsCount = toNullInt32(src.AwakeningsCount)
	dst.AwakeningsTotalDurationInMin = toNullInt32(src.AwakeningsTotalDurationInMin)
	dst.FinalWakeUpAt = src.FinalWakeUpAt
	dst.OutOfBedAt = toNullTime(src.OutOfBedAt)
	dst.SleepQuality = src.SleepQuality
	dst.Comments = toNullString(src.Comments)
}

func toNullTime(t *time.Time) sql.NullTime {
	if t != nil {
		return sql.NullTime{Time: t.UTC(), Valid: true}
	}
	return sql.NullTime{}
}

func fromNullTime(t sql.NullTime) *time.Time {
	if t.Valid {
		return &t.Time
	}
	return nil
}

func toNullInt32(i *int) sql.NullInt32 {
	if i != nil {
		return sql.NullInt32{Int32: int32(*i), Valid: true}
	}
	return sql.NullInt32{}
}

func fromNullInt32(i sql.NullInt32) *int {
	if i.Valid {
		val := int(i.Int32)
		return &val
	}
	return nil
}

func toNullInt64(i *int64) sql.NullInt64 {
	if i != nil {
		return sql.NullInt64{Int64: int64(*i), Valid: true}
	}
	return sql.NullInt64{}
}

func fromNullInt64(i sql.NullInt64) *int64 {
	if i.Valid {
		val := int64(i.Int64)
		return &val
	}
	return nil
}

func toNullString(s *string) sql.NullString {
	if s != nil {
		return sql.NullString{String: *s, Valid: true}
	}
	return sql.NullString{}
}

func fromNullString(s sql.NullString) *string {
	if s.Valid {
		return &s.String
	}
	return nil
}

func toPtr[T any](v T) *T {
	return &v
}
