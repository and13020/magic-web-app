package queries

const UserSchema = `
		CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			email varchar(225) not null unique,
			username varchar(225),
			password varchar(225) not null
		);
`

const TestSchema = `
CREATE TABLE IF NOT EXISTS test (
    id VARCHAR PRIMARY KEY,
    name VARCHAR)`

const CardSchema = `
CREATE TABLE IF NOT EXISTS card (
    id VARCHAR PRIMARY KEY,
    object VARCHAR,
    oracle_id VARCHAR,
    multiverse_ids TEXT,
    mtgo_id INTEGER,
    mtgo_foil_id INTEGER,
    tcgplayer_id INTEGER,
    cardmarket_id INTEGER,
    name VARCHAR,
    lang VARCHAR,
    released_at DATE,
    uri VARCHAR,
    scryfall_uri VARCHAR,
    layout VARCHAR,
    highres_image BOOLEAN,
    image_status VARCHAR,
    mana_cost VARCHAR,
    cmc DECIMAL(5,2),
    type_line VARCHAR,
    oracle_text TEXT,
    colors TEXT,
    color_identity TEXT,
    keywords TEXT,
    legalities JSONB,
    games TEXT,
    reserved BOOLEAN,
    game_changer BOOLEAN,
    foil BOOLEAN,
    nonfoil BOOLEAN,
    finishes TEXT,
    oversized BOOLEAN,
    promo BOOLEAN,
    reprint BOOLEAN,
    variation BOOLEAN,
    set_id VARCHAR,
    "set" VARCHAR,
    set_name VARCHAR,
    set_type VARCHAR,
    set_uri VARCHAR,
    set_search_uri VARCHAR,
    scryfall_set_uri VARCHAR,
    rulings_uri VARCHAR,
    prints_search_uri VARCHAR,
    collector_number VARCHAR,
    digital BOOLEAN,
    rarity VARCHAR,
    flavor_text TEXT,
    card_back_id VARCHAR,
    artist VARCHAR,
    artist_ids TEXT,
    illustration_id VARCHAR,
    border_color VARCHAR,
    frame VARCHAR,
    full_art BOOLEAN,
    textless BOOLEAN,
    booster BOOLEAN,
    story_spotlight BOOLEAN,
    edhrec_rank INTEGER,
    penny_rank INTEGER
);`

// not sure which size we'll want yet
// maybe normal and large?
const ImageUrisSchema = `
CREATE TABLE IF NOT EXISTS image_uris (
    card_id VARCHAR PRIMARY KEY REFERENCES cards(id) ON DELETE CASCADE,
    small VARCHAR,
    normal VARCHAR,
    large VARCHAR,
    png VARCHAR,
    art_crop VARCHAR,
    border_crop VARCHAR
);`

// needs to accept 3 types
// Legal, not legal, and limited (like black lotus in legacy)
const LegalitiesSchema = `
CREATE TABLE IF NOT EXISTS legalities (
    card_id VARCHAR PRIMARY KEY REFERENCES cards(id) ON DELETE CASCADE,
    standard VARCHAR,
    future VARCHAR,
    historic VARCHAR,
    timeless VARCHAR,
    gladiator VARCHAR,
    pioneer VARCHAR,
    modern VARCHAR,
    legacy VARCHAR,
    pauper VARCHAR,
    vintage VARCHAR,
    penny VARCHAR,
    commander VARCHAR,
    oathbreaker VARCHAR,
    standardbrawl VARCHAR,
    brawl VARCHAR,
    alchemy VARCHAR,
    paupercommander VARCHAR,
    duel VARCHAR,
    oldschool VARCHAR,
    premodern VARCHAR,
    predh VARCHAR
);`

const PricesSchema = `
CREATE TABLE IF NOT EXISTS prices (
    card_id VARCHAR PRIMARY KEY REFERENCES cards(id) ON DELETE CASCADE,
    usd VARCHAR,
    usd_foil VARCHAR,
    usd_etched VARCHAR,
    eur VARCHAR,
    eur_foil VARCHAR,
    tix VARCHAR
);`

// maybe keep 1 or 2 related uris?
// not sure on infinite articles/decks being useful
// probably exclude this table, just use purcahseUrisSchema below instead
const RelatedUrisSchema = `
CREATE TABLE IF NOT EXISTS related_uris (
    card_id VARCHAR PRIMARY KEY REFERENCES cards(id) ON DELETE CASCADE,
    gatherer VARCHAR,
    tcgplayer_infinite_articles VARCHAR,
    tcgplayer_infinite_decks VARCHAR,
    edhrec VARCHAR
);`

const PurchaseUrisSchema = `
CREATE TABLE IF NOT EXISTS purchase_uris (
    card_id VARCHAR PRIMARY KEY REFERENCES cards(id) ON DELETE CASCADE,
    tcgplayer VARCHAR,
    cardmarket VARCHAR,
    cardhoarder VARCHAR
);`
