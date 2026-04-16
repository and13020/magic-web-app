package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	q "magic/db_ops"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type CardRepository struct {
	db *sql.DB
}

type CardRepositoryInterface interface {
	GetCardsByName(name string) ([]ScryCard, error)
	GetRandomCard() (ScryCard, error)
	SaveCard(card ScryCard) error
	// UpdateCard(card Card) error
	// DeleteCard(id string) error
}

func NewCardRepository(db *sql.DB) *CardRepository {
	return &CardRepository{db: db}
}

// getCardsByName searches for cards with given string in their name
// if multiple results, print all
// if no results, return 404
// Documentation: https://scryfall.com/docs/api/cards/search
func (r *CardRepository) GetCardsByName(name string) ([]Card, error) {

	if len(name) > 999 {
		log.Fatal("Name is too long, must be less than 1000 characters")
	}

	// url encode the given name
	name = url.QueryEscape(name)
	req, err := http.NewRequest("GET", "https://api.scryfall.com/cards/search?q="+name, nil)
	if err != nil {
		return nil, err
	}

	// Create and set header, scryfall required values
	h := make(http.Header)
	h.Set("Accept", "application/json")
	h.Set("User-Agent", "myApp/1.0")
	req.Header = h

	c := http.Client{Timeout: 10 * time.Second}

	resp, err := c.Do(req)
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Status NOT success: %v", resp.StatusCode)
	} else if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	fmt.Printf("Status code: %v\n", resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Failed to read into body: ", err)
	}

	sCards := ScryCards{}
	cards := []Card{}

	err = json.Unmarshal(body, &sCards)
	if err != nil {
		return nil, err
	}

	// convert to our cards (safe for DB insertion)
	cards = convertToCard(sCards)

	for index, card := range cards {
		fmt.Printf("Card %v: %v\n", index+1, card.Name)
	}
	if len(cards) < 1 {
		fmt.Println("YOOOO we got no cards back")
	}

	return cards, nil
}

// GetRandomCard returns a random card from scryfall, converts it to Card, saves in DB then returns it
func (r *CardRepository) GetRandomCard() (Card, error) {
	c := ScryCard{}

	req, err := http.NewRequest("GET", "https://api.scryfall.com/cards/random", nil)
	if err != nil {
		return Card{}, err
	}

	// Create and set header, scryfall required values
	h := make(http.Header)
	h.Set("Accept", "application/json")
	h.Set("User-Agent", "myApp/1.0")
	req.Header = h

	client := http.Client{Timeout: 5 * time.Second}

	resp, err := client.Do(req)
	if err != nil {
		return Card{}, err
	}
	if resp.StatusCode != http.StatusOK || resp == nil {
		return Card{}, fmt.Errorf("Status NOT success or missing: %v\n", resp.StatusCode)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return Card{}, err
	}

	err = json.Unmarshal(body, &c)
	if err != nil {
		return Card{}, err
	}

	nc := convertToCard(ScryCards{Cards: []ScryCard{c}})

	return nc[0], nil
}

// SaveCard will passively store card data in our DB (cache)
// as other methods are called such as when any card is retrieved,
// it will be cached
func (r *CardRepository) SaveCard(c Card) error {

	fmt.Println("Entered save card")
	ctx, cancelCtx := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelCtx()

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		fmt.Println("Could not prepare txn")
		return err
	}

	// atomic txn - all card details are inserted, or none are
	defer tx.Rollback()

	for i, query := range q.InsertCardQueries {
		stmt, err := tx.Prepare(query)
		if err != nil {
			return err
		}
		defer stmt.Close()

		switch i {
		case 0: // card
			_, err = stmt.Exec(
				c.ID,
				c.Object,
				c.OracleID,
				c.MultiverseIDs,
				c.MtgoID,
				c.MtgoFoilID,
				c.TcgplayerID,
				c.CardmarketID,
				c.Name,
				c.Lang,
				c.ReleasedAt,
				c.Uri,
				c.ScryfallUri,
				c.Layout,
				c.HighresImage,
				c.ImageStatus,
				c.ManaCost,
				c.Cmc,
				c.TypeLine,
				c.OracleText,
				c.Colors,
				c.ColorIdentity,
				c.Keywords,
				c.Games,
				c.Reserved,
				c.GameChanger,
				c.Foil,
				c.Nonfoil,
				c.Finishes,
				c.Oversized,
				c.Promo,
				c.Reprint,
				c.Variation,
				c.SetID,
				c.Set,
				c.SetName,
				c.SetType,
				c.SetUri,
				c.SetSearchUri,
				c.ScryfallSetUri,
				c.RulingsUri,
				c.PrintsSearchUri,
				c.CollectorNumber,
				c.Digital,
				c.Rarity,
				c.FlavorText,
				c.CardBackID,
				c.Artist,
				c.ArtistIDs,
				c.IllustrationID,
				c.BorderColor,
				c.Frame,
				c.FullArt,
				c.Textless,
				c.Booster,
				c.StorySpotlight,
				c.EdhrecRank,
				c.PennyRank,
			)
		case 1: // image uri
			_, err = stmt.Exec(
				c.ID,
				c.ImageUris.Small,
				c.ImageUris.Normal,
				c.ImageUris.Large,
				c.ImageUris.Png,
				c.ImageUris.ArtCrop,
				c.ImageUris.BorderCrop,
			)

		case 2: // legalities
			_, err = stmt.Exec(
				c.ID,
				c.Legalities.Standard,
				c.Legalities.Future,
				c.Legalities.Historic,
				c.Legalities.Timeless,
				c.Legalities.Gladiator,
				c.Legalities.Pioneer,
				c.Legalities.Modern,
				c.Legalities.Legacy,
				c.Legalities.Pauper,
				c.Legalities.Vintage,
				c.Legalities.Penny,
				c.Legalities.Commander,
				c.Legalities.Oathbreaker,
				c.Legalities.StandardBrawl,
				c.Legalities.Brawl,
				c.Legalities.Alchemy,
				c.Legalities.PauperCommander,
				c.Legalities.Duel,
				c.Legalities.Oldschool,
				c.Legalities.Premodern,
				c.Legalities.Predh,
			)

		case 3: // prices
			_, err = stmt.Exec(
				c.ID,
				c.Prices.Usd,
				c.Prices.UsdFoil,
				c.Prices.UsdEtched,
				c.Prices.Eur,
				c.Prices.EurFoil,
				c.Prices.Tix,
			)
		case 4: // related uri
			_, err = stmt.Exec(
				c.ID,
				c.RelatedUris.Gatherer,
				c.RelatedUris.TcgplayerInfiniteArticles,
				c.RelatedUris.TcgplayerInfiniteDecks,
				c.EdhrecRank,
			)
		case 5: // purchase uri
			_, err = stmt.Exec(
				c.ID,
				c.PurchaseUris.Tcgplayer,
				c.PurchaseUris.Cardmarket,
				c.PurchaseUris.Cardhoarder,
			)
		}

		if err != nil {
			fmt.Printf("Failed to save query #%d DB: %v\n", i, query)
			return err
		}

	}

	tx.Commit()
	return nil
}

