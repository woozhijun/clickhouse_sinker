package impl

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestMMlogV3Parser_Parse(t *testing.T) {

	Convey("check tasks input", t, func() {

		body := "info 2019-05-31T09:56:07.247916+08:00 127.0.0.1 search 154280 853559408 _MMLOGV3_ - search_result {\"uid\": \"197775832\", \"func_name\": \"search_user_contact2\", \"abtest_key\": \"es_search\", \"req_lim\": 4, \"ptype\": \"\", \"words\": \"query=\\u5b89\\u5fc3\\u4fdd\\u9669&dict=\\u5b89\\u5fc3\\u4fdd\\u9669:1\", \"offset\": 0, \"now_in_sec\": \"1559267767\", \"result_len\": 4, \"event_type\": \"single\", \"url_prefix\": \"v2/s_uc\", \"event_id\": \"414523aa834711e9a139801844e1af5c\", \"d_query2type\": {\"\\u5b89\\u5fc3\\u4fdd\\u9669\": 1}, \"frm\": \"\", \"result_info\": [652695, 34066128, 573990, 55827243], \"req_id\": \"9deba9cfe75c4b21a13dcb5c114685f2\", \"now_str\": \"2019-05-31 09:56:07\", \"result_type\": \"regular\"}"
		p := &MMlogV3Parser{}
		parser := p.Parse([]byte(body))
		println(">>>>>." + parser.GetString("datetime"))
	})
}
