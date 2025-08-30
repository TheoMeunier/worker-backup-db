package databases

import (
	services "backup-dump-sql/services/s3"
	"backup-dump-sql/services/utils"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/JCoupalK/go-pgdump"
	"github.com/joho/godotenv"
)

type ServicePostgresqlImpl struct {
	database_name string
	host          string
	port          string
	username      string
	password      string
	sslmode       string
}

func ServiceImplPostgres() (*ServicePostgresqlImpl, error) {
	_ = godotenv.Load()

	return &ServicePostgresqlImpl{
		database_name: os.Getenv("DATABASE_NAME"),
		host:          os.Getenv("DATABASE_HOST"),
		port:          os.Getenv("DATABASE_PORT"),
		username:      os.Getenv("DATABASE_USER"),
		password:      os.Getenv("DATABASE_PASSWORD"),
		sslmode:       os.Getenv("POSTGRES_SSLMODE"),
	}, nil
}

func (s *ServicePostgresqlImpl) BackupPostgreSQL() error {
	port, _ := strconv.Atoi(s.port)

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		s.host, port, s.username, s.password, s.database_name)

	dumper := pgdump.NewDumper(psqlInfo, 50)

	// Créer un fichier temporaire
	currentTime := time.Now()
	tempFilename := filepath.Join(os.TempDir(), fmt.Sprintf("%s-%s.sql",
		s.database_name, currentTime.Format("20060102T150405")))

	fmt.Printf("Creating temporary backup file: %s\n", tempFilename)

	// Dump vers le fichier temporaire
	if err := dumper.DumpDatabase(tempFilename, &pgdump.TableOptions{
		TableSuffix: "",
		TablePrefix: "",
		Schema:      "",
	}); err != nil {
		return fmt.Errorf("error dumping database: %v", err)
	}

	defer os.Remove(tempFilename)

	// Lire le fichier
	fileData, err := os.ReadFile(tempFilename)
	if err != nil {
		return fmt.Errorf("error reading backup file: %v", err)
	}

	if len(fileData) == 0 {
		return fmt.Errorf("backup file is empty")
	}

	fmt.Printf("Original backup size: %d bytes\n", len(fileData))

	// Compresser avec gzip
	compressedData, err := utils.CompressGzip(fileData)
	if err != nil {
		return fmt.Errorf("error compressing backup: %v", err)
	}

	fmt.Printf("Compressed size: %d bytes (compression: %.1f%%)\n",
		len(compressedData), float64(len(compressedData))/float64(len(fileData))*100)

	// Créer le service S3
	s3Service, err := services.NewServiceS3WithR2()
	if err != nil {
		return fmt.Errorf("error creating S3 service: %v", err)
	}

	// Upload la version compressée
	if err := s3Service.UploadToS3(compressedData, s.database_name); err != nil {
		return fmt.Errorf("error uploading to S3: %v", err)
	}

	fmt.Printf("Compressed backup successfully uploaded to S3 for database: %s\n", s.database_name)
	return nil
}
