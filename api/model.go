package api

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

const DEFAULT_PAGE_SIZE = int64(100)
const MAX_PAGE_SIZE = int64(1000)
const MAX_COMMENT_LENGTH = 2048

type SleepQuality int

const (
	VeryPoorSleepQuality  SleepQuality = 1
	PoorSleepQuality      SleepQuality = 2
	AverageSleepQuality   SleepQuality = 3
	GoodSleepQuality      SleepQuality = 4
	ExcellentSleepQuality SleepQuality = 5
)

func (sq SleepQuality) String() string {
	switch sq {
	case VeryPoorSleepQuality:
		return "Very Poor"
	case PoorSleepQuality:
		return "Poor"
	case AverageSleepQuality:
		return "Average"
	case GoodSleepQuality:
		return "Good"
	case ExcellentSleepQuality:
		return "Excellent"
	default:
		return "Unknown"
	}
}

type SleepDiaryEntryDataDto struct {
	InBedAt                      *time.Time   `json:"get_in_bed_at,omitempty"`
	TriedToSleepAt               time.Time    `json:"tried_to_sleep_at"`
	SleepDelayInMin              *int         `json:"sleep_delay_in_min,omitempty"`
	AwakeningsCount              *int         `json:"awakenings_count,omitempty"`
	AwakeningsTotalDurationInMin *int         `json:"awakenings_total_duration_in_min,omitempty"`
	FinalWakeUpAt                time.Time    `json:"final_wake_up_at"`
	OutOfBedAt                   *time.Time   `json:"out_of_bed_at,omitempty"`
	SleepQuality                 SleepQuality `json:"sleep_quality"`
	Comments                     *string      `json:"comments,omitempty"`
}

func (dto *SleepDiaryEntryDataDto) Validate() []error {
	errors := validateTimeOrder(
		labeledTime{dto.InBedAt, "in_bed_at"},
		labeledTime{&dto.TriedToSleepAt, "tried_to_sleep_at"},
		labeledTime{&dto.FinalWakeUpAt, "final_wake_up_at"},
		labeledTime{dto.OutOfBedAt, "out_of_bed_at"},
	)
	if dto.TriedToSleepAt.IsZero() {
		errors = append(errors, fmt.Errorf("tried_to_sleep_at is required"))
	}
	if dto.FinalWakeUpAt.IsZero() {
		errors = append(errors, fmt.Errorf("final_wake_up_at is required"))
	}
	if dto.SleepQuality < VeryPoorSleepQuality || dto.SleepQuality > ExcellentSleepQuality {
		errors = append(errors, fmt.Errorf("sleep_quality should be between %d and %d", VeryPoorSleepQuality, ExcellentSleepQuality))
	}
	if dto.SleepDelayInMin != nil && *dto.SleepDelayInMin < 0 {
		errors = append(errors, fmt.Errorf("sleep_delay_in_min should be non-negative"))
	}
	if dto.AwakeningsCount != nil && *dto.AwakeningsCount < 0 {
		errors = append(errors, fmt.Errorf("awakenings_count should be non-negative"))
	}
	if dto.AwakeningsTotalDurationInMin != nil && *dto.AwakeningsTotalDurationInMin < 0 {
		errors = append(errors, fmt.Errorf("awakenings_total_duration_in_min should be non-negative"))
	}
	if dto.Comments != nil && len(*dto.Comments) > MAX_COMMENT_LENGTH {
		errors = append(errors, fmt.Errorf("comments should not exceed %d characters", MAX_COMMENT_LENGTH))
	}
	return errors
}

type SleepDiaryFilterDto struct {
	AccountUuids []string   `json:"account_uuids"`
	FromDate     *time.Time `json:"from_date,omitempty"`
	ToDate       *time.Time `json:"to_date,omitempty"`
	PageSize     int64      `json:"page_size"`
	PageNumber   int64      `json:"page_number"`
}

func (dto *SleepDiaryFilterDto) Validate() []error {
	errors := validateTimeOrder(
		labeledTime{dto.FromDate, "from_date"},
		labeledTime{dto.ToDate, "to_date"},
	)
	if len(dto.AccountUuids) == 0 {
		errors = append(errors, fmt.Errorf("account_uuids is required"))
	}
	for _, id := range dto.AccountUuids {
		if _, err := uuid.Parse(id); err != nil {
			errors = append(errors, fmt.Errorf("invalid UUID '%s'", id))
		}
	}
	if dto.PageSize < 1 {
		errors = append(errors, fmt.Errorf("page_size should be greater than 0"))
	}
	if dto.PageSize > MAX_PAGE_SIZE {
		errors = append(errors, fmt.Errorf("page_size should not exceed %d", MAX_PAGE_SIZE))
	}
	if dto.PageNumber < 1 {
		errors = append(errors, fmt.Errorf("page_number should be greater than 0"))
	}
	return errors
}

type UpdateSleepDiaryEntryDto struct {
	Version *int64 `json:"version,omitempty"`
	SleepDiaryEntryDataDto
}

func (dto *UpdateSleepDiaryEntryDto) Validate() []error {
	errors := dto.SleepDiaryEntryDataDto.Validate()
	return errors
}

type CreateSleepDiaryEntryDto struct {
	AccountUuid string `json:"account_uuid"`
	SleepDiaryEntryDataDto
}

func (dto *CreateSleepDiaryEntryDto) Validate() []error {
	errors := dto.SleepDiaryEntryDataDto.Validate()
	if dto.AccountUuid == "" {
		errors = append(errors, fmt.Errorf("account_uuid is required"))
	}
	return errors
}

type SleepDiaryEntryDto struct {
	Id          int64  `json:"id"`
	AccountUuid string `json:"account_uuid"`
	Version     int64  `json:"version"`
	SleepDiaryEntryDataDto
}

func (dto *SleepDiaryEntryDto) Validate() []error {
	errors := dto.SleepDiaryEntryDataDto.Validate()
	return errors
}

type PageDto[T any] struct {
	TotalCount int64 `json:"total_count"`
	PageSize   int64 `json:"page_size"`
	PageNumber int64 `json:"page_number"`
	Items      []T   `json:"items"`
}

type labeledTime struct {
	time  *time.Time
	label string
}

func validateTimeOrder(items ...labeledTime) []error {
	var prev *labeledTime
	errors := []error{}

	for _, i := range items {
		if i.time == nil {
			continue
		}
		if prev != nil && i.time.Before(*prev.time) {
			errors = append(errors, fmt.Errorf("%s should be before %s", prev.label, i.label))
		}
		prev = &i
	}

	return errors
}
