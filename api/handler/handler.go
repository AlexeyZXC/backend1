package handler

import (
	"context"
	"fmt"

	"github.com/AlexeyZXC/backend1/tree/CourseProject/app/repo/link"
)

type Handlers struct {
	hls *link.Links
}

func NewHandlers(ls *link.Links) *Handlers {
	r := &Handlers{
		hls: ls,
	}
	return r
}

func (ls *Handlers) CreateShortLink(ctx context.Context, longLink string) (*link.Link, error) {
	l, err := ls.hls.CreateShortLink(ctx, longLink)
	if err != nil {
		return nil, fmt.Errorf("err while creating short link: %w", err)
	}

	return l, nil
}

func (ls *Handlers) UpdateStat(ctx context.Context, shortLink int, ip string) error {
	err := ls.hls.UpdateStat(ctx, shortLink, ip)
	if err != nil {
		return fmt.Errorf("err while update stat for shortlink(%v), err: %w", shortLink, err)
	}

	return nil
}

func (ls *Handlers) GetStat(ctx context.Context, shortLink int) ([]link.Stat, error) {
	stat, err := ls.hls.GetStat(ctx, shortLink)
	if err != nil {
		return nil, fmt.Errorf("err while getstat for short link(%v), err: %w", shortLink, err)
	}

	return stat, nil
}
