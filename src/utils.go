package main

import (
    "reflect"
)

func inArray(val interface{}, array interface{}) bool {
    switch reflect.TypeOf(array).Kind() {
    case reflect.Slice:
        s := reflect.ValueOf(array)
        for i := 0; i < s.Len(); i++ {
            if reflect.DeepEqual(val, s.Index(i).Interface()) == true {
                return true
            }
        }
    }

    return false
}
