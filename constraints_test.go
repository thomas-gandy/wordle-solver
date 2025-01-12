package main

import (
	"testing"
)

func TestAtIndexConstraint(t *testing.T) {
	t.Run("can pass when on applicable index with specified letter", func(t *testing.T) {
		c := atIndexConstraint{letter: 'a', atIndex: 2}

		couldPass := c.couldPass([]byte("wra"), 2)

		if !couldPass {
			t.Errorf("constraint should have been able to pass, but it didn't")
		}
	})

	t.Run("cannot pass when on applicable index with different letter", func(t *testing.T) {
		c := atIndexConstraint{letter: 'a', atIndex: 2}

		couldPass := c.couldPass([]byte("wry"), 2)

		if couldPass {
			t.Errorf("constraint should not have been able to pass, but it did")
		}
	})

	t.Run("can pass when not on applicable index", func(t *testing.T) {
		c := atIndexConstraint{letter: 'a', atIndex: 2}

		couldPass := c.couldPass([]byte("wr"), 1)

		if !couldPass {
			t.Errorf("constraint should have been able to pass, but it didn't")
		}
	})
}

func TestNotAtIndexConstraint(t *testing.T) {
	t.Run("cannot pass when on applicable index with specified letter", func(t *testing.T) {
		c := notAtIndexConstraint{letter: 'a', notAtIndex: 2}

		couldPass := c.couldPass([]byte("wra"), 2)

		if couldPass {
			t.Errorf("constraint should not have been able to pass, but it did")
		}
	})

	t.Run("can pass when on applicable index with different letter", func(t *testing.T) {
		c := notAtIndexConstraint{letter: 'a', notAtIndex: 2}

		couldPass := c.couldPass([]byte("wry"), 2)

		if !couldPass {
			t.Errorf("constraint should have been able to pass, but it didn't")
		}
	})

	t.Run("can pass when not on applicable index", func(t *testing.T) {
		c := notAtIndexConstraint{letter: 'a', notAtIndex: 2}

		couldPass := c.couldPass([]byte("wr"), 1)

		if !couldPass {
			t.Errorf("constraint should have been able to pass, but it didn't")
		}
	})
}

func TestNotExistConstraint(t *testing.T) {
	t.Run("can pass with different letter", func(t *testing.T) {
		c := notExistConstraint{letterThatShouldNotExist: 'a'}

		couldPass := c.couldPass([]byte("e"), 0)

		if !couldPass {
			t.Error("expected to pass as different letter specified, but didn't")
		}
	})

	t.Run("cannot pass same different letter", func(t *testing.T) {
		c := notExistConstraint{letterThatShouldNotExist: 'a'}

		couldPass := c.couldPass([]byte("a"), 0)

		if couldPass {
			t.Error("expected to pass as same letter specified, but did")
		}
	})
}

