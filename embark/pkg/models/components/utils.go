package components

func copyTree(tree map[string]interface{}) map[string]interface{} {
	var cp = make(map[string]interface{})
	for k, v := range tree {
		switch v := v.(type) {
		case nil:
			continue
		case map[string]interface{}:
			cp[k] = copyTree(v)
		case []string:
			cp[k] = append([]string{}, v...)
		case []int:
			cp[k] = append([]int{}, v...)
		case []interface{}:
			cp[k] = append([]interface{}{}, v...)
		case Component:
			cp[k] = v.Copy()
		default:
			cp[k] = v
		}
	}
	return cp
}
