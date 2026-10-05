package repository

import (
	"context"
	"fmt"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
)

// сохраняет полученный файл в MinIO и возвращает сгенерированное имя объекта
func (r *Repository) UploadMedia(
	header *multipart.FileHeader,
) (string, error) {

	file, err := header.Open()
	if err != nil {
		return "", fmt.Errorf(
			"ошибка открытия файла: %w",
			err,
		)
	}

	defer file.Close()

	buffer := make([]byte, 512)

	readBytes, err := file.Read(buffer)
	if err != nil {
		return "", fmt.Errorf(
			"ошибка чтения файла: %w",
			err,
		)
	}

	contentType :=
		http.DetectContentType(
			buffer[:readBytes],
		)

	_, err = file.Seek(0, 0)

	if err != nil {
		return "", fmt.Errorf(
			"ошибка обработки файла: %w",
			err,
		)
	}

	extension :=
		strings.ToLower(
			filepath.Ext(header.Filename),
		)

	//генерируем новое имя латиницей
	filename :=
		uuid.New().String() + extension

	_, err = r.minio.PutObject(
		context.Background(),
		r.minioBucketName,
		filename,
		file,
		header.Size,
		minio.PutObjectOptions{
			ContentType: contentType,
		},
	)

	if err != nil {
		return "", fmt.Errorf(
			"ошибка загрузки файла в MinIO: %w",
			err,
		)
	}

	return filename, nil
}

// удаляет объект из MinIO, если создание услуги завершилось ошибкой
func (r *Repository) DeleteMedia(
	filename string,
) error {

	return r.minio.RemoveObject(
		context.Background(),
		r.minioBucketName,
		filename,
		minio.RemoveObjectOptions{},
	)
}
