package function

type BlockType string

const (
	BlockTypeProvider   BlockType = "provider"
	BlockTypeResource   BlockType = "resource"
	BlockTypeDataSource BlockType = "datasource"
)
