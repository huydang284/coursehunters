package main

import (
    "log"
    "runtime/debug"
)

func inSlice(search string, slice []string) bool {
    for _, value := range slice {
        if value == search {
            return true
        }
    }

    return false
}

func check(e error) {
    if e != nil {
        log.Fatalln(e)
        debug.PrintStack()
    }
}
