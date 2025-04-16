CREATE TABLE sleep_diary_entries (
    id BIGSERIAL PRIMARY KEY,
    account_uuid UUID NOT NULL,
    in_bed_at TIMESTAMPTZ NULL,
    tried_to_sleep_at TIMESTAMPTZ NOT NULL,
    sleep_delay_in_min INTEGER NULL,
    awakenings_count INTEGER NULL,
    awakenings_total_duration_in_min INTEGER NULL,
    final_wake_up_at TIMESTAMPTZ NOT NULL,
    out_of_bed_at TIMESTAMPTZ NULL,
    sleep_quality INTEGER NOT NULL,
    comments VARCHAR(2048) NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    version INTEGER NOT NULL DEFAULT 1
);

CREATE INDEX idx_sleep_diary_entries_account_uuid 
ON sleep_diary_entries (account_uuid);

CREATE INDEX idx_sleep_diary_entries_tried_to_sleep_at
ON sleep_diary_entries (tried_to_sleep_at);