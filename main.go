package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"slices"
	"strconv"
	"strings"
	"time"
)

func getLineFromUser(inputMessage string, scanner *bufio.Scanner) []byte {
	fmt.Printf("%s: ", inputMessage)
	scanner.Scan()

	return []byte(strings.ToLower(scanner.Text()))
}

func getNumberFromUser(inputMessage string, scanner *bufio.Scanner) int {
	fmt.Printf("%s: ", inputMessage)
	scanner.Scan()

	num, err := strconv.ParseInt(scanner.Text(), 10, 32)
	if err != nil {
		panic("could not parse input as a number")
	}

	return int(num)
}

func getNextWords(wordTree prefixTree, constraints []constraint, wordFrequencies map[string]int) []string {
	startTime := time.Now()
	possibleWords := wordTree.findPossibleWords(constraints)
	endTime := time.Now()
	elapsedTime := endTime.Sub(startTime)

	slices.SortFunc(possibleWords, func(a, b string) int {
		aFreq, exists := wordFrequencies[a]
		if !exists {
			aFreq = -1
		}

		bFreq, exists := wordFrequencies[b]
		if !exists {
			bFreq = -1
		}

		return bFreq - aFreq
	})

	fmt.Printf("found a total of %d words in %dms\n", len(possibleWords), elapsedTime.Milliseconds())
	return possibleWords
}

func createWordFrequenciesFromCsv(path string) map[string]int {
	file, err := os.Open(path)
	if err != nil {
		panic("failed to open word frequencies file")
	}
	startTime := time.Now()

	wordFrequencies := make(map[string]int, 1<<19)
	scanner := bufio.NewScanner(file)
	scanner.Scan()
	for scanner.Scan() {
		line := scanner.Text()
		elements := strings.Split(line, ",")
		word := elements[0]
		frequency, err := strconv.ParseInt(elements[1], 10, 64)
		if err != nil {
			panic("failed to parse word frequency")
		}

		wordFrequencies[word] = int(min(frequency, math.MaxInt-1))
	}

	endTime := time.Now()
	elapsedTime := endTime.Sub(startTime)
	fmt.Printf("initialised word frequency map (%dms)\n", elapsedTime.Milliseconds())

	return wordFrequencies
}

func printConstraints(constraints []constraint) {
	for _, c := range constraints {
		fmt.Println(c)
	}
}

func updateSession(session *sessionInfo, word, colors []byte) {
	if len(word) != len(colors) {
		return
	}

	constraints := make([]constraint, 0, 1<<4)
	constraints = append(constraints, &lengthConstraint{maxWordLength: len(word)})

	possiblyInvalidLetters := make(map[byte]struct{}, 26)
	letterFrequencyMap := make(map[byte]int, 26)
	sessionComplete := true

	for i, letter := range word {
		switch colors[i] {
		case 'g':
			constraints = append(constraints, &atIndexConstraint{letter: letter, atIndex: i})
			letterFrequencyMap[letter]++
		case 'y':
			constraints = append(constraints, &notAtIndexConstraint{letter: letter, notAtIndex: i})
			letterFrequencyMap[letter]++
			sessionComplete = false
		case 'b':
			possiblyInvalidLetters[letter] = struct{}{}
			sessionComplete = false
		}
	}
	session.completed = sessionComplete

	for letter, freq := range letterFrequencyMap {
		_, preciseCountKnown := possiblyInvalidLetters[letter]
		delete(possiblyInvalidLetters, letter)

		if preciseCountKnown {
			session.preciseLetters[letter] = freq
		} else {
			constraints = append(constraints, &minimumFrequencyConstraint{
				letter:              letter,
				minFrequency:        freq,
				wordLengthToApplyAt: len(word),
			})
		}
	}

	for letter := range possiblyInvalidLetters {
		session.invalidLetters[letter] = struct{}{}
	}

	for letter := range session.invalidLetters {
		constraints = append(constraints, &notExistConstraint{letterThatShouldNotExist: letter})
	}

	for letter, freq := range session.preciseLetters {
		constraints = append(constraints, &preciseFrequencyConstraint{
			letter:              letter,
			preciseFrequency:    freq,
			wordLengthToApplyAt: len(word),
		})
	}

	session.constraints = constraints
}

type sessionInfo struct {
	completed      bool
	constraints    []constraint
	invalidLetters map[byte]struct{}
	preciseLetters map[byte]int
}

func main() {
	wordFrequencies := createWordFrequenciesFromCsv("./word-frequencies.csv")
	wordTree := createPrefixTreeFromWordFile("./words.txt")

	inputScanner := bufio.NewScanner(os.Stdin)
	session := sessionInfo{
		constraints:    []constraint{},
		invalidLetters: make(map[byte]struct{}, 26),
		preciseLetters: make(map[byte]int, 26),
	}
	suggestedWord := getLineFromUser("enter your first word", inputScanner)

	for !session.completed {
		fmt.Printf("enter '%s'\n", suggestedWord)
		colors := getLineFromUser("color of each letter (G)reen, (Y)ellow, (B)lack", inputScanner)

		updateSession(&session, suggestedWord, colors)
		printConstraints(session.constraints)
		suggestedWords := getNextWords(wordTree, session.constraints, wordFrequencies)

		numToShow := min(5, len(suggestedWords))
		for i, possibleWord := range suggestedWords[:numToShow] {
			fmt.Printf("(%d) %s\n", i+1, possibleWord)
		}
		indexChosen := getNumberFromUser("index of word chosen", inputScanner)
		suggestedWord = []byte(suggestedWords[indexChosen-1])
	}
}
