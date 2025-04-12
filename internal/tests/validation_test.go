package tests

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/mabzd/snorlax/api"
)

func TestEmptyData(t *testing.T) {
	data := api.SleepDiaryEntryDataDto{}
	runValidationTests(t, data)
}

func TestInBedAtAndTriedToSleepAtNotInOrder(t *testing.T) {
	data := prepareSleepDiaryData()
	swap(data.InBedAt, &data.TriedToSleepAt)
	runValidationTests(t, data)
}

func TestInBedAtAndFinalWakeUpAtNotInOrder(t *testing.T) {
	data := prepareSleepDiaryData()
	swap(data.InBedAt, &data.FinalWakeUpAt)
	runValidationTests(t, data)
}

func TestInBedAtAndOutOfBedAtNotInOrder(t *testing.T) {
	data := prepareSleepDiaryData()
	swap(data.InBedAt, &data.FinalWakeUpAt)
	runValidationTests(t, data)
}

func TestTriedToSleepAtAndFinalWakeUpAtNotInOrder(t *testing.T) {
	data := prepareSleepDiaryData()
	swap(&data.TriedToSleepAt, &data.FinalWakeUpAt)
	runValidationTests(t, data)
}

func TestTriedToSleepAtAndOutOfBedAtNotInOrder(t *testing.T) {
	data := prepareSleepDiaryData()
	swap(&data.TriedToSleepAt, data.OutOfBedAt)
	runValidationTests(t, data)
}

func TestFinalWakeUpAtAndOutOfBedAtNotInOrder(t *testing.T) {
	data := prepareSleepDiaryData()
	swap(&data.FinalWakeUpAt, data.OutOfBedAt)
	runValidationTests(t, data)
}

func TestSleepQualityLowerBound(t *testing.T) {
	data := prepareSleepDiaryData()
	data.SleepQuality = api.VeryPoorSleepQuality - 1
	runValidationTests(t, data)
}

func TestSleepQualityUpperBound(t *testing.T) {
	data := prepareSleepDiaryData()
	data.SleepQuality = api.ExcellentSleepQuality + 1
	runValidationTests(t, data)
}

func TestNegativeSleepDelay(t *testing.T) {
	data := prepareSleepDiaryData()
	data.SleepDelayInMin = toPtr(-1)
	runValidationTests(t, data)
}

func TestNegativeAwakeningsCount(t *testing.T) {
	data := prepareSleepDiaryData()
	data.AwakeningsCount = toPtr(-1)
	runValidationTests(t, data)
}

func TestNegativeAwakeningsTotalDuration(t *testing.T) {
	data := prepareSleepDiaryData()
	data.AwakeningsTotalDurationInMin = toPtr(-1)
	runValidationTests(t, data)
}

func runValidationTests(t *testing.T, data api.SleepDiaryEntryDataDto) {
	runCreateAndAssertBadRequest(t, data)
	runUpdateAndAssertBadRequest(t, data)
}

func runCreateAndAssertBadRequest(t *testing.T, data api.SleepDiaryEntryDataDto) {
	addDto := api.AddSleepDiaryEntryDto{
		AccountUuid:            uuid.NewString(),
		SleepDiaryEntryDataDto: data,
	}

	resp := mustPost(t, "/sleep_diary/entries", addDto)
	defer resp.Body.Close()
	assertHttpStatusCode(t, http.StatusBadRequest, resp)
}

func runUpdateAndAssertBadRequest(t *testing.T, data api.SleepDiaryEntryDataDto) {
	addDto := api.AddSleepDiaryEntryDto{
		AccountUuid:            uuid.NewString(),
		SleepDiaryEntryDataDto: prepareSleepDiaryData(),
	}

	addResp := mustPost(t, "/sleep_diary/entries", addDto)
	defer addResp.Body.Close()
	assertHttpStatusCode(t, http.StatusCreated, addResp)
	createdEntry := mustDecode[api.SleepDiaryEntryDto](addResp.Body)

	updateDto := api.UpdateSleepDiaryEntryDto{
		Version:                &createdEntry.Version,
		SleepDiaryEntryDataDto: data,
	}

	updateResp := mustPut(t, fmt.Sprintf("/sleep_diary/entries/%v", createdEntry.Id), updateDto)
	defer updateResp.Body.Close()
	assertHttpStatusCode(t, http.StatusBadRequest, updateResp)
}
