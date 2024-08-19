CREATE TABLE memos (
  id UUID PRIMARY KEY, -- Memo ID
  user_id TEXT NOT NULL, -- User ID
  content TEXT NOT NULL, -- Memo content
  created_at TIMESTAMP NOT NULL -- Creation timestamp
);