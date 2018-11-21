package environment

type Environment struct {
	URLs string 	`env:"COUCHDB_CLEANER_URLS"`
	CleanInterval	int	`env:"COUCHDB_CLEANER_COMPACT_INTERVAL_MS"`
}