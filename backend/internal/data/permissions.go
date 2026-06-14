package data

import (
	"context"
	"slices"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Permissions []string

func (p Permissions) Include(code string) bool {
	return slices.Contains(p, code)
}

type PermissionsModel struct {
	DBpool *pgxpool.Pool
}

func (db *PermissionsModel) GetAllForUser(user_id int) (Permissions, error){

	qry := `select p.code from permissions p
		join user_permissions up on up.permission_id = p.permission_id
		join users users u on u.user_id = u.user_id
		where u.user_id = $1;`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	permissions := Permissions{}
	
	rows, err := db.DBpool.Query(ctx, qry, user_id)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var p string
		if err := rows.Scan(&p); err != nil {
			return nil, err
		}
		permissions = append(permissions, p)

	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	

	return permissions, err
}

func (db *PermissionsModel) AddForUser(user_id int, code string) error{

	qry := `insert into user_permission 
			select $1, permission_id from permissions where code = ANY($2)`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := db.DBpool.Exec(ctx, qry, user_id, code)

	return err
}

