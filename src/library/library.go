package main

/*
#include <git2.h>
#include <git2/sys/repository.h>
*/
import "C"
import (
	"metro"
	"reflect"
)

//export create
func create(directory *C.char, repo *C.git_repository) C.int {
	dirGo := C.GoString(directory)

	goRepo, _ := metro.Create(dirGo)

	value := reflect.ValueOf(goRepo)
	repo = value.FieldByName("ptr").Interface().(*C.git_repository)
	return 0
}

func main() {}
