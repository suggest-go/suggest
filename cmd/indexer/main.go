package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"flag"
	"github.com/alldroll/cdb"
	"github.com/alldroll/suggest"
	"github.com/alldroll/suggest/compression"
	"github.com/alldroll/suggest/dictionary"
	"github.com/alldroll/suggest/index"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"syscall"
	"time"
)

var (
	configPath string
	dict       string
	pidPath    string
)

const tmpSuffix = ".tmp"

func init() {
	flag.StringVar(&configPath, "config", "config.json", "config path")
	flag.StringVar(&dict, "dict", "", "reindex certain dict")
	flag.StringVar(&pidPath, "pid", "", "pid path")
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

	reindexed := false
	totalStart := time.Now()

	for _, config := range configs {
		if dict != "" && dict != config.Name {
			continue
		}

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
		reindexed = true
	}

	if !reindexed {
		log.Printf("There were not any reindex job")
		os.Exit(0)
	}

	log.Printf("Total time spent %s", time.Since(totalStart).String())

	tryToSendReindexSignal()
}

//
func buildDictionary(sourcePath string, config suggest.IndexDescription) dictionary.Dictionary {
	sourceFile, err := os.OpenFile(sourcePath, os.O_RDONLY, 0)
	if err != nil {
		log.Fatalf("cannot open source file %s", err)
	}

	destinationFile, err := os.OpenFile(
		config.GetDictionaryFile()+tmpSuffix,
		os.O_CREATE|os.O_RDWR|os.O_TRUNC,
		0644,
	)
	if err != nil {
		log.Fatal("Fail to create dictionary file %s", err)
	}

	cdbHandle := cdb.New()
	cdbWriter, err := cdbHandle.GetWriter(destinationFile)
	if err != nil {
		log.Fatal("Fail to create cdb writer %s", err)
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
			log.Fatalf("Fail to put record to cdb %s", err)
		}

		docId++
	}

	log.Printf("Number of string %d", docId)

	err = cdbWriter.Close()
	if err != nil {
		log.Fatalf("Fail to save cdb dictionary %s", err)
	}

	if err := os.Rename(config.GetDictionaryFile()+tmpSuffix, config.GetDictionaryFile()); err != nil {
		log.Fatalf("Error to rename file %s", err)
	}

	return dictionary.NewCDBDictionary(destinationFile)
}

//
func storeIndex(indices index.Indices, config suggest.IndexDescription) {
	writer := index.NewOnDiscIndicesWriter(
		compression.VBEncoder(),
		config.GetHeaderFile()+tmpSuffix,
		config.GetDocumentListFile()+tmpSuffix,
	)

	if err := writer.Save(indices); err != nil {
		log.Fatalf("Fail to save index %s", err)
	}

	if err := os.Rename(config.GetHeaderFile()+tmpSuffix, config.GetHeaderFile()); err != nil {
		log.Fatalf("Error to rename file %s", err)
	}

	if err := os.Rename(config.GetDocumentListFile()+tmpSuffix, config.GetDocumentListFile()); err != nil {
		log.Fatalf("Error to rename file %s", err)
	}
}

//
func tryToSendReindexSignal() {
	if pidPath == "" {
		return
	}

	d, err := ioutil.ReadFile(pidPath)
	if err != nil {
		log.Fatalf("error parsing pid from %s: %s", pidPath, err)
	}

	pid, err := strconv.Atoi(string(bytes.TrimSpace(d)))
	if err != nil {
		log.Fatalf("error parsing pid from %s: %s", pidPath, err)
	}

	if err := syscall.Kill(pid, syscall.SIGHUP); err != nil {
		log.Printf("Fail to send reindex signal to %d, %s", pid, err)
	}
}
