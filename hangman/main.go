package main

import (
	"errors"
	"fmt"
	"math/rand"
	"os"
	"time"
	"unicode"
)

var dictionary = []string{
	"Zombie",
	"Gopher",
	"Apple",
	"Programming",
	"Linux is Superior",
}

func main() {
	rand.Seed(time.Now().UnixNano())
	targetWord := getRandomWord()
	guessedLetters := initGuessedWords(targetWord)
	wrongGuesses := 0

	fmt.Println("Welcome to the Hangman!")
	fmt.Println("Please make a guess")

	for !checkIfAllFound(targetWord, guessedLetters) {
		printGameState(targetWord, guessedLetters)
		letter, err := getUserInput()

		if err != nil {
			fmt.Println(err)
			continue
		}

		if isPredictionCorrect(letter, guessedLetters) {
			guessedLetters[letter] = true
			fmt.Println("Well done! Keep going..")
		} else {
			hangman := readHangmanFile(wrongGuesses)
			wrongGuesses++
			fmt.Println(hangman)

			if wrongGuesses == 10 {
				fmt.Println("Whoops. You lost!")
				return
			}
		}
	}
	fmt.Println("\nCongrats!! You won!")
}

func getRandomWord() string {
	targetWord := dictionary[rand.Intn(len(dictionary))]
	return targetWord
}

func initGuessedWords(targetWord string) map[rune]bool {
	guessedLetters := map[rune]bool{}
	for _, char := range targetWord {
		guessedLetters[unicode.ToLower(char)] = false
	}
	return guessedLetters
}

func getUserInput() (rune, error) {
	var letter string
	fmt.Print("\n> ")
	fmt.Scan(&letter)

	if len(letter) > 1 {
		return -1, errors.New("Invalid value")
	}

	return rune(letter[0]), nil
}

func isPredictionCorrect(letter rune, guessedLetters map[rune]bool) bool {
	if _, ok := guessedLetters[letter]; ok {
		return true
	}
	return false
}

func checkIfAllFound(targetWord string, guessedLetters map[rune]bool) bool {
	for _, char := range targetWord {
		if !guessedLetters[unicode.ToLower(char)] {
			return false
		}
	}
	return true
}

func printGameState(targetWord string, guessedLetters map[rune]bool) {
	fmt.Println()
	for _, char := range targetWord {
		if char == ' ' {
			fmt.Print(" ")
		} else if guessedLetters[unicode.ToLower(char)] {
			fmt.Printf("%c", char)
		} else {
			fmt.Print("_")
		}
		fmt.Print(" ")
	}
	fmt.Println()
}

func readHangmanFile(state int) string {
	dat, err := os.ReadFile(fmt.Sprintf("./states/hangman%d", state))

	if err != nil {
		fmt.Println("\nAN ERROR OCCURED")
		fmt.Println("EXITING")
		panic(err)
	}

	return string(dat)
}
