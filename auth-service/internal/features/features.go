package features

import "github.com/Arclight-V/mtch/pkg/feature_list"

var Features = feature_list.Features{
	feature_list.FeatureKafka:      feature_list.FeatureDisabledByDefault,
	feature_list.VerifyCodeEnabled: feature_list.FeatureDisabledByDefault,
}
