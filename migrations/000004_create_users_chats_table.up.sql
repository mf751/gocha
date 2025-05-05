CREATE TABLE IF NOT EXISTS users_chats (
  chat_id UUID REFERENCES chats (id) ON DELETE CASCADE,
  user_id UUID REFERENCES users (id) ON DELETE CASCADE,
  is_admin BOOL NOT NULL,
  PRIMARY KEY (chat_id, user_id)
)
