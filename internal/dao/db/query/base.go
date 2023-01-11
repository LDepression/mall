package query

import "gorm.io/gorm"

func Paginate(pn, ps int64) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if pn == 0 {
			pn = 1
		}
		switch {
		case ps > 100:
			ps = 100
		case ps <= 0:
			ps = 10
		}
		offset := (pn - 1) * ps
		return db.Offset(int(offset)).Limit(int(ps))
	}
}
