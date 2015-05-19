package config

import (
	"io/ioutil"
	"strconv"
)

var (
	status_null        = 0
	status_word        = 1
	status_block_start = 2
	status_block_end   = 3
	status_newline     = 4
	status_comment     = 5
)

type itemType int

const (
	itemText       = 1
	itemBlockStart = 2
	itemBlockEnd   = 3
	itemEOF        = 4
)

type item struct {
	typ itemType
	val string
}

func newItem(typ itemType, val string) *item {
	return &item{
		typ: typ,
		val: val,
	}
}

type Config map[string]interface{}

func doParse(content string) Config {
	blocks := make([]Config, 0)
	items := lex(content)
	pairs := make([]*item, 0)
	root := make(Config)
	blocks = append(blocks, root)
	for i := range items {
		if i.typ == itemBlockStart {
			prev := blocks[len(blocks)-1]
			key := pairs[len(pairs)-1].val
			curr := make(Config)
			prev[key] = curr
			blocks = append(blocks, curr)
			pairs = pairs[0 : len(pairs)-1]
		} else if i.typ == itemBlockEnd {
			blocks = blocks[0 : len(blocks)-1]
		} else if i.typ == itemText {
			pairs = append(pairs, i)
			if len(pairs) == 2 {
				key := pairs[len(pairs)-2].val
				val := pairs[len(pairs)-1].val
				pairs = pairs[0 : len(pairs)-2]
				curr := blocks[len(blocks)-1]
				curr[key] = val
			}
		}
	}
	return root
}

func lex(content string) chan *item {
	items := make(chan *item)
	go func() {
		start := 0
		pos := 0
		line := 1
		size := len(content)
		status := status_null
		for pos < size {
			switch content[pos] {
			case ' ':
				if status == status_word {
					items <- newItem(itemText, content[start:pos])
					start = pos
					status = status_null
				} else if status == status_block_start || status == status_block_end {
					status = status_null
				}
				break
			case '\t':
				if status == status_word {
					items <- newItem(itemText, content[start:pos])
					start = pos
					status = status_null
				} else if status == status_block_start || status == status_block_end {
					status = status_null
				}
				break
			case '\n':
				line++
				if status == status_word {
					items <- newItem(itemText, content[start:pos])
					start = pos
					status = status_null
				} else if status == status_block_start || status == status_block_end {
					status = status_null
				} else if status == status_comment {
					status = status_null
				}
				break
			case '{':
				if status == status_word {
					items <- newItem(itemText, content[start:pos])
					start = pos
				}
				status = status_block_start
				items <- newItem(itemBlockStart, "")
				break
			case '}':
				status = status_block_end
				items <- newItem(itemBlockEnd, "")
				break
			case '#':
				status = status_comment
				break
			default:
				if status == status_null {
					start = pos
					status = status_word
				}
			}
			pos++
		}
		items <- newItem(itemEOF, "")
		close(items)
	}()
	return items
}

func (c Config) GetSection(key string) Config {
	val, present := c[key]
	if !present {
		return nil
	}
	return val.(Config)
}

func (c Config) Get(key string) interface{} {
	val, present := c[key]
	if !present {
		return nil
	}
	return val
}

func (c Config) GetStr(key string) string {
	val, present := c[key]
	if !present {
		return ""
	}
	return val.(string)
}

func (c Config) GetInt32(key string) int32 {
	val, err := strconv.ParseInt(c.GetStr(key), 10, 32)
	if err != nil {
		return 0
	}
	return int32(val)
}

func (c Config) GetInt16(key string) int16 {
	val, err := strconv.ParseInt(c.GetStr(key), 10, 16)
	if err != nil {
		return 0
	}
	return int16(val)
}

func (c Config) GetInt8(key string) int8 {
	val, err := strconv.ParseInt(c.GetStr(key), 10, 8)
	if err != nil {
		return 0
	}
	return int8(val)
}

func (c Config) GetUint32(key string) uint32 {
	val, err := strconv.ParseUint(c.GetStr(key), 10, 32)
	if err != nil {
		return 0
	}
	return uint32(val)
}

func (c Config) GetUint16(key string) uint16 {
	val, err := strconv.ParseUint(c.GetStr(key), 10, 16)
	if err != nil {
		return 0
	}
	return uint16(val)
}

func (c Config) GetUint8(key string) uint8 {
	val, err := strconv.ParseUint(c.GetStr(key), 10, 8)
	if err != nil {
		return 0
	}
	return uint8(val)
}

func (c Config) GetBool(key string) bool {
	val, err := strconv.ParseBool(c.GetStr(key))
	if err != nil {
		return false
	}
	return val
}

func Load(cfgfile string) (Config, error) {
	data, err := ioutil.ReadFile(cfgfile)
	if err != nil {
		return nil, err
	}
	return doParse(string(data)), nil
}
