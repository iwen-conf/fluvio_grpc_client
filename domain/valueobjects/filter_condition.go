package valueobjects

// FilterType 过滤类型
type FilterType string

const (
	FilterTypeKey     FilterType = "key"
	FilterTypeValue   FilterType = "value"
	FilterTypeHeader  FilterType = "header"
	FilterTypeOffset  FilterType = "offset"
)

// FilterOperator 过滤操作符
type FilterOperator string

const (
	FilterOperatorEq       FilterOperator = "eq"       // 等于
	FilterOperatorNe       FilterOperator = "ne"       // 不等于
	FilterOperatorGt       FilterOperator = "gt"       // 大于
	FilterOperatorGte      FilterOperator = "gte"      // 大于等于
	FilterOperatorLt       FilterOperator = "lt"       // 小于
	FilterOperatorLte      FilterOperator = "lte"      // 小于等于
	FilterOperatorContains FilterOperator = "contains" // 包含
	FilterOperatorRegex    FilterOperator = "regex"    // 正则表达式
)

// FilterCondition 过滤条件值对象
type FilterCondition struct {
	Type     FilterType
	Field    string         // 字段名（用于header类型）
	Operator FilterOperator
	Value    string
}

// NewFilterCondition 创建新的过滤条件
func NewFilterCondition(filterType FilterType, operator FilterOperator, value string) *FilterCondition {
	return &FilterCondition{
		Type:     filterType,
		Operator: operator,
		Value:    value,
	}
}

// NewHeaderFilter 创建头部过滤条件
func NewHeaderFilter(field string, operator FilterOperator, value string) *FilterCondition {
	return &FilterCondition{
		Type:     FilterTypeHeader,
		Field:    field,
		Operator: operator,
		Value:    value,
	}
}

// IsValid 验证过滤条件是否有效
func (fc *FilterCondition) IsValid() bool {
	if fc.Type == "" || fc.Operator == "" || fc.Value == "" {
		return false
	}
	
	// 头部过滤必须有字段名
	if fc.Type == FilterTypeHeader && fc.Field == "" {
		return false
	}
	
	return true
}

// String 返回过滤条件的字符串表示
func (fc *FilterCondition) String() string {
	if fc.Type == FilterTypeHeader {
		return string(fc.Type) + "." + fc.Field + " " + string(fc.Operator) + " " + fc.Value
	}
	return string(fc.Type) + " " + string(fc.Operator) + " " + fc.Value
}