package tests

import (
	"net/http"
	"testing"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"github.com/mabzd/snorlax/api"
)

func TestAddAndGetMinimalSleepDiary(t *testing.T) {
	runAddAndGetSleepDiary(
		t,
		api.AddSleepDiaryEntryDto{
			AccountUuid:            uuid.NewString(),
			SleepDiaryEntryDataDto: prepareMinimalSleepDiaryData(),
		})
}

func TestAddAndGetFullSleepDiary(t *testing.T) {
	runAddAndGetSleepDiary(
		t,
		api.AddSleepDiaryEntryDto{
			AccountUuid:            uuid.NewString(),
			SleepDiaryEntryDataDto: prepareSleepDiaryData(),
		})
}

func runAddAndGetSleepDiary(t *testing.T, addDto api.AddSleepDiaryEntryDto) {
	addResp := mustPost(t, "/sleep_diary/entries", addDto)
	defer addResp.Body.Close()
	assertHttpStatusCode(t, http.StatusCreated, addResp)
	createdEntry := mustDecode[api.SleepDiaryEntryDto](addResp.Body)
	retrievedEntry := mustGetEntryById(t, createdEntry.Id)
	assertEqualSleepDiaryEntryDto(t, createdEntry, retrievedEntry, true)
}
