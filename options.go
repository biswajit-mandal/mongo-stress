package main

import (
	"flag"
)

var (
	MongoIP                = "10.87.68.73"
	MongoPort              = "27017"
	MongoUserName          = ""
	MongoUserPasswd        = ""
	FlowsDBStr             = "flows_db"
	PreBuildJSON           = true
	ConcurrentMongoSession = 20
	DefSchemaFile          = "./schemas"
	DefGenCountPerSecond   = 2000
	DefDuration            = 30 /* 30 Seconds */
    RowTimeInMask          = 0x7fffff
    RowTimeInBits          = 23
)

type Stats struct {
	InsertErrorCount   int
	InsertSuccessCount int
}

var (
	GenCountPerSecond int
	SchemaFile        string
	Duration          int
)

func ParseArgs() {
	flag.StringVar(&SchemaFile, "schema_file", DefSchemaFile, "Schema file")
	flag.IntVar(&GenCountPerSecond, "count", DefGenCountPerSecond,
		"Generated count per second")
	flag.IntVar(&Duration, "duration", DefDuration,
		"Stress Test duration in seconds")
    flag.Parse()
}
