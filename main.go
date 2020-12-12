package main

import (
	"fmt"
	"strconv"

	"github.com/Tak1za/go-deck"
)

var gs GameState

// GameState struct
type GameState struct {
	Winner     map[string]deck.Card
	CardsDealt int
	Players    int
}

// Winner struct
type Winner struct {
	Player string
	Hand   deck.Card
}

// Round struct
type Round struct {
	cards  map[string][]deck.Card
	winner Winner
	dealt  map[string]deck.Card
}

func main() {
	var numPlayers int
	fmt.Println("Please enter the number of players: ")
	fmt.Scanf("%d", &numPlayers)

	gs.Players = numPlayers
	currentRound := dealCards(numPlayers)
	currentRound.startRound(numPlayers)
}

func dealCards(numPlayers int) *Round {
	newDeck := deck.New(deck.Shuffle)
	var round Round
	round.cards = make(map[string][]deck.Card)
	gs.CardsDealt = 52 / numPlayers
	for i := 1; i <= (52 / numPlayers); i++ {
		for j := i*numPlayers - numPlayers; j < i*numPlayers; j++ {
			playerID := "P" + strconv.Itoa(j%numPlayers)
			round.cards[playerID] = append(round.cards[playerID], newDeck[j])
			fmt.Printf("Dealt card: %+v, to: %s\n", newDeck[j], playerID)
		}
	}

	return &round
}

func (r *Round) startRound(numPlayers int) {
	fmt.Println("Round start...")
	r.dealt = make(map[string]deck.Card)
	for i := 0; i < numPlayers; i++ {
		playerIndex := "P" + strconv.Itoa(i)
		var cardPlayedNumber int
		fmt.Println(playerIndex + "'s chance: ")
		for j, card := range r.cards[playerIndex] {
			fmt.Println(strconv.Itoa(j) + ". " + card.String())
		}
		fmt.Println("Choose one: ")
		fmt.Scanf("%d", &cardPlayedNumber)
		cardPlayed := r.cards[playerIndex][cardPlayedNumber]
		r.dealt[playerIndex] = cardPlayed
	}

	fmt.Println("Cards played this round: ")
	for player, card := range r.dealt {
		fmt.Printf("Player: %s, Card: %+v\n", player, card)
	}
}
