package main

import (
	"fmt"
	"log"
	"os"
	
	"github.com/benlmaoujoud/vpchrono/aws"
)

func main() {
	
	region := "us-east-1"
	if len(os.Args) > 1 {
		region = os.Args[1]
	}
	
	sm, err := aws.NewSessionManager(region, "")
	if err != nil {
		log.Fatalf("Failed to create session manager: %v", err)
	}
	
	// Validate credentials
	if err := sm.ValidateCredentials(); err != nil {
		log.Fatalf("Invalid AWS credentials: %v", err)
	}
	
	// Get and print caller identity
	identity, err := sm.GetCallerIdentity()
	if err != nil {
		log.Fatalf("Failed to get caller identity: %v", err)
	}
	fmt.Printf("Authenticated as: %s\n", *identity.Arn)
	fmt.Printf("Account ID: %s\n", *identity.Account)
	
	// Get all VPCs
	vpcs, err := sm.GetAllVpcs()
	if err != nil {
		log.Fatalf("Failed to get VPCs: %v", err)
	}
	
	// Display VPC information
	fmt.Printf("\nFound %d VPCs in region %s:\n", len(vpcs), sm.GetRegion())
	for i, vpc := range vpcs {
		fmt.Printf("%d. VPC ID: %s\n", i+1, vpc.ID)
		fmt.Printf("   Name: %s\n", vpc.Name)
		fmt.Printf("   CIDR Block: %s\n", vpc.CidrBlock)
		fmt.Printf("   Is Default: %t\n", vpc.IsDefault)
		fmt.Println()
	}
	
	// Try to get the default VPC
	defaultVpc, err := sm.GetDefaultVpc()
	if err != nil {
		fmt.Printf("No default VPC found: %v\n", err)
	} else {
		fmt.Println("Default VPC Information:")
		fmt.Printf("   VPC ID: %s\n", defaultVpc.ID)
		fmt.Printf("   Name: %s\n", defaultVpc.Name)
		fmt.Printf("   CIDR Block: %s\n", defaultVpc.CidrBlock)
	}
}