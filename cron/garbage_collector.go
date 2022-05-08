package cron

import (
	"fmt"
	"os"
	"path"

	"github.com/lwinmgmg/http_data_store/environ"
)

var (
	GarbageChannel chan string      = make(chan string, 1000)
	env            *environ.Environ = environ.NewEnviron()
)

func GetGarbageChannelWriter() chan<- string {
	return GarbageChannel
}

func StartGarbageCollector() {
	go func(ch <-chan string) {
		folders_list := make([]string, 0, 4)
		folders_list = append(folders_list, env.HDS_FIRST_DIRS...)
		folders_list = append(folders_list, env.HDS_LAST_DIRS...)
		if env.HDS_GARBAGE_COLLECTOR_BK_ENABLE {
			folders_list = append(folders_list, env.HDS_BACKUP_DIR)
		}
		for vFile := range GarbageChannel {
			for _, vPath := range folders_list {
				if vPath != "" {
					real_path := path.Join(vPath, vFile)
					err := os.Remove(real_path)
					if err != nil {
						fmt.Printf("Error on delete file [%v] : %v\n", real_path, err)
					}
				}
			}
			fmt.Println(vFile, "is deleted successfully")
		}
	}(GarbageChannel)
}
