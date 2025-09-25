package repository

import (
	"context"
	"fmt"
)

type Row struct {
}

// StructScan fake StructScan
func (r *Row) StructScan(dest ...interface{}) error {
	return nil
}

type FakeDB struct {
}

func NewFakeDB() *FakeDB {
	return &FakeDB{}
}

func (f *FakeDB) GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	fmt.Println("Call GetContext")
	return nil
}

func (f *FakeDB) QueryRowxContext(ctx context.Context, query string, args ...interface{}) *Row {
	fmt.Println("Call QueryRowxContext")
	return &Row{}
}
