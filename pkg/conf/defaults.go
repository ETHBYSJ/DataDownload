package conf

var DatabaseConfig = &database{
	Type: 	"mysql",
	Port: 	3306,
}

var SystemConfig = &system{
	Debug: 	false,
	Listen: ":8080",
	SessionSecret: "file-manager",
	HashIDSalt: "file-manager",
	StorageRoot: "D:\\storage",
}
