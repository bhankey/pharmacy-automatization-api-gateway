CREATE OR REPLACE FUNCTION update_time_column()
    RETURNS TRIGGER AS $$
BEGIN
    NEW.update_time = now();
RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TABLE IF NOT EXISTS refresh_tokens (
                    id bigserial NOT NULL PRIMARY KEY,
                    user_id int REFERENCES users(id),
                    refresh_token text UNIQUE NOT NULL,
                    user_agent text NOT NULL,
                    ip text NOT NULL,
                    finger_print text NOT NULL,
                    is_available bool NOT NULL DEFAULT true,
                    creation_time timestamp NOT NULL DEFAULT NOW()
);
