package features

import "github.com/Arclight-V/mtch/pkg/feature_list"

const StoreUsersInDB = "store-codes-in-DB"

var Features = feature_list.Features{
	feature_list.FeatureKafka: feature_list.FeatureDisabledByDefault,
	StoreUsersInDB:            feature_list.FeatureDisabledByDefault,
}
