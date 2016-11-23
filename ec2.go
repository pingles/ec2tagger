package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"strings"
)

type Filter struct {
	NameContains string
	Tags         string
}

func describeInstancesInput(filter *Filter) *ec2.DescribeInstancesInput {
	return &ec2.DescribeInstancesInput{}
}

func instanceName(instance *ec2.Instance) string {
	for _, tag := range instance.Tags {
		if *tag.Key == "Name" {
			return *tag.Value
		}
	}
	return ""
}

func FindInstances(svc *ec2.EC2, filter *Filter) ([]*ec2.Instance, error) {
	in := describeInstancesInput(filter)
	resp, err := svc.DescribeInstances(in)
	if err != nil {
		return nil, err
	}

	instances := make([]*ec2.Instance, 0)
	for _, r := range resp.Reservations {
		for _, i := range r.Instances {
			n := instanceName(i)
			if strings.Contains(n, filter.NameContains) {
				instances = append(instances, i)
			}
		}
	}

	return instances, nil
}

func AddTag(svc *ec2.EC2, key, value string, instances []*ec2.Instance) error {
	resources := make([]*string, len(instances))
	for idx, instance := range instances {
		resources[idx] = instance.InstanceId
	}

	in := &ec2.CreateTagsInput{
		Resources: resources,
		Tags: []*ec2.Tag{
			&ec2.Tag{
				Key:   aws.String(key),
				Value: aws.String(value),
			},
		},
	}

	_, err := svc.CreateTags(in)

	return err
}
