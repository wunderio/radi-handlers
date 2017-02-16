package configwrapper

import (
	"strconv"
	"strings"

	log "github.com/Sirupsen/logrus"

	api_property "github.com/wunderkraut/radi-api/property"
)

// Convert Authorize string to int, invert it if asked to
func authStringToInt(authValue string, negate bool) int {
	negateMultiplier := 1
	if negate {
		negateMultiplier = -1
	}

	switch strings.ToLower(authValue) {
	case "allow":
		return 1 * negateMultiplier
	case "deny":
		return -1 * negateMultiplier
	default:
		return 0
	}
}

// Does a property match some string values
func authMatchProperty(prop api_property.Property, propStrings []string) bool {
	switch prop.Type() {
	case "bool":
		if len(propStrings) == 0 {
			propStrings = []string{"true"}
		}
		propValue := prop.Get().(bool)
		for _, stringValue := range propStrings {
			value, _ := ymlTool_Convert_ToBool(stringValue)
			if value == propValue {
				return true
			}
		}
	case "string":
		propValue := prop.Get().(string)
		for _, stringValue := range propStrings {
			if stringValue == propValue {
				return true
			}
		}
	case "stringslice":
		propValues := prop.Get().([]string)
		for _, stringValue := range propStrings {
			for _, propValue := range propValues {
				if stringValue == propValue {
					return true
				}
			}
		}
	case "int":
		propValue := prop.Get().(int64)
		for _, stringValue := range propStrings {
			value, _ := ymlTool_Convert_ToInt(stringValue)
			if value == propValue {
				return true
			}
		}
	default:
		log.WithFields(log.Fields{"Property": prop.Id(), "type": prop.Type()}).Error("Could not perform an authorization match on property as the auth type matching algorythm has not been written")
	}
	return false
}

/**
 * Some usefull tools to use when dealing with
 * data from YML sources.
 */

// Convert a YML string value to a bool
func ymlTool_Convert_ToBool(value string) (bool, error) {
	return strconv.ParseBool(value)
}

// Convert a YML string value to a bool
func ymlTool_Convert_ToInt(value string) (int64, error) {
	return strconv.ParseInt(value, 10, 64)
}
