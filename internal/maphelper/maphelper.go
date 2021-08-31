package maphelper

import (
	"fmt"
	"math"
)

func GetStringValue(data map[string]interface{}, key string) (string, bool) {
	if val, ok := data[key]; ok {
		if sVal, ok := val.(string); ok {
			return sVal, ok
		} else if fVal, ok := val.(float64); ok {
			return fmt.Sprintf("%f", fVal), ok
		} else if iVal, ok := val.(int); ok {
			return fmt.Sprintf("%d", iVal), ok
		} else if bVal, ok := val.(bool); ok {
			if bVal {
				return "true", ok
			} else {
				return "false", ok
			}
		}
	}

	return "", false
}

func GetBoolValue(data map[string]interface{}, key string) (bool, bool) {
	if val, ok := data[key]; ok {
		if bVal, ok := val.(bool); ok {
			return bVal, ok
		} else if iVal, ok := val.(int); ok {
			return iVal != 0, ok
		} else if fVal, ok := val.(float64); ok {
			return fVal != 0.00, ok
		}
	}

	return false, false
}

func GetFloatValue(data map[string]interface{}, key string) (float64, bool) {
	if val, ok := data[key]; ok {
		if fVal, ok := val.(float64); ok {
			return fVal, ok
		} else if iVal, ok := val.(int); ok {
			return float64(iVal), ok
		}
	}

	return 0.00, false
}

func GetIntValue(data map[string]interface{}, key string) (int, bool) {
	if val, ok := data[key]; ok {
		if iVal, ok := val.(int); ok {
			return iVal, ok
		} else if fVal, ok := val.(float64); ok {
			return int(math.Round(fVal)), ok
		}
	}

	return 0, false
}
