package service

import (
	"database/sql"
	"time"

	"github.com/mabzd/snorlax/api"
	"github.com/mabzd/snorlax/internal/utils"
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
		Version:   utils.ToNullInt64(dto.Version),
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
	dst.InBedAt = utils.FromNullTime(src.InBedAt)
	dst.TriedToSleepAt = src.TriedToSleepAt
	dst.SleepDelayInMin = utils.FromNullInt32(src.SleepDelayInMin)
	dst.AwakeningsCount = utils.FromNullInt32(src.AwakeningsCount)
	dst.AwakeningsTotalDurationInMin = utils.FromNullInt32(src.AwakeningsTotalDurationInMin)
	dst.FinalWakeUpAt = src.FinalWakeUpAt
	dst.OutOfBedAt = utils.FromNullTime(src.OutOfBedAt)
	dst.SleepQuality = src.SleepQuality
	dst.Comments = utils.FromNullString(src.Comments)
}

func assignDtoToEntry(src api.SleepDiaryEntryDataDto, dst *SleepDiaryEntry) {
	dst.InBedAt = utils.ToNullTime(src.InBedAt)
	dst.TriedToSleepAt = src.TriedToSleepAt
	dst.SleepDelayInMin = utils.ToNullInt32(src.SleepDelayInMin)
	dst.AwakeningsCount = utils.ToNullInt32(src.AwakeningsCount)
	dst.AwakeningsTotalDurationInMin = utils.ToNullInt32(src.AwakeningsTotalDurationInMin)
	dst.FinalWakeUpAt = src.FinalWakeUpAt
	dst.OutOfBedAt = utils.ToNullTime(src.OutOfBedAt)
	dst.SleepQuality = src.SleepQuality
	dst.Comments = utils.ToNullString(src.Comments)
}
