package gui

import (
	"fmt"
	"pdfprinter/domain"
	"pdfprinter/domain/models/application"
	"pdfprinter/reductor"
)

func GetModel() (*application.Application, error) {
	modelReductor, err := reductor.Instance().Model(domain.Application)
	if err != nil {
		return nil, fmt.Errorf("failed to get model from reductor: %w", err)
	}
	model, ok := modelReductor.(*application.Application)
	if !ok {
		return nil, fmt.Errorf("model is not of type *application.Application")
	}
	return model, nil
}

func SyncModel(app domain.Apper) (*application.Application, error) {
	modelReductor, err := reductor.Instance().Model(domain.Application)
	if err != nil {
		return nil, fmt.Errorf("failed to get model from reductor: %w", err)
	}
	model, ok := modelReductor.(*application.Application)
	if !ok {
		return nil, fmt.Errorf("model is not of type *application.Application")
	}
	if err := model.SyncToStore(app); err != nil {
		return nil, fmt.Errorf("failed to sync model from reductor: %w", err)
	}
	return model, nil
}
