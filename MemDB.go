package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/iancoleman/orderedmap"
)

type Cmd int

const (
	Get Cmd = iota
	Set
	Del
	Ext
	Unk
)

type Error int

func (e Error) Error() string {
	return "Empty command"
}

const (
	Empty Error = iota
)

type DB interface {
	Set(key string, value string) error

	Get(key string) (string, error)

	Del(key string) (string, error)
}

type memDB struct {
	values *orderedmap.OrderedMap
}

func (mem *memDB) Set(key, value string) error {
	mem.values.Set(key, value)
	return nil
}

func (mem *memDB) Get(key string) (string, error) {
	val, exists := mem.values.Get(key)
	if !exists {
		return "", errors.New("Key not found")
	}

	// Perform type assertion to convert interface{} to string
	byteValue, ok := val.(string)
	if !ok {
		return "", errors.New("Value is not of type string")
	}

	return byteValue, nil
}

func (mem *memDB) Del(key string) (string, error) {
	val, exists := mem.values.Get(string(key))
	if !exists {
		return "", errors.New("Key doesn't exist")
	}

	// Perform type assertion to convert interface{} to string
	byteValue, ok := val.(string)
	if !ok {
		return "", errors.New("Value is not of type string")
	}

	// Remove the key from the OrderedMap
	mem.values.Delete(string(key))

	return byteValue, nil
}

func NewInMem() *memDB {
	values := orderedmap.New()
	return &memDB{
		values,
	}
}

type Repl struct {
	db DB

	in  io.Reader
	out io.Writer
}

func (re *Repl) parseCmd(buf []byte) (Cmd, []string, error) {
	line := string(buf)
	elements := strings.Fields(line)
	if len(elements) < 1 {
		return Unk, nil, Empty
	}

	switch elements[0] {
	case "get":
		return Get, elements[1:], nil
	case "set":
		return Set, elements[1:], nil
	case "del":
		return Del, elements[1:], nil
	case "exit":
		return Ext, nil, nil
	default:
		return Unk, nil, nil
	}
}

func (re *Repl) Start() {
	scanner := bufio.NewScanner(re.in)

	for {
		fmt.Fprint(re.out, "> ")
		if !scanner.Scan() {
			break
		}
		buf := scanner.Text()
		cmd, elements, err := re.parseCmd([]byte(buf))
		if err != nil {
			fmt.Fprintf(re.out, "%s\n", err.Error())
			continue
		}
		switch cmd {
		case Get:
			if len(elements) != 1 {
				fmt.Fprintf(re.out, "Expected 1 argument, received: %d\n", len(elements))
				continue
			}
			v, err := re.db.Get(elements[0])
			if err != nil {
				fmt.Fprintln(re.out, err.Error())
				continue
			}
			fmt.Fprintln(re.out, v)
		case Set:
			if len(elements) != 2 {
				fmt.Printf("Expected 2 arguments, received: %d\n", len(elements))
				continue
			}
			err := re.db.Set(elements[0], elements[1])
			if err != nil {
				fmt.Fprintln(re.out, err.Error())
				continue
			}
		case Del:
			if len(elements) != 1 {
				fmt.Printf("Expected 1 argument, received: %d\n", len(elements))
				continue
			}
			v, err := re.db.Del(elements[0])
			if err != nil {
				fmt.Fprintln(re.out, err.Error())
				continue
			}
			fmt.Fprintln(re.out, v)
		case Ext:
			fmt.Fprintln(re.out, "Bye!")
			return
		case Unk:
			fmt.Fprintln(re.out, "Unknown command")
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(re.out, err.Error())
	} else {
		fmt.Fprintln(re.out, "Bye!")
	}
}

// func main() {
// 	db := NewInMem()
// 	repl := &Repl{
// 		db:  db,
// 		in:  os.Stdin,
// 		out: os.Stdout,
// 	}
// 	repl.Start()
// }
