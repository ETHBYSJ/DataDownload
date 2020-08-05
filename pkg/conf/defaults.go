package conf

var DatabaseConfig = &database{
	Type: 	"mysql",
	Port: 	3306,
}

var SystemConfig = &system{
	Host: "127.0.0.1",
	Debug: 	false,
	Listen: ":8080",
	Out: "8080",
	SessionSecret: "file-manager",
	HashIDSalt: "file-manager",
	StorageRoot: "D:\\storage",
}
