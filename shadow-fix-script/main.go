package main

import (
	"context"
	"encoding/json"
	i "fix-shadow-script/iot"
	"fmt"
)

func main() {
	// TEMP: get shadow for a single device
	thingName := "36sMjQdQgZkRqQWTz5F8ASMY4Ta_cbmb3d2ui4vzzh6z32"

	iotCore := i.NewIoTCore()

	shadow, err := iotCore.GetThingShadow(context.TODO(), thingName)
	if err != nil {
		fmt.Printf("Error getting shadow: %v\n", err)
		return
	}

	var pretty map[string]any
	if err := json.Unmarshal(shadow.Payload, &pretty); err != nil {
		fmt.Printf("Raw payload: %s\n", string(shadow.Payload))
		return
	}

	out, _ := json.MarshalIndent(pretty, "", "  ")
	fmt.Printf("Shadow for %s:\n%s\n", thingName, string(out))
}
