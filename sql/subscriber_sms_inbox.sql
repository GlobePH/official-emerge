DROP TABLE IF EXISTS subscriber_sms_inbox;
CREATE TABLE subscriber_sms_inbox (
	message_id         VARCHAR NOT NULL PRIMARY KEY,
	subscriber_number  VARCHAR NOT NULL,
	message            VARCHAR NOT NULL,
	received_at        TIMESTAMP NOT NULL,
	FOREIGN KEY (subscriber_number)
		REFERENCES subscribers (subscriber_number)
			ON UPDATE RESTRICT
			ON DELETE RESTRICT
);