func convertToCard(scryCards ScryCards) []Card {

	var cards []Card

	for _, sc := range scryCards.Cards {

		var mID strings.Builder
		for _, m := range sc.MultiverseIDs {
			mID.WriteString(strconv.Itoa(m))
		}

		c := Card{
			Object:          sc.Object,
			ID:              sc.ID,
			OracleID:        sc.OracleID,
			MultiverseIDs:   mID.String(),
			MtgoID:          sc.MtgoID,
			MtgoFoilID:      sc.MtgoFoilID,
			TcgplayerID:     sc.TcgplayerID,
			CardmarketID:    sc.CardmarketID,
			Name:            sc.Name,
			Lang:            sc.Lang,
			ReleasedAt:      sc.ReleasedAt,
			Uri:             sc.Uri,
			ScryfallUri:     sc.ScryfallUri,
			Layout:          sc.Layout,
			HighresImage:    sc.HighresImage,
			ImageStatus:     sc.ImageStatus,
			ImageUris:       sc.ImageUris,
			ManaCost:        sc.ManaCost,
			Cmc:             sc.Cmc,
			TypeLine:        sc.TypeLine,
			OracleText:      sc.OracleID,
			Colors:          strings.Join(sc.Colors, ","),
			ColorIdentity:   strings.Join(sc.ColorIdentity, ","),
			Keywords:        strings.Join(sc.Keywords, ","),
			Legalities:      sc.Legalities,
			Games:           strings.Join(sc.Games, ","),
			Reserved:        sc.Reserved,
			GameChanger:     sc.GameChanger,
			Foil:            sc.Foil,
			Nonfoil:         sc.Nonfoil,
			Finishes:        strings.Join(sc.Finishes, ","),
			Oversized:       sc.Oversized,
			Promo:           sc.Promo,
			Reprint:         sc.Reprint,
			Variation:       sc.Variation,
			SetID:           sc.SetID,
			Set:             sc.Set,
			SetName:         sc.SetName,
			SetType:         sc.SetType,
			SetUri:          sc.SetUri,
			SetSearchUri:    sc.SetSearchUri,
			ScryfallSetUri:  sc.ScryfallSetUri,
			RulingsUri:      sc.RulingsUri,
			PrintsSearchUri: sc.PrintsSearchUri,
			CollectorNumber: sc.CollectorNumber,
			Digital:         sc.Digital,
			Rarity:          sc.Rarity,
			FlavorText:      sc.FlavorText,
			CardBackID:      sc.CardBackID,
			Artist:          sc.Artist,
			ArtistIDs:       strings.Join(sc.ArtistIDs, ","),
			IllustrationID:  sc.IllustrationID,
			BorderColor:     sc.BorderColor,
			Frame:           sc.Frame,
			FullArt:         sc.FullArt,
			Textless:        sc.Textless,
			Booster:         sc.Booster,
			StorySpotlight:  sc.StorySpotlight,
			EdhrecRank:      sc.EdhrecRank,
			PennyRank:       sc.PennyRank,
			Prices:          sc.Prices,
			RelatedUris:     sc.RelatedUris,
			PurchaseUris:    sc.PurchaseUris,
		}
		cards = append(cards, c)
	}

	return cards
}
