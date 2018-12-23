package q50

import (
	"fmt"
	"strings"
	"testing"
	"time"
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
		err:      "Expected last ]",
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
	rawDate := "271018"
	rawTime := "035856"
	var sb strings.Builder

	fmt.Fprintf(&sb, "20%s-%s-%sT%s:%s:%s.000Z", rawDate[4:], rawDate[2:len(rawDate)-2], rawDate[0:2],
		rawTime[0:2], rawTime[2:len(rawTime)-2], rawTime[4:])

	dt, _ := time.Parse(time.RFC3339, sb.String())
	rt, _ := time.Parse(time.Kitchen, "3:04PM")

	lt := Location{
		Latitude:  00.312762,
		Longitude: 00.3385200,
	}

	msg := Message{
		ID:             1234567890,
		MessageType:    "LK",
		NetType:        "3G",
		BatteryPercent: 78,
		DeviceTime:     dt,
		ReceiveTime:    rt,
		Location:       lt,
	}
	b := []byte("[3G*1234567890*000D*LK,23201,0,78] [[[3G*1234567890*007D*CONFIG,TY:g36,UL:60,SY:0,CM:0,WT:0,HR:0,TB:1,CS:0,PP:0,AB:1,HH:1,TR:0,MO:0,FL:1,VD:0,DD:0,SD:0,XY:0,WF:0,WX:0,PH:0,RW:0,MT:1,][3G*1234567890*00CC*UD2,161018,060356,A,00.312705,N,00.3389767,E,3.00,116.8,0.0,6,88,2,9714,0,00000001,7,1,250,1,46612,1563,142,46612,1562,145,46612,6772,135,46612,1571,135,46612,1572,129,46612,6762,128,46612,8532,122,0,23.7][3G*1234567890*00CD*UD2,161018,060430,A,00.312618,N,00.3388017,E,0.00,348.2,0.0,6,100,2,9722,0,00000001,7,1,250,1,46612,1562,140,46612,1563,140,46612,6772,136,46612,1571,131,46612,1572,126,46612,6762,124,46612,6771,122,0,23.3][3G*1234567890*00CD*UD2,161018,060622,A,00.312600,N,00.3394767,E,1.72,30.7,0.0,6,83,2,9755,0,00000001,7,255,250,1,46612,1563,137,46612,1562,139,46612,6772,129,46612,1571,127,46612,6762,126,46612,1572,121,46612,6771,119,0,19.2][3G*1234567890*00CB*UD2,161018,060631,A,00.312638,N,00.3394183,E,0.00,27.5,0.0,7,93,2,9771,0,00000001,7,1,250,1,46612,1563,137,46612,1562,143,46612,6772,130,46612,1571,127,46612,1572,126,46612,6762,121,46612,8532,121,0,20.5][3G*1234567890*00CC*UD2,161018,060801,A,00.312762,N,00.3385200,E,2.29,56.9,0.0,4,100,2,9961,0,00000001,7,1,250,1,46612,1563,150,46612,1562,149,46612,1572,132,46612,1571,128,46612,6772,128,46612,1561,124,46612,6762,121,0,30.6][3G*1234567890*00CD*UD2,161018,061330,V,00.312762,N,00.3385200,E,0.00,0.0,0.0,0,89,2,10088,0,00000001,7,255,250,1,46612,1562,140,46612,1563,142,46612,6772,132,46612,1571,130,46612,6762,122,46612,6771,121,46612,1561,120,0,30.6][3G*1234567890*00CC*UD2,161018,061630,V,00.312762,N,00.3385200,E,0.00,0.0,0.0,0,100,2,10192,0,00000001,7,0,250,1,46612,1562,151,46612,1563,152,46612,1571,135,46612,8533,134,46612,1572,134,46612,1561,125,46612,6762,125,0,30.6][3G*1234567890*0092*UD2,271018,035247,V,00.312762,N,00.3385200,E,0.00,0.0,0.0,0,65,93,10360,0,00000000,3,255,250,1,46612,6761,133,46612,6762,132,46612,1571,121,0,30.6][3G*1234567890*0092*UD2,271018,035548,V,00.312762,N,00.3385200,E,0.00,0.0,0.0,0,45,93,10407,0,00000000,3,255,250,1,46612,6761,127,46612,6762,132,46612,1571,113,0,30.6][3G*1234567890*00B0*UD2,271018,035856,V,00.312762,N,00.3385200,E,0.00,0.0,0.0,0,71,93,10517,0,00000000,5,255,250,1,46612,6762,135,46612,9884,127,46612,1571,116,46612,6781,114,46612,1562,113,0,30.6]")
	message, _ := LastMessage(&b)
	message.ReceiveTime = rt

	if *message != msg {
		t.Error("broken message")
	}
}
