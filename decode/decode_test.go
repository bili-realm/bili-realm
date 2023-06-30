package decode

import (
	"fmt"
	"testing"
)

var TestStr = []byte(`
{
  "cmd": "DANMU_MSG",
  "info": [
    [
      0,
      1,
      25,
      5816798,
      1687934773234,
      2098934283,
      0,
      "ea6706c9",
      0,
      0,
      0,
      "",
      0,
      "{}",
      "{}",
      {
        "mode": 0,
        "show_player_type": 0,
        "extra": "{\"send_from_me\":false,\"mode\":0,\"color\":5816798,\"dm_type\":0,\"font_size\":25,\"player_mode\":1,\"show_player_type\":0,\"content\":\"哈哈哈哈哈\",\"user_hash\":\"3932620489\",\"emoticon_unique\":\"\",\"bulge_display\":0,\"recommend_score\":0,\"main_state_dm_color\":\"\",\"objective_state_dm_color\":\"\",\"direction\":0,\"pk_direction\":0,\"quartet_direction\":0,\"anniversary_crowd\":0,\"yeah_space_type\":\"\",\"yeah_space_url\":\"\",\"jump_to_url\":\"\",\"space_type\":\"\",\"space_url\":\"\",\"animation\":{},\"emots\":null,\"is_audited\":false,\"id_str\":\"45d97c722a2259be4ff079a616649bd748\"}"
      },
      {
        "activity_identity": "",
        "activity_source": 0,
        "not_show": 0
      },
      ""
    ],
    "哈哈哈哈哈",
    [
      8577714,
      "找不到答案",
      0,
      0,
      0,
      10000,
      1,
      ""
    ],
    [],
    [
      12,
      0,
      6406234,
      "\u003e50000",
      0
    ],
    [
      "",
      ""
    ],
    0,
    0,
    null,
    {
      "ts": 1687934773,
      "ct": "9596333F"
    },
    0,
    0,
    null,
    null,
    0,
    7,
    [
      7
    ]
  ],
  "dm_v2": "CiI0NWQ5N2M3MjJhMjI1OWJlNGZmMDc5YTYxNjY0OWJkNzQ4EAEYGSDeg+MCKghlYTY3MDZjOTIP5ZOI5ZOI5ZOI5ZOI5ZOIOPLPgoaQMUiL5OzoB2IAigEAmgEQCgg5NTk2MzMzRhC1ru+kBqIBhAEIssWLBBIP5om+5LiN5Yiw562U5qGIIkpodHRwczovL2kwLmhkc2xiLmNvbS9iZnMvZmFjZS9lNTM4ZjdmMWMzNzFiMTEwOGMzNDUyZWMxZjM2YjU0MjYxYTgwNzNiLmpwZziQTkABWgIIAWIPCAwQ2oCHAxoGPjUwMDAwagByAHoCCAeqAQQYjfQE"
}`)

func TestDecode(t *testing.T) {
	payload, err := Decode(TestStr)
	if err != nil {
		t.Fail()
	}
	t.Log(fmt.Sprintf("%#v", payload))
	if payload.Cmd != "DANMU_MSG" || payload.Danmaku == nil {
		t.Fail()
	}
}

func TestProcessDanmaku(t *testing.T) {

}
