package v3

import (
	"reflect"
	"sort"
	"testing"

	v3 "github.com/SigNoz/signoz/pkg/query-service/model/v3"
	. "github.com/smartystreets/goconvey/convey"
)

var testEnrichmentRequiredData = []struct {
	Name               string
	Params             v3.QueryRangeParamsV3
	EnrichmentRequired bool
}{
	{
		Name: "attribute enrichment not required",
		Params: v3.QueryRangeParamsV3{
			CompositeQuery: &v3.CompositeQuery{
				BuilderQueries: map[string]*v3.BuilderQuery{
					"test": {
						QueryName:  "test",
						Expression: "test",
						DataSource: v3.DataSourceLogs,
						AggregateAttribute: v3.AttributeKey{
							Key:      "test",
							Type:     v3.AttributeKeyTypeTag,
							DataType: v3.AttributeKeyDataTypeInt64,
						},
					},
				},
			},
		},
		EnrichmentRequired: true,
	},
	{
		Name: "attribute enrichment required",
		Params: v3.QueryRangeParamsV3{
			CompositeQuery: &v3.CompositeQuery{
				BuilderQueries: map[string]*v3.BuilderQuery{
					"test": {
						QueryName:  "test",
						Expression: "test",
						DataSource: v3.DataSourceLogs,
						AggregateAttribute: v3.AttributeKey{
							Key: "test",
						},
					},
				},
			},
		},
		EnrichmentRequired: true,
	},
	{
		Name: "filter enrichment required",
		Params: v3.QueryRangeParamsV3{
			CompositeQuery: &v3.CompositeQuery{
				BuilderQueries: map[string]*v3.BuilderQuery{
					"test": {
						QueryName:  "test",
						Expression: "test",
						DataSource: v3.DataSourceLogs,
						Filters: &v3.FilterSet{Operator: "AND", Items: []v3.FilterItem{
							{Key: v3.AttributeKey{Key: "user_name", DataType: v3.AttributeKeyDataTypeString, Type: v3.AttributeKeyTypeTag}, Value: "john", Operator: "="},
						}},
					},
				},
			},
		},
		EnrichmentRequired: true,
	},
	{
		Name: "filter enrichment required",
		Params: v3.QueryRangeParamsV3{
			CompositeQuery: &v3.CompositeQuery{
				BuilderQueries: map[string]*v3.BuilderQuery{
					"test": {
						QueryName:  "test",
						Expression: "test",
						DataSource: v3.DataSourceLogs,
						Filters: &v3.FilterSet{Operator: "AND", Items: []v3.FilterItem{
							{Key: v3.AttributeKey{Key: "user_name"}, Value: "john", Operator: "="},
						}},
					},
				},
			},
		},
		EnrichmentRequired: true,
	},
	{
		Name: "filter enrichment required required json",
		Params: v3.QueryRangeParamsV3{
			CompositeQuery: &v3.CompositeQuery{
				BuilderQueries: map[string]*v3.BuilderQuery{
					"test": {
						QueryName:  "test",
						Expression: "test",
						DataSource: v3.DataSourceLogs,
						Filters: &v3.FilterSet{Operator: "AND", Items: []v3.FilterItem{
							{Key: v3.AttributeKey{Key: "body.xyz", IsJSON: true, DataType: v3.AttributeKeyDataTypeString}, Value: "john", Operator: "="},
						}},
					},
				},
			},
		},
		EnrichmentRequired: true,
	},
	{
		Name: "groupBy enrichment required",
		Params: v3.QueryRangeParamsV3{
			CompositeQuery: &v3.CompositeQuery{
				BuilderQueries: map[string]*v3.BuilderQuery{
					"test": {
						QueryName:  "test",
						Expression: "test",
						DataSource: v3.DataSourceLogs,
						GroupBy:    []v3.AttributeKey{{Key: "userid", DataType: v3.AttributeKeyDataTypeString, Type: v3.AttributeKeyTypeTag}},
					},
				},
			},
		},
		EnrichmentRequired: true,
	},
	{
		Name: "groupBy enrichment required",
		Params: v3.QueryRangeParamsV3{
			CompositeQuery: &v3.CompositeQuery{
				BuilderQueries: map[string]*v3.BuilderQuery{
					"test": {
						QueryName:  "test",
						Expression: "test",
						DataSource: v3.DataSourceLogs,
						GroupBy:    []v3.AttributeKey{{Key: "userid"}},
					},
				},
			},
		},
		EnrichmentRequired: true,
	},
	{
		Name: "orderBy enrichment not required",
		Params: v3.QueryRangeParamsV3{
			CompositeQuery: &v3.CompositeQuery{
				BuilderQueries: map[string]*v3.BuilderQuery{
					"test": {
						QueryName:  "test",
						Expression: "test",
						DataSource: v3.DataSourceLogs,
						GroupBy:    []v3.AttributeKey{{Key: "userid", DataType: v3.AttributeKeyDataTypeString, Type: v3.AttributeKeyTypeTag}},
						OrderBy:    []v3.OrderBy{{ColumnName: "userid"}},
					},
				},
			},
		},
		EnrichmentRequired: true,
	},
	{
		Name: "orderBy enrichment required",
		Params: v3.QueryRangeParamsV3{
			CompositeQuery: &v3.CompositeQuery{
				BuilderQueries: map[string]*v3.BuilderQuery{
					"test": {
						QueryName:  "test",
						Expression: "test",
						DataSource: v3.DataSourceLogs,
						OrderBy:    []v3.OrderBy{{ColumnName: "userid"}},
					},
				},
			},
		},
		EnrichmentRequired: true,
	},
	{
		Name: "top level key",
		Params: v3.QueryRangeParamsV3{
			CompositeQuery: &v3.CompositeQuery{
				BuilderQueries: map[string]*v3.BuilderQuery{
					"test": {
						QueryName:  "test",
						Expression: "test",
						DataSource: v3.DataSourceLogs,
						GroupBy:    []v3.AttributeKey{{Key: "trace_id", Type: v3.AttributeKeyTypeUnspecified, DataType: v3.AttributeKeyDataTypeString, IsColumn: true}},
					},
				},
			},
		},
		EnrichmentRequired: false,
	},
	{
		Name: "orderBy enrichment required",
		Params: v3.QueryRangeParamsV3{
			CompositeQuery: &v3.CompositeQuery{
				BuilderQueries: map[string]*v3.BuilderQuery{
					"test": {
						QueryName:  "test",
						Expression: "test",
						DataSource: v3.DataSourceLogs,
						// here we have to fallback to trace_id attribute instead of column
						GroupBy: []v3.AttributeKey{{Key: "trace_id", DataType: v3.AttributeKeyDataTypeString, Type: v3.AttributeKeyTypeTag}},
						OrderBy: []v3.OrderBy{{ColumnName: "#SIGNOZ_VALUE", Order: "ASC"}},
					},
				},
			},
		},
		EnrichmentRequired: true,
	},
}

