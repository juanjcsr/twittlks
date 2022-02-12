package db

import (
	"bufio"
	"log"
	"os"

	"github.com/juanjcsr/twittlks/lks"
)

func ReadLineFromFile(filename string) (*[]lks.TuitLike, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(f)
	tlList := []lks.TuitLike{}
	for scanner.Scan() {
		line := scanner.Bytes()
		tl, err := lks.LineByteToTuitLike(line)
		if err != nil {
			log.Println(err)
			continue
		}
		tlList = append(tlList, *tl)
	}

	return &tlList, nil
}
