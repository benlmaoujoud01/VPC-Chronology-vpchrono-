package aws

import (
	"context"
	"fmt"
	
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

// VpcInfo contains basic information about a VPC
type VpcInfo struct {
	ID        string
	CidrBlock string
	IsDefault bool
	Name      string
	Tags      map[string]string
}

type Vpc struct {
	ID                   string            // The ID of the VPC
	Name                 string            // The name of the VPC
	Subnets              []Subnet          // A list of subnets in the VPC
	Tags                 map[string]string // The tags associated with the VPC
	CidrBlock            *string           // The primary IPv4 CIDR block for the VPC.
	CidrAssociations     []*string         // Information about the IPv4 CIDR blocks associated with the VPC.
	Ipv6CidrAssociations []*string         // Information about the IPv6 CIDR blocks associated with the VPC.
}
// GetAllVpcs retrieves all VPCs in the current region
func (sm *SessionManager) GetAllVpcs() ([]VpcInfo, error) {
	client := sm.GetEC2Client()
	
	result, err := client.DescribeVpcs(context.Background(), &ec2.DescribeVpcsInput{})
	if err != nil {
		return nil, fmt.Errorf("failed to describe VPCs: %w", err)
	}
	
	vpcs := make([]VpcInfo, 0, len(result.Vpcs))
	for _, vpc := range result.Vpcs {
		vpcInfo := VpcInfo{
			ID:        aws.ToString(vpc.VpcId),
			CidrBlock: aws.ToString(vpc.CidrBlock),
			IsDefault: aws.ToBool(vpc.IsDefault),
			Tags:      make(map[string]string),
		}
		
		// Extract tags
		for _, tag := range vpc.Tags {
			key := aws.ToString(tag.Key)
			value := aws.ToString(tag.Value)
			vpcInfo.Tags[key] = value
			
			// Set name if found
			if key == "Name" {
				vpcInfo.Name = value
			}
		}
		
		vpcs = append(vpcs, vpcInfo)
	}
	
	return vpcs, nil
}

// GetVpcById retrieves a specific VPC by ID
func (sm *SessionManager) GetVpcById(vpcId string) (*VpcInfo, error) {
	client := sm.GetEC2Client()
	
	result, err := client.DescribeVpcs(context.Background(), &ec2.DescribeVpcsInput{
		VpcIds: []string{vpcId},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to describe VPC %s: %w", vpcId, err)
	}
	
	if len(result.Vpcs) == 0 {
		return nil, fmt.Errorf("VPC %s not found", vpcId)
	}
	
	vpc := result.Vpcs[0]
	vpcInfo := &VpcInfo{
		ID:        aws.ToString(vpc.VpcId),
		CidrBlock: aws.ToString(vpc.CidrBlock),
		IsDefault: aws.ToBool(vpc.IsDefault),
		Tags:      make(map[string]string),
	}
	
	// Extract tags
	for _, tag := range vpc.Tags {
		key := aws.ToString(tag.Key)
		value := aws.ToString(tag.Value)
		vpcInfo.Tags[key] = value
		
		// Set name if found
		if key == "Name" {
			vpcInfo.Name = value
		}
	}
	
	return vpcInfo, nil
}

// GetDefaultVpc retrieves the default VPC in the current region
func (sm *SessionManager) GetDefaultVpc() (*VpcInfo, error) {
	client := sm.GetEC2Client()
	
	// Create filter for default VPC
	isDefaultFilter := types.Filter{
		Name:   aws.String("isDefault"),
		Values: []string{"true"},
	}
	
	result, err := client.DescribeVpcs(context.Background(), &ec2.DescribeVpcsInput{
		Filters: []types.Filter{isDefaultFilter},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to describe VPCs: %w", err)
	}
	
	if len(result.Vpcs) == 0 {
		return nil, fmt.Errorf("no default VPC found in region %s", sm.region)
	}
	
	vpc := result.Vpcs[0]
	vpcInfo := &VpcInfo{
		ID:        aws.ToString(vpc.VpcId),
		CidrBlock: aws.ToString(vpc.CidrBlock),
		IsDefault: aws.ToBool(vpc.IsDefault),
		Tags:      make(map[string]string),
	}
	
	// Extract tags
	for _, tag := range vpc.Tags {
		key := aws.ToString(tag.Key)
		value := aws.ToString(tag.Value)
		vpcInfo.Tags[key] = value
		
		// Set name if found
		if key == "Name" {
			vpcInfo.Name = value
		}
	}
	
	return vpcInfo, nil
}