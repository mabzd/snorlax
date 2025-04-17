![snorlax says hi](snorlax.png "Snorlax")

# Snorlax 

Snorlax is a showcase Go REST service that allows clients to track their sleep data.

Datapoints captured are modeled after [Consensus Sleep Diary (CSD)](https://pmc.ncbi.nlm.nih.gov/articles/PMC3250369/), a standardized questionnaire used to track sleep patterns and behaviors through daily self-reported entries. CSD aligns well with the core parameters required by the specification—sleep start and end time, and sleep quality.

| Field | Type | Description |
|-|-|-|
| `timezone` | text | The timezone the person slept in (optional, UTC if omitted). |
| `in_bed_at` | timestamp | The time the person got into bed (optional).|
| `tried_to_sleep_at` | timestamp **REQUIRED** | The time the person attempted to fall asleep. |
| `sleep_delay_in_min` | number | Number of minutes it took to fall asleep after trying (optional). |
| `awakenings_count` | number | Number of times the person woke up during the night (optional). |
| `awakenings_total_duration_in_min`| number  | Total duration of awakenings during the night in minutes (optional). |
| `final_wake_up_at` | timestamp **REQUIRED** | The time of the final awakening in the morning. |
| `out_of_bed_at` | timestamp | The time the person got out of bed for the day (optional). |
| `sleep_quality` | number (1-5) **REQUIRED** | Self-rated quality of sleep on a 1-5 scale (1=very poor, 5=excellent) |
| `comments` | text | Additional comments or notes about the sleep experience (optional). |

## How to Run

In the project root directory run
```
docker compose up
```
This runs the `snorlax-api` image hosting the API service on port :8080 and its dependency image `snorlax-db` hosting PostgreSQL database on port :5432.

## API Usage

### Create Entry
`POST /sleep_diary/entries`

Request
```
curl -X POST http://localhost:8080/sleep_diary/entries \
  -H "Content-Type: application/json" \
  -d '{
    "account_uuid": "c7f23d8a-5a10-4a1a-9c55-2f8c5d872f09",
    "timezone": "UTC",
    "in_bed_at": "2025-04-15T22:30:00Z",
    "tried_to_sleep_at": "2025-04-15T22:45:00Z",
    "sleep_delay_in_min": 15,
    "awakenings_count": 2,
    "awakenings_total_duration_in_min": 20,
    "final_wake_up_at": "2025-04-16T06:30:00Z",
    "out_of_bed_at": "2025-04-16T06:45:00Z",
    "sleep_quality": 4,
    "comments": "Woke up a couple of times, but overall decent sleep."
  }'
```

Response
```json
{
  "id": 1,
  "account_uuid": "c7f23d8a-5a10-4a1a-9c55-2f8c5d872f09",
  "version": 1,
  "timezone": "UTC",
  "in_bed_at": "2025-04-15T22:30:00Z",
  "tried_to_sleep_at": "2025-04-15T22:45:00Z",
  "sleep_delay_in_min": 15,
  "awakenings_count": 2,
  "awakenings_total_duration_in_min": 20,
  "final_wake_up_at": "2025-04-16T06:30:00Z",
  "out_of_bed_at": "2025-04-16T06:45:00Z",
  "sleep_quality": 4,
  "comments": "Woke up a couple of times, but overall decent sleep."
}
```

### Read Entry by ID
`GET /sleep_diary/entries/{id}`

Request
```
curl http://localhost:8080/sleep_diary/entries/1
```

Response
```json
{
  "id": 1,
  "account_uuid": "c7f23d8a-5a10-4a1a-9c55-2f8c5d872f09",
  "version": 1,
  "timezone": "UTC",
  "in_bed_at": "2025-04-15T22:30:00Z",
  "tried_to_sleep_at": "2025-04-15T22:45:00Z",
  "sleep_delay_in_min": 15,
  "awakenings_count": 2,
  "awakenings_total_duration_in_min": 20,
  "final_wake_up_at": "2025-04-16T06:30:00Z",
  "out_of_bed_at": "2025-04-16T06:45:00Z",
  "sleep_quality": 4,
  "comments": "Woke up a couple of times, but overall decent sleep."
}
```

### Query Entries
`GET /sleep_diary/entries?account_uuid={account_uuid}`

Allowed query parameters:

* `account_uuid` **REQUIRED** - returns entries for given account UUID. Multiple occurrences of this parameter is supported to retrieve data for many accounts.
* `date_from` - returns entries since given timestamp (inclusive). Each entry's timestamp is based on the tried_to_sleep_at attribute, which marks the start of the sleep attempt.
* `date_to` - returns entries up to given timestamp (exclusive).
* `page_size` - number of entries on each page (1-1000). Default is 100.
* `page_number` - used to iterate through pages. Pages are numbered from 1.

Request
```
curl http://localhost:8080/sleep_diary/entries?account_uuid=c7f23d8a-5a10-4a1a-9c55-2f8c5d872f09
```

Response
```json
{
  "total_count": 1,
  "page_size": 100,
  "page_number": 1,
  "items": [
    {
      "id": 1,
      "account_uuid": "c7f23d8a-5a10-4a1a-9c55-2f8c5d872f09",
      "version": 1,
      "timezone": "UTC",
      "in_bed_at": "2025-04-15T22:30:00Z",
      "tried_to_sleep_at": "2025-04-15T22:45:00Z",
      "sleep_delay_in_min": 15,
      "awakenings_count": 2,
      "awakenings_total_duration_in_min": 20,
      "final_wake_up_at": "2025-04-16T06:30:00Z",
      "out_of_bed_at": "2025-04-16T06:45:00Z",
      "sleep_quality": 4,
      "comments": "Woke up a couple of times, but overall decent sleep."
    }
  ]
}
```

### Update Entry by ID
`PUT /sleep_diary/entries/{id}`

Optimistic locking is supported if `version` attribute is specified: entry is updated only if given version has not changed, otherwise `409 Conflict` is returned. If the `version` attribute is omitted, the update will proceed and overwrite existing data unconditionally.

Request
```
curl -X PUT http://localhost:8080/sleep_diary/entries/1 \
  -H "Content-Type: application/json" \
  -d '{
    "version": 1,
    "timezone": "UTC",
    "in_bed_at": "2025-04-14T23:00:00Z",
    "tried_to_sleep_at": "2025-04-14T23:15:00Z",
    "sleep_delay_in_min": 25,
    "awakenings_count": 3,
    "awakenings_total_duration_in_min": 30,
    "final_wake_up_at": "2025-04-15T07:00:00Z",
    "out_of_bed_at": "2025-04-15T07:20:00Z",
    "sleep_quality": 3,
    "comments": "Restless night, kept waking up."
  }'
```

Response
```json
{
  "id": 1,
  "account_uuid": "c7f23d8a-5a10-4a1a-9c55-2f8c5d872f09",
  "version": 2,
  "timezone": "UTC",
  "in_bed_at": "2025-04-14T23:00:00Z",
  "tried_to_sleep_at": "2025-04-14T23:15:00Z",
  "sleep_delay_in_min": 25,
  "awakenings_count": 3,
  "awakenings_total_duration_in_min": 30,
  "final_wake_up_at": "2025-04-15T07:00:00Z",
  "out_of_bed_at": "2025-04-15T07:20:00Z",
  "sleep_quality": 3,
  "comments": "Restless night, kept waking up."
}
```

## Key Design & Implementation Decisions

### Run in Trusted Environment

The service is designed to run in a trusted environment, with authentication and account management handled upstream.

Rationale: separation of concerns and simpler design. The service is solely focused on the domain-specific logic.

### Timezones Support

API was designed to support timezones. Each timestamp in the entry is assumed to belong to specified `timezone`, even if timestamps themselves are with different offsets: eg. `"in_bed_at": "2025-04-17T11:31:27+0200"` provided with `"timezone": "America/New_York"` is interpreted in New York local time (i.e., 2025-04-17 05:31:27). Timezone information is persisted in the database.

Rationale: sleep and wake times are inherently tied to a user's local time, without proper timezone handling, collected data could be misaligned or misinterpreted, especially across regions. This datapoint allows for reliable analysis of collected data in the future.

### Separate HTTP and Service Layer

Codebase was designed with clear split between HTTP layer (`pkg/rest`) and service layer (`internal/service`) with service layer having no dependency on HTTP.

Rationale: domain logic remains agnostic of the delivery mechanism, allowing for clean architecture and flexibility across different interfaces.

### End-to-End Testing in Favor of Unit Testing

The service is tested end-to-end from HTTP layer to the persistence layer with help of `testconainers-go`. Tests can be run using Taskfile command `task e2e` (or by manually running `go test` inside tests module directory).

Rationale: end-to-end tests are often more effective than unit tests, as they validate the full request-handling flow and offer higher confidence with less overhead. Unit tests, while useful, can tightly couple tests to implementation details, making refactoring harder.

### Testing Code in Separate Module

Test code was offloaded to separate internal module (`internal/tests`).

Rationale: isolating test dependencies from service dependencies.

### Database Migrations

Separate database migration tool (`cmd/dbm`) takes care of applying migrations in order if it's detected that the schema needs an update. Proposed system is very simple but powerful: it maintains a table with set of files applied in SQL and compares it with the set of migration files embedded with the executable. 

Rationale: data often outlives the code that created it, maintaining a clear history of how the database evolves is crucial for long-term stability. Migration tool was designed to be a separate executable allowing to be reliably integrated into CI/CD pipelines or run on demand (i.e., `task upgrade-db`)