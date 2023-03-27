CREATE TABLE Users (
    id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
    user TINYTEXT NOT NULL,
    hash TEXT NOT NULL
);

CREATE TABLE Messages (
    id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
    recipient INTEGER NOT NULL REFERENCES Users(id),
    sender INTEGER NOT NULL REFERENCES Users(id),
    data TEXT NOT NULL,
    hash TEXT NOT NULL
);
