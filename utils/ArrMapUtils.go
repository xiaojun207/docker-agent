package utils

func MapInterfaceToString(m map[string]interface{}) map[string]string {
	res := map[string]string{}
	for s, i := range m {
		res[s] = i.(string)
	}
	return res
}

func ArrInterfaceToStr(arr []interface{}) []string {
	res := []string{}
	for _, a := range arr {
		res = append(res, a.(string))
	}
	return res
}