func TestEnrichmentRquired(t *testing.T) {
	for _, tt := range testEnrichmentRequiredData {
		Convey("testEnrichmentRequiredData", t, func() {
			res := EnrichmentRequired(&tt.Params)
			So(res, ShouldEqual, tt.EnrichmentRequired)
		})
	}
}

var testEnrichParamsData = []struct {
	Name   string
	Params v3.QueryRangeParamsV3
	Fields map[string]v3.AttributeKey
	Result v3.QueryRangeParamsV3
}{
	{
		Name: "Enriching query range v3 params",
		Params: v3.QueryRangeParamsV3{
			CompositeQuery: &v3.CompositeQuery{
				BuilderQueries: map[string]*v3.BuilderQuery{
					"test": {
						QueryName:  "test",
						Expression: "test",
						DataSource: v3.DataSourceLogs,
						AggregateAttribute: v3.AttributeKey{
							Key: "test",
						},
						Filters: &v3.FilterSet{Operator: "AND", Items: []v3.FilterItem{
							{Key: v3.AttributeKey{Key: "user_name"}, Value: "john", Operator: "="},
						}},
						GroupBy: []v3.AttributeKey{{Key: "trace_id"}},
						OrderBy: []v3.OrderBy{{ColumnName: "response_time"}},
					},
				},
			},
		},
		Fields: map[string]v3.AttributeKey{
			"test##tag##int64": {
				Key:      "test",
				Type:     v3.AttributeKeyTypeTag,
				DataType: v3.AttributeKeyDataTypeInt64,
			},
			"user_name##tag##string": {
				Key:      "user_name",
				Type:     v3.AttributeKeyTypeTag,
				DataType: v3.AttributeKeyDataTypeString,
			},
			"response_time##tag##int64": {
				Key:      "response_time",
				Type:     v3.AttributeKeyTypeTag,
				DataType: v3.AttributeKeyDataTypeInt64,
				IsColumn: true,
			},
		},
		Result: v3.QueryRangeParamsV3{
			CompositeQuery: &v3.CompositeQuery{
				BuilderQueries: map[string]*v3.BuilderQuery{
					"test": {
						QueryName:  "test",
						Expression: "test",
						DataSource: v3.DataSourceLogs,
						AggregateAttribute: v3.AttributeKey{
							Key:      "test",
							Type:     v3.AttributeKeyTypeTag,
							DataType: v3.AttributeKeyDataTypeInt64,
						},
						Filters: &v3.FilterSet{Operator: "AND", Items: []v3.FilterItem{
							{Key: v3.AttributeKey{Key: "user_name", Type: v3.AttributeKeyTypeTag, DataType: v3.AttributeKeyDataTypeString}, Value: "john", Operator: "="},
						}},
						GroupBy: []v3.AttributeKey{{Key: "trace_id", Type: v3.AttributeKeyTypeUnspecified, DataType: v3.AttributeKeyDataTypeString, IsColumn: true}},
						OrderBy: []v3.OrderBy{{ColumnName: "response_time", Key: "response_time", Type: v3.AttributeKeyTypeTag, DataType: v3.AttributeKeyDataTypeInt64, IsColumn: true}},
					},
				},
			},
		},
	},
	{
		Name: "Enriching query range v3 params with dot support",
		Params: v3.QueryRangeParamsV3{
			CompositeQuery: &v3.CompositeQuery{
				BuilderQueries: map[string]*v3.BuilderQuery{
					"test": {
						QueryName:  "test",
						Expression: "test",
						DataSource: v3.DataSourceLogs,
						AggregateAttribute: v3.AttributeKey{
							Key: "method.name",
						},
						Filters: &v3.FilterSet{Operator: "AND", Items: []v3.FilterItem{
							{Key: v3.AttributeKey{Key: "service.name"}, Value: "test", Operator: "="},
						}},
						GroupBy: []v3.AttributeKey{{Key: "host.name"}},
						OrderBy: []v3.OrderBy{{ColumnName: "host.name"}},
					},
				},
			},
		},
		Fields: map[string]v3.AttributeKey{
			"method.name##tag##string": {
				Key:      "method.name",
				Type:     v3.AttributeKeyTypeTag,
				DataType: v3.AttributeKeyDataTypeString,
				IsColumn: true,
			},
			"service.name##tag##string": {
				Key:      "service.name",
				Type:     v3.AttributeKeyTypeTag,
				DataType: v3.AttributeKeyDataTypeString,
			},
			"host.name##tag##string": {
				Key:      "host.name",
				Type:     v3.AttributeKeyTypeTag,
				DataType: v3.AttributeKeyDataTypeString,
			},
		},
		Result: v3.QueryRangeParamsV3{
			CompositeQuery: &v3.CompositeQuery{
				BuilderQueries: map[string]*v3.BuilderQuery{
					"test": {
						QueryName:  "test",
						Expression: "test",
						DataSource: v3.DataSourceLogs,
						AggregateAttribute: v3.AttributeKey{
							Key:      "method.name",
							Type:     v3.AttributeKeyTypeTag,
							DataType: v3.AttributeKeyDataTypeString,
							IsColumn: true,
						},
						Filters: &v3.FilterSet{Operator: "AND", Items: []v3.FilterItem{
							{Key: v3.AttributeKey{Key: "service.name", Type: v3.AttributeKeyTypeTag, DataType: v3.AttributeKeyDataTypeString}, Value: "test", Operator: "="},
						}},
						GroupBy: []v3.AttributeKey{{Key: "host.name", Type: v3.AttributeKeyTypeTag, DataType: v3.AttributeKeyDataTypeString}},
						OrderBy: []v3.OrderBy{{ColumnName: "host.name", Key: "host.name", Type: v3.AttributeKeyTypeTag, DataType: v3.AttributeKeyDataTypeString}},
					},
				},
			},
		},
	},
	{
		Name: "Don't enrich if other keys are non empty and not same",
		Params: v3.QueryRangeParamsV3{
			CompositeQuery: &v3.CompositeQuery{
				BuilderQueries: map[string]*v3.BuilderQuery{
					"test": {
						QueryName:  "test",
						Expression: "test",
						DataSource: v3.DataSourceLogs,
						AggregateAttribute: v3.AttributeKey{
							Key:      "test",
							Type:     v3.AttributeKeyTypeResource,
							DataType: v3.AttributeKeyDataTypeInt64,
						},
						Filters: &v3.FilterSet{Operator: "AND", Items: []v3.FilterItem{
							{Key: v3.AttributeKey{Key: "test", Type: v3.AttributeKeyTypeTag}, Value: "test", Operator: "="},
							{Key: v3.AttributeKey{Key: "test", DataType: v3.AttributeKeyDataTypeString}, Value: "test1", Operator: "="},
						}},
					},
				},
			},
		},
		Fields: map[string]v3.AttributeKey{
			"test##tag##string": {
				Key:      "test",
				Type:     v3.AttributeKeyTypeTag,
				DataType: v3.AttributeKeyDataTypeString,
				IsColumn: true,
			},
		},
		Result: v3.QueryRangeParamsV3{
			CompositeQuery: &v3.CompositeQuery{
				BuilderQueries: map[string]*v3.BuilderQuery{
					"test": {
						QueryName:  "test",
						Expression: "test",
						DataSource: v3.DataSourceLogs,
						AggregateAttribute: v3.AttributeKey{
							Key:      "test",
							Type:     v3.AttributeKeyTypeResource,
							DataType: v3.AttributeKeyDataTypeInt64,
						},
						Filters: &v3.FilterSet{Operator: "AND", Items: []v3.FilterItem{
							{Key: v3.AttributeKey{Key: "test", Type: v3.AttributeKeyTypeTag, DataType: v3.AttributeKeyDataTypeString, IsColumn: true}, Value: "test", Operator: "="},
							{Key: v3.AttributeKey{Key: "test", Type: v3.AttributeKeyTypeTag, DataType: v3.AttributeKeyDataTypeString, IsColumn: true}, Value: "test1", Operator: "="},
						}},
					},
				},
			},
		},
	},
	{
		Name: "Enrich if an attribute/resource attribute is materialized/dematerialized",
		Params: v3.QueryRangeParamsV3{
			CompositeQuery: &v3.CompositeQuery{
				BuilderQueries: map[string]*v3.BuilderQuery{
					"test": {
						QueryName:  "test",
						Expression: "test",
						DataSource: v3.DataSourceLogs,
						AggregateAttribute: v3.AttributeKey{
							Key:      "mat_resource",
							Type:     v3.AttributeKeyTypeResource,
							DataType: v3.AttributeKeyDataTypeInt64,
							IsColumn: true,
						},
						Filters: &v3.FilterSet{Operator: "AND", Items: []v3.FilterItem{
							{Key: v3.AttributeKey{Key: "mat_attr", Type: v3.AttributeKeyTypeTag, IsColumn: true}, Value: "test", Operator: "="},
							{Key: v3.AttributeKey{Key: "normal_attr", DataType: v3.AttributeKeyDataTypeString, IsColumn: false}, Value: "test1", Operator: "="},
						}},
					},
				},
			},
		},
		Fields: map[string]v3.AttributeKey{
			"mat_resource##resource##int64": {
				Key:      "mat_resource",
				Type:     v3.AttributeKeyTypeResource,
				DataType: v3.AttributeKeyDataTypeInt64,
				IsColumn: false,
			},
			"mat_attr##tag##string": {
				Key:      "mat_attr",
				Type:     v3.AttributeKeyTypeTag,
				DataType: v3.AttributeKeyDataTypeString,
				IsColumn: false,
			},
			"normal_attr##tag##string": {
				Key:      "normal_attr",
				Type:     v3.AttributeKeyTypeTag,
				DataType: v3.AttributeKeyDataTypeString,
				IsColumn: true,
			},
		},
		Result: v3.QueryRangeParamsV3{
			CompositeQuery: &v3.CompositeQuery{
				BuilderQueries: map[string]*v3.BuilderQuery{
					"test": {
						QueryName:  "test",
						Expression: "test",
						DataSource: v3.DataSourceLogs,
						AggregateAttribute: v3.AttributeKey{
							Key:      "mat_resource",
							Type:     v3.AttributeKeyTypeResource,
							DataType: v3.AttributeKeyDataTypeInt64,
							IsColumn: false,
						},
						Filters: &v3.FilterSet{Operator: "AND", Items: []v3.FilterItem{
							{Key: v3.AttributeKey{Key: "mat_attr", Type: v3.AttributeKeyTypeTag, DataType: v3.AttributeKeyDataTypeString, IsColumn: false}, Value: "test", Operator: "="},
							{Key: v3.AttributeKey{Key: "normal_attr", Type: v3.AttributeKeyTypeTag, DataType: v3.AttributeKeyDataTypeString, IsColumn: true}, Value: "test1", Operator: "="},
						}},
					},
				},
			},
		},
	},
}

