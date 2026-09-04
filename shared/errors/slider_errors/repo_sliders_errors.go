package slider_errors

import (
	"github.com/MamangRust/monolith-graphql-ecommerce-shared/errors"
)

var (
	ErrSliderNotFound           = errors.ErrNotFound.WithMessage("slider not found")
	ErrFindAllSliders           = errors.ErrInternal.WithMessage("failed to find all sliders")
	ErrFindActiveSliders        = errors.ErrInternal.WithMessage("failed to find active sliders")
	ErrFindTrashedSliders       = errors.ErrInternal.WithMessage("failed to find trashed sliders")
	ErrFindSliderByID           = errors.ErrInternal.WithMessage("failed to find slider by ID")
	ErrCreateSlider             = errors.ErrInternal.WithMessage("failed to create slider")
	ErrUpdateSlider             = errors.ErrInternal.WithMessage("failed to update slider")
	ErrTrashSlider              = errors.ErrInternal.WithMessage("failed to trash slider")
	ErrRestoreSlider            = errors.ErrInternal.WithMessage("failed to restore slider")
	ErrDeletePermanentSlider    = errors.ErrInternal.WithMessage("failed to permanently delete slider")
	ErrRestoreAllSlider         = errors.ErrInternal.WithMessage("failed to restore all sliders")
	ErrDeleteAllPermanentSlider = errors.ErrInternal.WithMessage("failed to permanently delete all sliders")
)
