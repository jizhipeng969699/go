package main

import "time"

type Message_Head struct {
	Version_num uint32//版本号
	Messagge_len uint32//消息长度
	Export_time time.Time//离开导出器导出时间   1970年1月1日UNIX时代以来以秒为单位表示 00:00 UTC
	Sequence_num uint32//序列号
	Observation_domain_id uint32//观察域ID
}


type Set_Base struct {
	Set_id uint32  //2 template msg  3 options msg  >=256 data msg   2字节
	Set_len uint32  // 包括set头 知道下一个 set的位置 2

}

type Template_Msg struct {
	Set_base *Set_Base //8 字节

	Template_id uint32 //用来表示data  他的dataset  2
	Field_count uint32 //字段数量    2

	Field_info []*Template_Filed //字段

	Padding_opt  uint32 //结尾  不处理
}

type Template_Filed struct {
	Information_Element_id uint32 //信息元素id    2字节
	Field_len   uint32 //字段的长度	2字节
	Enterprise_Number  uint32  //企业编号   4字节
}



type Options_Templte_Msg struct {
	Set_base *Set_Base		// 8 字节

	Template_id uint32 //用来表示data  他的dataset  2
	Field_count uint32 //字段数量    2

	Scope_field_count uint32 //范围字段数量  2字节 n
	Scope_filed []*Scope_Filed
	Scope_enterprise_num uint32 //4字节
	Option_filed []*Option_Filed
	Padding_opt uint32 //结尾 不处理


}
type Scope_Filed struct {
	Scope_info_element_id uint32//2
	Scop_field_len uint32		//2
}
type Option_Filed struct {
	Option_info_element_id uint32
	Option_field_len uint32
	Option_enterprise_num uint32 //4字节
}


type Data_Msg struct {
	//Set_id uint32 //就是templte 或者options 中templte 的id  2字节
	//Set_len uint32 //长度

	Set_base *Set_Base
	Field_value []*Field_value  //只有在templte 或者options 中才能解释 长度
	Padding_opt uint32 //结尾部处理

}

type Field_value struct {

}




//具体含义 参考： https://www.iana.org/assignments/ipfix/ipfix.xml
type Filed_info struct {
	Name string
	Abstract_Data_Type string
}
var Fiels_Info_Map =map[uint32]*Filed_info{
	1:&Filed_info{Name:"octetDeltaCount",Abstract_Data_Type:"unsigned64"},
	2:&Filed_info{Name:"octetDeltaCount",Abstract_Data_Type:"unsigned64"},
	3:&Filed_info{Name:"octetDeltaCount",Abstract_Data_Type:"unsigned64"},
	4:&Filed_info{Name:"octetDeltaCount",Abstract_Data_Type:"unsigned64"},
	5:&Filed_info{Name:"octetDeltaCount",Abstract_Data_Type:"unsigned64"},
	6:&Filed_info{Name:"octetDeltaCount",Abstract_Data_Type:"unsigned64"},
	7:&Filed_info{Name:"octetDeltaCount",Abstract_Data_Type:"unsigned64"},
	8:&Filed_info{Name:"octetDeltaCount",Abstract_Data_Type:"unsigned64"},
	9:&Filed_info{Name:"octetDeltaCount",Abstract_Data_Type:"unsigned64"},
	10:&Filed_info{Name:"octetDeltaCount",Abstract_Data_Type:"unsigned64"},
	11:&Filed_info{Name:"octetDeltaCount",Abstract_Data_Type:"unsigned64"},

}
