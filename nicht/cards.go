package main

// TODO:
// DONE: seperate library and game logic, very messy atm
// DONE: fix rand seed setup so I can optionally generate new shuffles/games each runtime.
// DONE: log game actions to file, for checking of game logic
// DONE: lots of things should be interfaces, for further game dev: Shuffler, Player, Scorer?
// implement naive AI, instead of taking random cards, take the card that increases personal score the most
// much further on: implement 'smart' AI, take the card that increases score delta by the most
// (vs. currently highest scoring player?)

import (
	"fmt"
	"math/rand"

	"github.com/frenata/gaga"
)

const (
	blue   = "Blue"
	green  = "Green"
	yellow = "Yellow"
	red    = "Red"

	special = "special"
	normal  = "normal"
)

var colors = [4]string{blue, green, yellow, red}
var vSpecial = [5]int{-1, -1, -1, 0, 2}
var vNormal = [10]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

var nichtPlayers = []*NichtPlayer{
	NewNichtPlayer("P1"),
	NewNichtPlayer("P2"),
	NewNichtPlayer("P3"),
	NewNichtPlayer("P4"),
}

// could load from a config file
func NewNichtDeck() *gaga.Deck {
	c := make([]gaga.Card, 60)
	n := 0
	for _, color := range colors {
		for _, i := range vNormal {
			c[n] = NewNichtCard(color, normal, i)
			n++
		}
		for _, i := range vSpecial {
			c[n] = NewNichtCard(color, special, i)
			n++
		}
	}
	return gaga.NewDeck(c)
}

type NichtCard struct {
	Color string
	Cat   string
	Value int
	play  *NichtPlayer
}

func NewNichtCard(color, cat string, value int) *NichtCard {
	c := new(NichtCard)
	c.Color = color
	c.Cat = cat
	c.Value = value

	return c
}

func (c *NichtCard) String() string {
	s := c.Color[:1]
	v := fmt.Sprint(c.Value)
	if c.Cat == special {
		return s + "x" + v
	} else { //if cat == normal {
		return s + v
	}
}

// On further thought, this abstraction isn't producing any benefit.
// NichtCard is merely fulfilling the interface mindlessly,
// the actual game logic in game.go is interacting with the struct
// directly.
// Having the Card interface inherently related to a Player might
// be worthwhile, but only if gaga.Deck also holds a []Player, so that
// by itself gaga.Deck could handle removing Player pointers to cards
// and reshuffling those cards into the deck. That seems a small benefit
// for an unwieldy extra abstraction.
/*
func (c *NichtCard) PlayedBy(p gaga.Player) gaga.Player {
	if p == nil {
		c.play = nil
	} else {
		c.play = p.(*NichtPlayer)
	}
	return c.play
}
*/

type NichtPlayer struct {
	Name  string
	Hand  []gaga.Card
	Table []gaga.Card
	Score int
}

func NewNichtPlayer(name string) *NichtPlayer {
	p := new(NichtPlayer)
	p.Name = name
	return p
}

func (p *NichtPlayer) PlayRand() *NichtCard {
	n := rand.Intn(len(p.Hand))
	c := p.Hand[n].(*NichtCard)
	p.Hand = append(p.Hand[:n], p.Hand[n+1:]...)
	c.play = p

	return c
}

func (p *NichtPlayer) AddTable(c *NichtCard) {
	p.Table = append(p.Table, c)
}

func (p *NichtPlayer) AddCard(c gaga.Card) {
	p.Hand = append(p.Hand, c.(*NichtCard))
}

func (p *NichtPlayer) PrintHand() string {
	var gc []gaga.Card
	for _, c := range p.Hand {
		gc = append(gc, c)
	}
	return gaga.PrintCards(gc)
}
func (p *NichtPlayer) PrintTable() string {
	var gc []gaga.Card
	for _, c := range p.Table {
		gc = append(gc, c)
	}
	return gaga.PrintCards(gc)
}

func (p *NichtPlayer) String() string {
	return p.Name
}
