package utils

import "reflect"

// InArray checks if element exist or not in array
// Ex. InArray("Dog", []string{"Cat", "Cow", "Dog"})
// Where "Dog" will be search string and second argument will be array in which we need to find
func InArray(val interface{}, array interface{}) (exists bool, index int) {
	exists = false
	index = -1

	switch reflect.TypeOf(array).Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(array)

		for i := 0; i < s.Len(); i++ {
			if reflect.DeepEqual(val, s.Index(i).Interface()) == true {
				index = i
				exists = true
				return
			}
		}
	}

	return
}