func TestMinimumFrequencyConstraint(t *testing.T) {
	t.Run("passes before applicable word length", func(t *testing.T) {
		word := "b"
		c := minimumFrequencyConstraint{
			letter:                   'a',
			minFrequency:             1,
			wordLengthToApplyAt:      len(word) + 1,
			currentLetterOccurrences: 0,
			occurrenceStates:         nil,
		}

		var couldPass bool
		for i := range word {
			couldPass = c.couldPass([]byte(word[:i+1]), len([]byte(word[:i+1]))-1)
		}

		if !couldPass {
			t.Error("expected to pass as before applicable word length, but did not")
		}
	})

	t.Run("passes at applicable word length when minimum met", func(t *testing.T) {
		word := "badaasad"
		c := minimumFrequencyConstraint{
			letter:                   'a',
			minFrequency:             2,
			wordLengthToApplyAt:      len(word),
			currentLetterOccurrences: 0,
			occurrenceStates:         nil,
		}

		var couldPass bool
		for i := range word {
			couldPass = c.couldPass([]byte(word[:i+1]), len([]byte(word[:i+1]))-1)
		}

		if !couldPass {
			t.Error("expected to pass as at applicable word length and min, but did not")
		}
	})

	t.Run("fails at applicable word length when minimum not met", func(t *testing.T) {
		word := "bb"
		c := minimumFrequencyConstraint{
			letter:                   'a',
			minFrequency:             1,
			wordLengthToApplyAt:      len(word),
			currentLetterOccurrences: 0,
			occurrenceStates:         nil,
		}

		var couldPass bool
		for i := range word {
			couldPass = c.couldPass([]byte(word[:i+1]), len([]byte(word[:i+1]))-1)
		}

		if couldPass {
			t.Error("expected to fail as at applicable word length without min, but did not")
		}
	})

	t.Run("state unaffected by sibling branch", func(t *testing.T) {
		c := minimumFrequencyConstraint{
			letter:              'a',
			minFrequency:        1,
			wordLengthToApplyAt: 4,
		}
		c.couldPass([]byte("b"), 0)
		c.saveState()
		c.couldPass([]byte("ba"), 1)
		c.restoreState()
		c.saveState()
		c.couldPass([]byte("bd"), 1)
		c.saveState()
		c.couldPass([]byte("bdi"), 2)
		c.saveState()

		couldPass := c.couldPass([]byte("bdie"), 3)

		if couldPass {
			t.Errorf("state incorrectly stored as %d occurrences instead of 0", c.currentLetterOccurrences)
		}
	})

	t.Run("at least one 'i' does not produce 'would'", func(t *testing.T) {
		c := minimumFrequencyConstraint{
			letter:              'i',
			minFrequency:        1,
			wordLengthToApplyAt: 5,
		}
		c.couldPass([]byte("w"), 0)
		c.couldPass([]byte("wo"), 1)
		c.couldPass([]byte("wou"), 2)
		c.couldPass([]byte("woul"), 3)

		couldPass := c.couldPass([]byte("would"), 4)

		if couldPass {
			t.Error("incorrectly passed word when it doesn't contain the letter 'i'")
		}
	})
}

func TestPreciseFrequencyConstraint(t *testing.T) {
	t.Run("passes before applicable word length", func(t *testing.T) {
		word := "b"
		c := preciseFrequencyConstraint{
			letter:                   'a',
			preciseFrequency:         1,
			wordLengthToApplyAt:      len(word) + 1,
			currentLetterOccurrences: 0,
			occurrenceStates:         nil,
		}

		var couldPass bool
		for i := range word {
			couldPass = c.couldPass([]byte(word[:i+1]), len([]byte(word[:i+1]))-1)
		}

		if !couldPass {
			t.Error("expected to pass as before applicable word length, but did not")
		}
	})

	t.Run("passes at applicable word length when precise met", func(t *testing.T) {
		word := "ba"
		c := preciseFrequencyConstraint{
			letter:                   'a',
			preciseFrequency:         1,
			wordLengthToApplyAt:      len(word),
			currentLetterOccurrences: 0,
			occurrenceStates:         nil,
		}

		var couldPass bool
		for i := range word {
			couldPass = c.couldPass([]byte(word[:i+1]), len([]byte(word[:i+1]))-1)
		}

		if !couldPass {
			t.Error("expected to pass as at applicable word length and min, but did not")
		}
	})

	t.Run("fails at applicable word length when precise not met", func(t *testing.T) {
		word := "bada"
		c := preciseFrequencyConstraint{
			letter:                   'a',
			preciseFrequency:         1,
			wordLengthToApplyAt:      len(word),
			currentLetterOccurrences: 0,
			occurrenceStates:         nil,
		}

		var couldPass bool
		for i := range word {
			couldPass = c.couldPass([]byte(word[:i+1]), len([]byte(word[:i+1]))-1)
		}

		if couldPass {
			t.Error("expected to fail at applicable word length, but did not")
		}
	})

	t.Run("state unaffected by sibling branch", func(t *testing.T) {
		c := preciseFrequencyConstraint{
			letter:              'a',
			preciseFrequency:    2,
			wordLengthToApplyAt: 4,
		}
		c.couldPass([]byte("b"), 0)
		c.saveState()
		c.couldPass([]byte("ba"), 1)
		c.restoreState()
		c.saveState()
		c.couldPass([]byte("bd"), 1)
		c.saveState()
		c.couldPass([]byte("bda"), 2)
		c.saveState()

		couldPass := c.couldPass([]byte("bdaa"), 3)

		if !couldPass {
			t.Errorf("state incorrectly stored as %d occurrences instead of 2", c.currentLetterOccurrences)
		}
	})
}
