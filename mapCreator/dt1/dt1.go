package dt1

import (
	"fmt"
	"game/interfaces"
	"game/tools"

	"log"
)

// BlockDataFormat represents the format of the block data
type BlockDataFormat int16

const (
	// BlockFormatRLE specifies the block format is RLE encoded
	BlockFormatRLE BlockDataFormat = 0

	// BlockFormatIsometric specifies the block format isometrically encoded
	BlockFormatIsometric BlockDataFormat = 1
)

const (
	numUnknownHeaderBytes = 260
	knownMajorVersion     = 7
	knownMinorVersion     = 6
	numUnknownTileBytes1  = 4
	numUnknownTileBytes2  = 4
	numUnknownTileBytes3  = 7
	numUnknownTileBytes4  = 12
	numUnknownBlockBytes  = 2
)

// DT1 represents a DT1 file.
type DT1 struct {
	majorVersion  int32
	minorVersion  int32
	numberOfTiles int32
	bodyPosition  int32
	Tiles         []Tile
}

// New creates a new DT1
func New() *DT1 {
	result := &DT1{
		majorVersion: knownMajorVersion,
		minorVersion: knownMinorVersion,
	}

	return result
}

// LoadDT1 loads a DT1 record
//nolint:funlen,gocognit,gocyclo // Can't reduce
func LoadDT1(fileData []byte) (*DT1, error) {
	result := &DT1{}
	br := tools.CreateStreamReader(fileData)
	var err error

	result.majorVersion, err = br.ReadInt32()
	if err != nil {
		return nil, err
	}

	result.minorVersion, err = br.ReadInt32()
	if err != nil {
		return nil, err
	}

	if result.majorVersion != knownMajorVersion || result.minorVersion != knownMinorVersion {
		const fmtErr = "expected to have a version of 7.6, but got %d.%d instead"
		return nil, fmt.Errorf(fmtErr, result.majorVersion, result.minorVersion)
	}

	br.SkipBytes(numUnknownHeaderBytes)

	result.numberOfTiles, err = br.ReadInt32()
	if err != nil {
		return nil, err
	}

	result.bodyPosition, err = br.ReadInt32()
	if err != nil {
		return nil, err
	}

	br.SetPosition(uint64(result.bodyPosition))

	//图片个数
	result.Tiles = make([]Tile, result.numberOfTiles)

	for tileIdx := range result.Tiles {
		tile := Tile{}

		tile.Direction, err = br.ReadInt32()
		if err != nil {
			return nil, err
		}

		tile.RoofHeight, err = br.ReadInt16()
		if err != nil {
			return nil, err
		}

		var matFlagBytes uint16

		matFlagBytes, err = br.ReadUInt16()
		if err != nil {
			return nil, err
		}

		tile.MaterialFlags = NewMaterialFlags(matFlagBytes)

		tile.Height, err = br.ReadInt32()
		if err != nil {
			return nil, err
		}

		tile.Width, err = br.ReadInt32()
		if err != nil {
			return nil, err
		}

		br.SkipBytes(numUnknownTileBytes1)

		tile.Type, err = br.ReadInt32()
		if err != nil {
			return nil, err
		}

		tile.Style, err = br.ReadInt32()
		if err != nil {
			return nil, err
		}

		tile.Sequence, err = br.ReadInt32()
		if err != nil {
			return nil, err
		}

		tile.RarityFrameIndex, err = br.ReadInt32()
		if err != nil {
			return nil, err
		}

		tile.unknown2, err = br.ReadBytes(numUnknownTileBytes2)
		if err != nil {
			return nil, err
		}

		for i := range tile.SubTileFlags {
			var subtileFlagBytes byte

			subtileFlagBytes, err = br.ReadByte()
			if err != nil {
				return nil, err
			}

			tile.SubTileFlags[i] = NewSubTileFlags(subtileFlagBytes)
		}

		br.SkipBytes(numUnknownTileBytes3)

		tile.blockHeaderPointer, err = br.ReadInt32()
		if err != nil {
			return nil, err
		}

		tile.blockHeaderSize, err = br.ReadInt32()
		if err != nil {
			return nil, err
		}

		var numBlocks int32

		numBlocks, err = br.ReadInt32()
		if err != nil {
			return nil, err
		}

		tile.Blocks = make([]Block, numBlocks)

		br.SkipBytes(numUnknownTileBytes4)

		result.Tiles[tileIdx] = tile
	}

	for tileIdx := range result.Tiles {
		tile := &result.Tiles[tileIdx]
		br.SetPosition(uint64(tile.blockHeaderPointer))

		for blockIdx := range tile.Blocks {
			result.Tiles[tileIdx].Blocks[blockIdx].X, err = br.ReadInt16()
			if err != nil {
				return nil, err
			}

			result.Tiles[tileIdx].Blocks[blockIdx].Y, err = br.ReadInt16()
			if err != nil {
				return nil, err
			}

			br.SkipBytes(numUnknownBlockBytes)

			result.Tiles[tileIdx].Blocks[blockIdx].GridX, err = br.ReadByte()
			if err != nil {
				return nil, err
			}

			result.Tiles[tileIdx].Blocks[blockIdx].GridY, err = br.ReadByte()
			if err != nil {
				return nil, err
			}

			result.Tiles[tileIdx].Blocks[blockIdx].format, err = br.ReadInt16()
			if err != nil {
				return nil, err
			}

			result.Tiles[tileIdx].Blocks[blockIdx].Length, err = br.ReadInt32()
			if err != nil {
				return nil, err
			}

			br.SkipBytes(numUnknownBlockBytes)

			result.Tiles[tileIdx].Blocks[blockIdx].FileOffset, err = br.ReadInt32()
			if err != nil {
				return nil, err
			}
		}

		for blockIndex, block := range tile.Blocks {
			br.SetPosition(uint64(tile.blockHeaderPointer + block.FileOffset))

			encodedData, err := br.ReadBytes(int(block.Length))
			if err != nil {
				return nil, err
			}

			tile.Blocks[blockIndex].EncodedData = encodedData
		}
	}

	return result, nil
}

func ImgIndexToRGBA(indexData []byte, palette interfaces.Palette) []byte {
	bytesPerPixel := 4
	colorData := make([]byte, len(indexData)*bytesPerPixel)

	for i := 0; i < len(indexData); i++ {
		// Index zero is hardcoded transparent regardless of palette
		if indexData[i] == 0 {
			continue
		}

		c, err := palette.GetColor(int(indexData[i]))
		if err != nil {
			log.Print(err)
		}

		colorData[i*bytesPerPixel] = c.R()
		colorData[i*bytesPerPixel+1] = c.G()
		colorData[i*bytesPerPixel+2] = c.B()
		colorData[i*bytesPerPixel+3] = c.A()
	}

	return colorData
}
