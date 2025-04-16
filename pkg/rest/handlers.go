package rest

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/mabzd/snorlax/api"
	"github.com/mabzd/snorlax/internal/service"
	"github.com/mabzd/snorlax/internal/utils"
)

func getSleepDiaryEntry(service *service.SleepDiaryService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := r.PathValue("id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			respondWithError(w, api.ERR_INVALID, "invalid ID format", err)
			return
		}

		dto, serviceErr := service.GetEntryById(id)
		if serviceErr != nil {
			respondWithApiError(w, serviceErr)
			return
		}

		respondWithJSON(w, http.StatusOK, dto)
	}
}

func getSleepDiaryEntries(service *service.SleepDiaryService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		accountUuids := query["account_uuid"]

		fromDate, err := parseTimeQueryParam(query.Get("from_date"))
		if err != nil {
			respondWithError(w, api.ERR_INVALID, "invalid from_date format", err)
			return
		}

		toDate, err := parseTimeQueryParam(query.Get("to_date"))
		if err != nil {
			respondWithError(w, api.ERR_INVALID, "invalid to_date format", err)
			return
		}

		pageSize, err := parseInt64QueryParam(query.Get("page_size"))
		if err != nil {
			respondWithError(w, api.ERR_INVALID, "invalid page_size format", err)
			return
		}

		pageNumber, err := parseInt64QueryParam(query.Get("page_number"))
		if err != nil {
			respondWithError(w, api.ERR_INVALID, "invalid page_number format", err)
			return
		}

		filter := api.SleepDiaryFilterDto{
			AccountUuids: accountUuids,
			FromDate:     fromDate,
			ToDate:       toDate,
			PageSize:     utils.WithDefault(pageSize, api.DEFAULT_PAGE_SIZE),
			PageNumber:   utils.WithDefault(pageNumber, 1),
		}

		entries, serviceErr := service.GetEntriesByFilter(filter)
		if serviceErr != nil {
			respondWithApiError(w, serviceErr)
			return
		}

		respondWithJSON(w, http.StatusOK, entries)
	}
}

func createSleepDiaryEntry(service *service.SleepDiaryService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			respondWithError(w, api.ERR_INVALID, "invalid request body", err)
			return
		}
		defer r.Body.Close()

		var dto api.CreateSleepDiaryEntryDto
		if err := json.Unmarshal(body, &dto); err != nil {
			respondWithError(w, api.ERR_INVALID, "invalid JSON format", err)
			return
		}

		result, serviceErr := service.CreateEntry(dto)
		if serviceErr != nil {
			respondWithApiError(w, serviceErr)
			return
		}

		respondWithJSON(w, http.StatusCreated, result)
	}
}

func updateSleepDiaryEntry(service *service.SleepDiaryService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := r.PathValue("id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			respondWithError(w, api.ERR_INVALID, "invalid ID format", err)
			return
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			respondWithError(w, api.ERR_INVALID, "invalid request body", err)
			return
		}
		defer r.Body.Close()

		var dto api.UpdateSleepDiaryEntryDto
		if err := json.Unmarshal(body, &dto); err != nil {
			respondWithError(w, api.ERR_INVALID, "invalid JSON format", err)
			return
		}

		result, serviceErr := service.UpdateEntry(id, dto)
		if serviceErr != nil {
			respondWithApiError(w, serviceErr)
			return
		}

		respondWithJSON(w, http.StatusOK, result)
	}
}

func parseTimeQueryParam(param string) (*time.Time, error) {
	if param == "" {
		return nil, nil
	}
	date, err := time.Parse(time.RFC3339, param)
	return &date, err
}

func parseInt64QueryParam(param string) (*int64, error) {
	if param == "" {
		return nil, nil
	}
	value, err := strconv.ParseInt(param, 10, 64)
	return &value, err
}

func respondWithApiError(w http.ResponseWriter, err api.Error) {
	dto := err.ToErrorDto()

	detailsMessage := ""
	if len(dto.Details) > 0 {
		detailsMessage = fmt.Sprintf("(%s)", strings.Join(dto.Details, ", "))
	}

	log.Printf("service error [%s]: %v %s\n", dto.Code, dto.Message, detailsMessage)
	respondWithJSON(w, toHttpError(dto.Code), dto)
}

func respondWithError(w http.ResponseWriter, code api.ErrorCode, message string, cause error) {
	log.Printf("Error(%v): %s, Cause: %v", code, message, cause)
	respondWithJSON(w, toHttpError(code), api.ErrorDto{Message: message, Code: code})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func toHttpError(code api.ErrorCode) int {
	switch code {
	case api.ERR_NOT_FOUND:
		return http.StatusNotFound
	case api.ERR_CONFLICT:
		return http.StatusConflict
	case api.ERR_INVALID:
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}
