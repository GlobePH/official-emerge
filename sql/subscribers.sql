CREATE TABLE subscribers (
	subscriber_number  VARCHAR NOT NULL PRIMARY KEY,
	access_token       VARCHAR NOT NULL UNIQUE,
	subscription_start DATE NOT NULL DEFAULT CURRENT_DATE,
	subscription_end   DATE
);
