package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
)

type node struct {
	value            uint8
	isWordTerminator bool
	children         map[uint8]*node
}

type prefixTree struct {
	root              *node
	longestWordLength int
}

func (tree *prefixTree) add(word string) {
	add(tree.root, word, len(word), 0)
}

func add(n *node, word string, length, i int) {
	if i == length {
		n.isWordTerminator = true
		return
	}

	currentChar := word[i]
	childNode, exists := n.children[currentChar]

	if exists {
		add(childNode, word, length, i+1)
	} else {
		newNode := &node{value: currentChar, children: make(map[uint8]*node)}
		n.children[currentChar] = newNode
		add(newNode, word, length, i+1)
	}
}

func (tree *prefixTree) findPossibleWords(constraints []constraint) []string {
	result := make([]string, 0, 512)
	currentWord := make([]byte, tree.longestWordLength)
	stateConstraints := make([]stateConstraint, 0, len(constraints))

	for _, c := range constraints {
		value, isStateConstraint := c.(stateConstraint)
		if isStateConstraint {
			stateConstraints = append(stateConstraints, value)
		}
	}

	for _, child := range tree.root.children {
		for _, constraintWithState := range stateConstraints {
			constraintWithState.saveState()
		}
		result = append(result, findPossibleWords(child, 0, constraints, stateConstraints, currentWord)...)
		for _, constraintWithState := range stateConstraints {
			constraintWithState.restoreState()
		}
	}

	return result
}

func findPossibleWords(n *node, i int, allConstraints []constraint, states []stateConstraint, word []byte) []string {
	word[i] = n.value
	constraintCompletesWord := false
	for _, currentConstraint := range allConstraints {
		if !currentConstraint.couldPass(word, i) {
			if currentConstraint.signalsWordEnd() && n.isWordTerminator {
				constraintCompletesWord = true
			} else {
				return []string{}
			}
		}
	}
	if constraintCompletesWord {
		return []string{string(word[:i+1])}
	}

	results := make([]string, 0, 512)
	for _, childNode := range n.children {
		for _, constraintWithState := range states {
			constraintWithState.saveState()
		}
		results = append(results, findPossibleWords(childNode, i+1, allConstraints, states, word)...)
		for _, constraintWithState := range states {
			constraintWithState.restoreState()
		}
	}

	return results
}

func createPrefixTreeFromWordFile(path string) prefixTree {
	file, err := os.Open(path)
	if err != nil {
		panic("failed to open words file")
	}
	startTime := time.Now()

	wordTree := prefixTree{root: &node{children: make(map[uint8]*node)}}
	scanner := bufio.NewScanner(file)
	longestWordLength := 0

	for scanner.Scan() {
		word := strings.ToLower(scanner.Text())
		if len(word) > longestWordLength {
			longestWordLength = len(word)
		}
		wordTree.add(word)
	}
	wordTree.longestWordLength = longestWordLength

	endTime := time.Now()
	elapsedTime := endTime.Sub(startTime)
	fmt.Printf("generated word prefix tree (%dms)\n", elapsedTime.Milliseconds())

	return wordTree
}
