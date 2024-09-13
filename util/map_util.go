package util

import "log"

func ConvertMap(res interface{}) (tmp map[string]interface{}) {
	// map 需要初始化一个出来
	tmp = make(map[string]interface{})
	log.Println("input res is : ", res)
	switch res.(type) {
	case nil:
		return tmp
	case map[string]interface{}:
		return res.(map[string]interface{})
	case map[interface{}]interface{}:
		log.Println("map[interface{}]interface{} res:", res)
		for k, v := range res.(map[interface{}]interface{}) {
			log.Println("loop:", k, v)
			switch k.(type) {
			case string:
				switch v.(type) {
				case map[interface{}]interface{}:
					log.Println("map[interface{}]interface{} v:", v)
					tmp[k.(string)] = ConvertMap(v)
					continue
				default:
					log.Printf("default v: %v %v \n", k, v)
					tmp[k.(string)] = v
				}

			default:
				continue
			}
		}
		return tmp
	default:
		// 暂时没遇到更复杂的数据
		log.Println("unknow data:", res)
	}
	return tmp
}
