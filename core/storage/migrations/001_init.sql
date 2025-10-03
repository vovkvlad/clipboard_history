-- Enable WAL mode for better concurrency and performance
PRAGMA journal_mode = WAL;
PRAGMA synchronous=NORMAL;

-- Enable FTS5 extension (should be available in modern SQLite)
PRAGMA foreign_keys = ON;

-- Create the main clipboard history table
CREATE TABLE IF NOT EXISTS clipboard_history (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  content TEXT NOT NULL,
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Create FTS5 virtual table for full-text search
CREATE VIRTUAL TABLE IF NOT EXISTS clipboard_history_fts USING fts5(
  content,
  content = clipboard_history,
  content_rowid = id
);

-- Create triggers to keep FTS5 table in sync with main table
CREATE TRIGGER IF NOT EXISTS clipboard_history_ai AFTER INSERT ON clipboard_history BEGIN
  INSERT INTO clipboard_history_fts(rowid, content) VALUES (new.id, new.content);
END;

CREATE TRIGGER IF NOT EXISTS clipboard_history_ad AFTER DELETE ON clipboard_history BEGIN
  INSERT INTO clipboard_history_fts(clipboard_history_fts, rowid, content) VALUES('delete', old.id, old.content);
END;

CREATE TRIGGER IF NOT EXISTS clipboard_history_au AFTER UPDATE ON clipboard_history BEGIN
  INSERT INTO clipboard_history_fts(clipboard_history_fts, rowid, content) VALUES('delete', old.id, old.content);
  INSERT INTO clipboard_history_fts(rowid, content) VALUES (new.id, new.content);
END;