func TestEnrichParams(t *testing.T) {
	for _, tt := range testEnrichParamsData {
		Convey("testEnrichmentRequiredData", t, func() {
			Enrich(&tt.Params, tt.Fields)
			So(tt.Params, ShouldResemble, tt.Result)
		})
	}
}

var testJSONFilterEnrichData = []struct {
	Name   string
	Filter v3.FilterItem
	Result v3.FilterItem
}{
	{
		Name: "array string",
		Filter: v3.FilterItem{
			Key: v3.AttributeKey{
				Key:      "body.requestor_list[*]",
				DataType: v3.AttributeKeyDataTypeUnspecified,
				Type:     v3.AttributeKeyTypeUnspecified,
			},
			Operator: "has",
			Value:    "index_service",
		},
		Result: v3.FilterItem{
			Key: v3.AttributeKey{
				Key:      "body.requestor_list[*]",
				DataType: v3.AttributeKeyDataTypeArrayString,
				Type:     v3.AttributeKeyTypeUnspecified,
				IsJSON:   true,
			},
			Operator: "has",
			Value:    "index_service",
		},
	},
	{
		Name: "int64",
		Filter: v3.FilterItem{
			Key: v3.AttributeKey{
				Key:      "body.intx",
				DataType: v3.AttributeKeyDataTypeUnspecified,
				Type:     v3.AttributeKeyTypeUnspecified,
			},
			Operator: "=",
			Value:    10,
		},
		Result: v3.FilterItem{
			Key: v3.AttributeKey{
				Key:      "body.intx",
				DataType: v3.AttributeKeyDataTypeInt64,
				Type:     v3.AttributeKeyTypeUnspecified,
				IsJSON:   true,
			},
			Operator: "=",
			Value:    10,
		},
	},
	{
		Name: "int64 string",
		Filter: v3.FilterItem{
			Key: v3.AttributeKey{
				Key:      "body.intx",
				DataType: v3.AttributeKeyDataTypeUnspecified,
				Type:     v3.AttributeKeyTypeUnspecified,
			},
			Operator: "=",
			Value:    "0",
		},
		Result: v3.FilterItem{
			Key: v3.AttributeKey{
				Key:      "body.intx",
				DataType: v3.AttributeKeyDataTypeInt64,
				Type:     v3.AttributeKeyTypeUnspecified,
				IsJSON:   true,
			},
			Operator: "=",
			Value:    int64(0),
		},
	},
	{
		Name: "float64",
		Filter: v3.FilterItem{
			Key: v3.AttributeKey{
				Key:      "body.float64[*]",
				DataType: v3.AttributeKeyDataTypeArrayFloat64,
				Type:     v3.AttributeKeyTypeUnspecified,
			},
			Operator: "!=",
			Value:    10.0,
		},
		Result: v3.FilterItem{
			Key: v3.AttributeKey{
				Key:      "body.float64[*]",
				DataType: v3.AttributeKeyDataTypeArrayFloat64,
				Type:     v3.AttributeKeyTypeUnspecified,
				IsJSON:   true,
			},
			Operator: "!=",
			Value:    10.0,
		},
	},
	{
		Name: "float64x",
		Filter: v3.FilterItem{
			Key: v3.AttributeKey{
				Key:      "body.float64x",
				DataType: v3.AttributeKeyDataTypeUnspecified,
				Type:     v3.AttributeKeyTypeUnspecified,
			},
			Operator: "!=",
			Value:    "10.0",
		},
		Result: v3.FilterItem{
			Key: v3.AttributeKey{
				Key:      "body.float64x",
				DataType: v3.AttributeKeyDataTypeFloat64,
				Type:     v3.AttributeKeyTypeUnspecified,
				IsJSON:   true,
			},
			Operator: "!=",
			Value:    10.0,
		},
	},
	{
		Name: "check IN",
		Filter: v3.FilterItem{
			Key: v3.AttributeKey{
				Key:      "body.attr",
				DataType: v3.AttributeKeyDataTypeUnspecified,
				Type:     v3.AttributeKeyTypeUnspecified,
			},
			Operator: "IN",
			Value:    []interface{}{"hello", "world"},
		},
		Result: v3.FilterItem{
			Key: v3.AttributeKey{
				Key:      "body.attr",
				DataType: v3.AttributeKeyDataTypeString,
				Type:     v3.AttributeKeyTypeUnspecified,
				IsJSON:   true,
			},
			Operator: "IN",
			Value:    []interface{}{"hello", "world"},
		},
	},
	{
		Name: "check NOT_IN",
		Filter: v3.FilterItem{
			Key: v3.AttributeKey{
				Key:      "body.attr",
				DataType: v3.AttributeKeyDataTypeUnspecified,
				Type:     v3.AttributeKeyTypeUnspecified,
			},
			Operator: "NOT_IN",
			Value:    []interface{}{10, 20},
		},
		Result: v3.FilterItem{
			Key: v3.AttributeKey{
				Key:      "body.attr",
				DataType: v3.AttributeKeyDataTypeInt64,
				Type:     v3.AttributeKeyTypeUnspecified,
				IsJSON:   true,
			},
			Operator: "NOT_IN",
			Value:    []interface{}{10, 20},
		},
	},
}

