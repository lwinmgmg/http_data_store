package views

type TempUrlQuery struct {
	File       string `form:"file" binding:"required"`
	FolderName string `form:"folder" binding:"required"`
	ExpireTime int64  `form:"expire" binding:"required"`
	HashKey    string `form:"hkey" binding:"required"`
	FileName   string `form:"filename"`
}
