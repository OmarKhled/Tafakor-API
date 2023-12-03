-- The initial setup for tafakor's DB

-- Contains all the verses (ayahs) that are scheduled for posting
CREATE TABLE IF NOT EXISTS verse (
    id varchar PRIMARY KEY NOT NULL,
    surah_number int8 NOT NULL,
    start SMALLINT NOT NULL, -- the start verse
    "end" SMALLINT NOT NULL, -- the end verse (same as start verse if only one ayah) | The longest surah in quran is 280 veres (Al Baqarah) so SMALL INT is more than enough
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Contains all the posted posts references the ayah used in each one
CREATE TABLE IF NOT EXISTS post (
    id uuid NOT NULL PRIMARY KEY,
    verse uuid NOT NULL references verse,
    published bool,
    footage varchar,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Contains all used stock footages references the post used in
CREATE TABLE IF NOT EXISTS stock_footage (
    id uuid PRIMARY KEY NOT NULL,
    post uuid NOT NULL references posting,
    provider varchar,
    provider_id varchar,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
)
-- stock_footage table is cruicial for ensuring that no stock footages are used in a post more than once (in a specific time interval)