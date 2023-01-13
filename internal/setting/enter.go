package setting

type group struct {
	Dao       mdao
	Log       elog
	Validator va
	Work      worker
	Maker     maker
	Es        es
}

var Group = new(group)

func InitAll() {
	Group.Dao.Init()
	Group.Log.Init()
	Group.Validator.InitTrans("zh")
	Group.Work.Init()
	Group.Maker.Init()
	Group.Es.InitEs()
}