func TestJsonEnrich(t *testing.T) {
	for _, tt := range testJSONFilterEnrichData {
		Convey(tt.Name, t, func() {
			res := jsonFilterEnrich(tt.Filter)
			So(res, ShouldResemble, tt.Result)
		})
	}
}

func TestJsonReplaceField(t *testing.T) {
	fields := map[string]v3.AttributeKey{
		"method.name": {
			Key:      "method.name",
			DataType: v3.AttributeKeyDataTypeString,
			Type:     v3.AttributeKeyTypeTag,
		},
		"status": {
			Key:      "status",
			DataType: v3.AttributeKeyDataTypeInt64,
			Type:     v3.AttributeKeyTypeTag,
		},
		"data.error": {
			Key:      "data.error",
			DataType: v3.AttributeKeyDataTypeString,
			Type:     v3.AttributeKeyTypeTag,
			IsColumn: true,
		},
	}
	var TestJsonReplaceFieldData = []struct {
		Name   string
		Filter v3.FilterItem
		Result v3.FilterItem
	}{
		{
			Name: "key in nested json",
			Filter: v3.FilterItem{
				Key: v3.AttributeKey{
					Key:      "body.method.name",
					DataType: v3.AttributeKeyDataTypeString,
					Type:     v3.AttributeKeyTypeUnspecified,
				},
				Operator: "has",
				Value:    "index_service",
			},
			Result: v3.FilterItem{
				Key: v3.AttributeKey{
					Key:      "method.name",
					DataType: v3.AttributeKeyDataTypeString,
					Type:     v3.AttributeKeyTypeTag,
					IsJSON:   false,
				},
				Operator: "has",
				Value:    "index_service",
			},
		},
		{
			Name: "key at top level",
			Filter: v3.FilterItem{
				Key: v3.AttributeKey{
					Key:      "body.status",
					DataType: v3.AttributeKeyDataTypeInt64,
					Type:     v3.AttributeKeyTypeUnspecified,
				},
				Operator: "=",
				Value:    10,
			},
			Result: v3.FilterItem{
				Key: v3.AttributeKey{
					Key:      "status",
					DataType: v3.AttributeKeyDataTypeInt64,
					Type:     v3.AttributeKeyTypeTag,
					IsJSON:   false,
				},
				Operator: "=",
				Value:    10,
			},
		},
		{
			Name: "key not present",
			Filter: v3.FilterItem{
				Key: v3.AttributeKey{
					Key:      "body.status.code",
					DataType: v3.AttributeKeyDataTypeInt64,
					Type:     v3.AttributeKeyTypeUnspecified,
				},
				Operator: "=",
				Value:    10,
			},
			Result: v3.FilterItem{
				Key: v3.AttributeKey{
					Key:      "body.status.code",
					DataType: v3.AttributeKeyDataTypeInt64,
					Type:     v3.AttributeKeyTypeUnspecified,
					IsJSON:   false,
				},
				Operator: "=",
				Value:    10,
			},
		},
		{
			Name: "key materialized",
			Filter: v3.FilterItem{
				Key: v3.AttributeKey{
					Key:      "body.data.error",
					DataType: v3.AttributeKeyDataTypeString,
					Type:     v3.AttributeKeyTypeUnspecified,
				},
				Operator: "=",
				Value:    10,
			},
			Result: v3.FilterItem{
				Key: v3.AttributeKey{
					Key:      "data.error",
					DataType: v3.AttributeKeyDataTypeString,
					Type:     v3.AttributeKeyTypeTag,
					IsJSON:   false,
					IsColumn: true,
				},
				Operator: "=",
				Value:    10,
			},
		},
	}
	for _, tt := range TestJsonReplaceFieldData {
		Convey(tt.Name, t, func() {
			res := jsonReplaceField(tt.Filter, fields)
			So(res, ShouldResemble, tt.Result)
		})
	}
}

