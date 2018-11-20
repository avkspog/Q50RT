package q50

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type Packet struct {
	ID       int
	Messages []*Message
}

type Message struct {
	MessageType    string
	NetType        string
	ID             int
	BatteryPercent uint8
	ReceiveTime    time.Time
	DeviceTime     time.Time
	Location
	errors []error
}

type Location struct {
	Latitude  float64
	Longitude float64
}

func Parse(data *[]byte) (packet *Packet, err error) {
	pack := new(Packet)
	text := strings.Trim(string(*data), " ")

	bktIndex := strings.Index(text, "[")

	if bktIndex == -1 || bktIndex != 0 {
		packet = nil
		err = errors.New("Expected [")
		return
	}

	fieldsBktFunc := func(r rune) bool {
		return r == '[' || r == ']'
	}
	f := strings.FieldsFunc(text, fieldsBktFunc)

	fieldsAstFunc := func(r rune) bool {
		return r == '*' || r == ','
	}

	currentTime := time.Now()

	for _, v := range f {
		messageFields := strings.FieldsFunc(v, fieldsAstFunc)

		if len(messageFields) < 7 {
			packet = nil
			err = errors.New("Broken message")
			return
		}

		message := new(Message)
		message.MessageType = messageFields[3]
		message.NetType = messageFields[0]
		message.ReceiveTime = currentTime
		id, err := strconv.Atoi(messageFields[1])
		if err == nil {
			pack.ID = id
			message.ID = id
		}

		switch message.MessageType {
		case "LK":
			parseLK(message, messageFields)
		case "UD":
			parseUD(message, messageFields)
		case "UD2":
			parseUD(message, messageFields)
		case "CONFIG":
			parseCONFIG(message, messageFields)
		}
		pack.addMessage(message)
	}

	packet = pack
	err = nil
	return
}

func parseLK(message *Message, messageFields []string) {
	//[3G*1234567890*000D*LK,23227,0,73]
	percent, err := strconv.ParseInt(messageFields[6], 10, 8)
	if err == nil {
		message.BatteryPercent = uint8(percent)
	}
}

func parseUD(message *Message, messageFields []string) {
	//[3G*1234567890*00A0*UD,051118,091654,V,00.000000,N,00.0000000,E,0.00,0.0,0.0,0,28,75,23282,0,00000008,4,255,250,1,46612,6762,122,46612,6761,128,46612,1562,117,46612,1561,113,0,36.6]
	rawDate := messageFields[4]
	rawTime := messageFields[5]
	var sb strings.Builder

	fmt.Fprintf(&sb, "20%s-%s-%sT%s:%s:%s.000Z", rawDate[4:], rawDate[2:len(rawDate)-2], rawDate[0:2],
		rawTime[0:2], rawTime[2:len(rawTime)-2], rawTime[4:])

	message.DeviceTime, _ = time.Parse(time.RFC3339, sb.String())

	if messageFields[8] == "N" {
		n, err := toFloat(messageFields[7])
		message.Latitude = n
		message.appendErr(err)
	}

	if messageFields[10] == "E" {
		n, err := toFloat(messageFields[9])
		message.Longitude = n
		message.appendErr(err)
	}
}

func parseUD2(message *Message, messageFields []string) {
	//[3G*1234567890*00CF*UD2,051118,090924,V,00.000000,N,00.0000000,E,0.00,0.0,0.0,0,100,77,23207,0,00000008,7,255,250,1,46612,6762,146,46612,6761,142,46612,6763,122,46612,1571,122,46612,1562,118,46612,1572,118,46612,9884,117,0,36.6]
}

func parseCONFIG(message *Message, messageFields []string) {
	//[3G*1234567890*007E*CONFIG,TY:g36,UL:300,SY:0,CM:0,WT:0,HR:0,TB:1,CS:0,PP:0,AB:1,HH:1,TR:0,MO:0,FL:1,VD:0,DD:0,SD:0,XY:0,WF:0,WX:0,PH:0,RW:0,MT:1,]
}

func toFloat(v string) (float64, error) {
	vt := strings.Trim(v, " ")
	if len(vt) == 0 || vt == "" {
		return 0, errors.New("Value is empty")
	}

	n, err := strconv.ParseFloat(vt, 64)
	if err != nil {
		return 0, err
	}
	return n, nil
}

func (m *Message) appendErr(err error) {
	if err != nil {
		m.errors = append(m.errors, err)
	}
}

func (p *Packet) addMessage(msg *Message) {
	if msg != nil {
		p.Messages = append(p.Messages, msg)
	}
}
