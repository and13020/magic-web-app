package queries

const InsertCard = `INSERT INTO card (
	id, 
	object, 
	oracle_id, 
	multiverse_ids, 
	mtgo_id, 
	mtgo_foil_id, 
	tcgplayer_id, 
	cardmarket_id,	
	name, 
	lang, 
	released_at, 
	uri, 
	scryfall_uri, 
	layout, 
	highres_image, 
	image_status, 
	mana_cost, 
	cmc, 
	type_line, 
	oracle_text, 
	colors, 
	color_identity, 
	keywords, 
	games, 
	reserved, 
	game_changer, 
	foil, 
	nonfoil, 
	finishes, 
	oversized, 
	promo, 
	reprint, 
	variation, 
	set_id, 
	"set", 
	set_name, 
	set_type, 
	set_uri, 
	set_search_uri, 
	scryfall_set_uri, 
	rulings_uri, 
	prints_search_uri, 
	collector_number, 
	digital, 
	rarity, 
	flavor_text, 
	card_back_id, 
	artist, 
	artist_ids, 
	illustration_id, 
	border_color, 
	frame, 
	full_art, 
	textless, 
	booster, 
	story_spotlight, 
	edhrec_rank, 
	penny_rank)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?,
		?, ?, ?, ?, ?, ?, ?, ?, ?, ?,
		?, ?, ?, ?, ?, ?, ?, ?, ?,
		?, ?, ?, ?, ?, ?, ?, ?, ?, ?,
		?, ?, ?, ?, ?, ?, ?, ?, ?, ?,
		?, ?, ?, ?, ?, ?, ?, ?, ?);`

const InsertImageUris = `INSERT INTO image_uris (
	card_id, small, normal, large, png, art_crop, border_crop)
	VALUES (?, ?, ?, ?, ?, ?, ?);`

// InsertLegalities has 22 fields including card_id
const InsertLegalities = `INSERT INTO legalities (
	card_id, 
	standard, 
	future, 
	historic, 
	timeless, 
	gladiator, 
	pioneer, 
	modern, 
	legacy, 
	pauper, 
	vintage, 
	penny, 
	commander, 
	oathbreaker, 
	standardbrawl, 
	brawl, 
	alchemy, 
	paupercommander, 
	duel, 
	oldschool, 
	premodern, 
	predh
) VALUES (
	?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);`

const InsertPrices = `INSERT INTO prices (
	card_id,
	usd,
	usd_foil,
	usd_etched,
	eur,
	eur_foil,
	tix
) VALUES (
	?, ?, ?, ?, ?, ?, ?);`

const InsertRelatedUris = `INSERT INTO related_uris (
	card_id,
	gatherer,
	tcgplayer_infinite_articles,
	tcgplayer_infinite_decks,
	edhrec
) VALUES (
	?, ?, ?, ?, ?);`

const InsertPurchaseUris = `INSERT INTO purchase_uris (
	card_id,
	tcgplayer,
	cardmarket,
	cardhoarder
) VALUES (
	?, ?, ?, ?);`

// InsertCardQueries contains insert queries in this order:
//
// 1. card
//
// 2. image uri
//
// 3. legality
//
// 4. price
//
// 5. related uri
//
// 6. purchased uri
var InsertCardQueries = []string{
	InsertCard,
	InsertImageUris,
	InsertLegalities,
	InsertPrices,
	InsertRelatedUris,
	InsertPurchaseUris,
}
