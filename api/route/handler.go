package route

import (
	"context"

	oas "github.com/onyxia-datalab/onyxia-onboarding/api/oas"
)

type MyHandler struct {
	oas.UnimplementedHandler
	onboardImpl func(ctx context.Context, req *oas.OnboardingRequest) (oas.OnboardRes, error)
}

func (h *MyHandler) Onboard(
	ctx context.Context,
	req *oas.OnboardingRequest,
) (oas.OnboardRes, error) {
	return h.onboardImpl(ctx, req)
}

var _ oas.Handler = (*MyHandler)(nil)
