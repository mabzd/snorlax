package tests

import (
	"net/http"
	"testing"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"github.com/mabzd/snorlax/api"
)

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

func runCreateAndGetSleepDiary(t *testing.T, dto api.CreateSleepDiaryEntryDto) {
	createResp := mustPost(t, "/sleep_diary/entries", dto)
	defer createResp.Body.Close()
	assertHttpStatusCode(t, http.StatusCreated, createResp)
	createdEntry := mustDecode[api.SleepDiaryEntryDto](createResp.Body)
	retrievedEntry := mustGetEntryById(t, createdEntry.Id)
	assertEqualEntryDto(t, createdEntry, retrievedEntry, true)
}
