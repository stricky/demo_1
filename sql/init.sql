CREATE TABLE messages_stat
(
    id        INTEGER PRIMARY KEY AUTOINCREMENT UNIQUE,
    timestamp TEXT,
    body      TEXT,
    hash      TEXT,
    status    TEXT
);

CREATE TABLE messages_queue
(
    id        INTEGER PRIMARY KEY AUTOINCREMENT UNIQUE,
    timestamp TEXT,
    body      TEXT,
    hash      TEXT
);

CREATE TABLE incoming_messages
(
    id        INTEGER PRIMARY KEY AUTOINCREMENT UNIQUE,
    timestamp TEXT,
    body      TEXT
);

INSERT INTO messages_queue (timestamp, body, hash) VALUES (datetime('now'), 'body1', 'rrrrrr');
INSERT INTO messages_queue (timestamp, body, hash) VALUES (datetime('now'), 'body2', 'yyyyyyy');
INSERT INTO messages_queue (timestamp, body, hash) VALUES (datetime('now'), 'body3', 'zzzzzzz');
