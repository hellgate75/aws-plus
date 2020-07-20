package io

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"gopkg.in/yaml.v2"
	"strings"
)

type Encoding string

const (
	EncodingJson Encoding = Encoding("json")
	EncodingXml Encoding = Encoding("xml")
	EncodingYaml Encoding = Encoding("yaml")
	EncodingGob Encoding = Encoding("gob")
	EncodingText Encoding = Encoding("text")
	EncodingUnknown Encoding = Encoding("")
)

func ParseEncoding(e string) Encoding {
	switch strings.ToLower(e) {
	case "json":
		return EncodingJson
	case "yaml":
		return EncodingYaml
	case "xml":
		return EncodingXml
	case "text":
		return EncodingText
	case "gob":
		return EncodingGob
	default:
		return EncodingUnknown
	}
}

func EncodeValue(in interface{}, enc Encoding) ([]byte, error) {
	var err error
	var out = make([]byte, 0)
	defer func() {
		if r := recover(); r != nil {
			err = errors.New(fmt.Sprintf("%v", r))
		}
	}()
	switch enc {
	case EncodingJson:
		out, err = json.Marshal(in)
	case EncodingXml:
		out, err = xml.Marshal(in)
	case EncodingYaml:
		out, err = yaml.Marshal(in)
	case EncodingGob:
		out, err = EncodeGobValue(in)
	case EncodingText:
		out, err = EncodeTextFormatSummary(in)
	default:
		err = errors.New(fmt.Sprintf("Unknown encoding format: %v", enc))
	}
	return out, err
}

func EncodeGobValue(in interface{}) ([]byte, error) {
	var err error
	var out = make([]byte, 0)
	defer func() {
		if r := recover(); r != nil {
			err = errors.New(fmt.Sprintf("%v", r))
		}
	}()
	var buff = bytes.NewBuffer([]byte{})
	err = gob.NewEncoder(buff).Encode(in)
	if err == nil {
		out = buff.Bytes()
		buff.Reset()
	}
	return out, err
}


func DecodeValue(out interface{}, in []byte, enc Encoding) error {
	var err error
	defer func() {
		if r := recover(); r != nil {
			err = errors.New(fmt.Sprintf("%v", r))
		}
	}()
	switch enc {
	case EncodingJson:
		err = json.Unmarshal(in, out)
	case EncodingXml:
		err = xml.Unmarshal(in, out)
	case EncodingYaml:
		err = yaml.Unmarshal(in, out)
	case EncodingGob:
		err = DecodeGobValue(in, out)
	default:
		err = errors.New(fmt.Sprintf("Unknown encoding format: %v", enc))
	}
	return err
}

func DecodeGobValue(in []byte, out interface{}) error  {
	var err error
	defer func() {
		if r := recover(); r != nil {
			err = errors.New(fmt.Sprintf("%v", r))
		}
	}()
	var buff = bytes.NewBuffer(in)
	err = gob.NewDecoder(buff).Decode(out)
	return err
}
