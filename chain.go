package markov

import (
	"math/rand"
	"time"
)

const begin = "___BEGIN__"
const end = "___END__"

// Chain is a Markov chain.
type Chain struct {
	corpus Corpus
	model  model
}

type state [2]string
type follows map[string]int

// model maps states to transition counts.
//
// Example corpus:
//   [["The", "quick", "fox."],
//    ["The", "lazy", "brown", "dog".]]
//
// Example (partial) model:
//   ["___BEGIN__", "___BEGIN__"]: { "The": 2 }
type model map[state]follows

// NewChain creates a new chain by building a model from the provided corpus.
func NewChain(corpus Corpus) Chain {
	return Chain{corpus, buildModel(corpus)}
}

// Walk returns a list representing a single run of the Markov model.
func (c *Chain) Walk() []string {
	sen := make([]string, 0)

	// The initial state is (___BEGIN__, ___BEGIN__).
	st := state{begin, begin}

	for {
		// Choose the next word.
		w := c.move(st)

		// If the next “word” is ___END__, return the sentence.
		if w == end {
			return sen
		}

		// Otheriwse, add the word to the sentence and update the state.
		sen = append(sen, w)
		st = state{st[1], w}
	}
}

func (c *Chain) move(s state) string {
	fs := c.model[s]
	tw := 0
	for _, w := range fs {
		tw += w
	}

	rand.Seed(time.Now().UnixNano())
	r := rand.Intn(tw)

	for c, w := range fs {
		r -= w
		if r <= 0 {
			return c
		}
	}

	panic("Something is fucked.")
}

// buildModel builds a representation of the Markov model.
func buildModel(corpus Corpus) model {
	m := make(model)

	for _, sen := range corpus {
		// Wrap the sentence with two ___BEGIN__ tokens and one ___END__ token.
		is := append(append([]string{begin, begin}, sen...), end)

		for i := range is {
			// Break if there's no more sentence left.
			if is[i+1] == end {
				break
			}

			// The current state is the current word and the next word.
			// If this state isn't represented in the model, add it.
			st := state{is[i], is[i+1]}
			if _, ok := m[st]; !ok {
				m[st] = make(follows)
			}

			// Update the transition map for the word that follows this state.
			f := is[i+2]
			if c, ok := m[st][f]; ok {
				m[st][f] = c + 1
			} else {
				m[st][f] = 1
			}
		}
	}

	return m
}
