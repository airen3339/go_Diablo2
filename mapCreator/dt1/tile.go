package dt1

// Tile is a representation of a map tile
type Tile struct {
	unknown2           []byte
	Direction          int32
	RoofHeight         int16
	MaterialFlags      MaterialFlags
	Height             int32
	Width              int32
	Type               int32
	Style              int32
	Sequence           int32
	RarityFrameIndex   int32
	SubTileFlags       [25]SubTileFlags
	blockHeaderPointer int32
	blockHeaderSize    int32
	Blocks             []Block
}

// GetSubTileFlags returns the tile flags for the given subtile
func (t *Tile) GetSubTileFlags(x, y int) bool {
	var subtileLookup = [5][5]int{
		{20, 21, 22, 23, 24},
		{15, 16, 17, 18, 19},
		{10, 11, 12, 13, 14},
		{5, 6, 7, 8, 9},
		{0, 1, 2, 3, 4},
	}

	return t.SubTileFlags[subtileLookup[y][x]].BlockWalk
}
