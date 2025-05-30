CREATE TABLE IF NOT EXISTS messages (
  id UUID PRIMARY KEY,
  sent TIMESTAMP(0) WITH TIME ZONE NOT NULL DEFAULT NOW(),
  chat_id UUID NOT NULL REFERENCES chats (id) ON DELETE CASCADE,
  user_id UUID NOT NULL REFERENCES users (id) ON DELETE CASCADE,
  content TEXT NOT NULL,
  type INT NOT NULL,
  deleted BOOL DEFAULT FALSE
)
