// The package provides core entity link and interfaces to work with it.
// Business logic.
package link

import (
	"context"
	"fmt"
)

// Stat contains statistic data for Link
type Stat struct {
	UserIP   string
	PassTime string
}

// Link contains the data for link and its statistics
type Link struct {
	ShortLink string
	LongLink  string
	StatData  []Stat
}

// LinkStore is interface that must be implemented by a real package used for data storing
type LinkStore interface {
	CreateShortLink(ctx context.Context, longLink string) (*Link, error)
	UpdateStat(ctx context.Context, shortLink int, ip string) error
	GetStat(ctx context.Context, shortLink int) ([]Stat, error)
	GetLongLink(ctx context.Context, shortLink int) (Link, error)
}

// Links wraps data storing package
type Links struct {
	lstore LinkStore
}

// NewLinks returns new Links object with embedded the interface
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

func (ls *Links) GetLongLink(ctx context.Context, shortLink int) (Link, error) {
	ll, err := ls.lstore.GetLongLink(ctx, shortLink)
	if err != nil {
		return Link{}, fmt.Errorf("get stat for short link(%v); error: %w", shortLink, err)
	}
	return ll, nil
}
