CREATE TABLE IF NOT EXISTS topics (
    id INTEGER PRIMARY KEY,
    key TEXT NOT NULL UNIQUE,
    title TEXT NOT NULL,
    comment TEXT NOT NULL DEFAULT '',
    status TEXT NOT NULL DEFAULT 'not_studied',
    last_reviewed_at TEXT,
    CONSTRAINT chk_status CHECK (
        status IN ('not_studied', 'familiar', 'learned', 'review')
      )
);

CREATE TABLE IF NOT EXISTS tasks (
    id INTEGER PRIMARY KEY,
    title TEXT NOT NULL,
    comment TEXT NOT NULL DEFAULT '',
    status TEXT NOT NULL DEFAULT 'not_studied',
    last_reviewed_at TEXT,
    CONSTRAINT chk_status CHECK (
        status IN ('not_studied', 'familiar', 'learned', 'review')
     )
);

CREATE TABLE IF NOT EXISTS goal (
    id INTEGER PRIMARY KEY CHECK ( id = 1 ),
    text TEXT NOT NULL DEFAULT '',
    comment TEXT NOT NULL DEFAULT ''
);

INSERT OR IGNORE INTO goal (id, text, comment)
VALUES (1, '', '');