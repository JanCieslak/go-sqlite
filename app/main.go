package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github/com/codecrafters-io/sqlite-starter-go/app/sqlite"
	"github/com/codecrafters-io/sqlite-starter-go/app/sqlite/varint"
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

		fmt.Printf("Type: %d\n", sqlite.PageType(pageType))
		fmt.Printf("Free block offset: %d\n", freeBlockOffset)
		fmt.Printf("number of tables: %d\n", numberOfCells)
		fmt.Printf("Content offset: %d\n", contentOffset)
		fmt.Printf("Fragmented free bytes: %d\n", fragmentedFreeBytes)

	case ".tables":
		//log.SetOutput(io.Discard)

		databaseFile, err := os.Open(databaseFilePath)
		if err != nil {
			return err
		}

		header := make([]byte, 100)
		_, err = databaseFile.Read(header)
		if err != nil {
			return err
		}

		var pageSize int16
		if err := binary.Read(bytes.NewReader(header[16:18]), binary.BigEndian, &pageSize); err != nil {
			return fmt.Errorf("failed to read page size integer, err = %w", err)
		}

		fmt.Printf("database page size: %v\n", pageSize)

		page := make([]byte, pageSize)
		n, err := databaseFile.Read(page)
		if err != nil {
			return err
		}

		buf := bytes.NewBuffer(page[:n])
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

		log.Printf("Type: %d\n", sqlite.PageType(pageType))
		log.Printf("Free block offset: %d\n", freeBlockOffset)
		fmt.Printf("number of tables: %d\n", numberOfCells)
		log.Printf("Content offset: %d\n", contentOffset)
		log.Printf("Fragmented free bytes: %d\n", fragmentedFreeBytes)

		offsets := make([]int16, numberOfCells)

		for i := 0; i < int(numberOfCells); i++ {
			var offset int16
			err = binary.Read(buf, binary.BigEndian, &offset)
			if err != nil {
				return err
			}

			offsets[i] = offset - 100
			log.Printf("Cell offset: %d\n", offset)
		}

		tables := make([]string, 0)

		for _, offset := range offsets {
			// TODO: figure out a size
			cell := bytes.NewBuffer(page[offset:])

			numOfBytes, err := varint.Decode(cell)
			if err != nil {
				return err
			}

			rowId, err := varint.Decode(cell)
			if err != nil {
				return err
			}

			payload := make([]byte, numOfBytes)
			err = binary.Read(cell, binary.BigEndian, payload)
			if err != nil {
				return err
			}

			log.Println("Num of bytes:", numOfBytes)
			log.Println("Row id:", rowId)
			log.Println("Payload:", string(payload))

			payloadBuf := bytes.NewReader(payload)

			headerBytes, err := varint.Decode(payloadBuf)
			if err != nil {
				return err
			}
			headerBytes--

			log.Println("Header bytes:", headerBytes)

			serials := make([]int64, 0)

			// todo: this shouldn't be a for loop (it depends on read varint's sizes)
			for i := 0; i < int(headerBytes); i++ {
				serial, err := varint.Decode(payloadBuf)
				if err != nil {
					return err
				}
				serials = append(serials, serial)
			}

			log.Println("Serials:", serials)

		loop:
			for i, serial := range serials {
				switch {
				case serial == 1:
					n, err := payloadBuf.ReadByte()
					if err != nil {
						return err
					}
					log.Printf("Integer: %d\n", n)
				case serial >= 13 && serial%2 == 1:
					text := make([]byte, (serial-13)/2)
					_, err = payloadBuf.Read(text)
					if err != nil {
						return err
					}
					if i == 1 {
						tables = append(tables, string(text))
						// TODO: temp fix: Err too much serials (additional 1 at the end)
						break loop
					}
					log.Printf("Text: %s\n", text)
				}
			}

			log.Println()
		}

		fmt.Printf("%v", tables)

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
