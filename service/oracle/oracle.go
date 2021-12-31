package oracle

import (
	"context"
	"encoding/base64"
	"time"

	"github.com/fox-one/pando/core"
	"github.com/fox-one/pando/service/oracle/dirtoracle"
	"github.com/pandodao/blst"
)

func New(
	oracles core.OracleStore,
) core.OracleService {
	return &oracleService{
		oracles: oracles,
	}
}

type oracleService struct {
	oracles core.OracleStore
}

func (s *oracleService) Parse(ctx context.Context, b []byte) (*core.Oracle, error) {
	var p dirtoracle.PriceData

	if err := p.UnmarshalBinary(b); err != nil {
		return nil, nil
	}

	oracle := &core.Oracle{
		CreatedAt: time.Unix(p.Timestamp, 0),
		AssetID:   p.AssetID,
		Current:   p.Price,
	}

	feeds, err := s.oracles.ListFeeds(ctx)
	if err != nil {
		return nil, err
	}

	var (
		pubs      []*blst.PublicKey
		governors []string
	)

	for idx, feed := range feeds {
		if p.Signature.Mask&(0x1<<(idx+1)) != 0 {
			bts, err := base64.StdEncoding.DecodeString(feed.PublicKey)
			if err != nil {
				continue
			}

			pub := blst.PublicKey{}
			if err := pub.FromBytes(bts); err != nil {
				continue
			}

			pubs = append(pubs, &pub)
			governors = append(governors, feed.UserID)
		}
	}

	if passed := blst.AggregatePublicKeys(pubs).Verify(p.Payload(), &p.Signature.Signature); passed {
		oracle.Governors = governors
	}

	return oracle, nil
}
