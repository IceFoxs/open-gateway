package model_test

import (
	"fmt"
	"github.com/IceFoxs/open-gateway/common"
	"github.com/IceFoxs/open-gateway/model"
	"github.com/cloudwego/hertz/pkg/common/json"
	"testing"
)

func TestAppMetadata(t *testing.T) {
	s := "{\n    \"appName\": \"fps-manager\",\n    \"methods\": [\n        \"FPS_MANAGER_GETSERVICECHANNELNONAMEPAIR\",\n        \"FPS_MANAGER_GETBANKCHANNELIDNAMEPAIR\",\n        \"FPS_MANAGER_GETBANKSERVICEPAYAPPLY\",\n        \"FPS_MANAGER_GENERATERSAKEYS\",\n        \"FPS_MANAGER_GETBANKREFUNDSERVICEAPPLYBYID\",\n        \"FPS_MANAGER_GETSCIDNAMEPAIR\",\n        \"FPS_MANAGER_GETTRADECHANNELCONFIG\",\n        \"FPS_MANAGER_EDITTRADECHANNELCONFIG\",\n        \"FPS_MANAGER_GETBANKPAYAPPLY\",\n        \"FPS_MANAGER_GETBANKSERVICEPAYAPPLYBYID\",\n        \"FPS_MANAGER_GETGATEWAYCHANNELCONFIG\",\n        \"FPS_MANAGER_DELETEBANKCHANNELCONFIG\",\n        \"FPS_MANAGER_DELETEGATEWAYCHANNELCONFIG\",\n        \"FPS_MANAGER_GETBANKREFUNDAPPLY\",\n        \"FPS_MANAGER_GETACCOUNTRECONRECORD\",\n        \"FPS_MANAGER_GETREQUESTREFUNDSERVICETRADE\",\n        \"FPS_MANAGER_GETSERVICEFACTORINGNONAMEPAIR\",\n        \"FPS_MANAGER_GETREQUESTFUNDTRADE\",\n        \"FPS_MANAGER_GETACCOUNTFUNDSRECEIPT\",\n        \"FPS_MANAGER_EDITBANKCHANNELCONFIG\",\n        \"FPS_MANAGER_ADDSERVICECHANNELCONFIG\",\n        \"FPS_MANAGER_GETCLEARINGFACTORINGFUNDS\",\n        \"FPS_MANAGER_ADDGATEWAYCHANNELCONFIG\",\n        \"FPS_MANAGER_GETBANKREFUNDDETAIL\",\n        \"FPS_MANAGER_GETRECEIVABLEBILLRECONRECORD\",\n        \"FPS_MANAGER_ADDBANKCHANNELCONFIG\",\n        \"FPS_MANAGER_GETCLEARINGBANKFUNDS\",\n        \"FPS_MANAGER_CONFREFRESH\",\n        \"FPS_MANAGER_GETREQUESTREFUNDTRADE\",\n        \"FPS_MANAGER_DELETESERVICECHANNELCONFIG\",\n        \"FPS_MANAGER_GETBANKPAYDETAIL\",\n        \"FPS_MANAGER_GETACCOUNTBILLUPLOAD\",\n        \"FPS_MANAGER_GETSERVICECHANNELCONFIG\",\n        \"FPS_MANAGER_DELETETRADECHANNELCONFIG\",\n        \"FPS_MANAGER_GETACCOUNTERRORBILL\",\n        \"FPS_MANAGER_EDITSERVICECHANNELCONFIG\",\n        \"FPS_MANAGER_GETBANKREFUNDSERVICEAPPLY\",\n        \"FPS_MANAGER_EDITGATEWAYCHANNELCONFIG\",\n        \"FPS_MANAGER_ADDTRADECHANNELCONFIG\",\n        \"FPS_MANAGER_GETRECEIVABLEERRORBILLRRECORD\",\n        \"FPS_MANAGER_GETBANKCHANNELCONFIG\",\n        \"FPS_MANAGER_GETDEBETPAYCOLLECT\"\n    ]\n}"
	var appMetadata model.AppMetadata
	err := json.Unmarshal([]byte(s), &appMetadata)
	if err != nil {
		return
	}
	fmt.Printf("appMetadata:%s", common.ToJSON(appMetadata))
}
