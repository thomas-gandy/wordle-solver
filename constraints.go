package main

import "strconv"

type constraint interface {
	couldPass(word []byte, i int) bool
	signalsWordEnd() bool
}

type stateConstraint interface {
	constraint
	saveState()
	restoreState()
}

type atIndexConstraint struct {
	letter  uint8
	atIndex int
}

type notAtIndexConstraint struct {
	letter     uint8
	notAtIndex int
}

type notExistConstraint struct {
	letterThatShouldNotExist uint8
}

type minimumFrequencyConstraint struct {
	letter                   uint8
	minFrequency             int
	wordLengthToApplyAt      int
	currentLetterOccurrences int
	occurrenceStates         []int
}
type preciseFrequencyConstraint struct {
	letter                   uint8
	preciseFrequency         int
	wordLengthToApplyAt      int
	currentLetterOccurrences int
	occurrenceStates         []int
}

type lengthConstraint struct {
	maxWordLength int
}

func (c *atIndexConstraint) couldPass(word []byte, i int) bool {
	return c.atIndex != i || c.letter == word[i]
}

func (c *atIndexConstraint) signalsWordEnd() bool {
	return false
}

func (c *atIndexConstraint) String() string {
	return "'" + string(rune(c.letter)) + "' must be at index " + strconv.Itoa(c.atIndex)
}

func (c *notAtIndexConstraint) couldPass(word []byte, i int) bool {
	return c.notAtIndex != i || c.letter != word[i]
}

func (c *notAtIndexConstraint) signalsWordEnd() bool {
	return false
}

func (c *notAtIndexConstraint) String() string {
	return "'" + string(rune(c.letter)) + "' must not be at index " + strconv.Itoa(c.notAtIndex)
}

func (c *notExistConstraint) couldPass(word []byte, i int) bool {
	return c.letterThatShouldNotExist != word[i]
}

func (c *notExistConstraint) signalsWordEnd() bool {
	return false
}

func (c *notExistConstraint) String() string {
	return "'" + string(rune(c.letterThatShouldNotExist)) + "' must not exist"
}

func (c *minimumFrequencyConstraint) couldPass(word []byte, i int) bool {
	if c.letter == word[i] {
		c.currentLetterOccurrences++
	}

	if c.wordLengthToApplyAt == i+1 {
		return c.currentLetterOccurrences >= c.minFrequency
	}

	return true
}

func (c *minimumFrequencyConstraint) signalsWordEnd() bool {
	return false
}

func (c *minimumFrequencyConstraint) saveState() {
	c.occurrenceStates = append(c.occurrenceStates, c.currentLetterOccurrences)
}

func (c *minimumFrequencyConstraint) restoreState() {
	c.currentLetterOccurrences = c.occurrenceStates[len(c.occurrenceStates)-1]
	c.occurrenceStates = c.occurrenceStates[:len(c.occurrenceStates)-1]
}

func (c *minimumFrequencyConstraint) String() string {
	l := string(rune(c.letter))
	o := strconv.Itoa(c.minFrequency)
	return "'" + l + "' must occur at least " + o + " time(s)"
}

func (c *preciseFrequencyConstraint) couldPass(word []byte, i int) bool {
	if c.letter == word[i] {
		c.currentLetterOccurrences++
	}

	if c.wordLengthToApplyAt == i+1 {
		return c.currentLetterOccurrences == c.preciseFrequency
	}

	return true
}

func (c *preciseFrequencyConstraint) signalsWordEnd() bool {
	return false
}

func (c *preciseFrequencyConstraint) saveState() {
	c.occurrenceStates = append(c.occurrenceStates, c.currentLetterOccurrences)
}

func (c *preciseFrequencyConstraint) restoreState() {
	c.currentLetterOccurrences = c.occurrenceStates[len(c.occurrenceStates)-1]
	c.occurrenceStates = c.occurrenceStates[:len(c.occurrenceStates)-1]
}

func (c *preciseFrequencyConstraint) String() string {
	l := string(rune(c.letter))
	o := strconv.Itoa(c.preciseFrequency)
	return "'" + l + "' must occur exactly " + o + " time(s)"
}

func (c *lengthConstraint) couldPass(_ []byte, i int) bool {
	return i+1 < c.maxWordLength
}

func (c *lengthConstraint) signalsWordEnd() bool {
	return true
}

func (c *lengthConstraint) String() string {
	return "the word must have a length of " + strconv.Itoa(c.maxWordLength)
}
