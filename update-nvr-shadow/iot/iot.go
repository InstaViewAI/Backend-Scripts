package aws

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

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

// UpdateNVRShadow updates exit_dly and nvr_sec_state for a thing
func (ait *IoTCore) UpdateNVRShadow(thingName string) error {
	fmt.Printf("Updating the Shadow for thing: %s\n", thingName)

	payload := ShadowPayloadV2{
		State: StateV2{
			Desired: map[string]any{
				"exit_dly":     30,
				"nvr_sec_state": "Disarmed",
			},
		},
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal shadow payload: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = ait.UpdateThingShadow(ctx, thingName, payloadBytes)
	if err != nil {
		return fmt.Errorf("failed to update shadow for %s: %w", thingName, err)
	}

	fmt.Printf("Successfully updated shadow for thing: %s\n", thingName)
	return nil
}

// UpdateNVRsubShadow will update the shadow to subscription field
func (ait *IoTCore) UpdateNVRSubShadow(thingName string) error {
	fmt.Printf("Updating the Subs Shadow for thing: %s\n", thingName)
	payload := ShadowPayloadV2{
		State: StateV2{
			Desired:map[string]any{
				"subscription": map[string]any{
					"supported_features": []map[string]string{
						// {
						// 	"id":    "pro_monitoring",
						// 	"value": "",
						// },
						{
							"id":    "cloud_storage",
							"value": "30",
						},
					},
				},
			},
		},
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal shadow payload: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err = ait.UpdateThingShadow(ctx, thingName, payloadBytes)
	if err != nil {
		return fmt.Errorf("failed to update shadow for %s: %w", thingName, err)
	}

	fmt.Printf("Subscription shadow updated successfully for thing: %s\n", thingName)
	return nil
}


func (ait *IoTCore) ThingExists(thingName string) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := ait.IotClient.DescribeThing(ctx, &iot.DescribeThingInput{
		ThingName: aws.String(thingName),
	})
	return err == nil
}


