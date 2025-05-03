CREATE TABLE IF NOT EXISTS messages
(
    id        SERIAL PRIMARY KEY,
    user_id   INT          NOT NULL,
    username  VARCHAR(255) NOT NULL,
    content   TEXT         NOT NULL,
    timestamp TIMESTAMP    NOT NULL
);
