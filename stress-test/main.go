package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"reflect"
	"strconv"
	"strings"
	"time"
)

var (
	SchemaKeyStr = "tableindex"
	KeySeperator = ";"
	BuiltDataCnt = 0
)

type SchemaData struct {
	Collection string                 `json:"collection_name"`
	Data       map[string]interface{} `json:"data"`
	Timestamp  int64                  `json:"Timestamp"`
	TableIndex string                 `json:"table_index"`
    T1         int64                    `json:"T1"`
    T2         int64                    `json:"T2"`
}

var mg MongoDBHandler
var stats Stats
var signalCh chan os.Signal

func getKeyVal(key string, val interface{}, idx int, maxidx int) interface{} {
	new_val := val.(string)
	if "string" == new_val || "uuid" == new_val {
		return fmt.Sprintf("%v%d", key, idx)
	} else if "int" == new_val || "double" == new_val {
		return rand.Intn(maxidx)
	} else {
		log.Fatalln("Support not added for type ", new_val)
	}
	return nil
}

func getKeyStrValByType(val interface{}) string {
	var keyVal string
	v := reflect.ValueOf(val)
	switch v.Kind() {
	case reflect.String:
		keyVal = val.(string)
	case reflect.Int:
		keyVal = strconv.Itoa(val.(int))
	}
	return keyVal
}

func buildRow(schemaName string, schema map[string]interface{}, idx int,
	maxidx int) SchemaData {
	var tree = make(map[string]interface{})
	for key, val := range schema {
		keysSplitted := strings.Split(key, ".")
		if key == SchemaKeyStr {
			tree[key] = val
			continue
		}
		tree[key] = getKeyVal(keysSplitted[len(keysSplitted)-1], val, idx,
			maxidx)
	}
	keys := strings.Split(tree[SchemaKeyStr].(string), KeySeperator)
	var keyVal string
	for i := 0; i < len(keys); i++ {
		keyVal += ";" + keys[i] + "=" + getKeyStrValByType(tree[keys[i]])
	}
    timeStamp := time.Now().UnixNano() / 1000000
    t1 := timeStamp & int64(RowTimeInMask)
    RowTimeInBitsUInt := uint(RowTimeInBits)
    t2 := timeStamp >> RowTimeInBitsUInt
	return SchemaData{Data: tree, Collection: schemaName,
		Timestamp: time.Now().UnixNano() / 1000000, TableIndex: keyVal[1:],
        T1: t1, T2: t2}
}

func generateData() {
	var (
		err error
	)
	b, err := ioutil.ReadFile(SchemaFile)
	if err != nil {
		log.Fatalln("Schema file read error:", err)
	}
	var schemas map[string]interface{}
	err = json.Unmarshal(b, &schemas)
	if err != nil {
		log.Fatalln("JSON unmarshal error:", err)
	}
	logDbCollectionCount()
	schemaCnt := len(schemas)
	actCount := (GenCountPerSecond / schemaCnt) * schemaCnt
	log.Printf("Started Generating data for %d count", actCount)
	rate := time.Second / time.Duration(GenCountPerSecond)
	throttle := time.Tick(rate)
	log.Println("Getting rate as:", rate)
	if PreBuildJSON == true {
		preBuildJSONAndPushToDB(schemas, actCount, throttle)
	} else {
		postBuildJSONAndPushToDB(schemas, actCount, throttle)
	}
}

func preBuildJSONAndPushToDB(schemas map[string]interface{}, count int,
	throttle <-chan time.Time) {
	actCount := count * Duration
	resultJSON := make([]SchemaData, 0, actCount*len(schemas))
	log.Println("Pre build json started at:", time.Now(), len(resultJSON))
	for i := 0; i < actCount; i++ {
		for schemaName, schema := range schemas {
			newSchema := schema.(map[string]interface{})
			data := buildRow(schemaName, newSchema, i, count)
			resultJSON = append(resultJSON, data)
		}
	}
	log.Println("Pre build json ended at:", time.Now(), len(resultJSON))
	log.Println("\n\n")
	log.Println("Stress Test Started at:", time.Now())
	for {
		for i := range resultJSON {
			<-throttle
			if i >= len(resultJSON)-1 {
			    log.Println("Reached MAX", i, len(resultJSON))
				preShutdownProcess()
				os.Exit(1)
			}
			go func(i int) {
			    collection := resultJSON[i].Collection
			    createRowInDB(&mg, collection, resultJSON[i])
			}(i)
		}
	}
}

func logDbCollectionCount() {
	b, err := ioutil.ReadFile(SchemaFile)
	if err != nil {
		log.Fatalln("Schema file read error:", err)
	}

	var schemas map[string]interface{}
	err = json.Unmarshal(b, &schemas)
	if err != nil {
		log.Fatalln("JSON unmarshal error:", err)
	}

	for key, _ := range schemas {
		count, err := getCollectionCout(&mg, key)
		if err != nil {
			log.Printf("Getting Count error %v for table %s", err, key)
		}
		log.Printf("Count %v for Table %v", key, count)
	}
}

func postBuildJSONAndPushToDB(schemas map[string]interface{}, count int,
	throttle <-chan time.Time) {
	for {
		<-throttle
		//go func() {
		buildRowsAndPushToDB(schemas, count)
		//}()
	}
}

func buildRowsAndPushToDB(schemas map[string]interface{}, count int) {
	for i := 0; i < count; i++ {
		for schemaName, schema := range schemas {
			newSchema := schema.(map[string]interface{})
			data := buildRow(schemaName, newSchema, i, count)
			createRowInDB(&mg, schemaName, data)
		}
	}
}

func preShutdownProcess() {
	log.Println("Stress Test Ended at:", time.Now())
	statData, _ := json.Marshal(stats)
	log.Println("Getting Stats data as:", string(statData))
	logDbCollectionCount()
	close(signalCh)
}

func main() {
	signalCh = make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt)
	ParseArgs()
	DBSetup(&mg)
	go func() {
		for {
			select {
			case sig := <-signalCh:
				log.Println("Interrupt is detected", sig)
				preShutdownProcess()
				os.Exit(1)
			}
		}
	}()
	generateData()
}
