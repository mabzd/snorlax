package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"github.com/mabzd/snorlax/api"
	"github.com/mabzd/snorlax/internal/config"
	"github.com/mabzd/snorlax/pkg/dbm"
	"github.com/mabzd/snorlax/pkg/rest"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"gotest.tools/v3/assert"
)

var srv *httptest.Server

func TestMain(m *testing.M) {
	log.SetOutput(io.Discard)
	ctx := context.Background()

	DB_USER := "postgres"
	DB_PASS := "postgres"
	DB_NAME := "testdb"

	containerReq := testcontainers.ContainerRequest{
		Image:        "postgres:15",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_USER":     DB_USER,
			"POSTGRES_PASSWORD": DB_PASS,
			"POSTGRES_DB":       DB_NAME,
		},
		WaitingFor: wait.ForListeningPort("5432/tcp").WithStartupTimeout(30 * time.Second),
	}

	dbContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: containerReq,
		Started:          true,
	})
	if err != nil {
		log.Fatalf("Could not start container: %s", err)
	}
	defer func() {
		_ = dbContainer.Terminate(ctx)
	}()

	port, _ := dbContainer.MappedPort(ctx, "5432")
	cfg := config.Config{
		ApiPort:            "8080",
		DbHost:             "localhost",
		DbPort:             port.Port(),
		DbUser:             DB_USER,
		DbPass:             DB_PASS,
		DbName:             DB_NAME,
		ServerTimeoutInSec: 5,
	}

	dbm.UpgradeDatabaseIfNeeded(cfg)
	handler := rest.NewServerHandler(cfg)
	srv = httptest.NewServer(handler)
	defer srv.Close()
	code := m.Run()
	os.Exit(code)
}

func newMinimalRandomEntryData() api.SleepDiaryEntryDataDto {
	tz := newRandomTimezone()
	sleepAt := time.Now().Add(time.Duration(-rand.Intn(24)) * time.Hour)
	return api.SleepDiaryEntryDataDto{
		Timezone:       tz.String(),
		TriedToSleepAt: sleepAt.In(tz),
		FinalWakeUpAt:  sleepAt.In(tz).Add(time.Duration(rand.Intn(4)+6) * time.Hour),
		SleepQuality:   api.ExcellentSleepQuality,
	}
}

func newRandomEntryData() api.SleepDiaryEntryDataDto {
	sleepAt := time.Now().Add(time.Duration(-rand.Intn(24)) * time.Hour)
	return newRandomEntryDataForSleepAt(sleepAt)
}

func newRandomEntryDataForSleepAt(sleepAt time.Time) api.SleepDiaryEntryDataDto {
	tz := newRandomTimezone()
	wakeUpAt := sleepAt.Add(time.Duration(8+rand.Intn(4)) * time.Hour)
	return api.SleepDiaryEntryDataDto{
		Timezone:                     tz.String(),
		InBedAt:                      toPtr(sleepAt.In(tz).Add(time.Duration(-rand.Intn(60)) * time.Minute)),
		TriedToSleepAt:               sleepAt.In(tz),
		SleepDelayInMin:              toPtr(rand.Intn(60)),
		AwakeningsCount:              toPtr(rand.Intn(5)),
		AwakeningsTotalDurationInMin: toPtr(rand.Intn(60)),
		FinalWakeUpAt:                wakeUpAt.In(tz),
		OutOfBedAt:                   toPtr(wakeUpAt.In(tz).Add(time.Duration(60) * time.Minute)),
		SleepQuality:                 api.SleepQuality(rand.Intn(5) + 1),
		Comments:                     toPtr("Good sleep"),
	}
}

func newRandomTimezone() *time.Location {
	timezones := []string{"America/New_York", "Europe/London", "Asia/Tokyo", "Australia/Sydney", "America/Los_Angeles"}
	randomTimezone := timezones[rand.Intn(len(timezones))]
	location, err := time.LoadLocation(randomTimezone)
	if err != nil {
		log.Fatalf("Failed to load location: %v", err)
	}
	return location
}

