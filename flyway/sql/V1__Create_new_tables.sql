CREATE TABLE messages (
                          message_offset INT CONSTRAINT messages_primary_key PRIMARY KEY,
                          topic VARCHAR,
                          partition INT,
                          consumer INT,
                          key VARCHAR,
                          value VARCHAR,
                          created_at TIMESTAMP NOT NULL default now()
)