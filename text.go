package markov

import (
	"log"
	"math"
	"regexp"
	"strings"
)

// A Text is a fucking dick.
type Text struct {
	Chain
	sentences    [][]string
	rejoinedText string
}

type Corpus [][]string

// NewText creates a new text.
func NewText(text string) Text {
	s := splitSentences(splitText(text))
	t := joinText(s)
	c := NewChain(s)

	return Text{Chain: c, rejoinedText: t, sentences: s}
}

// MakeSentence makes a sentence.
func (t *Text) MakeSentence() string {
	for i := 0; i < 10; i++ {
		s := t.Chain.Walk()
		if t.testOutputSentence(s) {
			return strings.Join(s, " ")
		}
	}

	return ""
}

// joinText joins sentences into a reconstruction of the input text.
func joinText(sentences [][]string) string {
	t := make([]string, len(sentences))

	for i, words := range sentences {
		t[i] = strings.Join(words, " ")
	}

	return strings.Join(t, " ")
}

// splitText splits the input text into sentences.
func splitText(text string) []string {
	r, err := regexp.Compile(`([\w\.'’&\]\)]+[\.\?!])([‘’“”'\"\)\]]*)(\s+)[^a-z\-–—]`)
	if err != nil {
		log.Fatal(err)
	}

	ms := r.FindAllStringSubmatchIndex(text, -1)
	sens := make([]string, len(ms)+1)

	for i, m := range ms {
		if i == 0 {
			// The first sentence begins at the beginning, of course.
			sens[i] = text[0:m[3]]
		} else {
			// Sentences except the first begin after the end of the previous.
			sens[i] = text[ms[i-1][7]:m[3]]
		}

		if i+1 == len(ms) {
			// The last sentence ends at the end, of course.
			sens[i+1] = text[m[7]:len(text)]
		}
	}

	return sens
}

// splitSentences splits sentences into lists of words.
func splitSentences(sentences []string) [][]string {
	r, err := regexp.Compile("\\s+")
	if err != nil {
		log.Fatal(err)
	}

	runs := make([][]string, 0, len(sentences))

	for _, s := range sentences {
		if testInputSentence(s) {
			runs = append(runs, r.Split(s, -1))
		}
	}

	return runs
}

// testInputSentence is a sentence filter that rejects sentences which
// “contain the type of punctuation that would look strange on its own.”
func testInputSentence(sentence string) bool {
	if len(strings.TrimSpace(sentence)) == 0 {
		return false
	}

	r, err := regexp.Compile(`(^')|('$)|\s'|'\s|[\"(\(\)\[\])]`)
	if err != nil {
		log.Fatal(err)
	}

	if r.FindString(sentence) != "" {
		return false
	}

	return true
}

// testOutputSentence is a sentence filter that rejects sentences which
// “too closely match the original text, namely those that contain any
// identical sequence of words of [an arbitrary maths-based] length”.
func (t *Text) testOutputSentence(sentence []string) bool {
	// I have no idea how this formula works.
	l := float64(len(sentence))
	m := math.Min(15, math.Floor(0.7*l))
	c := int(math.Max(1, l-m))

	for i := 0; i < c; i++ {
		// Split the sentence into `c` “windows” of length `m`,
		// (regardless of how we derived `c` and `m` in the first place).
		g := sentence[i : i+int(m)+1]

		// Reject the sentence if this part is present in the text.
		s := strings.Join(g, " ")
		if strings.Contains(t.rejoinedText, s) {
			return false
		}
	}

	return true
}
