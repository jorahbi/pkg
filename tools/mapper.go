package mapper

import (
	"encoding/json"
	"reflect"
	"strings"
	"sync"
)

// ----------------------------
// Cache for MapToStruct
// ----------------------------

var fieldCache sync.Map     // map[reflect.Type]map[string]int
var fieldPathCache sync.Map // map[reflect.Type]map[string][]int

// MapToStruct maps a flat map to a struct pointer using json tag or field name
func MapToStruct(m map[string]interface{}, target interface{}) {
	v := reflect.ValueOf(target)
	if v.Kind() != reflect.Ptr || v.Elem().Kind() != reflect.Struct {
		return
	}
	v = v.Elem()
	t := v.Type()

	fieldMap, ok := fieldCache.Load(t)
	if !ok {
		newMap := make(map[string]int)
		for i := 0; i < t.NumField(); i++ {
			f := t.Field(i)
			tag := f.Tag.Get("json")
			name := strings.Split(tag, ",")[0]
			if name == "" {
				name = f.Name
			}
			newMap[name] = i
		}
		fieldCache.Store(t, newMap)
		fieldMap = newMap
	}

	fields := fieldMap.(map[string]int)
	for k, val := range m {
		if idx, found := fields[k]; found {
			field := v.Field(idx)
			if field.CanSet() {
				setField(field, val)
			}
		}
	}
}

func setField(field reflect.Value, val interface{}) {
	if val == nil {
		return
	}
	valValue := reflect.ValueOf(val)
	if valValue.Type().AssignableTo(field.Type()) {
		field.Set(valValue)
		return
	}

	bytes, err := json.Marshal(val)
	if err != nil {
		return
	}
	newVal := reflect.New(field.Type()).Interface()
	if err := json.Unmarshal(bytes, newVal); err == nil {
		field.Set(reflect.ValueOf(newVal).Elem())
	}
}

// ----------------------------
// StructToMap: depth-first, stack-based
// ----------------------------

// StructToMap converts struct to map[string]interface{} with nested support (non-recursive)
func StructToMap(input interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	stack := []struct {
		prefix string
		val    reflect.Value
	}{}

	v := reflect.ValueOf(input)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		return result
	}

	stack = append(stack, struct {
		prefix string
		val    reflect.Value
	}{"", v})

	for len(stack) > 0 {
		// pop
		n := len(stack) - 1
		item := stack[n]
		stack = stack[:n]

		t := item.val.Type()
		for i := 0; i < t.NumField(); i++ {
			field := t.Field(i)
			if field.PkgPath != "" {
				continue
			}
			name := field.Name
			if tag := field.Tag.Get("json"); tag != "" {
				s := strings.Split(tag, ",")[0]
				if s != "" && s != "-" {
					name = s
				}
			}

			key := item.prefix + name
			val := item.val.Field(i)
			if val.Kind() == reflect.Struct && field.Anonymous == false {
				// push nested struct onto stack
				stack = append(stack, struct {
					prefix string
					val    reflect.Value
				}{prefix: key + ".", val: val})
			} else {
				result[key] = val.Interface()
			}
		}
	}
	return result
}

// MapToStructNested maps a map with dot-notation keys to a nested struct
func MapToStructNested(data map[string]interface{}, target interface{}) {
	v := reflect.ValueOf(target)
	if v.Kind() != reflect.Ptr || v.Elem().Kind() != reflect.Struct {
		return
	}
	v = v.Elem()
	t := v.Type()

	paths, ok := fieldPathCache.Load(t)
	if !ok {
		newPaths := make(map[string][]int)
		type stackItem struct {
			path  string
			t     reflect.Type
			index []int
		}
		stack := []stackItem{{"", t, nil}}

		for len(stack) > 0 {
			n := len(stack) - 1
			item := stack[n]
			stack = stack[:n]

			for i := 0; i < item.t.NumField(); i++ {
				field := item.t.Field(i)
				if field.PkgPath != "" {
					continue
				}
				tag := field.Tag.Get("json")
				name := strings.Split(tag, ",")[0]
				if name == "" {
					name = field.Name
				}
				fullPath := name
				if item.path != "" {
					fullPath = item.path + "." + name
				}
				indexPath := append(item.index, i)
				ft := field.Type
				if ft.Kind() == reflect.Ptr {
					ft = ft.Elem()
				}
				if ft.Kind() == reflect.Struct && field.Anonymous == false {
					stack = append(stack, stackItem{fullPath, ft, indexPath})
				}
				newPaths[fullPath] = indexPath
			}
		}
		fieldPathCache.Store(t, newPaths)
		paths = newPaths
	}
	pathMap := paths.(map[string][]int)

	for key, val := range data {
		if path, found := pathMap[key]; found {
			fv := v
			ft := t
			for i, idx := range path {
				field := ft.Field(idx)
				if fv.Kind() == reflect.Ptr {
					fv = fv.Elem()
				}
				f := fv.Field(idx)

				if i < len(path)-1 {
					// make sure nested struct is addressable and set
					if f.Kind() == reflect.Ptr && f.IsNil() {
						ptr := reflect.New(f.Type().Elem())
						f.Set(ptr)
					}
					if f.Kind() == reflect.Struct {
						// ok
					} else if f.Kind() == reflect.Ptr && f.Elem().Kind() == reflect.Struct {
						// ok
					} else {
						break
					}
					fv = f
					ft = field.Type
					if ft.Kind() == reflect.Ptr {
						ft = ft.Elem()
					}
				} else {
					if f.CanSet() {
						setField(f, val)
					}
				}
			}
		}
	}
}

// columns, _ := rows.Columns()
// values := make([]interface{}, len(columns))
// result := make([]map[string]interface{}, 0)

// for rows.Next() {
//     ptrs := make([]interface{}, len(columns))
//     for i := range values {
//         ptrs[i] = &values[i]
//     }

//     if err := rows.Scan(ptrs...); err != nil {
//         return err
//     }

//     row := make(map[string]interface{})
//     for i, col := range columns {
//         row[col] = values[i]
//     }
//     result = append(result, row)
// }

// jsonBytes, err := json.Marshal(result)
