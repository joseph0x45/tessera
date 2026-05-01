package db

import "github.com/joseph0x45/sad"

var migrations = []sad.Migration{
	{
		Version: 1,
		Name:    "apps",
		SQL: `
      create table apps (
        id text not null primary key,
        name text not null unique
      );
    `,
	},
	{
		Version: 2,
		Name:    "users",
		SQL: `
      create table users (
        id text not null primary key,
        app_id text not null references apps.id,
        name text not null,
        password text not null,
        unique(id, name)
      );
    `,
	},
}
