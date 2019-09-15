package main

import (
	"encoding/binary"
	"errors"
	"log"
	"net"
	"time"
)

var Options_map = map[uint32]*Options_Templte_Msg{}
var Templte_map = map[uint32]*Template_Msg{}

const (
	Message_Header_Len = 16 //消息头长度
)

func Parse(c net.Conn) {
	defer func() {
		err := c.Close()
		if err != nil {
			log.Println("conn close err :", err)
		}
	}()

	for {

		msg_head_sli := make([]byte, Message_Header_Len)
		n, err := c.Read(msg_head_sli)
		if err != nil || n != Message_Header_Len {
			log.Println("conn read msg head err :", err, "读到的字节数：", n)
			return
		}
		msg_head := parse_msg_head(msg_head_sli[:])

		msg_data_sli := make([]byte, msg_head.Messagge_len-Message_Header_Len)
		n, err = c.Read(msg_data_sli)
		if err != nil {
			log.Println("conn read data err :", err)
			return
		}

		parse_msg(msg_data_sli[:])

	}
}

func binary_BigEndian_Uint32(data []byte) uint32 {
	return binary.BigEndian.Uint32(data)
}

func parse_msg_head(tmp []byte) *Message_Head {
	if tmp == nil || len(tmp) != Message_Header_Len {
		return nil
	}

	return &Message_Head{
		Version_num:           binary_BigEndian_Uint32(tmp[:2]),
		Messagge_len:          binary_BigEndian_Uint32(tmp[2:4]),
		Export_time:           time.Unix(int64(binary_BigEndian_Uint32(tmp[4:8])), 0).UTC(),
		Sequence_num:          binary_BigEndian_Uint32(tmp[8:12]),
		Observation_domain_id: binary_BigEndian_Uint32(tmp[12:16]),
	}
}

func parse_msg(tmp []byte) error {
	if tmp == nil || len(tmp) == 0 {
		return errors.New("data info err ")
	}

	set_base := &Set_Base{
		Set_id:  binary_BigEndian_Uint32(tmp[:2]),
		Set_len: binary_BigEndian_Uint32(tmp[2:4]),
	}
	tmp = tmp[4:]

	switch set_base.Set_id {
	case 2:
		field_count := binary_BigEndian_Uint32(tmp[6:8])

		temp := &Template_Msg{
			Set_base:    set_base,
			Template_id: binary_BigEndian_Uint32(tmp[4:6]),
			Field_count: field_count,
			Field_info:  parse_template_field(field_count, tmp[8:]),
			Padding_opt: nil,
		}

		Templte_map[temp.Template_id] = temp
		tmp = tmp[temp.Set_base.Set_len:]

	case 3:
		field_count := binary_BigEndian_Uint32(tmp[6:8])
		opt := &Options_Templte_Msg{
			Set_base:             set_base,
			Template_id:          binary_BigEndian_Uint32(tmp[4:6]),
			Field_count:          field_count,
			Scope_field_count:    binary_BigEndian_Uint32(tmp[8:10]),
			Scope_filed:          nil,
			Scope_enterprise_num: nil,
			Option_filed:         nil,
		}
		Options_map[opt.Template_id] = opt
		tmp = tmp[opt.Set_base.Set_len:]

	default:
		if set_base.Set_id >= 256 {
			if tv, tok := Templte_map[set_base.Set_id]; tok {
				data_msg := &Data_Msg{
					Set_base:    set_base,
					Field_value: parse_template_field_value(tv.Field_info),
					Padding_opt: nil,
				}

			} else if ov, ook := Options_map[set_base.Set_id]; ook {
				data_msg := &Data_Msg{
					Set_base:    set_base,
					Field_value: parse_options_field_value(ov.),
					Padding_opt: nil,
				}
			}
		}
	}

	err := parse_msg(tmp)
	if err!=nil{
		log.Fatal(err)
		return err
	}
	return nil

}

func parse_template_field(num uint32, tmp []byte) []*Template_Filed {
	if num <= 0 || tmp == nil || len(tmp) == 0 {
		return nil
	}

	tmp_sli := []*Template_Filed{}

	for i := 0; i < int(num); i++ {
		tmp_sli = append(tmp_sli, &Template_Filed{
			Information_Element_id: binary_BigEndian_Uint32(tmp[:2]),
			Field_len:              binary_BigEndian_Uint32(tmp[2:4]),
			Enterprise_Number:      binary_BigEndian_Uint32(tmp[4:8]),
		})
		tmp = tmp[8:]
	}

	return tmp_sli

}

//TODO
func parse_options_field_value(temp *Template_Filed) []*Field_value {
	return nil
}

func parse_template_field_value(temp []*Template_Filed) []*Field_value {
	return nil
}
