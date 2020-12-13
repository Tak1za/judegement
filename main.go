package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"

	"github.com/Tak1za/go-deck"
)

var gs GameState
var scoreData [][]string

// GameState struct
type GameState struct {
	Winner  map[string]deck.Card
	Players int
}

// Winner struct
type Winner struct {
	Player string
	Hand   deck.Card
}

// Round struct
type Round struct {
	cardsToBeDealt int
	number         int
	cards          map[string][]deck.Card
	winner         Winner
	dealt          map[string]deck.Card
	handEstimate   map[string]int
	trump          deck.Suit
}

// WinnerHand struct
type WinnerHand struct {
	Player string
	Hand   deck.Card
}

func main() {
	scoreFile, err := os.Create("./scoring.csv")
	if err != nil {
		panic(err)
	}
	writer := csv.NewWriter(scoreFile)

	var numPlayers int
	fmt.Println("Please enter the number of players: ")
	fmt.Scanf("%d", &numPlayers)

	scoreData = make([][]string, (52/numPlayers)+2)

	gs.Players = numPlayers
	addHeader(writer)

	for i := 0; i < 52/numPlayers; i++ {
		var round Round
		round.cardsToBeDealt = (52/numPlayers - i)
		fmt.Println("Cards to be dealt: ", round.cardsToBeDealt)
		round.dealCards(writer)
		round.getHandEstimate(writer)
		maxEstimatePlayer := "P0"
		maxEstimate := round.handEstimate[maxEstimatePlayer]
		for playerID, estimate := range round.handEstimate {
			if estimate > maxEstimate {
				maxEstimate = estimate
				maxEstimatePlayer = playerID
			}
		}
		fmt.Printf("What do you want as the trump %s?: \n", maxEstimatePlayer)
		for index, suit := range []string{"Spades", "Diamonds", "Clubs", "Hearts"} {
			fmt.Printf("%s. %s\n", strconv.Itoa(index), suit)
		}
		fmt.Scanf("%d", &round.trump)
		round.startRound()
		order := []string{"P0", "P1", "P2", "P3", "P4"}
		roundWinner := round.dealWinner(order)
		fmt.Println("Deal Winner: ", roundWinner)
	}
}

func addHeader(w *csv.Writer) {
	scoreData[0] = append(scoreData[0], "Cards")
	for i := 0; i < gs.Players; i++ {
		playerID := "P" + strconv.Itoa(i)
		scoreData[0] = append(scoreData[0], playerID)
	}
	err := w.Write(scoreData[0])
	if err != nil {
		panic(err)
	}
	w.Flush()
	err = w.Error()
	if err != nil {
		panic(err)
	}
}

func (r *Round) dealCards(w *csv.Writer) {
	newDeck := deck.New(deck.Shuffle)
	r.cards = make(map[string][]deck.Card)
	for i := 1; i <= r.cardsToBeDealt; i++ {
		for j := i*gs.Players - gs.Players; j < i*gs.Players; j++ {
			playerID := "P" + strconv.Itoa(j%gs.Players)
			r.cards[playerID] = append(r.cards[playerID], newDeck[j])
		}
	}
}

func (r *Round) getHandEstimate(w *csv.Writer) {
	var handEstimate int
	r.handEstimate = make(map[string]int, gs.Players)
	scoringRow := (52 / gs.Players) + 1 - r.cardsToBeDealt
	scoreData[scoringRow] = append(scoreData[scoringRow], fmt.Sprint(r.cardsToBeDealt))
	for i := 0; i < gs.Players; i++ {
		playerID := "P" + strconv.Itoa(i)
		fmt.Printf("How many hands can you make %s?: ", playerID)
		fmt.Scanf("%d", &handEstimate)
		r.handEstimate[playerID] = handEstimate
		scoreData[scoringRow] = append(scoreData[scoringRow], fmt.Sprint(handEstimate))
	}
	err := w.Write(scoreData[scoringRow])
	w.Flush()
	err = w.Error()
	if err != nil {
		panic(err)
	}
}

func (r *Round) startRound() {
	fmt.Println("Round start...")
	r.dealt = make(map[string]deck.Card)
	for i := 0; i < gs.Players; i++ {
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

func (r *Round) dealWinner(order []string) WinnerHand {
	winnerCard := r.dealt[order[0]]
	winnerPlayer := order[0]
	for i := 0; i < gs.Players-1; i++ {
		if (winnerCard.Suit == r.dealt[order[i+1]].Suit && winnerCard.Rank < r.dealt[order[i+1]].Rank) || (winnerCard.Suit != r.dealt[order[i+1]].Suit && r.dealt[order[i+1]].Suit == r.trump) {
			winnerCard = r.dealt[order[i+1]]
			winnerPlayer = order[i+1]
		}
	}

	return WinnerHand{
		Player: winnerPlayer,
		Hand:   winnerCard,
	}
}
