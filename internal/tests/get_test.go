package tests

import (
	"fmt"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/mabzd/snorlax/api"
	"github.com/stretchr/testify/assert"
)

func TestGetOneEntryByFilter(t *testing.T) {
	createdEntry := mustCreateRandomEntry(t)
	page := mustGetEntriesByQuery(t, fmt.Sprintf("?account_uuid=%s", createdEntry.AccountUuid))
	assertDefaultPageEqual(t, 1, []api.SleepDiaryEntryDto{createdEntry}, page)
}

func TestGetNoEntriesByFilter(t *testing.T) {
	unknownUuid := uuid.NewString()
	page := mustGetEntriesByQuery(t, fmt.Sprintf("?account_uuid=%s", unknownUuid))
	assertDefaultPageEqual(t, 0, []api.SleepDiaryEntryDto{}, page)
}

func TestGetMultipleEntriesByFilter(t *testing.T) {
	data := mustCreateTestData(t)
	page := mustGetEntriesByQuery(t, fmt.Sprintf("?account_uuid=%s", data.UuidA))
	assertDefaultPageEqual(t, 4, data.A, page)
}
func TestPageSizeOne(t *testing.T) {
	data := mustCreateTestData(t)
	page := mustGetEntriesByQuery(t, fmt.Sprintf("?account_uuid=%s&page_size=1", data.UuidA))
	assertPageEqual(t, 4, 1, 1, data.A[:1], page)
}

func TestPageSizeTwoPageTwo(t *testing.T) {
	data := mustCreateTestData(t)
	page := mustGetEntriesByQuery(t, fmt.Sprintf("?account_uuid=%s&page_size=2&page_number=2", data.UuidA))
	assertPageEqual(t, 4, 2, 2, data.A[2:4], page)
}

func TestPageSizeThreePageTwo(t *testing.T) {
	data := mustCreateTestData(t)
	page := mustGetEntriesByQuery(t, fmt.Sprintf("?account_uuid=%s&page_size=3&page_number=2", data.UuidA))
	assertPageEqual(t, 4, 3, 2, data.A[3:], page)
}

func TestPageNumberExceeding(t *testing.T) {
	data := mustCreateTestData(t)
	page := mustGetEntriesByQuery(t, fmt.Sprintf("?account_uuid=%s&page_size=10&page_number=100", data.UuidA))
	assertPageEqual(t, 4, 10, 100, []api.SleepDiaryEntryDto{}, page)
}

func TestMultipleAccounts(t *testing.T) {
	data := mustCreateTestData(t)
	page := mustGetEntriesByQuery(t, fmt.Sprintf("?account_uuid=%s&account_uuid=%s", data.UuidB, data.UuidC))
	assertDefaultPageEqual(t, 2, append(data.B, data.C...), page)
}

/*
*   0  1  2  3
* --|--|--|--|->
*         \____ (inclusive)
* should return 2, 3
 */
func TestDateFrom(t *testing.T) {
	data := mustCreateTestData(t)
	page := mustGetEntriesByQuery(t, fmt.Sprintf("?account_uuid=%s&from_date=%s", data.UuidA, getDate(data.A[2], 0)))
	assertDefaultPageEqual(t, 2, data.A[2:], page)
}

/*
*   0  1  2  3
* --|--|--|--|->
*          \____ +1s
* should return 3
 */
func TestDateFromShiftSecondUp(t *testing.T) {
	data := mustCreateTestData(t)
	page := mustGetEntriesByQuery(t, fmt.Sprintf("?account_uuid=%s&from_date=%s", data.UuidA, getDate(data.A[2], 1)))
	assertDefaultPageEqual(t, 1, data.A[3:], page)
}

/*
*   0  1  2  3
* --|--|--|--|->
* ________| (exclusive)
* should return 0, 1
 */
func TestDateTo(t *testing.T) {
	data := mustCreateTestData(t)
	page := mustGetEntriesByQuery(t, fmt.Sprintf("?account_uuid=%s&to_date=%s", data.UuidA, getDate(data.A[2], 0)))
	assertDefaultPageEqual(t, 2, data.A[:2], page)
}

/*
*   0  1  2  3
* --|--|--|--|->
* _________| +1s
* should return 0, 1, 2
 */
func TestDateToShiftSecondUp(t *testing.T) {
	data := mustCreateTestData(t)
	page := mustGetEntriesByQuery(t, fmt.Sprintf("?account_uuid=%s&to_date=%s", data.UuidA, getDate(data.A[2], 1)))
	assertDefaultPageEqual(t, 3, data.A[:3], page)
}

/*
*   0  1  2  3
* --|--|--|--|->
*      \__| (inclusive; exclusive)
* should return 1
 */
func TestFromAndTo(t *testing.T) {
	data := mustCreateTestData(t)
	page := mustGetEntriesByQuery(t, fmt.Sprintf(
		"?account_uuid=%s&from_date=%s&to_date=%s",
		data.UuidA,
		getDate(data.A[1], 0),
		getDate(data.A[2], 0)))

	assertDefaultPageEqual(t, 1, data.A[1:2], page)
}

func TestNoAccountUuid(t *testing.T) {
	runGetEntriesByQueryAndAssertBadRequest(t, "?")
}

func TestInvalidAccountUuid(t *testing.T) {
	runGetEntriesByQueryAndAssertBadRequest(t, "?account_uuid=invalid")
}

func TestZeroPageSize(t *testing.T) {
	runGetEntriesByQueryAndAssertBadRequest(t, fmt.Sprintf("?account_uuid=%s&page_size=0", uuid.NewString()))
}

