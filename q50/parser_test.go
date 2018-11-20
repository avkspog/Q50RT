package q50

import (
	"testing"
)

type testStruct struct {
	TestName       string
	MessageType    string
	Message        string
	MessageCount   int
	ID             int
	BatteryPercent uint8
	Latitude       float64
	Longitude      float64
	err            string
}

var packetTests = []testStruct{
	{
		TestName: "Packet ID",
		Message:  "[3G*1234567890*000D*LK,23201,0,78]",
		ID:       1234567890,
		err:      "packet.ID != test.ID",
	},
	{
		TestName: "Broken message 1",
		Message:  "3G*1234567890*000D*LK,23201,0,78]",
		err:      "Expected [",
	},
	{
		TestName: "Broken message 2",
		Message:  "ghjad*3G*sf dkfjdf* hjf[ adsf*LK* adsf asdf] df asdffd623]",
		err:      "Expected [",
	},
	{
		TestName: "Broken message 3",
		Message:  "[ghjad*3G*sf dkfjdf* hjf[ adsf*LK* adsf asdf] df asdffd623",
		err:      "Broken message",
	},
	{
		TestName:     "Messages count 1",
		Message:      "[3G*1234567890*000D*LK,23201,0,78]",
		MessageCount: 1,
		err:          "Messages count != 1",
	},
	{
		TestName: "Broken message 4",
		Message:  "[3G*1234567890*000D*LK00000078]",
		err:      "Broken message",
	},
}

func TestPacket(t *testing.T) {
	for _, test := range packetTests {
		t.Run(test.TestName, func(t *testing.T) {
			b := []byte(test.Message)
			packet, err := Parse(&b)

			if err != nil {
				if err.Error() != test.err {
					t.Error(err)
				}
			}

			if test.ID != 0 && packet.ID != test.ID {
				t.Error(test.err)
			}

			if test.MessageCount != 0 && test.MessageCount != len(packet.Messages) {
				t.Error(test.err)
			}
		})
	}
}

func TestMessages(t *testing.T) {

}
