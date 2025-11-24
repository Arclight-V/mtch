package features

import "github.com/Arclight-V/mtch/pkg/feature_list"

const StoreCodesInDB = "store-codes-in-DB"

var Features = feature_list.Features{
	feature_list.FeatureKafka: feature_list.FeatureDisabledByDefault,
	StoreCodesInDB:            feature_list.FeatureEnabledByDefault,
}
