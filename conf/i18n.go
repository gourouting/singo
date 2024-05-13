package conf

import (
	"os"
	"strings"

	yaml "gopkg.in/yaml.v2"
)

// Dictionary 字典
var Dictionary *map[interface{}]interface{}

// LoadLocales 读取国际化文件
func LoadLocales(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	m := make(map[interface{}]interface{})
	err = yaml.Unmarshal([]byte(data), &m)
	if err != nil {
		return err
	}

	Dictionary = &m

	return nil
}

// T 翻译
func T(key string) string {
	dic := *Dictionary
	keys := strings.Split(key, ".")
	for index, path := range keys {
		// 如果到达了最后一层，寻找目标翻译
		if len(keys) == (index + 1) {
			for k, v := range dic {
				if k, ok := k.(string); ok {
					if k == path {
						if value, ok := v.(string); ok {
							return value
						}
					}
				}
			}
			return path
		}
		// 如果还有下一层，继续寻找
		for k, v := range dic {
			if ks, ok := k.(string); ok {
				if ks == path {
					if dic, ok = v.(map[interface{}]interface{}); ok == false {
						return path
					}
				}
			} else {
				return ""
			}
		}
	}

	return ""
}
