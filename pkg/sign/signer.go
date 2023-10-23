package sign

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/credentials/stscreds"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/aws/smithy-go/logging"
	"github.com/google/uuid"
	"github.com/valyala/bytebufferpool"
)

const (
	ServiceName = "execute-api"
)

type SignConfig struct {
	AccessKeyID string //AWS IAM User Access Key Id
	SecretKey   string //AWS IAM User Secret Key
	Region      string //AWS Region
	RoleArn     string //AWS IAM Role ARN
	MaxRetry    int
	LogMode     []aws.ClientLogMode
	Logger      logging.Logger
}

type AWSSigV4Signer struct {
	cfg                    SignConfig
	awsConfig              aws.Config
	aws4Signer             *v4.Signer
	awsCredentials         *aws.Credentials
	awsStsAssumeRoleOutput *sts.AssumeRoleOutput
}

func NewSigner(cfg SignConfig) *AWSSigV4Signer {

	s := AWSSigV4Signer{
		cfg: cfg,
	}

	s.aws4Signer = v4.NewSigner(
		func(s *v4.SignerOptions) {
			s.DisableURIPathEscaping = true
		},
	)

	return &s
}

func (s *AWSSigV4Signer) autoRefreshCredentials(ctx context.Context) error {

	roleSessionName := uuid.New().String()

	client := sts.NewFromConfig(s.awsConfig)

	role, err := client.AssumeRole(ctx, &sts.AssumeRoleInput{
		RoleArn:         aws.String(s.cfg.RoleArn),
		RoleSessionName: aws.String(roleSessionName),
	})
	if err != nil {
		return err
	}

	s.awsStsAssumeRoleOutput = role

	appCreds := stscreds.NewAssumeRoleProvider(client, s.cfg.RoleArn)

	cred, err := appCreds.Retrieve(ctx)
	if err != nil {
		return err
	}

	s.awsCredentials = &cred

	return nil
}

// Deprecated: Util 20231231. Use new_signer instead.
func (s *AWSSigV4Signer) Sign(r *http.Request) error {

	var clientLogMode aws.ClientLogMode
	defaultMode := aws.LogRetries | aws.LogRequestWithBody | aws.LogResponse

	ctx := r.Context()

	if len(s.cfg.LogMode) > 0 {
		for _, mode := range s.cfg.LogMode {
			clientLogMode |= mode
		}
	} else {
		clientLogMode = defaultMode
	}

	credPvd := credentials.NewStaticCredentialsProvider(s.cfg.AccessKeyID, s.cfg.SecretKey, "")

	awsCfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion(s.cfg.Region),
		config.WithCredentialsProvider(credPvd),
		config.WithRetryMaxAttempts(s.cfg.MaxRetry),
		config.WithClientLogMode(clientLogMode),
		config.WithLogger(s.cfg.Logger),
	)
	if err != nil {
		return nil
	}

	s.awsConfig = awsCfg

	if s.awsCredentials == nil || s.awsCredentials.Expired() {
		if err := s.autoRefreshCredentials(ctx); err != nil {
			return err
		}
	}

	r.Header.Add("X-Amz-Security-Token", *s.awsStsAssumeRoleOutput.Credentials.SessionToken)

	bp := bytebufferpool.Get()
	defer bytebufferpool.Put(bp)
	bp.Reset()

	if r.Body != nil {
		_, err := bp.ReadFrom(r.Body)
		if err != nil {
			return err
		}
		r.ContentLength = int64(bp.Len())
		r.Body = io.NopCloser(bytes.NewReader(bp.Bytes()))
	}

	t := time.Now().UTC()
	sum := sha256.Sum256(bp.Bytes())
	payloadHash := hex.EncodeToString(sum[:])

	return s.aws4Signer.SignHTTP(ctx, *s.awsCredentials, r, payloadHash, ServiceName, s.cfg.Region, t)
}
