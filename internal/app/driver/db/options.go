// Copyright (c) 2020 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package db

import "github.com/harunnryd/tempokerja/config"

const (
	// MysqlDialectParam ...
	MysqlDialectParam = "mysql"
	// PgsqlDialectParam ...
	PgsqlDialectParam = "postgres"
)

// Option ...
type Option func(*DB)

// WithConfig ....
func WithConfig(config config.Config) Option {
	return func(db *DB) {
		db.config = config
	}
}
