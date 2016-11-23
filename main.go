package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"gopkg.in/alecthomas/kingpin.v2"
	"os"
	"text/tabwriter"
)

var (
	region  = kingpin.Flag("region", "EC2 Region").Default("eu-west-1").String()
	name    = kingpin.Flag("name", "Name filter").String()
	filters = kingpin.Flag("tags", "Tag Filter (k=v). e.g. foo=bar").StringMap()
	dryRun  = kingpin.Flag("dry-run", "Don't update tags. Shows which instances will be updated.").Bool()

	tagKey   = kingpin.Arg("key", "Key for new tag").String()
	tagValue = kingpin.Arg("value", "Value for new tag.").String()
)

func main() {
	kingpin.Parse()

	sess, err := session.NewSession()
	if err != nil {
		fmt.Println("error creating aws session:", err.Error())
		os.Exit(1)
	}

	svc := ec2.New(sess, &aws.Config{Region: aws.String(*region)})

	filter := &Filter{
		NameContains: *name,
	}
	instances, err := FindInstances(svc, filter)
	if err != nil {
		fmt.Println("error finding instances:", err.Error())
		os.Exit(1)
	}

	if !*dryRun {
		err = AddTag(svc, *tagKey, *tagValue, instances)
		if err != nil {
			fmt.Println("error updating tags:", err.Error())
			os.Exit(1)
		}
	}

	w := tabwriter.NewWriter(os.Stdout, 15, 0, 1, ' ', 0)
	fmt.Fprintln(w, "INSTANCE\tNAME")
	for _, instance := range instances {
		fmt.Fprintf(w, "%s\t%s\n", *instance.InstanceId, instanceName(instance))
	}

	w.Flush()
}
