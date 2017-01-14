package cddb

import (
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

func createQueryCmd(cmdArray []string) (queryCmd QueryCmd, ok bool) {
	if len(cmdArray) < 4 {
		return QueryCmd{}, false
	}

	var err error

	queryCmd.discID = cmdArray[0]
	queryCmd.trackCount, err = strconv.Atoi(cmdArray[1])
	if err != nil {
		return QueryCmd{}, false
	}

	if len(cmdArray[2:len(cmdArray)-1]) != queryCmd.trackCount {
		return QueryCmd{}, false
	}

	queryCmd.offsets = make([]int, queryCmd.trackCount+1)

	for i := 0; i < queryCmd.trackCount; i++ {
		offset, err := strconv.Atoi(cmdArray[i+2])
		if err != nil {
			return QueryCmd{}, false
		}
		queryCmd.offsets[i] = offset
	}

	queryCmd.totalSeconds, err = strconv.Atoi(cmdArray[len(cmdArray)-1])
	if err != nil {
		return QueryCmd{}, false
	}

	queryCmd.offsets[len(queryCmd.offsets)-1] = queryCmd.totalSeconds * 75

	if queryCmd.offsets[0] == 0 {
		for i := range queryCmd.offsets {
			queryCmd.offsets[i] += 150
		}
	}

	return queryCmd, true
}

func createReadCmd(cmdArray []string) (readCmd ReadCmd, ok bool) {
	if len(cmdArray) != 2 {
		return ReadCmd{}, false
	}

	readCmd.category = cmdArray[0]
	readCmd.discID = cmdArray[1]

	return readCmd, true
}
