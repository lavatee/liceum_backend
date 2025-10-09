CREATE TABLE events (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL UNIQUE,
    description VARCHAR(500),
);

CREATE TABLE event_blocks (
    id SERIAL PRIMARY KEY,
    event_id INT NOT NULL,
    start_date TIMESTAMP NOT NULL,
    end_date TIMESTAMP NOT NULL,
    name VARCHAR(255),
    description VARCHAR(500),
    link VARCHAR(255) NOT NULL
);

ALTER TABLE event_blocks ADD FOREIGN KEY event_id REFERENCES events(id);