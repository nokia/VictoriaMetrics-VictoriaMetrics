package logger

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
)

var (
	fieldType      = "type"
	fieldClogMsg   = "message"
	fieldClogTime  = "time"
	fieldSystem    = "system"
	fieldSystemId  = "systemid"
	fieldHost      = "host"
	fieldTz        = "timezone"
	fieldContainer = "container"
	fieldTypeValue = "log"
	fieldExtension = "extension"

	logExtensions = flag.String("loggerExt", "", "Key value pairs that need to be added to the emitting logs, it needs to be json string")
)

func getEnvOrDefault(key string, defValue string) string {
	val, set := os.LookupEnv(key)
	if set {
		return val
	} else {
		return defValue
	}
}

func getExtensionAsString() string {
	if *logExtensions != "" {
		jsonString := *logExtensions
		// Parse the JSON string into a map
		var myMap map[string]interface{}
		err := json.Unmarshal([]byte(jsonString), &myMap)
		if err != nil {
			panic(fmt.Errorf("error parsing Extensions:%q", err))
		}

		// Convert the map back to a JSON string
		jsonBytes, err := json.Marshal(myMap)
		if err != nil {
			panic(fmt.Errorf("error converting to string %q", err))

		}
		return string(jsonBytes)
	}
	return ""
}

// This shall be called after initTimeZone as this is referring to that variable
func getClogMessage() string {
	fieldSystemValue := getEnvOrDefault("SYSTEM", "BCMT")
	fieldSystemIdValue := getEnvOrDefault("SYSTEM_ID", "BCMT_ID")
	fieldHostValue := getEnvOrDefault("HOST", "localhost.default")
	fieldContainerValue := getEnvOrDefault("CONTAINER_NAME", os.Args[0])
	extAsString := getExtensionAsString()
	var clogMsg string
	if *logExtensions != "" {
		clogMsg = fmt.Sprintf(
			"%q:%q,%q:%q,%q:%q,%q:%q,%q:%q,%q:%q,%q:%s",
			fieldType, fieldTypeValue,
			fieldSystem, fieldSystemValue,
			fieldSystemId, fieldSystemIdValue,
			fieldHost, fieldHostValue,
			fieldContainer, fieldContainerValue,
			fieldTz, timezone,
			fieldExtension, extAsString)
	} else {
		clogMsg = fmt.Sprintf(
			"%q:%q,%q:%q,%q:%q,%q:%q,%q:%q,%q:%q",
			fieldType, fieldTypeValue,
			fieldSystem, fieldSystemValue,
			fieldSystemId, fieldSystemIdValue,
			fieldHost, fieldHostValue,
			fieldContainer, fieldContainerValue,
			fieldTz, timezone)

	}

	return clogMsg
}
