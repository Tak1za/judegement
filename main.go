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
	Winner  string
	Players int
}

// Round struct
type Round struct {
	cardsToBeDealt int
	number         int
	cards          map[string][]deck.Card
	dealt          map[string]deck.Card
	handEstimate   map[string]int
	trump          deck.Suit
	handCounts     map[string]int
	winners        []string
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
		round.number = i + 1
		round.cardsToBeDealt = (52/numPlayers - i)
		fmt.Println("Cards to be dealt: ", round.cardsToBeDealt)
		round.dealCards(writer)
		maxEstimatePlayer := round.getHandEstimate(writer)
		fmt.Printf("What do you want as the trump %s?: \n", maxEstimatePlayer)
		for index, suit := range []string{"Spades", "Diamonds", "Clubs", "Hearts"} {
			fmt.Printf("%s. %s\n", strconv.Itoa(index), suit)
		}
		fmt.Scanf("%d", &round.trump)
		for round.trump >= 4 || round.trump < 0 {
			fmt.Printf("Invalid trump. Please select again: ")
			fmt.Scanf("%d", &round.trump)
		}
		fmt.Println("Trump for this round is: ", round.trump)

		round.startRound()

		round.findRoundWinners()
		round.updateScores(writer)
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

func (r *Round) getHandEstimate(w *csv.Writer) string {
	r.handEstimate = make(map[string]int, gs.Players)
	scoreData[r.number] = append(scoreData[r.number], fmt.Sprint(r.cardsToBeDealt))
	maxEstimate := -1
	maxEstimatePlayer := "P0"
	handEstimateSum := 0
	for i := 0; i < gs.Players; i++ {
		playerID := "P" + strconv.Itoa(i)
		fmt.Printf("How many hands can you make %s?: ", playerID)
		var handEstimate int
		fmt.Scanf("%d", &handEstimate)

		for handEstimate > r.cardsToBeDealt || handEstimate < 0 {
			fmt.Printf("Invalid estimate. Estimate again: ")
			fmt.Scanf("%d", &handEstimate)
		}
		handEstimateSum += handEstimate

		// find the player with maximum hands, will be trump setter
		if handEstimate > maxEstimate {
			maxEstimate = handEstimate
			maxEstimatePlayer = playerID
		}
		r.handEstimate[playerID] = handEstimate
		scoreData[r.number] = append(scoreData[r.number], fmt.Sprint(handEstimate))
	}
	err := w.Write(scoreData[r.number])
	w.Flush()
	err = w.Error()
	if err != nil {
		panic(err)
	}

	return maxEstimatePlayer
}

func (r *Round) startRound() {
	fmt.Println("Round start...")
	r.handCounts = make(map[string]int)
	for k := 0; k < r.cardsToBeDealt; k++ {
		r.dealt = make(map[string]deck.Card, gs.Players)
		for i := 0; i < gs.Players; i++ {
			playerIndex := "P" + strconv.Itoa(i)
			var cardPlayedNumber int
			fmt.Println(playerIndex + "'s chance: ")
			for j, card := range r.cards[playerIndex] {
				fmt.Println(strconv.Itoa(j) + ". " + card.String())
			}
			fmt.Println("Choose one: ")
			fmt.Scanf("%d", &cardPlayedNumber)
			for cardPlayedNumber > r.cardsToBeDealt-k-1 || cardPlayedNumber < 0 {
				fmt.Printf("Invalid card. Please select again: ")
				fmt.Scanf("%d", &cardPlayedNumber)
			}
			cardPlayed := r.cards[playerIndex][cardPlayedNumber]
			r.dealt[playerIndex] = cardPlayed
			r.cards[playerIndex] = removeFromPlayerDeck(r.cards[playerIndex], cardPlayedNumber)
		}

		fmt.Println("Cards played this hand: ")
		for player, card := range r.dealt {
			fmt.Printf("Player: %s, Card: %+v\n", player, card)
		}

		order := []string{"P0", "P1", "P2", "P3", "P4"}
		handWinner := r.handWinner(order)
		fmt.Println("Hand Winner: ", handWinner)

		r.handCounts[handWinner.Player]++
	}

	fmt.Printf("Hand Counts: %+v\n", r.handCounts)
}

// not making use of pointer to the playerDeck for sake of clean code
func removeFromPlayerDeck(playerDeck []deck.Card, index int) []deck.Card {
	copy(playerDeck[index:], playerDeck[index+1:])
	playerDeck[len(playerDeck)-1] = deck.Card{}
	playerDeck = playerDeck[:len(playerDeck)-1]
	return playerDeck
}

func (r *Round) handWinner(order []string) WinnerHand {
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

func (r *Round) findRoundWinners() {
	r.winners = make([]string, 0)
	for player, estimate := range r.handEstimate {
		if estimate == r.handCounts[player] {
			r.winners = append(r.winners, player)
		}
	}
}

func (r *Round) updateScores(w *csv.Writer) {
	if len(r.winners) <= 0 {
		for i := 0; i < gs.Players; i++ {
			scoreData[r.number] = append(scoreData[r.number], fmt.Sprint(0))
		}
	} else {
		for i := 0; i < gs.Players; i++ {
			var updated bool = false
			for _, winner := range r.winners {
				winnerID := winner[1:]
				winnerIDInt, _ := strconv.Atoi(winnerID)
				if i == winnerIDInt {
					updated = true
					scoreData[r.number][i] = "1" + scoreData[r.number][i]
					break
				}
			}
			if !updated {
				scoreData[r.number][i] = fmt.Sprint(0)
			}
		}
	}
}
