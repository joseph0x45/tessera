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
        app_id text not null references apps(id) on delete cascade,
        name text not null,
        password text not null,
        unique(id, name)
      );
    `,
	},
	{
		Version: 3,
		Name:    "metadata",
		SQL: `
      create table app_metadata (
        key text not null primary key,
        value text not null
      );
    `,
	},
	{
		Version: 4,
		Name:    "seed",
		SQL: `
    insert into app_metadata (
      key, value
    )
    values
      ('admin_password', '$2y$10$F2WPrx9uIsY1QvC205x.euRny62xAzsdOlnT/2smlVxY/uvnHkz7K'),
      ('signing_secret', 'signing_secret');
    `,
	},
	{
		Version: 5,
		Name:    "sessions",
		SQL: `
      create table sessions (
        id text not null primary key,
        session_user_id text not null references users(id) on delete cascade
      );
    `,
	},
}
