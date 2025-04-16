package tests

import (
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"github.com/mabzd/snorlax/api"
)

func TestCreateInvalidAccountUuid(t *testing.T) {
	dto := api.CreateSleepDiaryEntryDto{
		AccountUuid:            "invalid",
		SleepDiaryEntryDataDto: newMinimalRandomEntryData(),
	}

	createResp := mustPost(t, "/sleep_diary/entries", dto)
	defer createResp.Body.Close()
	assertHttpStatusCode(t, http.StatusBadRequest, createResp)
}

func TestCreateAndGetMinimalSleepDiary(t *testing.T) {
	runCreateAndGetSleepDiary(
		t,
		api.CreateSleepDiaryEntryDto{
			AccountUuid:            uuid.NewString(),
			SleepDiaryEntryDataDto: newMinimalRandomEntryData(),
		})
}

func TestCreateAndGetFullSleepDiary(t *testing.T) {
	runCreateAndGetSleepDiary(
		t,
		api.CreateSleepDiaryEntryDto{
			AccountUuid:            uuid.NewString(),
			SleepDiaryEntryDataDto: newRandomEntryData(),
		})
}

func TestCreateTimesGetConvertedToTargetTimezone(t *testing.T) {
	dto := api.CreateSleepDiaryEntryDto{
		AccountUuid:            uuid.NewString(),
		SleepDiaryEntryDataDto: newRandomEntryData(),
	}

	tz, err := time.LoadLocation("Europe/Warsaw")
	if err != nil {
		t.Fatal("Could not load timezone")
	}

	dto.Timezone = tz.String()
	createdDto := mustCreateEntry(t, dto)

	assertValuesEqualTimeMsPrec(t, toPtr(dto.InBedAt.In(tz)), createdDto.InBedAt, "InBedAt")
	assertValuesEqualTimeMsPrec(t, toPtr(dto.TriedToSleepAt.In(tz)), &createdDto.TriedToSleepAt, "TriedToSleepAt")
	assertValuesEqualTimeMsPrec(t, toPtr(dto.FinalWakeUpAt.In(tz)), &createdDto.FinalWakeUpAt, "FinalWakeUpAt")
	assertValuesEqualTimeMsPrec(t, toPtr(dto.OutOfBedAt.In(tz)), createdDto.OutOfBedAt, "OutOfBedAt")
}

func runCreateAndGetSleepDiary(t *testing.T, dto api.CreateSleepDiaryEntryDto) {
	createResp := mustPost(t, "/sleep_diary/entries", dto)
	defer createResp.Body.Close()
	assertHttpStatusCode(t, http.StatusCreated, createResp)
	createdEntry := mustDecode[api.SleepDiaryEntryDto](createResp.Body)
	retrievedEntry := mustGetEntryById(t, createdEntry.Id)
	assertEqualEntryDto(t, createdEntry, retrievedEntry, true)
}
