package worker

import "github.com/mauriciofsnts/hermes/internal/storage"

func StartWorker() {
	storage := storage.NewStorage()

	storage.Read()
}
