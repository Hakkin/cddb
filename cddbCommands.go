package cddb

import (
	"fmt"
	"log"
	"strconv"
)

type QueryCmd struct {
	discID       string
	trackCount   int
	offsets      []int
	totalSeconds int
	language     string
	country      string
}

type ReadCmd struct {
	category string
	discID   string
	language string
	country  string
}

var syntaxError error = fmt.Errorf("%v", cddbStatus(500, "Command syntax error", true))

func logSyntaxError(cmdArray []string) {
	log.Println("syntax error:")
	log.Println(cmdArray)
}

func createQueryCmd(cmdArray []string) (queryCmd QueryCmd, err error) {
	if len(cmdArray) < 4 {
		logSyntaxError(cmdArray)
		return QueryCmd{}, syntaxError
	}

	queryCmd.discID = cmdArray[0]
	queryCmd.trackCount, err = strconv.Atoi(cmdArray[1])
	if err != nil {
		logSyntaxError(cmdArray)
		return QueryCmd{}, syntaxError
	}

	if len(cmdArray[2:len(cmdArray)-1]) != queryCmd.trackCount {
		logSyntaxError(cmdArray)
		return QueryCmd{}, syntaxError
	}

	for i := 0; i < queryCmd.trackCount; i++ {
		offset, err := strconv.Atoi(cmdArray[i+2])
		if err != nil {
			logSyntaxError(cmdArray)
			return QueryCmd{}, syntaxError
		}
		queryCmd.offsets = append(queryCmd.offsets, offset)
	}
	
	queryCmd.totalSeconds, err = strconv.Atoi(cmdArray[len(cmdArray)-1])
	if err != nil {
		logSyntaxError(cmdArray)
		return QueryCmd{}, syntaxError
	}
	
	queryCmd.offsets = append(queryCmd.offsets, queryCmd.totalSeconds*75)

	if queryCmd.offsets[0] == 0 {
		for i := range queryCmd.offsets {
			queryCmd.offsets[i] += 150
		}
	}

	return queryCmd, nil
}

func createReadCmd(cmdArray []string) (readCmd ReadCmd, err error) {
	if len(cmdArray) != 2 {
		logSyntaxError(cmdArray)
		return ReadCmd{}, syntaxError
	}

	readCmd.category = cmdArray[0]
	readCmd.discID = cmdArray[1]

	return readCmd, nil
}
