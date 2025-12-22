package aws

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/iot"
	"github.com/aws/aws-sdk-go-v2/service/iotdataplane"
)

// Using thing name variable in policy
// https://docs.aws.amazon.com/iot/latest/developerguide/thing-policy-variables.html

// IoTCore is a wrapper for iot package
type IoTCore struct {
	IotClient     *iot.Client
	IotDataClient *iotdataplane.Client
}

// variables are used to create a singleton instance of the iot.Client. at package level
var (
	iotOnce sync.Once
	iotCore *IoTCore
)

// NewIoTCore create a function to create instance of AwsIoTClient
func NewIoTCore() *IoTCore {
	var (
		err error
		cfg aws.Config
	)
	iotOnce.Do(
		func() {
			cfg, err = config.LoadDefaultConfig(context.TODO())
			if err != nil {
				panic(fmt.Errorf("failed to read aws config: %w", err))
			}
			// creating instance of iot.Client using config
			iotSvc := iot.NewFromConfig(cfg)
			// for shadow functions
			iotDataSvc := iotdataplane.NewFromConfig(cfg)

			iotCore = &IoTCore{
				IotClient:     iotSvc,
				IotDataClient: iotDataSvc,
			}
		},
	)
	if err != nil {
		panic(fmt.Errorf("failed to initializing aws service: %w", err))
	}

	return iotCore
}

func (ait *IoTCore) GetThingShadow(ctx context.Context, thingName string) (*iotdataplane.GetThingShadowOutput, error) {
	getThingShadowInput := &iotdataplane.GetThingShadowInput{
		ThingName: aws.String(thingName),
	}

	getThingShadowOutput, err := ait.IotDataClient.GetThingShadow(ctx, getThingShadowInput)
	if err != nil {
		return nil, err
	}

	return getThingShadowOutput, nil
}

func (ait *IoTCore) UpdateThingShadow(ctx context.Context, thingName string, payload []byte) (*iotdataplane.UpdateThingShadowOutput, error) {
	updateThingShadowInput := &iotdataplane.UpdateThingShadowInput{
		ThingName: aws.String(thingName),
		Payload:   payload,
	}

	updateThingShadowOutput, err := ait.IotDataClient.UpdateThingShadow(ctx, updateThingShadowInput)
	if err != nil {
		return nil, err
	}

	return updateThingShadowOutput, err
}

type ShadowPayloadV2 struct {
	State StateV2 `json:"state"`
}

type StateV2 struct {
	Desired map[string]any `json:"desired,omitempty"`
}

func (ait *IoTCore) UpdateBitRateToEnum(thingName string) error {

	fmt.Printf("Reading bitrate for thing: %s\n", thingName)

	sh, err := ait.GetThingShadow(context.TODO(), thingName)
	if err != nil {
		return fmt.Errorf("failed to get thing shadow for %s: %w, %v", thingName, err, sh)
	}

	currShadow := &ShadowPayloadV2{}
	if err := json.Unmarshal(sh.Payload, currShadow); err != nil {
		return fmt.Errorf("failed to unmarshal shadow payload for %s: %w", thingName, err)
	}

	if currShadow.State.Desired != nil && currShadow.State.Desired["0xFC1A:0x00"] != nil {
		// Assuming bitrate is a string that needs to be converted to an enum
		bitrateValue, ok := currShadow.State.Desired["0xFC1A:0x00"].(string)
		if !ok {
			print("thing name: ", thingName, " does not have a valid bitrate value")

			updateBitrateShadow := &ShadowPayloadV2{
				State: StateV2{
					Desired: map[string]any{
						"0xFC1A:0x00": "Auto", // Defaulting
						// to "Auto" if the value is not a string
					},
				},
			}

			payload, err := json.Marshal(updateBitrateShadow)
			if err != nil {
				fmt.Printf("Failed to marshal shadow payload for %s: %v\n", thingName, err)
			}

			_, err = ait.UpdateThingShadow(context.TODO(), thingName, payload)
			if err != nil {
				fmt.Printf("Failed to update shadow for %s: %v\n", thingName, err)
			} else {
				fmt.Printf("Updated bitrate for thing: %s to Auto\n", thingName)
			}

			return nil

		} else {
			fmt.Printf("thing name: %s has bitrate value: %v\n\n", thingName, bitrateValue)
		}
	}

	return nil
}
