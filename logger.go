package main

import (
	"log"
	"os"
)

func InitLogger() {
	log.SetOutput(os.Stdout)
	log.SetFlags(log.LstdFlags)
}
