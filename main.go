package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

func main() {
	region := flag.String("region", "eu-west-1", "AWS region")
	key := flag.String("key", "", "SSH key for which instances will be shown")
	flag.Parse()

	svc := ec2.New(session.New(), &aws.Config{Region: region})

	resp, err := svc.DescribeInstances(&ec2.DescribeInstancesInput{})
	if err != nil {
		fmt.Printf("error running DescribeInstances in %s: %s\n", *region, err)
		os.Exit(1)
	}

	keys := make(map[string][]string)

	for _, reservation := range resp.Reservations {
		for _, instance := range reservation.Instances {
			if _, ok := keys[*instance.KeyName]; !ok {
				keys[*instance.KeyName] = []string{}
			}

			keys[*instance.KeyName] = append(keys[*instance.KeyName], *instance.InstanceId)
		}
	}

	for k, v := range keys {
		fmt.Printf("SSH key %s used by %d instances\n", k, len(v))
	}

	if *key != "" {
		if _, ok := keys[*key]; !ok {
			fmt.Printf("SSH key %s not found", *key)
			os.Exit(1)
		}

		fmt.Printf("instance IDs for key %s:\n", *key)
		for _, id := range keys[*key] {
			fmt.Printf("\t%s\n", id)
		}
	}
}
