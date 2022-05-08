package environ

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

var env *Environ

type Environ struct {
	//server
	HDS_PORT   int
	HDS_WORKER int

	//database
	HDS_DB_HOST      string
	HDS_DB_PORT      int
	HDS_DB_USER      string
	HDS_DB_PASSWORD  string
	HDS_DB_NAME      string
	HDS_TABLE_PREFIX string

	//logger
	HDS_LOG_FILENAME          string
	HDS_LOG_CONSOLE_ONLY      bool
	HDS_LOG_ROTATION_COUNT    int
	HDS_LOG_MAX_PENDING_COUNT int
	HDS_LOG_LEVEL             string

	//just for initial
	HDS_TOKEN_KEY             string
	HDS_TOKEN_AUDIENCE        string
	HDS_TOKEN_CHANGE_ALWAYS   bool
	HDS_TOKEN_DEFAULT_TIMEOUT int // in second

	HDS_ADMIN_PASSWORD string

	//filestore
	HDS_FIRST_DIRS                  []string
	HDS_LAST_DIRS                   []string
	HDS_BACKUP_DIR                  string
	HDS_GARBAGE_COLLECTOR_BK_ENABLE bool
}

func NewEnviron() *Environ {
	LoadEnvFile()
	bkdir := getEnviron("HDS_BACKUP_DIR", "")
	if bkdir == "" {
		panic("HDS_BACKUP_DIR can't be empty!")
	}
	return &Environ{
		HDS_PORT:   intParser(getEnviron("HDS_PORT", 8000), 8000),
		HDS_WORKER: intParser(getEnviron("HDS_WORKER", 2), 2),

		HDS_DB_HOST:      getEnviron("HDS_DB_HOST", "localhost"),
		HDS_DB_PORT:      intParser(getEnviron("HDS_DB_PORT", 3306), 3306),
		HDS_DB_USER:      getEnviron("HDS_DB_USER", "root"),
		HDS_DB_PASSWORD:  getEnviron("HDS_DB_PASSWORD", "password"),
		HDS_DB_NAME:      getEnviron("HDS_DB_NAME", "hds"),
		HDS_TABLE_PREFIX: getEnviron("HDS_TABLE_PREFIX", "hds"),

		HDS_LOG_FILENAME:          getEnviron("HDS_LOG_FILENAME", "/tmp/hds.log"),
		HDS_LOG_CONSOLE_ONLY:      boolParser(getEnviron("HDS_LOG_CONSOLE_ONLY", true), true),
		HDS_LOG_ROTATION_COUNT:    intParser(getEnviron("HDS_LOG_ROTATION_COUNT", 0), 0),
		HDS_LOG_MAX_PENDING_COUNT: intParser(getEnviron("HDS_LOG_MAX_PENDING_COUNT", 1000), 1000),
		HDS_LOG_LEVEL:             getEnviron("HDS_LOG_LEVEL", "INFO"),

		HDS_TOKEN_KEY:             getEnviron("HDS_TOKEN_KEY", "secret"),
		HDS_TOKEN_AUDIENCE:        getEnviron("HDS_TOKEN_AUDIENCE", "admin"),
		HDS_TOKEN_CHANGE_ALWAYS:   boolParser(getEnviron("HDS_TOKEN_CHANGE_ALWAYS", false), false),
		HDS_TOKEN_DEFAULT_TIMEOUT: intParser(getEnviron("HDS_TOKEN_DEFAULT_TIMEOUT", 24*60*60), 24*60*60),
		HDS_ADMIN_PASSWORD:        getEnviron("HDS_ADMIN_PASSWORD", "admin"),

		HDS_FIRST_DIRS:                  stringListParser(getEnviron("HDS_FIRST_DIRS", ""), "HDS_FIRST_DIRS"),
		HDS_LAST_DIRS:                   stringListParser(getEnviron("HDS_LAST_DIRS", ""), "HDS_LAST_DIRS"),
		HDS_BACKUP_DIR:                  bkdir,
		HDS_GARBAGE_COLLECTOR_BK_ENABLE: boolParser(getEnviron("HDS_GARBAGE_COLLECTOR_BK_ENABLE", true), true),
	}
}

func GetAllEnv() *Environ {
	if env == nil {
		env = NewEnviron()
	}
	return env
}

func getEnviron(key string, default_val any) string {
	if envData, ok := os.LookupEnv(key); ok {
		return envData
	}
	return fmt.Sprintf("%v", default_val)
}

func stringListParser(input string, name string) []string {
	if input == "" {
		panic(fmt.Sprintf("%v can't be empty!", name))
	}
	stringList := strings.Split(input, ",")
	output := make([]string, len(stringList))
	for k, v := range stringList {
		str := strings.TrimSpace(v)
		str = strings.Trim(str, "\"")
		output[k] = str
	}
	return output
}

func boolParser(input string, default_val bool, mesg ...interface{}) bool {
	output, err := strconv.ParseBool(input)
	if err != nil {
		fmt.Println(mesg...)
		output = default_val
	}
	return output
}

func intParser(input string, default_val int, mesg ...interface{}) int {
	output, err := strconv.Atoi(input)
	if err != nil {
		fmt.Println(mesg...)
		output = default_val
	}
	return output
}

func checkEqual(input string) int {
	count := 0
	for _, v := range input {
		if v == '=' {
			count += 1
		}
	}
	return count
}

func splitWithFirstEqual(input string) []string {
	output := []string{}
	tmp_str := []rune{}
	count := 0
	for _, v := range input {
		if v == '=' && count == 0 {
			output = append(output, string(tmp_str))
			count += 1
			tmp_str = []rune{}
			continue
		}
		tmp_str = append(tmp_str, v)
	}
	output = append(output, string(tmp_str))
	return output
}

func setEnvByStrList(strList []string) {
	newStrList := []string{}
	for _, v := range strList {
		str := strings.TrimSpace(v)
		str = strings.Trim(str, "\"")
		newStrList = append(newStrList, str)
	}
	os.Setenv(newStrList[0], newStrList[1])
}

func LoadEnvFile() {
	envPath := getEnviron("HDS_ENV_PATH", ".env") //default to .env file
	reader, err := os.Open(envPath)
	if err == nil {
		defer reader.Close()
		bufioScanner := bufio.NewScanner(reader)
		for bufioScanner.Scan() {
			line := strings.TrimSpace(bufioScanner.Text())
			if len(line) <= 0 || line[0] == '#' {
				continue
			}
			fmt.Println(line)
			var strList []string = []string{}
			if checkEqual(line) == 1 {
				strList = strings.Split(line, "=")
			} else if checkEqual(line) > 1 {
				strList = splitWithFirstEqual(line)
			}
			setEnvByStrList(strList)
		}
	} else {
		fmt.Println(envPath, "does not exist")
	}
}
