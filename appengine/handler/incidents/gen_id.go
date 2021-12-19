package incidents

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/speps/go-hashids"
)

func NewRandomID() string {
	rand.Seed(time.Now().UnixNano())
	return NewID(11, strconv.Itoa(rand.Int()))
}

func NewID(digit int, seed ...string) string {
	if digit == 1 {
		panic(fmt.Sprintf("digit=%d. digit must be greater than 1", digit))
	}

	hd := hashids.NewData()
	hd.MinLength = digit
	hd.Salt = strings.Join(seed, "")

	h, err := hashids.NewWithData(hd)
	if err != nil {
		panic(err)
	}

	id, err := h.Encode([]int{42})
	if err != nil {
		panic(err)
	}

	return id
}
