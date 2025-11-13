package feature_list

import (
	"context"
	"errors"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/open-feature/go-sdk/openfeature"
)

type FeatureState bool

const (
	FeatureDisabledByDefault FeatureState = false
	FeatureEnabledByDefault  FeatureState = true
)

// Features a type for storing features and their default states
type Features map[string]FeatureState

// defaultFor get the default state of the feature or an error if it doesn't exist
func (f Features) defaultFor(name string) (FeatureState, error) {
	v, ok := f[name]
	if !ok {
		return FeatureDisabledByDefault, errors.New("feature not found")
	}

	return v, nil
}

// Feature wrapper above the feature name and its default state
type Feature struct {
	Name  string
	State FeatureState
}

// FeatureList represents the state of features
type FeatureList struct {
	provider openfeature.FeatureProvider
	client   *openfeature.Client

	logger   log.Logger
	features Features
}

// IsEnabled check if feature is enabled
func (fl *FeatureList) IsEnabled(name string) bool {
	defaultState, err := fl.features.defaultFor(name)
	if err != nil {
		level.Error(fl.logger).Log("msg", "error checking feature", "name", name, "err", err)
		return false
	}

	feature := Feature{Name: name, State: defaultState}

	featureEnabled, err := fl.client.BooleanValue(
		context.Background(), feature.Name, bool(feature.State), openfeature.EvaluationContext{},
	)
	if err != nil {
		level.Error(fl.logger).Log("msg", "failed to load feature state", "name", name, "err", err)
		return bool(feature.State)
	}

	return featureEnabled

}

// NewFeatureList create a new FeatureList that wrapped around *openfeature.Client
func NewFeatureList(provider openfeature.FeatureProvider, domain string, logger log.Logger, features Features) (*FeatureList, error) {

	if err := openfeature.SetProviderAndWait(provider); err != nil {
		// If a provider initialization error occurs, log it and exit
		level.Error(logger).Log("msg", "failed to set the OpenFeature provider", "err", err)
		return nil, err
	}

	featureList := &FeatureList{
		provider: provider,
		// Initialize OpenFeature client
		client:   openfeature.NewClient(domain),
		logger:   logger,
		features: features,
	}

	return featureList, nil
}
