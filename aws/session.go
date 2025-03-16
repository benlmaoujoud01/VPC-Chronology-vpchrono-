package aws

import (
	"context"
	"fmt"
	
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/sts"
)

// SessionManager handles AWS authentication and session management
type SessionManager struct {
	region  string
	profile string
	config  aws.Config
}

// NewSessionManager creates a new session manager
func NewSessionManager(region, profile string) (*SessionManager, error) {
	sm := &SessionManager{
		region:  region,
		profile: profile,
	}

	// Load the AWS SDK configuration
	cfg, err := sm.loadConfig()
	if err != nil {
		return nil, err
	}

	sm.config = cfg
	return sm, nil
}

// loadConfig loads the AWS SDK configuration with the specified region and profile
func (sm *SessionManager) loadConfig() (aws.Config, error) {
	ctx := context.Background()
	var opts []func(*config.LoadOptions) error

	// Add region to options
	opts = append(opts, config.WithRegion(sm.region))

	// Add profile to options if specified
	if sm.profile != "" {
		opts = append(opts, config.WithSharedConfigProfile(sm.profile))
	}

	// Load the configuration
	cfg, err := config.LoadDefaultConfig(ctx, opts...)
	if err != nil {
		return aws.Config{}, fmt.Errorf("failed to load AWS configuration: %w", err)
	}

	return cfg, nil
}

// GetEC2Client returns a new EC2 client
func (sm *SessionManager) GetEC2Client() *ec2.Client {
	return ec2.NewFromConfig(sm.config)
}

// GetS3Client returns a new S3 client
func (sm *SessionManager) GetS3Client() *s3.Client {
	return s3.NewFromConfig(sm.config)
}

// ValidateCredentials checks if the current credentials are valid
func (sm *SessionManager) ValidateCredentials() error {
	stsClient := sts.NewFromConfig(sm.config)
	_, err := stsClient.GetCallerIdentity(context.Background(), &sts.GetCallerIdentityInput{})
	if err != nil {
		return fmt.Errorf("failed to validate AWS credentials: %w", err)
	}
	return nil
}

// GetCallerIdentity returns information about the authenticated user
func (sm *SessionManager) GetCallerIdentity() (*sts.GetCallerIdentityOutput, error) {
	stsClient := sts.NewFromConfig(sm.config)
	identity, err := stsClient.GetCallerIdentity(context.Background(), &sts.GetCallerIdentityInput{})
	if err != nil {
		return nil, fmt.Errorf("failed to get caller identity: %w", err)
	}
	return identity, nil
}

// GetRegion returns the configured AWS region
func (sm *SessionManager) GetRegion() string {
	return sm.region
}

// UseRegion creates a new session manager with the specified region
func (sm *SessionManager) UseRegion(region string) (*SessionManager, error) {
	return NewSessionManager(region, sm.profile)
}