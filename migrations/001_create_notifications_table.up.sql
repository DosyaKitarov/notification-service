CREATE TABLE notifications (
    id SERIAL PRIMARY KEY,
    user_id VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    type VARCHAR(255) NOT NULL,
    notification_channel JSONB NOT NULL, -- Changed to JSONB to store JSON arrays
    metadata JSONB NOT NULL, -- Added metadata column to store additional details
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);