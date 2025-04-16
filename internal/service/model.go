package service

import (
	"database/sql"
	"time"

	"github.com/mabzd/snorlax/api"
)

type SleepDiaryEntry struct {
	Id                           int64
	AccountUuid                  string
	Timezone                     string
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

func fromCreateSleepDiaryEntryDto(dto api.CreateSleepDiaryEntryDto) SleepDiaryEntry {
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

func toSleepDiaryEntryDto(entry SleepDiaryEntry) (api.SleepDiaryEntryDto, error) {
	dto := api.SleepDiaryEntryDto{
		Id:          entry.Id,
		AccountUuid: entry.AccountUuid,
		Version:     entry.Version.Int64,
	}
	err := assignEntryToDto(entry, &dto.SleepDiaryEntryDataDto)
	return dto, err
}

func assignEntryToDto(src SleepDiaryEntry, dst *api.SleepDiaryEntryDataDto) error {
	tz, err := time.LoadLocation(src.Timezone)
	if err != nil {
		return err
	}

	dst.Timezone = src.Timezone
	dst.InBedAt = fromNullTime(src.InBedAt, tz)
	dst.TriedToSleepAt = src.TriedToSleepAt.In(tz)
	dst.SleepDelayInMin = fromNullInt32(src.SleepDelayInMin)
	dst.AwakeningsCount = fromNullInt32(src.AwakeningsCount)
	dst.AwakeningsTotalDurationInMin = fromNullInt32(src.AwakeningsTotalDurationInMin)
	dst.FinalWakeUpAt = src.FinalWakeUpAt.In(tz)
	dst.OutOfBedAt = fromNullTime(src.OutOfBedAt, tz)
	dst.SleepQuality = src.SleepQuality
	dst.Comments = fromNullString(src.Comments)
	return nil
}

func assignDtoToEntry(src api.SleepDiaryEntryDataDto, dst *SleepDiaryEntry) {
	dst.Timezone = src.Timezone
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
		return sql.NullTime{Time: *t, Valid: true}
	}
	return sql.NullTime{}
}

func fromNullTime(nt sql.NullTime, tz *time.Location) *time.Time {
	if nt.Valid {
		t := nt.Time.In(tz)
		return &t
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
