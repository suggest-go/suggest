package cmd

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"syscall"
	"time"

	"github.com/spf13/cobra"

	"github.com/alldroll/suggest/pkg/compression"
	"github.com/alldroll/suggest/pkg/dictionary"
	"github.com/alldroll/suggest/pkg/index"
	"github.com/alldroll/suggest/pkg/suggest"
)

var (
	dict string
	host string
)

func init() {
	indexCmd.Flags().StringVarP(&dict, "dict", "d", "", "reindex certain dict")
	indexCmd.Flags().StringVarP(&host, "host", "", "", "host to send reindex request")

	rootCmd.AddCommand(indexCmd)
}

var indexCmd = &cobra.Command{
	Use:   "indexer -c [config file]",
	Short: "builds indexes and dictionaries",
	Long:  `builds indexes and dictionaries and send signal to reload suggest-service`,
	RunE: func(cmd *cobra.Command, args []string) error {
		log.SetPrefix("indexer: ")
		log.SetFlags(0)

		configs, err := readConfigs()

		if err != nil {
			return err
		}

		reindexed := false
		totalStart := time.Now()

		for _, config := range configs {
			if dict != "" && dict != config.Name {
				continue
			}

			if err := indexJob(config); err != nil {
				return err
			}

			reindexed = true
		}

		if !reindexed {
			log.Printf("There were not any reindex job")
			return nil
		}

		log.Printf("Total time spent %s", time.Since(totalStart).String())

		if pidPath != "" {
			if err := tryToSendReindexSignal(); err != nil {
				return err
			}
		}

		if host != "" {
			if err := tryToSendReindexRequest(); err != nil {
				return err
			}
		}

		return nil
	},
}

// readConfigs retrieves a list of index descriptions from the configPath
func readConfigs() ([]suggest.IndexDescription, error) {
	f, err := os.Open(configPath)

	if err != nil {
		return nil, fmt.Errorf("Could not open config file %s", err)
	}

	defer f.Close()
	configs, err := suggest.ReadConfigs(f)

	if err != nil {
		return nil, fmt.Errorf("Invalid config file format %s", err)
	}

	return configs, nil
}

// indexJob performs building a dictionary, a search index for the given index description
func indexJob(config suggest.IndexDescription) error {
	log.Printf("Start process '%s' config", config.Name)

	// create a cdb dictionary
	log.Printf("Building a dictionary...")
	start := time.Now()

	dictReader, err := newDictionaryReader(config)

	if err != nil {
		return err
	}

	dict, err := dictionary.BuildCDBDictionary(dictReader, config.GetDictionaryFile())

	if err != nil {
		return err
	}

	log.Printf("Time spent %s", time.Since(start))

	// create a search index
	log.Printf("Creating a search index...")
	start = time.Now()

	if err = buildIndex(dict, config); err != nil {
		return err
	}

	log.Printf("Time spent %s", time.Since(start))
	log.Printf("End process\n\n")

	return nil
}

// dictionaryReader is an adapter, that implements dictionary.Iterable for bufio.Scanner
type dictionaryReader struct {
	lineScanner *bufio.Scanner
}

// Iterate iterates through each line of the corresponding dictionary
func (dr *dictionaryReader) Iterate(iterator dictionary.Iterator) error {
	docID := dictionary.Key(0)

	for dr.lineScanner.Scan() {
		if err := iterator(docID, dr.lineScanner.Text()); err != nil {
			return err
		}

		docID++
	}

	return dr.lineScanner.Err()
}

// newDictionaryReader creates an adapter to Iterable interface, that scans all lines
// from the SourcePath and creates pairs of <DocID, Value>
func newDictionaryReader(config suggest.IndexDescription) (dictionary.Iterable, error) {
	f, err := os.Open(config.SourcePath)

	if err != nil {
		return nil, fmt.Errorf("Could not open a source file %s", err)
	}

	scanner := bufio.NewScanner(f)

	return &dictionaryReader{
		lineScanner: scanner,
	}, nil
}

// buildIndex builds a search index by using the given config and the dictionary
// and persists it on FS
func buildIndex(dict dictionary.Dictionary, config suggest.IndexDescription) error {
	directory, err := index.NewFSDirectory(config.OutputPath)

	if err != nil {
		return err
	}

	indexWriter := index.NewIndexWriter(
		directory,
		config.CreateWriterConfig(),
		compression.VBEncoder(),
	)

	alphabet := config.CreateAlphabet()
	cleaner := index.NewCleaner(alphabet.Chars(), config.Pad, config.Wrap)
	generator := index.NewGenerator(config.NGramSize)

	if err = index.BuildIndex(dict, indexWriter, generator, cleaner); err != nil {
		return err
	}

	return nil
}

// tryToSendReindexSignal sends a SIGHUP signal to the pid
func tryToSendReindexSignal() error {
	d, err := ioutil.ReadFile(pidPath)

	if err != nil {
		return fmt.Errorf("error parsing pid from %s: %s", pidPath, err)
	}

	pid, err := strconv.Atoi(string(bytes.TrimSpace(d)))

	if err != nil {
		return fmt.Errorf("error parsing pid from %s: %s", pidPath, err)
	}

	if err := syscall.Kill(pid, syscall.SIGHUP); err != nil {
		return fmt.Errorf("Fail to send reindex signal to %d, %s", pid, err)
	}

	return nil
}

// tryToSendReindexRequest sends a http request with a reindex purpose
func tryToSendReindexRequest() error {
	resp, err := http.Post(host, "text/plain", nil)

	if err != nil {
		return fmt.Errorf("Fail to send reindex request to %s, %s", host, err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return fmt.Errorf("Fail to read response body %s", err)
	}

	if string(body) != "OK" {
		return fmt.Errorf("Something goes wrong with reindex request, %s", body)
	}

	return nil
}
