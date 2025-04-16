![snorlax says hi](snorlax.png "Snorlax")

# Snorlax 

Snorlax is a showcase Go REST service that allows clients to track their sleep data.

Datapoints captured are modeled after [Consensus Sleep Diary (CSD)](https://pmc.ncbi.nlm.nih.gov/articles/PMC3250369/), a standardized questionnaire used to track sleep patterns and behaviors through daily self-reported entries:

| Field | Type | Description |
|-|-|-|
| `in_bed_at` | timestamp | The time the person got into bed (optional).|
| `tried_to_sleep_at` | timestamp **REQUIRED** | The time the person attempted to fall asleep. |
| `sleep_delay_in_min` | number | Number of minutes it took to fall asleep after trying (optional). |
| `awakenings_count` | number | Number of times the person woke up during the night (optional). |
| `awakenings_total_duration_in_min`| number  | Total duration of awakenings during the night in minutes (optional). |
| `final_wake_up_at` | timestamp **REQUIRED** | The time of the final awakening in the morning. |
| `out_of_bed_at` | timestamp | The time the person got out of bed for the day (optional). |
| `sleep_quality` | number, 1-5 **REQUIRED** | Self-rated quality of sleep on a 1-5 scale (1=very poor, 5=excellent) |
| `comments` | text | Additional comments or notes about the sleep experience (optional). |

## How to run

In the project root directory run
```
docker-compose up
```
This runs the `snorlax-api` image hosting the API service on port :8080 and it's dependency image `snorlax-db` hosting PostgreSQL database on port :5432.

## API usage

### Create entry
`POST /sleep_diary/entries`

Request
```
curl -X POST http://localhost:8080/sleep_diary/entries \
  -H "Content-Type: application/json" \
  -d '{
    "account_uuid": "c7f23d8a-5a10-4a1a-9c55-2f8c5d872f09",
    "get_in_bed_at": "2025-04-15T22:30:00Z",
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
  "get_in_bed_at": "2025-04-15T22:30:00Z",
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

### Read entry by ID
`GET /sleep_diary/entries/{id}`

Request
```
curl localhost:8080/sleep_diary/entries/1
```

Response
```json
{
  "id": 1,
  "account_uuid": "c7f23d8a-5a10-4a1a-9c55-2f8c5d872f09",
  "version": 1,
  "get_in_bed_at": "2025-04-15T22:30:00Z",
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

### Query entries
`GET /sleep_diary/entries?account_uuid={account_uuid}`

Allowed query parameters:

* `account_uuid` **REQUIRED** - returns entries for given account UUID. Multiple ocurrences of this parameter is supported to retrieve data for many accounts.
* `date_from` - returns entries since given timestamp (inclusive). Each entry's timestamp is based on the tried_to_sleep_at attribute, which marks the start of the sleep attempt.
* `date_to` - returns entries up to given timestamp (exclusive).
* `page_size` - number of entries on each page (1-1000). Default is 100.
* `page_number` - used to interate through pages. Pages are numbered from 1.

Request
```
curl localhost:8080/sleep_diary/entries?account_uuid=c7f23d8a-5a10-4a1a-9c55-2f8c5d872f09
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
      "get_in_bed_at": "2025-04-15T22:30:00Z",
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

### Update entry by ID
`PUT /sleep_diary/entries/{id}`

Optimisic locking is supported if `version` attribute is specified: entry is updated only if given version has not changed, otherwise `409 Conflict` is returned. If the `version` attribute is omitted, the update will proceed and overwrite existing data unconditionally.

Request
```
curl -X PUT http://localhost:8080/sleep_diary/entries/1 \
  -H "Content-Type: application/json" \
  -d '{
    "version": 1,
    "get_in_bed_at": "2025-04-14T23:00:00Z",
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
  "get_in_bed_at": "2025-04-14T23:00:00Z",
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