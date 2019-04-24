package core

import (
	"math/rand"
	"reflect"
	"time"
)

func IsInstanceOf(objectPtr, typePtr interface{}) bool {
	return reflect.TypeOf(objectPtr) == reflect.TypeOf(typePtr)
}

func Srand() {
	rand.Seed(time.Now().UnixNano())
}

func Rand(min int, max int) int {
	return rand.Intn(max-min) + min
}
