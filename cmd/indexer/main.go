package main

import (
	"bufio"
	"encoding/binary"
	"flag"
	"github.com/alldroll/cdb"
	"github.com/alldroll/suggest"
	"github.com/alldroll/suggest/compression"
	"github.com/alldroll/suggest/dictionary"
	"github.com/alldroll/suggest/index"
	"log"
	"os"
	"time"
)

var (
	configPath string
)

func init() {
	flag.StringVar(&configPath, "config", "config.json", "config path")
}

//
func buildDictionary(sourcePath string, config suggest.IndexDescription) dictionary.Dictionary {
	sourceFile, err := os.OpenFile(sourcePath, os.O_RDONLY, 0)
	if err != nil {
		log.Fatalf("cannot open source file %s", err)
	}

	destinationFile, err := os.OpenFile(
		config.GetDictionaryFile(),
		os.O_CREATE|os.O_RDWR|os.O_TRUNC,
		0644,
	)
	if err != nil {
		log.Fatal(err)
	}

	cdbHandle := cdb.New()
	cdbWriter, err := cdbHandle.GetWriter(destinationFile)
	if err != nil {
		log.Fatal(err)
	}

	var (
		docId   uint32 = 0
		key            = make([]byte, 4)
		scanner        = bufio.NewScanner(sourceFile)
	)

	for scanner.Scan() {
		binary.LittleEndian.PutUint32(key, docId)
		err = cdbWriter.Put(key, scanner.Bytes())
		if err != nil {
			log.Fatal(err)
		}

		docId++
	}

	log.Printf("Number of string %d", docId)

	err = cdbWriter.Close()
	if err != nil {
		log.Fatal(err)
	}

	return dictionary.NewCDBDictionary(destinationFile)
}

//
func storeIndex(indices index.Indices, config suggest.IndexDescription) {
	writer := index.NewOnDiscIndicesWriter(compression.VBEncoder(), config.GetHeaderFile(), config.GetDocumentListFile())
	err := writer.Save(indices)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	log.SetPrefix("indexer: ")
	log.SetFlags(0)

	flag.Parse()

	f, err := os.Open(configPath)
	if err != nil {
		log.Fatalf("could not open config file %s", err)
	}

	defer f.Close()
	configs, err := suggest.ReadConfigs(f)
	if err != nil {
		log.Fatalf("invalid config file format %s", err)
	}

	totalStart := time.Now()

	for _, config := range configs {
		log.Printf("Start process '%s' config", config.Name)

		alphabet := config.CreateAlphabet()
		cleaner := index.NewCleaner(alphabet.Chars(), config.Pad, config.Wrap)
		generator := index.NewGenerator(config.NGramSize, alphabet)

		log.Printf("Building dictionary...")

		start := time.Now()
		dictionary := buildDictionary(config.SourcePath, config)
		log.Printf("Time spent %s", time.Since(start))

		// create index in memory
		indexer := index.NewIndexer(config.NGramSize, generator, cleaner)
		log.Printf("Creating index...")

		start = time.Now()
		indices := indexer.Index(dictionary)
		log.Printf("Time spent %s", time.Since(start))

		// store index on disc
		log.Printf("Storing index...")

		start = time.Now()
		storeIndex(indices, config)
		log.Printf("Time spent %s", time.Since(start))

		log.Printf("End process\n\n")
	}

	log.Printf("Total time spent %s", time.Since(totalStart).String())
}
