package util

import (
	"fmt"
	"github.com/IceFoxs/open-gateway/common"
	hessian "github.com/apache/dubbo-go-hessian2"
	"github.com/cloudwego/hertz/pkg/common/json"
	"log"
)

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

func ConvertHessianMap(m map[string]interface{}) map[string]hessian.Object {
	hmap := make(map[string]hessian.Object)
	for k, v := range m {
		hmap[k] = hessian.Object(v)
	}
	return hmap
}

func JsonStringToMap(jsonStr string) (map[string]interface{}, error) {
	var data map[string]interface{}
	err := json.Unmarshal([]byte(jsonStr), &data)
	if err != nil {
		return nil, err
	}
	fmt.Printf("JsonStringToMap:%s", common.ToJSON(data))
	return data, nil
}
