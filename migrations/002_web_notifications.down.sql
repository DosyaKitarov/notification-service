
ALTER TABLE email_notifications RENAME TO notifications;

ALTER TABLE notifications ADD COLUMN notification_channel VARCHAR(255);

DROP TABLE IF EXISTS web_notifications;