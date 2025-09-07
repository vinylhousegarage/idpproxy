package secretmgr

import (
	"context"
	"errors"

	sm "cloud.google.com/go/secretmanager/apiv1"
	smpb "cloud.google.com/go/secretmanager/apiv1/secretmanagerpb"
)

type Getter struct {
	client *sm.Client
}

func NewGetter(ctx context.Context) (*Getter, error) {
	c, err := sm.NewClient(ctx)
	if err != nil {
		return nil, errors.Join(ErrInitFailed, err)
	}

	return &Getter{client: c}, nil
}

func (g *Getter) Get(ctx context.Context, name string) (string, error) {
	resp, err := g.client.AccessSecretVersion(ctx, &smpb.AccessSecretVersionRequest{
		Name: name,
	})
	if err != nil {
		return "", errors.Join(ErrAccessFailed, err)
	}

	return string(resp.Payload.Data), nil
}

func (g *Getter) Close() error { return g.client.Close() }
