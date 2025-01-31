package domain

type OnboardingRequest struct {
	Group string
}

type OnboardingService interface {
	Onboard(req OnboardingRequest) error
}
