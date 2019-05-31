package seqno

import (
	"github.com/jinzhu/gorm"
	lua "github.com/yuin/gopher-lua"
)

// DocumentStandard 单据标准表
type DocumentStandard struct {
	ID          int64  `json:"id",gorm:"auto-increment"`
	LogicID     string `json:"logic_ID"`     //类型
	FormateCode string `json:"formate_Code"` //格式代码
	MarkCode    string `json:"mark_Code"`    //标识代码
	ScriptCode  string `json:"script_Code"`  //脚本代码
	StepNum     int64  `json:"step_Num"`     //步长
	Remark      string `json:"remark"`       //备注
}

// DocumentGainNumber 年度单据发号表
type DocumentGainNumber struct {
	ID               int64 `json:"id",gorm:"auto-increment"`
	DsID             int64
	DocumentStandard DocumentStandard `gorm:"ForeignKey:DsID;association_foreignkey:ID"`
	LogicID          string           `json:"logic_ID"`
	FormateCode      string           `json:"formate_Code"`
	MarkCode         string           `json:"mark_Code"`
	CurrentNum       int64            `json:"current_Num"`
	StepNum          int64            `json:"step_Num"`
}

// DocNo 操作代理
type DocNo struct {
	conn             *gorm.DB
	currentNum       int64
	logicID          string
	formateCode      string
	markCode         string //lua代码
	step             int64
	remark           string
	gainMark         string //通过lua生成出来的
	documentStandard *DocumentStandard
}

// NewNoGenerator 初始化新代理
func NewNoGenerator(db *gorm.DB, logicID string) *DocNo {
	return &DocNo{
		conn:        db,
		currentNum:  0,
		logicID:     logicID,
		formateCode: "%s%06d",
		markCode:    "",
		step:        1,
		remark:      "",
	}
}

// InitTable 初始化表
func InitTable(db *gorm.DB) {
	db.AutoMigrate(&DocumentStandard{})
	db.AutoMigrate(&DocumentGainNumber{})
}

// Step 步长
func (s *DocNo) Step(step int64) *DocNo {
	s.step = step
	return s
}

// StartWith 起始数
func (s *DocNo) StartWith(start int64) *DocNo {
	s.currentNum = start
	return s
}

// FormateCode 格式代码
func (s *DocNo) FormateCode(format string) *DocNo {
	s.formateCode = format
	return s
}

// MarkCode 标识代码
func (s *DocNo) MarkCode(markCode string) *DocNo {
	s.markCode = markCode
	return s
}

// Remark 备注
func (s *DocNo) Remark(remark string) *DocNo {
	s.remark = remark
	return s
}

// Next 返回序列号
func (s *DocNo) Next() (int64, error) {
	return s.next()
}

// 通过lua生成mark
func (s *DocNo) gainMarkCode(L *lua.LState) string {
	return ""
}

// 通过lua生成 Elements
func (s *DocNo) gainElements(L *lua.LState) []string {
	newArr := make([]string, 0)
	return newArr
}

// 返回
func (s *DocNo) next() (int64, error) {
	L := lua.NewState()
	defer L.Close()
	s = s.findDocumentStandard()
	var gainMark = ""
	if s.markCode != "" {
		gainMark = s.gainMarkCode(L)
	}

	s.gainMark = gainMark

	dgn := s.generateNextSeqNumber()

	nextSeq := dgn.CurrentNum + dgn.StepNum

	s.currentNum = nextSeq

	// elements := s.gainElements(L)

	s.conn.Model(dgn).Update("current_Num", s.currentNum)

	return s.currentNum, nil

}

func (s *DocNo) findDocumentStandard() *DocNo {
	var ds DocumentStandard
	query := s.conn.First(&ds, "logic_ID = ?", s.logicID)

	if query.Error != nil { //没有找到，新建一个
		documentStandard := &DocumentStandard{
			LogicID:     s.logicID,
			FormateCode: s.formateCode,
			MarkCode:    s.markCode,
			StepNum:     s.step,
			Remark:      s.remark,
		}
		s.conn.Create(documentStandard)
		s.documentStandard = documentStandard
	} else {
		s.formateCode = ds.FormateCode
		s.markCode = ds.MarkCode
		s.step = ds.StepNum
		s.step = ds.StepNum
		s.remark = ds.Remark
		s.documentStandard = &ds
	}

	return s
}

func (s *DocNo) generateNextSeqNumber() *DocumentGainNumber {
	var dgn DocumentGainNumber

	var condition = map[string]interface{}{
		"logic_ID": s.logicID,
		"ds_id":    s.documentStandard.ID,
	}

	if s.gainMark != "" {
		condition["mark_Code"] = s.gainMark
	}

	query := s.conn.Where(condition).First(&dgn)

	if query.Error != nil { //没有找到，新建一个
		documentGainNumber := &DocumentGainNumber{
			LogicID:          s.logicID,
			FormateCode:      s.formateCode,
			MarkCode:         s.markCode,
			DocumentStandard: *s.documentStandard,
			DsID:             s.documentStandard.ID,
			StepNum:          s.step,
			CurrentNum:       s.currentNum,
		}
		s.conn.Create(documentGainNumber)
		return documentGainNumber
	}

	return &dgn
}
