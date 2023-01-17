package model

import (
	"context"
	"gorm.io/gorm"
	"gorm.io/gorm/utils"
	"mall/internal/global"
	"strconv"
)

type Category struct {
	BaseModel
	Name             string      `gorm:"type:varchar(20);not null" json:"name"`
	ParentCategoryID int32       `json:"parent"`
	ParentCategory   *Category   `json:"-"`
	SubCategory      []*Category `gorm:"foreignKey:ParentCategoryID;references:ID" json:"sub_category"`
	Level            int32       `gorm:"type:int;not null;default:1" json:"level"`
	IsTab            bool        `gorm:"default:false;not null" json:"is_tab"`
}

type Brand struct {
	BaseModel
	Name string `gorm:"column:name;type:varchar(20);not null"`
	Logo string `gorm:"type:varchar(200);default:'';not null"`
}

type GoodsCategoryBrand struct {
	BaseModel
	CategoryID int32 `gorm:"type:int;index:idx_category_brand,unique"`
	Category   Category

	BrandID int32 `gorm:"type:int;index:idx_category_brand,unique"`
	Brand   Brand //这里不能写成Brands,因为默认是用brand_id去进行查询的
}
type Good struct {
	BaseModel

	CategoryID int32 `gorm:"type:int;not null"`
	Category   Category
	BrandID    int32 `gorm:"type:int;not null"`
	Brand      Brand

	OnSale   bool `gorm:"default:false;not null"`
	ShipFree bool `gorm:"default:false;not null"`
	IsNew    bool `gorm:"default:false;not null"`
	IsHot    bool `gorm:"default:false;not null"`

	Name            string   `gorm:"type:varchar(50);not null"`
	GoodsSn         string   `gorm:"type:varchar(50);not null"`
	ClickNum        int32    `gorm:"type:int;default:0;not null"`
	SoldNum         int32    `gorm:"type:int;default:0;not null"`
	FavNum          int32    `gorm:"type:int;default:0;not null"`
	MarketPrice     float32  `gorm:"not null"`
	ShopPrice       float32  `gorm:"not null"`
	GoodsBrief      string   `gorm:"type:varchar(100);not null"`
	Images          GormList `gorm:"type:varchar(1000);not null"`
	DescImages      GormList `gorm:"type:varchar(1000);not null"`
	GoodsFrontImage string   `gorm:"type:varchar(200);not null"`
}

//func (g *Good) BeforeCreate(tx *gorm.DB) (err error) {
//	u.UUID = uuid.New()
//
//	if !u.IsValid() {
//		err = errors.New("can't save invalid data")
//	}
//	return
//}

func (g *Good) AfterCreate(tx *gorm.DB) (err error) {
	esModel := EsGoods{
		ID:          g.ID,
		CategoryID:  g.CategoryID,
		BrandsID:    g.BrandID,
		OnSale:      g.OnSale,
		ShipFree:    g.ShipFree,
		IsNew:       g.IsNew,
		IsHot:       g.IsHot,
		Name:        g.Name,
		ClickNum:    g.ClickNum,
		SoldNum:     g.SoldNum,
		FavNum:      g.FavNum,
		MarketPrice: g.MarketPrice,
		GoodsBrief:  g.GoodsBrief,
		ShopPrice:   g.ShopPrice,
	}
	_, err = global.EsClient.Index().Index(EsGoods{}.IndexName()).BodyJson(esModel).Id(utils.ToString(g.ID)).Do(context.Background())
	if err != nil {
		return err
	}
	return nil
}
func (g *Good) AfterDelete(tx *gorm.DB) (err error) {
	_, err = global.EsClient.Delete().Index(EsGoods{}.IndexName()).Id(strconv.Itoa(int(g.ID))).Do(context.Background())
	if err != nil {
		return err
	}
	return nil
}

func (g *Good) AfterUpdate(tx *gorm.DB) (err error) {
	esModel := EsGoods{
		ID:          g.ID,
		CategoryID:  g.CategoryID,
		BrandsID:    g.BrandID,
		OnSale:      g.OnSale,
		ShipFree:    g.ShipFree,
		IsNew:       g.IsNew,
		IsHot:       g.IsHot,
		Name:        g.Name,
		ClickNum:    g.ClickNum,
		SoldNum:     g.SoldNum,
		FavNum:      g.FavNum,
		MarketPrice: g.MarketPrice,
		GoodsBrief:  g.GoodsBrief,
		ShopPrice:   g.ShopPrice,
	}

	_, err = global.EsClient.Update().Index(esModel.IndexName()).
		Doc(esModel).Id(strconv.Itoa(int(g.ID))).Do(context.Background())
	if err != nil {
		return err
	}
	return nil
}
