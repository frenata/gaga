// Package deck defines some useful interfaces and functions for manipulating card decks in games.
package deck

import (
	"math/rand"
	"time"
)

/* Card interface defines a object that represents a card in a card game
   The implementing struct should be able to represent the Card via the String()
   function.
   Cards can be collected and shuffled in a Deck and be held by Players.

   TODO: Should implementations store information about what Player has played a Card?
*/
type Card interface {
	//PlayedBy(Player) Player
	String() string
}

// Player defines an object that can be given Cards, and can print the cards it holds
// as a String.
// TODO: String() should be removed or renamed to Hand() string or Cards() string,
// the same implementations kept, to better represent how it is being implemented.
type Player interface {
	AddCard(c Card)
	String() string
}

// PrintCards takes a slice of Cards and returns a single string with each Card printed
// and seperated by spaces.
// For use as a helper function for structs that implement Card,
func PrintCards(stack []Card) string {
	var s string

	for _, c := range stack {
		if s != "" {
			s = s + " " + c.String()
		} else {
			s = c.String()
		}
	}
	return s
}

// Deck collects Cards into a usable collection, suitable for shuffling and dealing.
//	cards - all the cards that have been shuffled but not yet dealt to Players.
//	discards - all the cards that have been returned to the deck
type Deck struct {
	cards    []Card //cards currently in the deck
	discards []Card //cards currently in the deck
}

// New accepts a slice of Cards and creates a new Deck ready for use.
// Cards and Shuffled are both populated with equivalent slices.
// Shuffle must be called explicitly after a NewDeck is returned to shuffle the cards.
func New(cards []Card) *Deck {
	d := new(Deck)
	d.cards = make([]Card, len(cards))
	copy(d.cards, cards)
	d.discards = make([]Card, 0, len(d.cards))
	return d
}

// Cards returns a list of the cards in the deck.
func (d *Deck) Cards() []Card {
	cards := make([]Card, len(d.cards))
	copy(cards, d.cards)
	return cards
}

// Discards returns a list of the cards in the deck.
func (d *Deck) Discards() []Card {
	discards := make([]Card, len(d.discards))
	copy(discards, d.discards)
	return discards
}

// String prints all the *shuffled* cards in the deck.
func (d *Deck) String() string {
	return PrintCards(d.cards)
}

// Shuffle takes a seed and randomizes the cards contained in the Shuffled slice.
// TODO: refactor to remove seed argument. Allow user to set the seed once and directly,
// via a seperate method.
func (d *Deck) Shuffle(seed int) {
	rnd := deckSeed(seed)

	var toshuffle []Card
	toshuffle = append(toshuffle, d.cards...)
	toshuffle = append(toshuffle, d.discards...)

	d.discards = make([]Card, 0, len(toshuffle))

	shuffled := make([]Card, len(toshuffle))
	r := rnd.Perm(len(toshuffle))
	j := 0
	for _, i := range r {
		shuffled[j] = toshuffle[i]
		j++
	}
	d.cards = shuffled
}

// ReturnCards takes a slice of cards (from a Player?) and adds them back into Shuffled.
// TODO: This may be better represented by a seperate slice, Discards?
// Likewise, since the implication is that these cards were previously dealt by
// this deck, there should be error checking to verify: if a card was not previously
// dealt, an error should be returned. And if it was, remove it from Dealt slice.
func (d *Deck) Discard(cards []Card) {
	for _, c := range cards {
		d.discards = append(d.discards, c)
	}
}

// DealAll takes a slice of Players and deals cards to each in turn until
// none are left.
// This does *not* guarantee equal dealing.
// Returns number of cards dealt.
func (d *Deck) DealAll(players []Player) (n int) {
	for {
		for _, p := range players {
			if d.Deal(p) {
				n++
			} else { // if Deal fails, no more shuffled
				return n
			}
		}
	}
}

// Deal adds a card from the shuffled cards to a Player.
// Returns true if a card was dealt, false if not. (no cards left)
func (d *Deck) Deal(p Player) bool {
	if len(d.cards) > 0 {
		p.AddCard(d.cards[0])
		d.cards = d.cards[1:]
		return true
	} else {
		return false
	}
}

// Seed sets a random seed for Deck shuffling and dealing. Passing -1 uses the current
// time, anything else is static and suitable for testing, etc.
// TODO: Rewrite to provide user access and permanent struct seed.
func deckSeed(seed int) *rand.Rand {
	s64 := int64(seed)

	if s64 == -1 {
		s64 = time.Now().UnixNano()
	}
	s := rand.NewSource(s64)
	r := rand.New(s)
	return r
}

// PopCard attempts to remove a Card from a slice of Cards, returns true if successful.
// **Not currently being used**
func popCard(c Card, s []Card) bool {
	for i, v := range s {
		if c == v {
			s = append(s[:i], s[i+1:]...)
			return true
		}
	}
	return false
}
