package image

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4"
	"github.com/mittz/roleplay-webapp-portal/database"
)

type ImageHash struct {
	Name string `json:"name"`
	Hash string `json:"hash"`
}

type BulkImageHashes struct {
	ImageHashes []ImageHash `json:"image_hashes"`
}

func NewImageHash() ImageHash {
	return ImageHash{}
}

func NewBulkImageHashes() BulkImageHashes {
	return BulkImageHashes{}
}

func (b BulkImageHashes) OverrideDatabase() error {
	dbPool := database.GetDatabaseConnection()

	queryTableCreation := `
	DROP TABLE IF EXISTS image_hashes;
	CREATE TABLE image_hashes (
		name character varying(40) NOT NULL,
		hash character varying(100) NOT NULL,
		PRIMARY KEY(name)
	);
	GRANT ALL ON image_hashes TO PUBLIC;
	`

	if _, err := dbPool.Exec(context.Background(), queryTableCreation); err != nil {
		return fmt.Errorf("Table recreation failed: %v\n", err)
	}

	if _, err := dbPool.CopyFrom(
		context.Background(),
		pgx.Identifier{"image_hashes"},
		[]string{"name", "hash"},
		pgx.CopyFromSlice(len(b.ImageHashes), func(i int) ([]interface{}, error) {
			return []interface{}{
				b.ImageHashes[i].Name,
				b.ImageHashes[i].Hash,
			}, nil
		}),
	); err != nil {
		return fmt.Errorf("Bulk insertion failed: %v\n", err)
	}

	return nil
}