func TestMoreThanMaxPageSize(t *testing.T) {
	runGetEntriesByQueryAndAssertBadRequest(t, fmt.Sprintf("?account_uuid=%s&page_size=%v", uuid.NewString(), api.MAX_PAGE_SIZE+1))
}

func TestInvalidPageSize(t *testing.T) {
	runGetEntriesByQueryAndAssertBadRequest(t, fmt.Sprintf("?account_uuid=%s&page_size=invalid", uuid.NewString()))
}

func TestZeroPageNumber(t *testing.T) {
	runGetEntriesByQueryAndAssertBadRequest(t, fmt.Sprintf("?account_uuid=%s&page_number=0", uuid.NewString()))
}

func TestInvalidPageNumber(t *testing.T) {
	runGetEntriesByQueryAndAssertBadRequest(t, fmt.Sprintf("?account_uuid=%s&page_number=invalid", uuid.NewString()))
}

func TestDateFromInvalid(t *testing.T) {
	runGetEntriesByQueryAndAssertBadRequest(t, fmt.Sprintf("?account_uuid=%s&from_date=invalid", uuid.NewString()))
}

func TestDateToInvalid(t *testing.T) {
	runGetEntriesByQueryAndAssertBadRequest(t, fmt.Sprintf("?account_uuid=%s&to_date=invalid", uuid.NewString()))
}

func TestDateFromAndToNotInOrder(t *testing.T) {
	date := time.Now().Format(time.RFC3339)
	runGetEntriesByQueryAndAssertBadRequest(t, fmt.Sprintf("?account_uuid=%s&from_date=%s&to_date=%s", uuid.NewString(), date, date))
}

type testData struct {
	UuidA string
	UuidB string
	UuidC string
	A     []api.SleepDiaryEntryDto
	B     []api.SleepDiaryEntryDto
	C     []api.SleepDiaryEntryDto
}

func mustCreateTestData(t *testing.T) testData {
	now := time.Now()
	a := uuid.NewString()
	b := uuid.NewString()
	c := uuid.NewString()
	a1 := mustCreateEntry(t, api.AddSleepDiaryEntryDto{AccountUuid: a, SleepDiaryEntryDataDto: prepateSleepDiaryDataForSleepAt(now)})
	a2 := mustCreateEntry(t, api.AddSleepDiaryEntryDto{AccountUuid: a, SleepDiaryEntryDataDto: prepateSleepDiaryDataForSleepAt(now.Add(24 * time.Hour))})
	a3 := mustCreateEntry(t, api.AddSleepDiaryEntryDto{AccountUuid: a, SleepDiaryEntryDataDto: prepateSleepDiaryDataForSleepAt(now.Add(48 * time.Hour))})
	a4 := mustCreateEntry(t, api.AddSleepDiaryEntryDto{AccountUuid: a, SleepDiaryEntryDataDto: prepateSleepDiaryDataForSleepAt(now.Add(72 * time.Hour))})
	b1 := mustCreateEntry(t, api.AddSleepDiaryEntryDto{AccountUuid: b, SleepDiaryEntryDataDto: prepateSleepDiaryDataForSleepAt(now.Add(1 * time.Second))})
	c1 := mustCreateEntry(t, api.AddSleepDiaryEntryDto{AccountUuid: c, SleepDiaryEntryDataDto: prepateSleepDiaryDataForSleepAt(now.Add(2 * time.Second))})

	return testData{
		UuidA: a,
		UuidB: b,
		UuidC: c,
		A:     []api.SleepDiaryEntryDto{a1, a2, a3, a4},
		B:     []api.SleepDiaryEntryDto{b1},
		C:     []api.SleepDiaryEntryDto{c1},
	}
}

func mustGetEntriesByQuery(t *testing.T, query string) api.PageDto[api.SleepDiaryEntryDto] {
	pageResp := mustGet(t, fmt.Sprintf("/sleep_diary/entries%s", query))
	defer pageResp.Body.Close()
	assertHttpStatusCode(t, http.StatusOK, pageResp)
	return mustDecode[api.PageDto[api.SleepDiaryEntryDto]](pageResp.Body)
}

func runGetEntriesByQueryAndAssertBadRequest(t *testing.T, query string) {
	pageResp := mustGet(t, fmt.Sprintf("/sleep_diary/entries%s", query))
	defer pageResp.Body.Close()
	assertHttpStatusCode(t, http.StatusBadRequest, pageResp)
}

func assertDefaultPageEqual(t *testing.T, total int64, items []api.SleepDiaryEntryDto, actual api.PageDto[api.SleepDiaryEntryDto]) {
	assertPageEqual(t, total, api.DEFAULT_PAGE_SIZE, api.DEFAULT_PAGE_NUMBER, items, actual)
}

func assertPageEqual(t *testing.T, total, pageSize, pageNumber int64, items []api.SleepDiaryEntryDto, actual api.PageDto[api.SleepDiaryEntryDto]) {
	assert.Equal(t, total, actual.TotalCount)
	assert.Equal(t, pageSize, actual.PageSize)
	assert.Equal(t, pageNumber, actual.PageNumber)
	assert.Equal(t, len(items), len(actual.Items))
	for i := 0; i < len(items); i++ {
		assertEqualSleepDiaryEntryDto(t, items[i], actual.Items[i], true)
	}
}

func getDate(item api.SleepDiaryEntryDto, shift int) string {
	return url.QueryEscape(item.TriedToSleepAt.Add(time.Duration(shift) * time.Second).Format(time.RFC3339))
}
