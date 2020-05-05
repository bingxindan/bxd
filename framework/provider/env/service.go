package env

import (
	"bufio"
	"bytes"
	"github.com/pkg/errors"
	"io"
	"os"
	"path"
	"strconv"
)

type BxdEnv struct {
	folder string
	maps map[string]string
}

func NewBxdEnv(params ...interface{}) (interface{}, error) {
	if len(params) != 1 {
		return nil, errors.New("NewBxdEnv param error")
	}

	folder := params[0].(string)
	file := path.Join(folder, ".env")
	_, err := os.Stat(file)
	if err != nil || os.IsNotExist(err) {
		return nil, errors.New("file " + file + " not exist:" + err.Error())
	}

	fi, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer fi.Close()

	bxdEnv := &BxdEnv{
		folder:folder,
		maps:map[string]string{},
	}
	br := bufio.NewReader(fi)
	for {
		line, _, c := br.ReadLine()
		if c == io.EOF {
			break
		}
		s := bytes.SplitN(line, []byte{'='}, 2)
		if len(s) < 2 {
			continue
		}
		key := string(s[0])
		val := string(s[1])
		bxdEnv.maps[key] = val
	}

	return bxdEnv, nil
}

func (en *BxdEnv) Get(key string) string {
	if val, ok := en.maps[key]; ok {
		return val
	}
	return ""
}

func (en *BxdEnv) AppEnv() string {
	return en.Get("APP_ENV")
}

func (en *BxdEnv) AppDebug() bool {
	b, err := strconv.ParseBool(en.Get("APP_DEBUG"))
	if err == nil {
		return b
	}
	return false
}

func (en *BxdEnv) AppURL() string {
	return en.Get("APP_URL")
}

func (en *BxdEnv) IsExist(key string) bool {
	_, ok := en.maps[key]
	return ok
}

func (en *BxdEnv) All() map[string]string {
	return en.maps
}
