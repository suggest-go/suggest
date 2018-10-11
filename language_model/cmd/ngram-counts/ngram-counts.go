package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"github.com/alldroll/suggest/alphabet"
	lm "github.com/alldroll/suggest/language_model"
	"io"
	"log"
	"os"
	"strings"
)

type countTable = map[uint32]*count

type count struct {
	higher map[uint32]*count
	count  uint32
}

type indexer struct {
	table  map[lm.Word]uint32
	holder []lm.Word
}

func (i *indexer) getOrCreate(ngram lm.Word) uint32 {
	index, ok := i.table[ngram]

	if !ok {
		index = uint32(len(i.holder))
		i.table[ngram] = index
		i.holder = append(i.holder, ngram)
	}

	return index
}

func (i *indexer) find(index uint32) (lm.Word, error) {
	if uint32(len(i.holder)) <= index {
		return "", errors.New("index is not exists")
	}

	return i.holder[int(index)], nil
}

const (
	startSymbol = "<S>"
	endSymbol   = "</S>"
)

var (
	sourcePath string
)

func init() {
	flag.StringVar(&sourcePath, "source", "", "source path")
	//flag.StringVar(&sourcePath, "output", "", "source path")
}

func out(writer io.Writer, indexer *indexer, table countTable, level int, grams lm.NGram) {
	for index, counts := range table {
		nGram, err := indexer.find(index)
		if err != nil {
			panic(err)
		}

		if level == 0 {
			fmt.Fprintf(writer, "%s\t%d\n", strings.Join(append(grams, nGram), " "), counts.count)
		} else if level > 0 && counts.higher != nil {
			out(writer, indexer, counts.higher, level-1, append(grams, nGram))
		}
	}
}

func main() {
	log.SetPrefix("indexer: ")
	log.SetFlags(0)

	flag.Parse()

	sourceFile, err := os.Open(sourcePath)
	if err != nil {
		log.Fatalf("could not open source file %s", err)
	}

	defer sourceFile.Close()

	retriever := lm.NewSentenceRetriver(
		lm.NewTokenizer(alphabet.NewEnglishAlphabet()),
		bufio.NewReader(sourceFile),
		alphabet.NewSimpleAlphabet([]rune{'.', '?', '!'}),
	)

	nGramOrder := uint8(3)
	table := countTable{}
	indexer := &indexer{
		table:  map[lm.Word]uint32{},
		holder: []lm.Word{},
	}

	generators := []lm.Generator{}

	for i := uint8(1); i <= nGramOrder; i++ {
		generators = append(
			generators,
			lm.NewGenerator(
				i,
				startSymbol,
				endSymbol,
			),
		)
	}

	for {
		sentence := retriever.Retrieve()
		if sentence == nil {
			break
		}

		if len(sentence) == 0 {
			continue
		}

		// TODO process it concurrent
		for i := 0; i < int(nGramOrder); i++ {
			generator := generators[i]
			nGramsSet := generator.Generate(sentence)

			for _, nGrams := range nGramsSet {
				t := table

				for j, nGram := range nGrams {
					index := indexer.getOrCreate(nGram)

					if t[index] == nil {
						t[index] = &count{}
					}

					if j == i {
						t[index].count += 1
					}

					if j < i {
						if t[index].higher == nil {
							t[index].higher = countTable{}
						}

						t = t[index].higher
					}
				}
			}
		}
	}

	grams := lm.NGram{}
	for i := 0; i < int(nGramOrder); i++ {
		f, err := os.OpenFile(fmt.Sprintf("%d-gm", i+1), os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			panic(err)
		}

		defer f.Close()

		buf := bufio.NewWriter(f)
		grams = grams[:0]
		out(buf, indexer, table, i, grams)
		buf.Flush()
	}
}
