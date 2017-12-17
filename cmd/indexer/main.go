package main

import (
	"flag"
	"os"
	"log"
	"github.com/alldroll/suggest"
	"github.com/alldroll/cdb"
	"fmt"
	"bufio"
	"encoding/binary"
	"time"
	"io"
	"io/ioutil"
	"encoding/json"
)

type IndexConfig struct {
	Name string `json:"name"`
	NGramSize int `json:"nGramSize"`
	SourcePath string `json:"source"`
	OutputPath string `json:"output"`
	Alphabet []string `json:"alphabet"`
	Pad string `json:"pad"`
	//Wrap [2]string `json:"wrap"`
	Wrap string `json:"wrap"`
}

func readConfig(reader io.Reader) ([]IndexConfig, error) {
	var configs []IndexConfig

	data, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(data, &configs)
	if err != nil {
		return nil, err
	}

	return configs, nil
}

var (
	configPath string
	alphabetMap map[string]suggest.Alphabet
)

func init() {
	flag.StringVar(&configPath, "config", "config.json", "config path")

	alphabetMap = map[string]suggest.Alphabet {
		"english": suggest.NewEnglishAlphabet(),
		"russian": suggest.NewRussianAlphabet(),
		"numbers": suggest.NewNumberAlphabet(),
	}
}

//
func getAlphabet(list []string) suggest.Alphabet {
	alphabets := make([]suggest.Alphabet, 0)

	for _, symbols := range list {
		if alphabet, ok := alphabetMap[symbols]; ok {
			alphabets = append(alphabets, alphabet)
			continue
		}

		alphabets = append(alphabets, suggest.NewSimpleAlphabet([]rune(symbols)))
	}

	return suggest.NewCompositeAlphabet(alphabets)
}

//
func buildDictionary(name, sourcePath, outputPath string) suggest.Dictionary {
	sourceFile, err := os.OpenFile(sourcePath, os.O_RDONLY, 0)
	if err != nil {
		log.Fatalf("cannot open source file %s", err)
	}

	destinationFile, err := os.OpenFile(
		fmt.Sprintf("%s/%s.cdb", outputPath, name),
		os.O_CREATE | os.O_RDWR | os.O_TRUNC,
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
		docId uint32 = 0
		key = make([]byte, 4)
		scanner = bufio.NewScanner(sourceFile)
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

	return suggest.NewCDBDictionary(destinationFile)
}

//
func storeIndex(name string, outputPath string, index suggest.Index) {
	key := make([]byte, 4)
	cdbHandle := cdb.New()

	for length, table := range index {
		if table == nil {
			continue
		}

		destinationFile, err := os.OpenFile(
			fmt.Sprintf("%s/%s.%d.cdb", outputPath, name, length),
			os.O_CREATE | os.O_WRONLY | os.O_TRUNC,
			0644,
		)
		if err != nil {
			log.Fatal(err)
		}

		cdbWriter, err := cdbHandle.GetWriter(destinationFile)
		if err != nil {
			log.Fatal(err)
		}

		for docId, postingList := range table {
			if postingList == nil {
				continue
			}

			value := make([]byte, len(postingList) * 4)
			for i, x := range postingList {
				binary.LittleEndian.PutUint32(value[4*i:], uint32(x))
			}

			binary.LittleEndian.PutUint32(key, uint32(docId))

			err = cdbWriter.Put(key, value)
			if err != nil {
				log.Fatal(err)
			}
		}

		err = cdbWriter.Close()
		if err != nil {
			log.Fatal(err)
		}

		destinationFile.Close()
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
	configs, err := readConfig(f)
	if err != nil {
		log.Fatalf("invalid config file format %s", err)
	}

	totalStart := time.Now()

	for _, config := range configs {
		log.Printf("Start process '%s' config", config.Name)

		alphabet := getAlphabet(config.Alphabet)
		cleaner := suggest.NewCleaner(alphabet.Chars(), config.Pad, config.Wrap)
		generator := suggest.NewGenerator(config.NGramSize, alphabet)

		log.Printf("Building dictionary...")

		start := time.Now()
		dictionary := buildDictionary(config.Name, config.SourcePath, config.OutputPath)
		log.Printf("Time spent %s", time.Since(start))

		// create index in memory
		indexer := suggest.NewIndexer(config.NGramSize, generator, cleaner)
		log.Printf("Creating index...")

		start = time.Now()
		index := indexer.Index(dictionary)
		log.Printf("Time spent %s", time.Since(start))

		// store index on disc
		log.Printf("Storing index...")

		start = time.Now()
		storeIndex(config.Name, config.OutputPath, index)
		log.Printf("Time spent %s", time.Since(start))

		log.Printf("End process\n\n")
	}

	log.Printf("Total time spent %s", time.Since(totalStart).String())
}
