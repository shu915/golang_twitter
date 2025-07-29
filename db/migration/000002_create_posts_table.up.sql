CREATE TABLE posts (
  id SERIAL PRIMARY KEY,
  user_id INT NOT NULL,
  content VARCHAR(140) NOT NULL,
  created_at TIMESTAMPTZ DEFAULT now()
);