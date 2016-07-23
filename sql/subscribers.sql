DROP TABLE IF EXISTS subscribers;
CREATE TABLE subscribers (
	subscriber_number  VARCHAR NOT NULL PRIMARY KEY,
	access_token       VARCHAR NOT NULL UNIQUE,
	subscription_start TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	subscription_end   TIMESTAMP
);
