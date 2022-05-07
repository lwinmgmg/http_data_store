package views

type FolderCreate struct {
	Name string `json:"name" binding:"required"`
	Key  string `json:"key"`
}

type FolderRead struct {
	ID     uint   `json:"id"`
	Name   string `json:"name"`
	Key    string `json:"key"`
	UserID uint   `json:"user_id"`
}

type FolderUpdate struct {
	Name *string
	Key  *string
}

func (folder *FolderUpdate) ToMap() map[string]any {
	resMap := make(map[string]any, 2)
	if folder.Name != nil {
		resMap["name"] = *folder.Name
	}
	if folder.Key != nil {
		resMap["key"] = *folder.Key
	}
	return resMap
}
