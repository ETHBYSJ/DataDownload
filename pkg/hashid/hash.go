package hashid

import (
	"errors"
	"github.com/speps/go-hashids"
	"go-file-manager/pkg/conf"
)

const (
	UserID = iota
	FileID
	FolderID
)

var (
	ErrTypeNotMatch = errors.New("ID类型不匹配")
)

func HashEncode(v []int) (string, error) {
	hd := hashids.NewData()
	hd.Salt = conf.SystemConfig.HashIDSalt
	h, err := hashids.NewWithData(hd)
	if err != nil {
		return "", err
	}
	id, err := h.Encode(v)
	if err != nil {
		return "", err
	}
	return id, nil
}

func HashDecode(raw string) ([]int, error) {
	hd := hashids.NewData()
	hd.Salt = conf.SystemConfig.HashIDSalt
	h, err := hashids.NewWithData(hd)
	if err != nil {
		return []int{}, err
	}
	return h.DecodeWithError(raw)
}

func HashID(id uint, t int) string {
	v, _ := HashEncode([]int{int(id), t})
	return v
}

func DecodeHashID(id string, t int) (uint, error) {
	v, _ := HashDecode(id)
	if len(v) != 2 || v[1] != t {
		return 0, ErrTypeNotMatch
	}
	return uint(v[0]), nil
}