var testParseStrValueData = []struct {
	Name       string
	Operator   v3.FilterOperator
	Value      interface{}
	ResultType string
	Result     interface{}
}{
	{
		Name:       "bool",
		Value:      "true",
		Operator:   v3.FilterOperatorEqual,
		ResultType: "bool",
		Result:     true,
	},
	{
		Name:       "int",
		Value:      "10",
		Operator:   v3.FilterOperatorNotEqual,
		ResultType: "int64",
		Result:     10,
	},
	{
		Name:       "float",
		Value:      "10.0",
		Operator:   v3.FilterOperatorGreaterThan,
		ResultType: "float64",
		Result:     10.0,
	},
	{
		Name:       "string",
		Value:      "hello",
		Operator:   v3.FilterOperatorLessThan,
		ResultType: "string",
		Result:     "hello",
	},
}

func TestParseStrValue(t *testing.T) {
	for _, tt := range testParseStrValueData {
		Convey(tt.Name, t, func() {
			vtype, value := parseStrValue(tt.Value.(string), tt.Operator)
			So(vtype, ShouldEqual, tt.ResultType)
			So(value, ShouldEqual, tt.Result)
		})
	}
}

func TestGetFieldNames(t *testing.T) {
	type args struct {
		compositeQuery *v3.CompositeQuery
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "nil composite query",
			args: args{
				compositeQuery: nil,
			},
			want: []string{},
		},
		{
			name: "empty builder queries",
			args: args{
				compositeQuery: &v3.CompositeQuery{
					BuilderQueries: map[string]*v3.BuilderQuery{},
				},
			},
			want: []string{},
		},
		{
			name: "aggregate attribute only",
			args: args{
				compositeQuery: &v3.CompositeQuery{
					BuilderQueries: map[string]*v3.BuilderQuery{
						"test": {
							QueryName:  "test",
							Expression: "test",
							DataSource: v3.DataSourceLogs,
							AggregateAttribute: v3.AttributeKey{
								Key: "response_time",
							},
						},
					},
				},
			},
			want: []string{"response_time"},
		},
		{
			name: "filters only",
			args: args{
				compositeQuery: &v3.CompositeQuery{
					BuilderQueries: map[string]*v3.BuilderQuery{
						"test": {
							QueryName:  "test",
							Expression: "test",
							DataSource: v3.DataSourceLogs,
							Filters: &v3.FilterSet{
								Operator: "AND",
								Items: []v3.FilterItem{
									{Key: v3.AttributeKey{Key: "user_name"}, Value: "john", Operator: "="},
									{Key: v3.AttributeKey{Key: "status"}, Value: "error", Operator: "!="},
								},
							},
						},
					},
				},
			},
			want: []string{"user_name", "status"},
		},
		{
			name: "JSON filters",
			args: args{
				compositeQuery: &v3.CompositeQuery{
					BuilderQueries: map[string]*v3.BuilderQuery{
						"test": {
							QueryName:  "test",
							Expression: "test",
							DataSource: v3.DataSourceLogs,
							Filters: &v3.FilterSet{
								Operator: "AND",
								Items: []v3.FilterItem{
									{Key: v3.AttributeKey{Key: "body.user_name"}, Value: "john", Operator: "="},
									{Key: v3.AttributeKey{Key: "body.status.code"}, Value: 200, Operator: "="},
									{Key: v3.AttributeKey{Key: "body.requestor_list[*]"}, Value: "index_service", Operator: "has"},
								},
							},
						},
					},
				},
			},
			// requestor_list[*] is not required though
			want: []string{"user_name", "status.code", "requestor_list[*]"},
		},
		{
			name: "groupBy only",
			args: args{
				compositeQuery: &v3.CompositeQuery{
					BuilderQueries: map[string]*v3.BuilderQuery{
						"test": {
							QueryName:  "test",
							Expression: "test",
							DataSource: v3.DataSourceLogs,
							GroupBy: []v3.AttributeKey{
								{Key: "service_name"},
								{Key: "host_name"},
							},
						},
					},
				},
			},
			want: []string{"service_name", "host_name"},
		},
		{
			name: "orderBy only",
			args: args{
				compositeQuery: &v3.CompositeQuery{
					BuilderQueries: map[string]*v3.BuilderQuery{
						"test": {
							QueryName:  "test",
							Expression: "test",
							DataSource: v3.DataSourceLogs,
							OrderBy: []v3.OrderBy{
								{ColumnName: "timestamp", Order: "DESC"},
								{ColumnName: "response_time", Order: "ASC"},
							},
						},
					},
				},
			},
			want: []string{"timestamp", "response_time"},
		},
		{
			name: "all field types combined",
			args: args{
				compositeQuery: &v3.CompositeQuery{
					BuilderQueries: map[string]*v3.BuilderQuery{
						"test": {
							QueryName:  "test",
							Expression: "test",
							DataSource: v3.DataSourceLogs,
							AggregateAttribute: v3.AttributeKey{
								Key: "response_time",
							},
							Filters: &v3.FilterSet{
								Operator: "AND",
								Items: []v3.FilterItem{
									{Key: v3.AttributeKey{Key: "user_name"}, Value: "john", Operator: "="},
									{Key: v3.AttributeKey{Key: "body.status"}, Value: "error", Operator: "!="},
								},
							},
							GroupBy: []v3.AttributeKey{
								{Key: "service_name"},
								{Key: "user_name"},
							},
							OrderBy: []v3.OrderBy{
								{ColumnName: "timestamp", Order: "DESC"},
							},
						},
					},
				},
			},
			want: []string{"response_time", "user_name", "status", "service_name", "timestamp"},
		},
		{
			name: "multiple builder queries",
			args: args{
				compositeQuery: &v3.CompositeQuery{
					BuilderQueries: map[string]*v3.BuilderQuery{
						"query1": {
							QueryName:  "query1",
							Expression: "query1",
							DataSource: v3.DataSourceLogs,
							AggregateAttribute: v3.AttributeKey{
								Key: "field1",
							},
						},
						"query2": {
							QueryName:  "query2",
							Expression: "query2",
							DataSource: v3.DataSourceLogs,
							Filters: &v3.FilterSet{
								Operator: "AND",
								Items: []v3.FilterItem{
									{Key: v3.AttributeKey{Key: "field2"}, Value: "value", Operator: "="},
								},
							},
						},
						"query3": {
							QueryName:  "query3",
							Expression: "different_expression",
							DataSource: v3.DataSourceLogs,
							AggregateAttribute: v3.AttributeKey{
								Key: "field3",
							},
						},
					},
				},
			},
			want: []string{"field1", "field2"},
		},
		{
			name: "empty aggregate attribute key and filters",
			args: args{
				compositeQuery: &v3.CompositeQuery{
					BuilderQueries: map[string]*v3.BuilderQuery{
						"test": {
							QueryName:  "test",
							Expression: "test",
							DataSource: v3.DataSourceLogs,
							AggregateAttribute: v3.AttributeKey{
								Key: "",
							},
							Filters: &v3.FilterSet{
								Operator: "AND",
								Items: []v3.FilterItem{
									{Key: v3.AttributeKey{Key: "user_name"}, Value: "john", Operator: "="},
								},
							},
							GroupBy: []v3.AttributeKey{
								{Key: "service_name"},
							},
						},
					},
				},
			},
			want: []string{"user_name", "service_name"},
		},
		{
			name: "JSON filter with array suffix - no use",
			args: args{
				compositeQuery: &v3.CompositeQuery{
					BuilderQueries: map[string]*v3.BuilderQuery{
						"test": {
							QueryName:  "test",
							Expression: "test",
							DataSource: v3.DataSourceLogs,
							Filters: &v3.FilterSet{
								Operator: "AND",
								Items: []v3.FilterItem{
									{Key: v3.AttributeKey{Key: "body.tags[*]"}, Value: "important", Operator: "has"},
									{Key: v3.AttributeKey{Key: "body.nested.array[*]"}, Value: "value", Operator: "has"},
								},
							},
						},
					},
				},
			},
			want: []string{"tags[*]", "nested.array[*]"},
		},
		{
			name: "mixed JSON and regular filters",
			args: args{
				compositeQuery: &v3.CompositeQuery{
					BuilderQueries: map[string]*v3.BuilderQuery{
						"test": {
							QueryName:  "test",
							Expression: "test",
							DataSource: v3.DataSourceLogs,
							Filters: &v3.FilterSet{
								Operator: "AND",
								Items: []v3.FilterItem{
									{Key: v3.AttributeKey{Key: "user_name"}, Value: "john", Operator: "="},
									{Key: v3.AttributeKey{Key: "body.status.abc"}, Value: "error", Operator: "!="},
									{Key: v3.AttributeKey{Key: "service_name"}, Value: "api", Operator: "="},
								},
							},
						},
					},
				},
			},
			want: []string{"user_name", "status.abc", "service_name"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetFieldNames(tt.args.compositeQuery)
			// Sort both slices since order doesn't matter
			sort.Strings(got)
			sort.Strings(tt.want)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetFieldNames() = %v, want %v", got, tt.want)
			}
		})
	}
}
