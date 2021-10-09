package shapeziocompression

import (
	"encoding/hex"
	"encoding/json"
	"hash/crc32"
	"os"
	"strings"

	LZString "github.com/Lazarus/lz-string-go"
)

const _keyStrUriSafe = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+-$"
const _charmap = "!#%&'()*+,-./:;<=>?@[]^_`{|}~¥¦§¨©ª«¬\u00ac\u00ad®¯°±²³´µ¶·¸¹º»¼½¾¿ÀÁÂÃÄÅÆÇÈÉÊËÌÍÎÏÐÑÒÓÔÕÖ×ØÙÚÛÜÝÞßàáâãäåæçèéêëìíîïðñòóôõö÷øùúûüýþÿABCDEFGHIJKLMNOPQRSTUVWXYZ"
const _crcPrefix = "crc32---------------------------"

var compressionCache = map[int]string{}
var decompressionCache = map[string]int{}

type ObjectInternalJSON struct {
	Keys   []string               `json:"keys"`
	Values []string               `json:"values"`
	Data   map[string]interface{} `json:"data"`
}

func UnmarshalObjectInternalJSON(data []byte) (ObjectInternalJSON, error) {
	var r ObjectInternalJSON
	err := json.Unmarshal(data, &r)
	return r, err
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func indexOfArray(slice []string, item string) int {
	for i := range slice {
		if slice[i] == item {
			return i
		}
	}
	return -1
}

func padStart(a string, b string, l int) string {
	for i := len(a); i < l; i++ {
		a = b + a
	}
	return a
}

func computeCrc(a []byte) string {
	crc := crc32.NewIEEE()
	crc.Write(a)
	return _crcPrefix + padStart(hex.EncodeToString(crc.Sum(nil)), "0", 8)
}

func CompressFile(file1 string, file2 string) {
	dat, err := os.ReadFile(file1)
	check(err)
	var o1 map[string]interface{}
	err = json.Unmarshal(dat, &o1)
	check(err)
	k1 := []string{}
	v1 := []string{}
	o := CompressObjectInternal(o1, &k1, &v1)
	if d1, ok := o.(map[string]interface{}); ok {
		o2 := ObjectInternalJSON{Data: d1, Keys: k1, Values: v1}
		c, err := json.Marshal(o2)
		check(err)
		err = os.WriteFile("1.bin", []byte(c), 0666)
		check(err)
		b := computeCrc(c)
		a := []byte{}
		a = append(a, b...)
		a = append(a, c...)
		err = os.WriteFile("2.bin", a, 0666)
		check(err)
		z := LZString.Compress(string(a), _keyStrUriSafe)
		z2 := []byte{0x01}
		z2 = append(z2, []byte(z)...)
		err = os.WriteFile(file2, z2, 0666)
		check(err)
	}
}

func DecompressFile(file1 string, file2 string) {
	dat, err := os.ReadFile(file1)
	check(err)
	a, err := LZString.Decompress(string(dat[1:]), _keyStrUriSafe)
	check(err)
	b := a[40:]
	obj, err := UnmarshalObjectInternalJSON([]byte(b))
	check(err)
	o := DecompressObjectInternal(obj.Data, obj.Keys, obj.Values)
	o1, err := json.Marshal(o)
	check(err)
	err = os.WriteFile(file2, o1, 0666)
	check(err)
}

func DecompressObjectInternal(obj interface{}, keys []string, values []string) interface{} {
	if o, ok := obj.([]interface{}); ok {
		var result []interface{}
		for i := 0; i < len(o); i++ {
			result = append(result, DecompressObjectInternal(o[i], keys, values))
		}
		return result
	} else if o, ok := obj.(map[string]interface{}); ok {
		result := make(map[string]interface{})
		for k := range o {
			realIndex := DecompressInt(k)
			for realIndex < len(keys) {
				realIndex += len(keys)
			}
			for realIndex >= len(keys) {
				realIndex -= len(keys)
			}
			result[keys[realIndex]] = DecompressObjectInternal(o[k], keys, values)
		}
		return result
	} else if o, ok := obj.(string); ok {
		realIndex := DecompressInt(o)
		return values[realIndex]
	}
	return obj
}

func DecompressInt(s string) int {
	if v, ok := decompressionCache[s]; ok {
		return v
	}

	result := 0
	for i := len(s) - 1; i >= 0; i-- {
		result = result*len(_charmap) + strings.Index(_charmap, string(s[i]))
	}
	result -= 1
	decompressionCache[s] = result
	return result
}

func CompressObjectInternal(obj interface{}, keys *[]string, values *[]string) interface{} {
	if o, ok := obj.([]interface{}); ok {
		var result []interface{}
		for i := 0; i < len(o); i++ {
			result = append(result, CompressObjectInternal(o[i], keys, values))
		}
		return result
	} else if o, ok := obj.(map[string]interface{}); ok {
		result := make(map[string]interface{})
		for k := range o {
			index := indexOfArray(*keys, k)
			if index == -1 {
				index = len(*keys)
				*keys = append(*keys, k)
			}
			result[CompressInt(index)] = CompressObjectInternal(o[k], keys, values)
		}
		return result
	} else if o, ok := obj.(string); ok {
		index := indexOfArray(*values, o)
		if index == -1 {
			index = len(*values)
			*values = append(*values, o)
		}
		return CompressInt(index)
	}
	return obj
}

func CompressInt(a int) string {
	a += 1

	cache_key := a

	if v, ok := compressionCache[cache_key]; ok {
		return v
	}

	result := ""
	for a > 0 {
		result += string(_charmap[a%len(_charmap)])
		a = a / len(_charmap)
	}

	compressionCache[cache_key] = result
	return result
}
