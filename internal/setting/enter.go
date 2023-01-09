package setting

type group struct {
	Dao       mdao
	Log       log
	Validator va
	Work      worker
	Maker     maker
}

var Group = new(group)

func InitAll() {
	Group.Dao.Init()
	Group.Log.Init()
	Group.Validator.InitTrans("zh")
	Group.Work.Init()
	Group.Maker.Init()
}
