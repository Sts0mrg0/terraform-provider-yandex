// Code generated by sdkgen. DO NOT EDIT.

// nolint
package postgresql

import (
	"context"

	"google.golang.org/grpc"

	"github.com/yandex-cloud/go-genproto/yandex/cloud/mdb/postgresql/v1"
	"github.com/yandex-cloud/go-genproto/yandex/cloud/operation"
)

// DatabaseServiceClient is a postgresql.DatabaseServiceClient with
// lazy GRPC connection initialization.
type DatabaseServiceClient struct {
	getConn func(ctx context.Context) (*grpc.ClientConn, error)
}

var _ postgresql.DatabaseServiceClient = &DatabaseServiceClient{}

// Create implements postgresql.DatabaseServiceClient
func (c *DatabaseServiceClient) Create(ctx context.Context, in *postgresql.CreateDatabaseRequest, opts ...grpc.CallOption) (*operation.Operation, error) {
	conn, err := c.getConn(ctx)
	if err != nil {
		return nil, err
	}
	return postgresql.NewDatabaseServiceClient(conn).Create(ctx, in, opts...)
}

// Delete implements postgresql.DatabaseServiceClient
func (c *DatabaseServiceClient) Delete(ctx context.Context, in *postgresql.DeleteDatabaseRequest, opts ...grpc.CallOption) (*operation.Operation, error) {
	conn, err := c.getConn(ctx)
	if err != nil {
		return nil, err
	}
	return postgresql.NewDatabaseServiceClient(conn).Delete(ctx, in, opts...)
}

// Get implements postgresql.DatabaseServiceClient
func (c *DatabaseServiceClient) Get(ctx context.Context, in *postgresql.GetDatabaseRequest, opts ...grpc.CallOption) (*postgresql.Database, error) {
	conn, err := c.getConn(ctx)
	if err != nil {
		return nil, err
	}
	return postgresql.NewDatabaseServiceClient(conn).Get(ctx, in, opts...)
}

// List implements postgresql.DatabaseServiceClient
func (c *DatabaseServiceClient) List(ctx context.Context, in *postgresql.ListDatabasesRequest, opts ...grpc.CallOption) (*postgresql.ListDatabasesResponse, error) {
	conn, err := c.getConn(ctx)
	if err != nil {
		return nil, err
	}
	return postgresql.NewDatabaseServiceClient(conn).List(ctx, in, opts...)
}

// Update implements postgresql.DatabaseServiceClient
func (c *DatabaseServiceClient) Update(ctx context.Context, in *postgresql.UpdateDatabaseRequest, opts ...grpc.CallOption) (*operation.Operation, error) {
	conn, err := c.getConn(ctx)
	if err != nil {
		return nil, err
	}
	return postgresql.NewDatabaseServiceClient(conn).Update(ctx, in, opts...)
}