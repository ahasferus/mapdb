package main

import (
	"fmt"
	"regexp"
	"strings"
)

type database struct {
	storage map[string]string
}

type args struct {
	key   string
	value string
}

func newDatabase() *database {
	db := new(database)
	db.storage = make(map[string]string)
	return db
}

func (db *database) set(args *args, reply *string) {
	db.storage[args.key] = args.value
	*reply = ""
}

func (db *database) get(args *args, reply *string) {
	var exist bool
	*reply, exist = db.storage[args.key]
	if !exist {
		*reply = fmt.Sprintf("Requsted key %s does not exist", args.key)
	}
}

func (db *database) del(args *args, reply *string) {
	var exist bool
	_, exist = db.storage[args.key]
	if !exist {
		*reply = fmt.Sprintf("Requsted key %s does not exist", args.key)
	} else {
		delete(db.storage, args.key)
		*reply = ""
	}
}

func (db *database) keys(pattern string, reply *string) {
	var keys []string
	for k, _ := range db.storage {
		matched, err := regexp.MatchString(pattern, k)
		if err != nil {
			*reply = fmt.Sprintf("%s", err)
			return
		}
		if matched {
			keys = append(keys, k)
		}
	}
	*reply = strings.Join(keys, ", ")
}
