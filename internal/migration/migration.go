package migration

//
//import (
//	"context"
//	"io/fs"
//	"os"
//	"path/filepath"
//
//	"github.com/andreevym/gofermart/internal/logger"
//	//"github.com/andreevym/gofermart/internal/storage/postgres"
//	"go.uber.org/zap"
//)
//
//func MigrateStorage(ctx context.Context, client *postgres.Client) error {
//	return filepath.Walk("migrations", func(path string, info fs.FileInfo, err error) error {
//		if !info.IsDir() {
//			logger.Logger().Info("apply migration", zap.String("path", path))
//			bytes, err := os.ReadFile(path)
//			if err != nil {
//				return err
//			}
//
//			err = client.ApplyMigration(ctx, string(bytes))
//			if err != nil {
//				return err
//			}
//		}
//
//		return nil
//	})
//}
