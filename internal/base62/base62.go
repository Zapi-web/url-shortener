package base62

import (
	"errors"
	"math"

	"github.com/Zapi-web/url-shortener/internal/domain"
)

var (
	ErrInvalidCharacter = errors.New("invalid base62 character")
	ErrOverflow         = errors.New("uint64 overflow")
)

type Base62Encoder struct{}

const base62Alphabet = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

var decodeMap [256]int8

func init() {
	for i := 0; i < len(decodeMap); i++ {
		decodeMap[i] = -1
	}

	for i := 0; i < len(base62Alphabet); i++ {
		decodeMap[base62Alphabet[i]] = int8(i)
	}
}

func New() *Base62Encoder {
	return &Base62Encoder{}
}

func (e *Base62Encoder) Encode(num uint64) string {
	if num == 0 {
		return "0"
	}

	var buf [11]byte
	i := len(buf)

	for num > 0 {
		i--
		buf[i] = base62Alphabet[num%62]
		num /= 62
	}

	return string(buf[i:])
}

func (e *Base62Encoder) Decode(str string) (uint64, error) {
	if str == "" {
		return 0, domain.ErrInputisEmpty
	}

	var res uint64

	for i := 0; i < len(str); i++ {
		val := decodeMap[str[i]]

		if val == -1 {
			return 0, ErrInvalidCharacter
		}

		if res > (math.MaxUint64-uint64(val))/62 {
			return 0, ErrOverflow
		}

		res = res*62 + uint64(val)
	}

	return res, nil
}
