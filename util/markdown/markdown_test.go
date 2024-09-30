package markdown_test

import (
	"fmt"
	"github.com/IceFoxs/open-gateway/util/markdown"
	"testing"
)

func TestMarkdown(t *testing.T) {
	s := markdown.NewMarkdown(markdown.WithLi([]string{"1", "2", "3"}), markdown.WithText(markdown.WrapColor(markdown.CONTENT_COLOR, "1231231321"))).Builder()
	fmt.Println(s)
	s = markdown.NewMarkdown(markdown.WithTitle3(markdown.WrapFontColor(markdown.NOTICE_COLOR, "网关服务注册重复文件名")), markdown.WithText(markdown.WrapFontColor(markdown.ERROR_COLOR, "方法名：FPS_SSSS_SSS重复出现，具体应用如下："))).Builder()
	fmt.Println(s)

}
