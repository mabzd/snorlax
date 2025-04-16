package service

import (
	"database/sql"
	"log"

	"github.com/mabzd/snorlax/api"
	"github.com/mabzd/snorlax/internal/config"
	"github.com/mabzd/snorlax/internal/database"
)

type SleepDiaryService struct {
	db *sql.DB
}

func NewSleepDiaryService(cfg config.Config) *SleepDiaryService {
	db, err := database.InitDB(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	return &SleepDiaryService{
		db: db,
	}
}

func (s *SleepDiaryService) GetEntryById(id int64) (api.SleepDiaryEntryDto, api.Error) {
	entry, err := getSleepDiaryEntryById(s.db, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return api.SleepDiaryEntryDto{}, api.NewError("entry not found", api.ERR_NOT_FOUND)
		}
		log.Printf("Reading entry by ID %d failed: %v\n", id, err)
		return api.SleepDiaryEntryDto{}, api.NewError("read failed", api.ERR_UNKNOWN)
	}

	dto, err := toSleepDiaryEntryDto(entry)
	if err != nil {
		log.Printf("Converting entry %d to DTO failed: %v\n", id, err)
		return api.SleepDiaryEntryDto{}, api.NewError("conversion failed", api.ERR_UNKNOWN)
	}

	return dto, nil
}

func (s *SleepDiaryService) GetEntriesByFilter(filter api.SleepDiaryFilterDto) (api.PageDto[api.SleepDiaryEntryDto], api.Error) {
	errs := filter.Validate()
	if len(errs) > 0 {
		return api.PageDto[api.SleepDiaryEntryDto]{}, api.NewValidationError("invalid filter data", errs)
	}

	count, err := countSleepDiaryEntriesByFilter(s.db, filter)
	if err != nil {
		log.Printf("Counting entries by filter %v failed: %v\n", filter, err)
		return api.PageDto[api.SleepDiaryEntryDto]{}, api.NewError("count failed", api.ERR_UNKNOWN)
	}

	entries, err := getSleepDiaryEntriesByFilter(s.db, filter)
	if err != nil {
		log.Printf("Reading entries by filter %v failed: %v\n", filter, err)
		return api.PageDto[api.SleepDiaryEntryDto]{}, api.NewError("read failed", api.ERR_UNKNOWN)
	}

	items := make([]api.SleepDiaryEntryDto, len(entries))
	for i, entry := range entries {
		dto, err := toSleepDiaryEntryDto(entry)
		if err != nil {
			log.Printf("Converting entry %d to DTO failed: %v\n", entry.Id, err)
			return api.PageDto[api.SleepDiaryEntryDto]{}, api.NewError("conversion failed", api.ERR_UNKNOWN)
		}
		items[i] = dto
	}

	return api.PageDto[api.SleepDiaryEntryDto]{
		TotalCount: count,
		PageSize:   filter.PageSize,
		PageNumber: filter.PageNumber,
		Items:      items,
	}, nil
}

func (s *SleepDiaryService) CreateEntry(dto api.CreateSleepDiaryEntryDto) (api.SleepDiaryEntryDto, api.Error) {
	errs := dto.Validate()
	if len(errs) > 0 {
		return api.SleepDiaryEntryDto{}, api.NewValidationError("invalid create data", errs)
	}

	entry := fromCreateSleepDiaryEntryDto(dto)
	createdEntry, err := insertSleepDiaryEntry(s.db, entry)
	if err != nil {
		log.Printf("Inserting entry %v failed: %v\n", dto, err)
		return api.SleepDiaryEntryDto{}, api.NewError("insert failed", api.ERR_UNKNOWN)
	}

	createdDto, err := toSleepDiaryEntryDto(createdEntry)
	if err != nil {
		log.Printf("Converting entry %d to DTO failed: %v\n", createdEntry.Id, err)
		return api.SleepDiaryEntryDto{}, api.NewError("conversion failed", api.ERR_UNKNOWN)
	}

	return createdDto, nil
}

func (s *SleepDiaryService) UpdateEntry(id int64, dto api.UpdateSleepDiaryEntryDto) (api.SleepDiaryEntryDto, api.Error) {
	errs := dto.Validate()
	if len(errs) > 0 {
		return api.SleepDiaryEntryDto{}, api.NewValidationError("invalid update data", errs)
	}

	entry := fromUpdateSleepDiaryEntryDto(dto)
	entry.Id = id
	updatedEntry, err := updateSleepDiaryEntry(s.db, entry)
	if err != nil {
		if err == ErrConflict {
			return api.SleepDiaryEntryDto{}, api.NewError("version conflict", api.ERR_CONFLICT)
		}
		if err == sql.ErrNoRows {
			return api.SleepDiaryEntryDto{}, api.NewError("entry not found", api.ERR_NOT_FOUND)
		}
		log.Printf("Updating entry %v failed: %v", dto, err)
		return api.SleepDiaryEntryDto{}, api.NewError("update failed", api.ERR_UNKNOWN)
	}

	updatedDto, err := toSleepDiaryEntryDto(updatedEntry)
	if err != nil {
		log.Printf("Converting entry %d to DTO failed: %v\n", updatedDto.Id, err)
		return api.SleepDiaryEntryDto{}, api.NewError("conversion failed", api.ERR_UNKNOWN)
	}

	return updatedDto, nil
}
