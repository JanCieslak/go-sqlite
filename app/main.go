package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github/com/codecrafters-io/sqlite-starter-go/app/sqlite/btree"
	"log"
	"os"
	// Available if you need it!
	// "github.com/xwb1989/sqlparser"
)

func RunCommand(databaseFilePath string, command string) error {
	switch command {
	case ".dbinfo":
		databaseFile, err := os.Open(databaseFilePath)
		if err != nil {
			return err
		}

		header := make([]byte, 100)
		_, err = databaseFile.Read(header)
		if err != nil {
			return err
		}

		var pageSize uint16
		if err := binary.Read(bytes.NewReader(header[16:18]), binary.BigEndian, &pageSize); err != nil {
			return fmt.Errorf("failed to read page size integer, err = %w", err)
		}

		fmt.Printf("database page size: %v\n", pageSize)

		rest := make([]byte, 10000)
		n, err := databaseFile.Read(rest)
		if err != nil {
			return err
		}

		buf := bytes.NewBuffer(rest[:n])
		pageType, err := buf.ReadByte()
		if err != nil {
			return err
		}

		var freeBlockOffset int16
		err = binary.Read(buf, binary.BigEndian, &freeBlockOffset)
		if err != nil {
			return err
		}

		var numberOfCells int16
		err = binary.Read(buf, binary.BigEndian, &numberOfCells)
		if err != nil {
			return err
		}

		var contentOffset int16
		err = binary.Read(buf, binary.BigEndian, &contentOffset)
		if err != nil {
			return err
		}

		fragmentedFreeBytes, err := buf.ReadByte()
		if err != nil {
			return err
		}

		fmt.Printf("Type: %d\n", btree.PageType(pageType))
		fmt.Printf("Free block offset: %d\n", freeBlockOffset)
		fmt.Printf("number of tables: %d\n", numberOfCells)
		fmt.Printf("Content offset: %d\n", contentOffset)
		fmt.Printf("Fragmented free bytes: %d\n", fragmentedFreeBytes)

	default:
		return fmt.Errorf("unknown command %s", command)
	}

	return nil
}

func main() {
	databaseFilePath := os.Args[1]
	command := os.Args[2]

	if err := RunCommand(databaseFilePath, command); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
