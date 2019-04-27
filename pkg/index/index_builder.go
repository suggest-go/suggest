package index

import "github.com/alldroll/suggest/pkg/dictionary"

// BuildIndex performs building and storing the given dictionary data
func BuildIndex(dict dictionary.Dictionary, writer *Writer, generator Generator, cleaner Cleaner) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err, _ = r.(error)
		}
	}()

	err = dict.Iterate(func(key dictionary.Key, value dictionary.Value) error {
		word := cleaner.CleanAndWrap(value)

		return writer.AddDocument(key, generator.Generate(word))
	})

	if err != nil {
		return err
	}

	if err = writer.Commit(); err != nil {
		return err
	}

	return nil
}
