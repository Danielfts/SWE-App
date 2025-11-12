CREATE TABLE stocks (
    id INT PRIMARY KEY DEFAULT unique_rowid(),
    ticker varchar,
    target_from decimal(15,2),
    target_to decimal(15,2),
    company varchar,
    action varchar,
    brokerage varchar,
    rating_from varchar,
    rating_to varchar,
    time timestamptz
);