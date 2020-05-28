package main

import (
	"fmt"
	"log"

	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/compute/v1"
)

func main() {

	project := "travel-prod-76683"
	zone := []string{"asia-southeast1-a", "asia-southeast1-b", "asia-southeast1-c"}
	instance := []string{"travel-order-postgres-01", "travel-order-postgres-02"}
	var labels []string

	ctx := context.Background()

	c, err := google.DefaultClient(ctx, compute.CloudPlatformScope)
	if err != nil {
		log.Fatal(err)
	}

	computeService, err := compute.New(c)
	if err != nil {
		log.Println(err)
	}

	for i := 0; i < len(instance); i++ {
		for j := 0; j < len(zone); j++ {

			resp, err := computeService.Instances.Get(project, zone[j], instance[i]).Context(ctx).Do()

			if err != nil {
				log.Println(err)

			} else {
				labels = append(labels, resp.LabelFingerprint)
				fmt.Printf("%#v\n", resp.LabelFingerprint)
			}
		}
	}

	for _, label := range labels {
		fmt.Println(label)
	}

	for i := 0; i < len(instance); i++ {
		for j := 0; j < len(zone); j++ {

			rb := &compute.InstancesSetLabelsRequest{
				Labels:           map[string]string{"env": "prod", "tribe": "travel", "squad": "travel", "type": "db"},
				LabelFingerprint: labels[i],
			}
			resp, err := computeService.Instances.SetLabels(project, zone[j], instance[i], rb).Context(ctx).Do()

			if err != nil {
				log.Println(err)
			}
			fmt.Printf("%#v\n", resp)
		}
	}

}
