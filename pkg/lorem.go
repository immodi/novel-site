// pkg/lorem.go
package pkg

import (
	"math/rand"
	"strings"
	"time"
)

var loremWords = []string{
	"lorem", "ipsum", "dolor", "sit", "amet", "consectetur",
	"adipiscing", "elit", "sed", "do", "eiusmod", "tempor",
	"incididunt", "ut", "labore", "et", "dolore", "magna", "aliqua",
}

func init() {
	rand.New(rand.NewSource((time.Now().UnixNano())))
}

func randomWord() string {
	return loremWords[rand.Intn(len(loremWords))]
}

// LoremSentence generates a sentence with N–M words.
func LoremSentence(minWords, maxWords int) string {
	n := rand.Intn(maxWords-minWords+1) + minWords
	words := make([]string, n)
	for i := range words {
		words[i] = randomWord()
	}
	s := strings.Join(words, " ")
	return strings.ToUpper(s[:1]) + s[1:] + "."
}

// LoremParagraph generates a paragraph with N–M sentences.
func LoremParagraph(minSent, maxSent int) string {
	n := rand.Intn(maxSent-minSent+1) + minSent
	sentences := make([]string, n)
	for i := range sentences {
		sentences[i] = LoremSentence(4, 12)
	}
	return strings.Join(sentences, " ")
}

// LoremText generates multiple paragraphs separated by newlines.
func LoremText(paragraphs int) string {
	paras := make([]string, paragraphs)
	for i := range paras {
		paras[i] = LoremParagraph(3, 6)
	}
	return strings.Join(paras, "\n\n")
}
