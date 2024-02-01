-- +migrate Up

CREATE TABLE IF NOT EXISTS event_logs (
    uuid CHAR (50) PRIMARY KEY,
    level VARCHAR (50) NOT NULL,
    caller VARCHAR (50) NOT NULL,
    host VARCHAR (50) NOT NULL,
    status INT NOT NULL,
    method VARCHAR (50) NOT NULL,
    user_id VARCHAR (50) NOT NULL,
    message VARCHAR (128) NOT NULL,
    payload TEXT NOT NULL,
    created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP
);