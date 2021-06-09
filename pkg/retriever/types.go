package retriever

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/ssm"
)

type SSMClient interface {
	GetParameter(ctx context.Context, params *ssm.GetParameterInput, optFns ...func(*ssm.Options)) (*ssm.GetParameterOutput, error)
}
