package pbupdate

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

func CopyValuesFromPaths[S any, D any](paths []string, source S, dest D) (*D, error) {
	srcMap, err := objectToMap(source)
	if err != nil {
		return nil, err
	}

	destMap, err := objectToMap(dest)
	if err != nil {
		return nil, err
	}

	for _, path := range paths {
		err = copyFromPath(path, srcMap, destMap)
		if err != nil {
			return nil, err
		}
	}

	destStr, err := json.Marshal(destMap)
	if err != nil {
		return nil, err
	}

	var result D
	err = json.Unmarshal(destStr, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func copyFromPath(path string, src map[string]any, dest map[string]any) error {
	switch {
	case src == nil:
		return errors.New("source object was nil, please pass in a valid object")
	case dest == nil:
		return errors.New("dest object was nil, please pass in a valid object")
	case len(path) == 0:
		return errors.New("path has a length of zero, please pass in a valid path")
	}

	srcValue, err := readJsonPath(path, src)
	if err != nil {
		return err
	}

	return updateJsonValue(path, dest, srcValue)
}

func readJsonPath(path string, object map[string]any) (any, error) {
	if len(path) == 0 {
		return nil, errors.New("readJsonPath: received an empty string for the path")
	}

	key, remainingKeys := splitDotPath(path)

	value, ok := object[key]
	if !ok {
		return nil, errors.New(fmt.Sprintf("readJsonPath: could not find key %s in object %v", key, object))
	}

	// If there are more paths to go down
	if len(remainingKeys) > 0 {
		if next, ok := value.(map[string]any); ok {
			return readJsonPath(remainingKeys, next)
		}

		return nil, errors.New(fmt.Sprintf("readJsonPath: value %v was not assignable to map[string]interface{} at path %s", value, path))
	}

	return value, nil
}

func updateJsonValue(path string, object map[string]any, newValue any) error {
	if len(path) == 0 {
		return errors.New("updateJsonValue: received an empty string for the path")
	}

	key, remainingKeys := splitDotPath(path)

	value, ok := object[key]
	if !ok {
		return errors.New(fmt.Sprintf("updateJsonValue: could not find key %s in object %v", key, object))
	}

	// If there are more paths to go down
	if len(remainingKeys) > 0 {
		if next, ok := value.(map[string]any); ok {
			return updateJsonValue(remainingKeys, next, newValue)
		}

		return errors.New(fmt.Sprintf("updateJsonValue: value %v was not assignable to map[string]interface{} at path %s", value, path))
	}

	object[key] = newValue
	return nil
}

func splitDotPath(path string) (string, string) {
	split := strings.Split(path, ".")
	if len(split) == 0 {
		return "", ""
	}

	return split[0], strings.Join(split[1:], ".")
}

func objectToMap(object any) (map[string]any, error) {
	jsonStr, err := json.Marshal(object)
	if err != nil {
		return nil, err
	}

	var jsonMap map[string]any
	err = json.Unmarshal(jsonStr, &jsonMap)
	if err != nil {
		return nil, err
	}

	return jsonMap, nil
}
