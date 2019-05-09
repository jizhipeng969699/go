package models

import (
	"github.com/astaxie/beego/orm"
	"time"
)
import _ "github.com/go-sql-driver/mysql"

type Thirduser struct {
	Id   int
	Name string
	Pwd  string

	Txt []*Txt `orm:"rel(m2m)"`
}
type Txt struct {
	Id        int       `orm:"pk;auto"`
	Title     string    `orm:"size(50)"`
	Content   string    `orm:"size(500)"`
	Img       string    `orm:"null"`
	Time      time.Time `orm:"type(datetime);auto_now_add"`
	Readcount int       `orm:"defalut(0)"`

	Txttype *Txttype `orm:"rel(fk);null;on_delete(set_null)"` //文章和文章类型是一对多关系

	Thirduser []*Thirduser `orm:"reverse(many)"` //
}
type Txttype struct {
	Id       int
	Typename string `orm:"unique"`

	Txt []*Txt `orm:"reverse(many)"`
}

func init() {
	orm.RegisterDataBase("default", "mysql", "root:woshini88@tcp(127.0.0.1:3306)/third")
	orm.RegisterModel(new(Thirduser), new(Txt), new(Txttype))
	orm.RunSyncdb("default", false, true)

}
