package handler

import (
	"os"
)


type Handler interface {
	Read(buf []byte) (n int, err os.Error)
	Write(buf []byte) (n int, err os.Error)
	Close()
}