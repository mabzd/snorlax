package tests

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/mabzd/snorlax/api"
	"github.com/stretchr/testify/assert"
)

func TestUpdateEntry(t *testing.T) {
	createdEntry := mustCreateRandomEntry(t)

	updateDto := api.UpdateSleepDiaryEntryDto{
		Version:                &createdEntry.Version,
		SleepDiaryEntryDataDto: newRandomEntryData(),
	}

	updateResp := mustPut(t, fmt.Sprintf("/sleep_diary/entries/%v", createdEntry.Id), updateDto)
	defer updateResp.Body.Close()
	assertHttpStatusCode(t, http.StatusOK, updateResp)
	updatedEntry := mustDecode[api.SleepDiaryEntryDto](updateResp.Body)

	assert.Equal(t, createdEntry.Version+1, updatedEntry.Version)
	retrievedEntry := mustGetEntryById(t, createdEntry.Id)
	assertEqualEntryDto(t, retrievedEntry, updatedEntry, false)
}

func TestUpdateMinimalEntry(t *testing.T) {
	createdEntry := mustCreateRandomEntry(t)

	updateDto := api.UpdateSleepDiaryEntryDto{
		Version:                &createdEntry.Version,
		SleepDiaryEntryDataDto: newMinimalRandomEntryData(),
	}

	updateResp := mustPut(t, fmt.Sprintf("/sleep_diary/entries/%v", createdEntry.Id), updateDto)
	defer updateResp.Body.Close()
	assertHttpStatusCode(t, http.StatusOK, updateResp)
	updatedEntry := mustDecode[api.SleepDiaryEntryDto](updateResp.Body)

	assert.Equal(t, createdEntry.Version+1, updatedEntry.Version)
	retrievedEntry := mustGetEntryById(t, createdEntry.Id)
	assertEqualEntryDto(t, retrievedEntry, updatedEntry, false)
}

func TestUpdateNonExistingEntry(t *testing.T) {
	updateDto := api.UpdateSleepDiaryEntryDto{
		SleepDiaryEntryDataDto: newRandomEntryData(),
	}

	updateResp := mustPut(t, fmt.Sprintf("/sleep_diary/entries/%v", 99999999999999999), updateDto)
	defer updateResp.Body.Close()
	assertHttpStatusCode(t, http.StatusNotFound, updateResp)
}

func TestUpdateWithWrongVersion(t *testing.T) {
	createdEntry := mustCreateRandomEntry(t)

	updateDto := api.UpdateSleepDiaryEntryDto{
		Version:                toPtr(createdEntry.Version + 1),
		SleepDiaryEntryDataDto: newRandomEntryData(),
	}

	updateResp := mustPut(t, fmt.Sprintf("/sleep_diary/entries/%v", createdEntry.Id), updateDto)
	defer updateResp.Body.Close()
	assertHttpStatusCode(t, http.StatusConflict, updateResp)
}

func TestUpdateWithoutVersion(t *testing.T) {
	createdEntry := mustCreateRandomEntry(t)

	updateDto := api.UpdateSleepDiaryEntryDto{
		SleepDiaryEntryDataDto: newRandomEntryData(),
	}

	updateResp := mustPut(t, fmt.Sprintf("/sleep_diary/entries/%v", createdEntry.Id), updateDto)
	defer updateResp.Body.Close()
	assertHttpStatusCode(t, http.StatusOK, updateResp)
	updatedEntry := mustDecode[api.SleepDiaryEntryDto](updateResp.Body)

	assert.Equal(t, createdEntry.Version+1, updatedEntry.Version)
	retrievedEntry := mustGetEntryById(t, createdEntry.Id)
	assertEqualEntryDto(t, retrievedEntry, updatedEntry, false)
}
