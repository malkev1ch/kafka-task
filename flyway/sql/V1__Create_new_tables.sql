CREATE TABLE messages (
                          message_offset INT,
                          partition INT,
                          topic VARCHAR,
                          consumer INT,
                          key VARCHAR,
                          value VARCHAR,
                          created_at TIMESTAMP NOT NULL default now()
)