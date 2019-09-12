package main

import (
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"time"
)

var Message_len uint32

func Parse(c net.Conn) {
	tmp := make([]byte, Message_Header_Len)

	n, err := c.Read(tmp)
	if err != nil || n != Message_Header_Len {
		log.Fatal("read err :", n, err)
	}
	fmt.Println("read count : ", n)

	vsersion_num := binary_BigEndian_Uint32(tmp[:2])
	message_len := binary_BigEndian_Uint32(tmp[2:4])

	Message_len = message_len - Message_Header_Len - Set_Len

	export_time := time.Unix(int64(binary_BigEndian_Uint32(tmp[4:8])), 0).UTC()
	sequence_number := binary_BigEndian_Uint32(tmp[8:12])
	observation_domain_id := binary_BigEndian_Uint32(tmp[12:16])

	fmt.Println(vsersion_num)
	fmt.Println(message_len)
	fmt.Println(export_time)
	fmt.Println(sequence_number)
	fmt.Println(observation_domain_id)

	set := make([]byte, Set_Len)
	sn, err := c.Read(set)
	if err!=nil||sn!=Set_Len {
		log.Fatal("read set len err :",sn,err)
	}
	fmt.Println("read set count :",sn)

	set_id := binary_BigEndian_Uint32(set[:2])
	set_len := binary_BigEndian_Uint32(set[2:4])
	tempalte_id:=binary_BigEndian_Uint32(set[4:6])
	field_count := binary_BigEndian_Uint32(set[6:8])
	src_ip :=binary_BigEndian_Uint32(set[8:10])

}

func binary_BigEndian_Uint32(data []byte)uint32{
	return binary.BigEndian.Uint32(data)
}
