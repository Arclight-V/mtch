package feature_list

import (
	"testing"
)

func TestFeatures_defaultFor_Existing(t *testing.T) {
	f := Features{FeatureKafka: FeatureDisabledByDefault}

	state, err := f.defaultFor(FeatureKafka)
	if err != nil {
		t.Fatal("unexpected error", err)
	}
	if state != FeatureDisabledByDefault {
		t.Fatalf("unexpected state, got %v, want %v", state, FeatureDisabledByDefault)
	}
}

func TestFeatures_defaultFor_Missing(t *testing.T) {
	f := Features{FeatureKafka: FeatureDisabledByDefault}

	state, err := f.defaultFor("NoopFeature")
	if err == nil {
		t.Fatal("unexpected error", err)
	}
	if state != FeatureDisabledByDefault {
		t.Fatalf("unexpected state, got %v, want %v", state, FeatureDisabledByDefault)
	}
}
