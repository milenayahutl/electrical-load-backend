package repository

import (
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

//подключения к хранилищам, которые используются приложением
type Repository struct {
	db *gorm.DB

	minio           *minio.Client
	minioBucketName string
}

//настройки подключения к PostgreSQL и MinIO
type RepositorySettings struct {
	PostgresDSN string

	MinioEndpoint   string
	MinioAccessKey  string
	MinioSecretKey  string
	MinioBucketName string
}

func New(
	settings *RepositorySettings,
) (*Repository, error) {

	//PostgreSQL
	db, err := gorm.Open(
		postgres.Open(settings.PostgresDSN),
		&gorm.Config{},
	)

	if err != nil {
		return nil, err
	}

	//MinIO.
	minioClient, err := minio.New(
		settings.MinioEndpoint,
		&minio.Options{
			Creds: credentials.NewStaticV4(
				settings.MinioAccessKey,
				settings.MinioSecretKey,
				"",
			),
			Secure: false,
		},
	)

	if err != nil {
		return nil, err
	}

	return &Repository{
		db:              db,
		minio:           minioClient,
		minioBucketName: settings.MinioBucketName,
	}, nil
}
