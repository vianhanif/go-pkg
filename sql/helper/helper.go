package helper

import (
	"fmt"
	"net/http"
)

// QueryFilter : represent a query value with its corresponding column and operation
type QueryFilter struct {
	Key       string
	Operation string
	Column    string
	Value     string
	Default   string
	Exec      string
}

// appendFilter : sql helper to build query
func appendFilter(where, key string, column string, operator string, id int) (string, int) {
	value := fmt.Sprintf(`$%d`, id)
	if operator == "" {
		operator = "="
	}
	where = fmt.Sprintf(`%s AND %s %s %s`, where, column, operator, value)
	if id == 1 {
		where = fmt.Sprintf(`WHERE %s %s %s`, column, operator, value)
	}
	return where, id + 1
}

// BuildFilter conbine all parameters into a query statement
func BuildFilter(params ...QueryFilter) (string, []interface{}) {
	where := ""
	args := []interface{}{}
	id := 1
	for _, param := range params {
		switch {
		case param.Operation == "order":
			where = fmt.Sprintf(`%s ORDER BY %s %s`, where, param.Column, param.Value)
		case param.Operation == "limit":
			where = fmt.Sprintf(`%s LIMIT %s`, where, param.Value)
		case param.Operation == "offset":
			where = fmt.Sprintf(`%s OFFSET %s`, where, param.Value)
		case param.Exec != "":
			if where != "" {
				where = fmt.Sprintf(`%s AND %s`, where, param.Exec)
			} else {
				where = fmt.Sprintf("WHERE %s", param.Exec)
			}
		default:
			column := param.Column
			if param.Column == "" {
				column = fmt.Sprintf(`"%s"`, param.Key)
			}
			where, id = appendFilter(where, param.Key, column, param.Operation, id)
			if param.Default != "" && param.Value == "" {
				param.Value = param.Default
			}
			args = append(args, param.Value)
		}
	}
	return where, args
}

func contains(required []string, key string) bool {
	for _, item := range required {
		if item == key {
			return true
		}
	}
	return false
}

// GetQueries : assigning requested QueryFilter values
func GetQueries(r *http.Request, filters []QueryFilter, required []string) ([]QueryFilter, error) {
	var keys []QueryFilter
	var err error
	for index, el := range filters {
		value := r.URL.Query()[el.Key]
		if contains(required, el.Key) && len(value) == 0 {
			if index == 0 {
				err = fmt.Errorf("params_required : %s is empty", el.Key)
			}
			break
		} else {
			var val string
			switch {
			case len(value) > 0 || el.Value != "":
				val = el.Value
				if len(value) > 0 && el.Operation != "offset" && el.Operation != "limit" {
					val = value[0]
				}
			case el.Default != "" && (len(value) == 0 || el.Value == ""):
				val = el.Default
			}
			if val != "" {
				keys = append(keys,
					QueryFilter{
						Key:       el.Key,
						Operation: el.Operation,
						Column:    el.Column,
						Value:     val,
					})
			}
		}
	}
	return keys, err
}
