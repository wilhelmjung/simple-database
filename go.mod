module main

go 1.15

replace db => ./src/

require (
	db v0.0.0-00010101000000-000000000000
	github.com/stretchr/testify v1.6.1 // indirect
)
