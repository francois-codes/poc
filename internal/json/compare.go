package json

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/wI2L/jsondiff"
)

func CompareJSONFiles(file1, file2 string) (string, error) {
	// Lire les fichiers
	json1, err := os.ReadFile(file1)
	if err != nil {
		return "", fmt.Errorf("error reading file1: %w", err)
	}
	json2, err := os.ReadFile(file2)
	if err != nil {
		return "", fmt.Errorf("error reading file2: %w", err)
	}

	// Unmarshal en `any` (interface{})
	var v1, v2 any
	if err = json.Unmarshal(json1, &v1); err != nil {
		return "", fmt.Errorf("error parsing file1: %w", err)
	}
	if err = json.Unmarshal(json2, &v2); err != nil {
		return "", fmt.Errorf("error parsing file2: %w", err)
	}

	// Comparer avec jsondiff
	diff, _ := jsondiff.Compare(v1, v2)
	return diff.String(), nil
}

/*
func CompareDBJSONWithFile(ctx context.Context, repo db.VersionRepository, datamodelID int64, filepath string) (string, error) {
	// 1. Lire le JSON depuis la base (champ json.RawMessage)
	dbVersion, err := repo.GetLatestByDatamodelID(ctx, datamodelID)
	if err != nil {
		return "", fmt.Errorf("failed to get version from DB: %w", err)
	}

	// 2. Lire le JSON depuis le fichier
	fileJSON, err := os.ReadFile(filepath)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %w", err)
	}

	// 3. Unmarshal les deux
	var dbData, fileData any
	if err := json.Unmarshal(dbVersion.Json, &dbData); err != nil {
		return "", fmt.Errorf("invalid JSON from DB: %w", err)
	}
	if err := json.Unmarshal(fileJSON, &fileData); err != nil {
		return "", fmt.Errorf("invalid JSON from file: %w", err)
	}

	// 4. Diff avec jsondiff
	opts := jsondiff.DefaultConsoleOptions()
	diff, err := jsondiff.Compare(dbData, fileData, &opts)
	if err != nil {
		return "", fmt.Errorf("error computing diff: %w", err)
	}

	return diff.String(), nil
}

*/
