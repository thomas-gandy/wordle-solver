package main

import (
	"fmt"
	"testing"
)

func TestPrefixTreeCreate(t *testing.T) {

}

func TestPrefixTreeAddWord(t *testing.T) {

}

func TestPrefixTreeFindWord(t *testing.T) {

}

func TestSomething(t *testing.T) {
	wordFrequencies := createWordFrequenciesFromCsv("./word-frequencies.csv")
	wordTree := createPrefixTreeFromWordFile("./words.txt")

	session := sessionInfo{constraints: []constraint{}, invalidLetters: make(map[byte]struct{}, 26)}

	// 1
	suggestedWord := []byte("arise")
	fmt.Printf("enter '%s'\n", suggestedWord)
	colors := []byte("bbybb")

	updateSession(&session, suggestedWord, colors)
	printConstraints(session.constraints)
	suggestedWord = getNextWords(wordTree, session.constraints, wordFrequencies)

	// 2
	fmt.Printf("enter '%s'\n", suggestedWord)
	colors = []byte("bgybb")

	updateSession(&session, suggestedWord, colors)
	printConstraints(session.constraints)
	suggestedWord = getNextWords(wordTree, session.constraints, wordFrequencies)
}
