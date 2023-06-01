package repository

import (
	"database/sql"
	"errors"
	"log"
	"product-go/helper/failerror"
	"product-go/helper/random"
	"product-go/helper/timeout"
	"product-go/model"
	"product-go/publisher"
	"time"
)

type repository struct {
	db   *sql.DB
	sent publisher.Publisher
}

func NewRepository(db *sql.DB, sent publisher.Publisher) Repositorier {
	return &repository{
		db:   db,
		sent: sent,
	}
}

func (repo *repository) GetProduct(req model.ProductSearch) ([]model.Product, error) {
	ctx, cancel := timeout.NewCtxTimeout()
	defer cancel()

	querySearchProduct := `select p.id, p.store_id, p.category_id, p.size_id, p.color_id ,p.name, p.brand, p.subtitle, p.description ,p.unit_price ,p.status ,p.stock ,p.sku ,p.weight ,p.created_at ,p.updated_at from products p  join categories c on c.id = p.category_id join product_sizes s on s.id = p.size_id join product_colors co on co.id = p.color_id where c.name like '%' || $1 || '%' and p.name like '%' || $2 || '%' and p.brand like '%' || $3 || '%'`

	result, err := repo.db.QueryContext(ctx, querySearchProduct, req.Category, req.Name, req.Brand)
	failerror.FailError(err, "error query")

	var data = []model.Product{}

	for result.Next() {
		var temp model.Product
		result.Scan(&temp.Id, &temp.StoreID, &temp.CategoryID, &temp.SizeID, &temp.ColorID, &temp.Name, &temp.Brand, &temp.Subtitle, &temp.Description, &temp.UnitPrice, &temp.Status, &temp.Stock, &temp.Sku, &temp.Weight, &temp.Created_at, &temp.Update_at)
		data = append(data, temp)
	}

	return data, nil
}

func (repo *repository) ShowProduct(id int) (model.Product, error) {
	ctx, cancel := timeout.NewCtxTimeout()
	defer cancel()

	query := `select  id, store_id, category_id, size_id, color_id , name, subtitle , description , unit_price , status , stock ,sku , weight , created_at , updated_at from products p where id = $1`

	result, err := repo.db.QueryContext(ctx, query, id)
	failerror.FailError(err, "error prepare")

	var temp = model.Product{}
	for result.Next() {
		result.Scan(&temp.Id, &temp.StoreID, &temp.CategoryID, &temp.SizeID, &temp.ColorID, &temp.Name, &temp.Subtitle, &temp.Description, &temp.UnitPrice, &temp.Status, &temp.Stock, &temp.Sku, &temp.Weight, &temp.Created_at, &temp.Update_at)
	}

	if temp.Id <= 0 {
		return temp, errors.New("product not found")
	}

	return temp, nil
}

func (repo *repository) CreateProduct(req []model.ProductReq) ([]model.Product, error) {
	ctx, cancel := timeout.NewCtxTimeout()
	defer cancel()

	var sent []model.Product

	for _, v := range req {
		inrandom := model.Product{
			StoreID:     v.StoreID,
			CategoryID:  v.CategoryID,
			SizeID:      v.SizeID,
			ColorID:     v.ColorID,
			Name:        v.Name,
			Brand:       v.Brand,
			Subtitle:    v.Subtitle,
			Description: v.Description,
			UnitPrice:   v.UnitPrice,
			Status:      false,
			Stock:       v.Stock,
			Sku:         random.NewRandom().RandomString(),
			Weight:      v.Weight,
		}

		sent = append(sent, inrandom)
	}

	err := repo.sent.Public(sent, "create_product")
	if err != nil {
		return nil, errors.New("failed publisher")
	}

	time.Sleep(3 * time.Second)

	var resultss []model.Product
	query := `select  id, store_id, category_id, size_id, color_id , name, brand,subtitle , description , unit_price , status , stock ,sku , weight , created_at , updated_at from products p where sku = $1`

	stmt, err := repo.db.PrepareContext(ctx, query)
	failerror.FailError(err, "error prepare")

	for _, v := range sent {

		result, err := stmt.QueryContext(ctx, v.Sku)
		failerror.FailError(err, "error prepare")

		var temp model.Product
		for result.Next() {
			result.Scan(&temp.Id, &temp.StoreID, &temp.CategoryID, &temp.SizeID, &temp.ColorID, &temp.Name, &temp.Brand, &temp.Subtitle, &temp.Description, &temp.UnitPrice, &temp.Status, &temp.Stock, &temp.Sku, &temp.Weight, &temp.Created_at, &temp.Update_at)
		}
		if temp.Id == 0 {
			continue
		}
		resultss = append(resultss, temp)
	}

	if resultss == nil {
		return []model.Product{}, errors.New("error create product")
	}

	return resultss, nil
}

func (repo *repository) UpdateProduct(req model.ProductUpd) (model.Product, error) {
	ctx, cancel := timeout.NewCtxTimeout()
	defer cancel()

	err := repo.sent.Public(req, "update_product")
	if err != nil {
		return model.Product{}, errors.New("failed publisher")
	}

	log.Println(req)
	time.Sleep(1 * time.Second)

	query := `select id, store_id, category_id, size_id, color_id , name, brand,subtitle , description , unit_price , status , stock ,sku , weight , created_at , updated_at from products p where id = $1`

	result, err := repo.db.QueryContext(ctx, query, req.Id)
	failerror.FailError(err, "error prepare")

	var temp model.Product
	for result.Next() {
		result.Scan(&temp.Id, &temp.StoreID, &temp.CategoryID, &temp.SizeID, &temp.ColorID, &temp.Name, &temp.Brand, &temp.Subtitle, &temp.Description, &temp.UnitPrice, &temp.Status, &temp.Stock, &temp.Sku, &temp.Weight, &temp.Created_at, &temp.Update_at)
	}

	if temp.Id == 0 {
		return model.Product{}, errors.New("error create product")
	}

	return temp, nil
}

func (repo *repository) DeleteProduct(id int) (int, error) {

	ctx, cancel := timeout.NewCtxTimeout()
	defer cancel()

	var idCheck int
	queryCheck := `select id from products where id = $1`
	err := repo.db.QueryRowContext(ctx, queryCheck, id).Scan(&idCheck)
	failerror.FailError(err, "error exec")

	if idCheck == 0 {
		return 0, errors.New("product tidak ditemukan")
	}

	query := `DELETE FROM products WHERE id = $1`
	_, err = repo.db.ExecContext(ctx, query, idCheck)
	failerror.FailError(err, "error exec")

	return idCheck, nil
}
