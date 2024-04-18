package service

import (
	"context"
	"javacode/internal/config"

	"github.com/pkg/errors"
)

func (s *servImpl) ChangingSum(ctx context.Context, data *TrasactionInfo) (int, error) {

	if data.Role == config.HeaderInc {
		sum, err := s.rep.SumInc(ctx, data.Sum)
		if err != nil {
			return 0, errors.Wrap(err, "occurred error ChangingSum")
		}
		return sum, nil
	}
	if data.Role == config.HeaderDec {
		sum, err := s.rep.SumDec(ctx, data.Sum)
		if err != nil {
			return 0, errors.Wrap(err, "occurred error ChangingSum")
		}
		return sum, nil
	}

	return 0, errors.Wrap(errors.New("bad Role"), "occurred error ChangingSum")
}
