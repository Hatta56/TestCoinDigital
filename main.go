package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

type Player struct {
	ID       int
	Dice     []int
	Score    int
	HasDice  bool
	Next     *Player
	Previous *Player
}

func main() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Masukkan jumlah pemain: ")
	numPlayersStr, _ := reader.ReadString('\n')
	numPlayers, _ := strconv.Atoi(strings.TrimSpace(numPlayersStr))

	fmt.Print("Masukkan jumlah dadu: ")
	numDiceStr, _ := reader.ReadString('\n')
	numDice, _ := strconv.Atoi(strings.TrimSpace(numDiceStr))

	players := createPlayers(numPlayers, numDice)
	setupPlayerConnections(players)

	fmt.Println("==================")
	fmt.Printf("Pemain = %d, Dadu = %d\n", numPlayers, numDice)
	fmt.Println("==================")

	round := 1
	for {
		fmt.Printf("Giliran %d lempar dadu:\n", round)

		// Lemparkan dadu
		for _, player := range players {
			player.rollDice()
			fmt.Printf("%s (%d): %v\n", player.getName(), player.Score, player.Dice)
		}

		// Evaluasi hasil lemparan
		for _, player := range players {
			player.evaluateDice()
		}

		// Cek pemain yang masih memiliki dadu
		remainingPlayers := make([]*Player, 0)
		for _, player := range players {
			if player.HasDice {
				remainingPlayers = append(remainingPlayers, player)
			}
		}

		// Cetak hasil evaluasi
		fmt.Println("Setelah evaluasi:")
		for _, player := range players {
			if player.HasDice {
				fmt.Printf("%s (%d): %v\n", player.getName(), player.Score, player.Dice)
			} else {
				fmt.Printf("Pemain #%d (%d): _ (Berhenti bermain karena tidak memiliki dadu)\n", player.ID, player.Score)
			}
		}
		fmt.Println("==================")

		// Cek kondisi berhenti
		if len(remainingPlayers) == 1 {
			break
		}

		round++
	}

	fmt.Println("Game berakhir karena hanya", players[0].getName(), "yang memiliki dadu.")
	winner := getWinner(players)
	fmt.Println("Game dimenangkan oleh", winner.getName(), "karena memiliki", winner.Score, "poin lebih banyak dari pemain lainnya.")
}

func createPlayers(numPlayers, numDice int) []*Player {
	players := make([]*Player, numPlayers)
	for i := 0; i < numPlayers; i++ {
		players[i] = &Player{
			ID:      i + 1,
			Dice:    make([]int, numDice),
			Score:   0,
			HasDice: true,
		}
	}
	return players
}

func setupPlayerConnections(players []*Player) {
	numPlayers := len(players)
	for i := 0; i < numPlayers; i++ {
		players[i].Next = players[(i+1)%numPlayers]
		players[i].Previous = players[(i+numPlayers-1)%numPlayers]
	}
}

func (p *Player) rollDice() {
	for i := range p.Dice {
		p.Dice[i] = roll()
	}
}

func (p *Player) evaluateDice() {
	newDice := make([]int, 0)
	for _, dice := range p.Dice {
		switch dice {
		case 6:
			p.Score++
			// newDice := make([]int, 0, len(p.Dice))
			// for _, v := range p.Dice {
			// 	if v != dice {
			// 		p.Dice = append(p.Dice, v)
			// 	}
			// }
			p.Dice = newDice
		case 1:
			p.getNextPlayer().giveDice()
		default:
			newDice = append(newDice, dice)
		}
	}
	p.Dice = newDice
	p.HasDice = len(p.Dice) > 0
}

func (p *Player) giveDice() {
	p.Dice = append(p.Dice, 1)
}

func (p *Player) getNextPlayer() *Player {
	return p.Next
}

func (p *Player) getName() string {
	return fmt.Sprintf("Pemain #%d", p.ID)
}

func getWinner(players []*Player) *Player {
	winner := players[0]
	for _, player := range players[1:] {
		if player.Score > winner.Score {
			winner = player
		}
	}
	return winner
}

func roll() int {
	return rand.Intn(6) + 1
}

func init() {
	rand.Seed(time.Now().UnixNano())
}
