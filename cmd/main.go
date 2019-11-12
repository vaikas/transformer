package main

import (
	"context"
	"fmt"
	"log"
	"os"

	cloudevents "github.com/cloudevents/sdk-go"
	"github.com/kelseyhightower/envconfig"
)

type envConfig struct {
	// Port on which to listen for cloudevents
	Port int    `envconfig:"RCV_PORT" default:"8080"`
	Path string `envconfig:"RCV_PATH" default:"/"`
	Step string `envconfig:"STEP" default:"0"`
}

func main() {
	var env envConfig
	if err := envconfig.Process("", &env); err != nil {
		log.Printf("[ERROR] Failed to process env var: %s", err)
		os.Exit(1)
	}
	os.Exit(_main(os.Args[1:], env))
}

type Example struct {
	Sequence int    `json:"id"`
	Message  string `json:"message"`
}

type stepper struct {
	step string
}

func NewStepper(step string) *stepper {
	return &stepper{step: step}
}

func (s *stepper) gotEvent(ctx context.Context, event cloudevents.Event, resp *cloudevents.EventResponse) error {
	fmt.Printf("Got Event Context: %+v\n", event.Context)
	data := &Example{}
	if err := event.DataAs(data); err != nil {
		fmt.Printf("Got Data Error: %s\n", err.Error())
	}
	fmt.Printf("Got Data: %+v\n", data)
	fmt.Printf("Got Transport Context: %+v\n", cloudevents.HTTPTransportContextFrom(ctx))
	fmt.Printf("----------------------------\n")

	responseData := Example{
		Sequence: data.Sequence,
		// Just tack our step number to the Message to demo changing the event as it traverses
		// the sequence.
		Message: fmt.Sprintf("%s - Handled by %s", data.Message, s.step),
	}

	r := cloudevents.NewEvent()
	r.SetSource(fmt.Sprintf("/transformer/%s", s.step))
	r.SetType("samples.http.mod3")
	r.SetID(event.Context.GetID())
	r.SetData(responseData)
	r.SetDataContentType(cloudevents.ApplicationJSON)
	resp.RespondWith(200, &r)
	return nil
}

func _main(args []string, env envConfig) int {
	ctx := context.Background()

	t, err := cloudevents.NewHTTPTransport(
		cloudevents.WithPort(env.Port),
		cloudevents.WithPath(env.Path),
	)
	if err != nil {
		log.Fatalf("failed to create transport: %s", err.Error())
	}
	c, err := cloudevents.NewClient(t,
		cloudevents.WithUUIDs(),
		cloudevents.WithTimeNow(),
	)
	if err != nil {
		log.Fatalf("failed to create client: %s", err.Error())
	}

	stepper := NewStepper(env.Step)
	if err := c.StartReceiver(ctx, stepper.gotEvent); err != nil {
		log.Fatalf("failed to start receiver: %s", err.Error())
	}

	log.Printf("listening on :%d%s\n", env.Port, env.Path)
	<-ctx.Done()

	return 0
}
