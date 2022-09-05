package link

import (
	"context"
	"fmt"
	"time"
)

type Stat struct {
	UserIP   string
	PassTime time.Time
}

type Link struct {
	ShortLink int
	LongLink  string
	StatData  []Stat
}

type LinkStore interface {
	CreateShortLink(ctx context.Context, longLink string) (*Link, error)
	UpdateStat(ctx context.Context, shortLink int, ip string) error
	GetStat(ctx context.Context, shortLink int) ([]Stat, error)
}

type Links struct {
	lstore LinkStore
}

func NewLinks(lstore LinkStore) *Links {
	return &Links{
		lstore: lstore,
	}
}

func (ls *Links) CreateShortLink(ctx context.Context, longLink string) (*Link, error) {
	link, err := ls.lstore.CreateShortLink(ctx, longLink)
	if err != nil {
		return nil, fmt.Errorf("create short link error: %w", err)
	}
	return link, nil
}

func (ls *Links) UpdateStat(ctx context.Context, shortLink int, ip string) error {
	err := ls.lstore.UpdateStat(ctx, shortLink, ip)
	if err != nil {
		return fmt.Errorf("update short link error: %w", err)
	}
	return nil
}

func (ls *Links) GetStat(ctx context.Context, shortLink int) ([]Stat, error) {
	stat, err := ls.lstore.GetStat(ctx, shortLink)
	if err != nil {
		return nil, fmt.Errorf("get stat for short link(%v); error: %w", shortLink, err)
	}

	return stat, nil
}
