package diag

import (
	"github.com/xjlgod/nebula-database-diagnostic/pkg/config"
	"log"
)

func Run(conf *config.Config) {
	log.Println(conf.String())
}
