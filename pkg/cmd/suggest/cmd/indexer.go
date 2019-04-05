package cmd

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"syscall"
	"time"

	"github.com/spf13/cobra"

	"github.com/alldroll/cdb"
	"github.com/alldroll/suggest/pkg/compression"
	"github.com/alldroll/suggest/pkg/dictionary"
	"github.com/alldroll/suggest/pkg/index"
	"github.com/alldroll/suggest/pkg/suggest"
)

var (
	dict string
	host string
)

const tmpSuffix = ".tmp"

func init() {
	indexCmd.Flags().StringVarP(&dict, "dict", "d", "", "reindex certain dict")
	indexCmd.Flags().StringVarP(&host, "host", "", "", "host to send reindex request")

	rootCmd.AddCommand(indexCmd)
}

var indexCmd = &cobra.Command{
	Use:   "indexer -c [config file]",
	Short: "builds indexes and dictionaries",
	Long:  `builds indexes and dictionaries and send signal to reload fresh data`,
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

//
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

//
func indexJob(config suggest.IndexDescription) error {
	log.Printf("Start process '%s' config", config.Name)

	alphabet := config.CreateAlphabet()
	cleaner := index.NewCleaner(alphabet.Chars(), config.Pad, config.Wrap)
	generator := index.NewGenerator(config.NGramSize)

	// create dictionary
	log.Printf("Building dictionary...")

	start := time.Now()

	dictionary, err := buildDictionary(config.SourcePath, config)
	if err != nil {
		return err
	}

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

	if err := storeIndex(indices, config); err != nil {
		return err
	}

	log.Printf("Time spent %s", time.Since(start))

	log.Printf("End process\n\n")

	return nil
}

//
func buildDictionary(sourcePath string, config suggest.IndexDescription) (dictionary.Dictionary, error) {
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
		return nil, fmt.Errorf("Fail to create dictionary file %s", err)
	}

	cdbHandle := cdb.New()
	cdbWriter, err := cdbHandle.GetWriter(destinationFile)
	if err != nil {
		return nil, fmt.Errorf("Fail to create cdb writer %s", err)
	}

	var (
		docID   uint32
		key     = make([]byte, 4)
		scanner = bufio.NewScanner(sourceFile)
	)

	for scanner.Scan() {
		binary.LittleEndian.PutUint32(key, docID)
		err = cdbWriter.Put(key, scanner.Bytes())
		if err != nil {
			return nil, fmt.Errorf("Fail to put record to cdb %s", err)
		}

		docID++
	}

	log.Printf("Number of string %d", docID)

	err = cdbWriter.Close()
	if err != nil {
		return nil, fmt.Errorf("Fail to save cdb dictionary %s", err)
	}

	if err := os.Rename(config.GetDictionaryFile()+tmpSuffix, config.GetDictionaryFile()); err != nil {
		return nil, fmt.Errorf("Error to rename file %s", err)
	}

	return dictionary.NewCDBDictionary(destinationFile), nil
}

//
func storeIndex(indices index.Indices, config suggest.IndexDescription) error {
	writer := index.NewOnDiscIndicesWriter(
		compression.VBEncoder(),
		config.GetHeaderFile()+tmpSuffix,
		config.GetDocumentListFile()+tmpSuffix,
	)

	if err := writer.Save(indices); err != nil {
		return fmt.Errorf("Fail to save index %s", err)
	}

	if err := os.Rename(config.GetHeaderFile()+tmpSuffix, config.GetHeaderFile()); err != nil {
		return fmt.Errorf("Error to rename file %s", err)
	}

	if err := os.Rename(config.GetDocumentListFile()+tmpSuffix, config.GetDocumentListFile()); err != nil {
		return fmt.Errorf("Error to rename file %s", err)
	}

	return nil
}

//
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

//
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
