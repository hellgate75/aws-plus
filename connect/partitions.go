package connect

import (
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go/aws/endpoints"
)

// Describes an AWS Partition information
type AwsPartition struct{
	Id			string
	DnsSuffix	string
	partition	endpoints.Partition
	regions		map[string]endpoints.Region
	services	map[string]endpoints.Service
}

// Gets the list of Endpoint Regions names
func (ap *AwsPartition) GetRegionNames() []string {
	out := make([]string, 0)
	for k, _ := range  ap.regions {
		out = append(out, k)
	}
	return out
}

// Gets a specific Endpoint Region
func (ap *AwsPartition) GetRegion(name string) (*endpoints.Region, error) {
	if r, ok := ap.regions[name]; ok {
		return &r, nil
	}
	return nil, errors.New(fmt.Sprintf("Region %s not found", name))
}

// Gets the list of Endpoint Services names
func (ap *AwsPartition) ServiceNames() []string {
	out := make([]string, 0)
	for k, _ := range  ap.services {
		out = append(out, k)
	}
	return out
}

// Gets a specific Endpoint Service
func (ap *AwsPartition) GetService(id string) (*endpoints.Service, error) {
	if s, ok := ap.services[id]; ok {
		return &s, nil
	}
	return nil, errors.New(fmt.Sprintf("Service %s not found", id))
}

// Acquires Endpoint Resolver for region and service with options
func (ap *AwsPartition) EndpointFor(service, region string, opts ...func(options *endpoints.Options)) (endpoints.ResolvedEndpoint, error) {
	return ap.partition.EndpointFor(service, region, opts...)
}

var awsPartitions = make(map[string]AwsPartition)

var partitionIds = make([]string, 0)

var fulfilled bool

func GetPartitionIds() []string {
	return partitionIds
}

func GetAwsPartition(id string) (*AwsPartition, error) {
	if p, ok := awsPartitions[id]; ok {
		return &p, nil
	}
	return nil, errors.New(fmt.Sprintf("Partition %s not found", id))
}

func InitPartitions() {
	if fulfilled {
		return
	}
	partitions := endpoints.DefaultResolver().(endpoints.EnumPartitions).Partitions()
	for _, partition := range partitions {
		partitionIds = append(partitionIds, partition.ID())
		awsPartitions[partition.ID()]=AwsPartition{
			Id: partition.ID(),
			DnsSuffix: partition.DNSSuffix(),
			partition: partition,
			regions: partition.Regions(),
			services: partition.Services(),
		}
	}
	fulfilled = true
}
