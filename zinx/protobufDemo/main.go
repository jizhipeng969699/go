package main

import (
	"encoding/base64"
	"fmt"
	"github.com/golang/protobuf/proto"
	"zinx/protobufDemo/pb"
)

//
//type AA struct {
//	Name                 string   `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
//	Age                  int32    `protobuf:"varint,2,opt,name=age,proto3" json:"age,omitempty"`
//	XXX_NoUnkeyedLiteral struct{} `json:"-"`
//	XXX_unrecognized     []byte   `json:"-"`
//	XXX_sizecache        int32    `json:"-"`
//}
const (
	base64Table = "123QRSTUabcdVWXYZHijKLAWDCABDstEFGuvwxyzGHIJklmnopqr234560178912"
)

func main() {
	enc := base64.NewEncoding(base64Table)
	str := enc.EncodeToString([]byte("我是你爸爸"))

	person := &pb.Person{
		Name:   str,
		Age:    18,
		Emails: []string{"aaa", "bbb", "ccc"},
		Phones: []*pb.PhoneNumber{
			&pb.PhoneNumber{
				Number: "111",
			},
			&pb.PhoneNumber{
				Number: "222",
			},
			&pb.PhoneNumber{
				Number: "333",
			},
		},
		Data: &pb.Person_School{
			School: "ooo",
		},
	}

	olddata, err := proto.Marshal(person)
	if err != nil {
		panic(err)
	}

	fmt.Println(person.Name)
	fmt.Println(person.Age)
	fmt.Println(person.Emails)
	fmt.Println(person.Phones)

	fmt.Println("--------------------------------")
	//aa:=new(AA)
	per := new(pb.Person)
	err = proto.Unmarshal(olddata, per)
	if err != nil {
		panic(err)
	}

	//base64.NewEncoding(base64Table)
	by, err := enc.DecodeString(per.Name)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(by))
	fmt.Println(per.Age)
	fmt.Println(per.Emails)
	fmt.Println(per.Phones)
	fmt.Println(per.GetSchool())
}