func assertEqualEntryDto(t *testing.T, expected api.SleepDiaryEntryDto, actual api.SleepDiaryEntryDto, compareVersion bool) {
	assertValuesEqual(t, &expected.Id, &actual.Id, "Id")
	assertValuesEqual(t, &expected.AccountUuid, &actual.AccountUuid, "AccountUuid")
	assertValuesEqualTimeMsPrec(t, expected.InBedAt, actual.InBedAt, "InBedAt")
	assertValuesEqualTimeMsPrec(t, &expected.TriedToSleepAt, &actual.TriedToSleepAt, "TriedToSleepAt")
	assertValuesEqual(t, expected.SleepDelayInMin, actual.SleepDelayInMin, "SleepDelayInMin")
	assertValuesEqual(t, expected.AwakeningsCount, actual.AwakeningsCount, "AwakeningsCount")
	assertValuesEqual(t, expected.AwakeningsTotalDurationInMin, actual.AwakeningsTotalDurationInMin, "AwakeningsTotalDurationInMin")
	assertValuesEqualTimeMsPrec(t, &expected.FinalWakeUpAt, &actual.FinalWakeUpAt, "FinalWakeUpAt")
	assertValuesEqualTimeMsPrec(t, expected.OutOfBedAt, actual.OutOfBedAt, "OutOfBedAt")
	assertValuesEqual(t, &expected.SleepQuality, &actual.SleepQuality, "SleepQuality")
	assertValuesEqual(t, expected.Comments, actual.Comments, "Comments")

	if compareVersion {
		assertValuesEqual(t, &expected.Version, &actual.Version, "Version")
	}
}

func assertValuesEqual[T comparable](t *testing.T, expected, actual *T, msg string) {
	if expected == nil && actual == nil {
		return
	}
	if expected == nil || actual == nil {
		t.Fatalf("Expected and actual do not match: expected=%v, actual=%v", expected, actual)
	}
	assert.Equal(t, *expected, *actual, msg)
}

func assertValuesEqualTimeMsPrec(t *testing.T, expected *time.Time, actual *time.Time, msg string) {
	if expected == nil && actual == nil {
		return
	}
	if expected == nil || actual == nil {
		t.Fatalf("Expected and actual times do not match (%s): expected=%v, actual=%v", msg, expected, actual)
	}
	assert.Equal(
		t,
		expected.Truncate(time.Millisecond).Format(time.RFC3339),
		actual.Truncate(time.Millisecond).Format(time.RFC3339),
		msg)
}

func mustCreateRandomEntry(t *testing.T) api.SleepDiaryEntryDto {
	return mustCreateEntry(t, api.CreateSleepDiaryEntryDto{
		AccountUuid:            uuid.NewString(),
		SleepDiaryEntryDataDto: newRandomEntryData(),
	})
}

func mustCreateEntry(t *testing.T, dto api.CreateSleepDiaryEntryDto) api.SleepDiaryEntryDto {
	resp := mustPost(t, "/sleep_diary/entries", dto)
	defer resp.Body.Close()
	assertHttpStatusCode(t, http.StatusCreated, resp)
	return mustDecode[api.SleepDiaryEntryDto](resp.Body)
}

func mustGetEntryById(t *testing.T, id int64) api.SleepDiaryEntryDto {
	resp := mustGet(t, fmt.Sprintf("/sleep_diary/entries/%v", id))
	defer resp.Body.Close()
	assertHttpStatusCode(t, http.StatusOK, resp)
	return mustDecode[api.SleepDiaryEntryDto](resp.Body)
}

func assertHttpStatusCode(t *testing.T, expectedStatusCode int, resp *http.Response) {
	if expectedStatusCode != resp.StatusCode {
		body, _ := io.ReadAll(resp.Body)
		t.Fatalf(
			"Expected status code %d, got %d. Response body: '%s'",
			expectedStatusCode,
			resp.StatusCode,
			strings.TrimSpace(string(body)))
	}
}

func mustPost(t *testing.T, path string, payload interface{}) *http.Response {
	return mustSendJson(t, http.MethodPost, path, payload)
}

func mustPut(t *testing.T, path string, payload interface{}) *http.Response {
	return mustSendJson(t, http.MethodPut, path, payload)
}

func mustGet(t *testing.T, path string) *http.Response {
	resp, err := http.Get(srv.URL + path)
	if err != nil {
		t.Fatalf("Failed to make GET request: %v", err)
	}
	return resp
}

func mustSendJson(t *testing.T, method string, path string, payload interface{}) *http.Response {
	url := srv.URL + path
	json := mustMashal(payload)
	req, err := http.NewRequest(method, url, bytes.NewReader(json))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	return resp
}

func mustMashal(v interface{}) []byte {
	data, err := json.Marshal(v)
	if err != nil {
		log.Fatalf("Failed to marshal: %v", err)
	}
	return data
}

func mustDecode[T any](r io.Reader) T {
	var v T
	if err := json.NewDecoder(r).Decode(&v); err != nil {
		log.Panicf("failed to decode JSON: %v", err)
	}
	return v
}

func swap(t1, t2 *time.Time) {
	*t1, *t2 = *t2, *t1
}

func toPtr[T any](value T) *T {
	return &value
}
