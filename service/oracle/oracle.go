package oracle

import (
	"context"
	"encoding/base64"
	"errors"
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
		return nil, err
	}

	oracle, err := s.oracles.Find(ctx, p.AssetID)
	if err != nil {
		return nil, err
	}

	if oracle.Threshold == 0 {
		return nil, errors.New("oracle with zero threshold is not allowed to be poke")
	}

	oracle = &core.Oracle{
		CreatedAt: time.Unix(p.Timestamp, 0),
		AssetID:   p.AssetID,
		Current:   p.Price,
	}

	feeds, err := s.oracles.ListFeeds(ctx)
	if err != nil {
		return nil, err
	}

	var pubs []*blst.PublicKey
	for idx, feed := range feeds {
		if p.Signature.Mask&(0x1<<(idx+1)) != 0 {
			bts, err := base64.StdEncoding.DecodeString(feed.PublicKey)
			if err != nil {
				return nil, err
			}

			pub := blst.PublicKey{}
			if err := pub.FromBytes(bts); err != nil {
				return nil, err
			}

			pubs = append(pubs, &pub)
			oracle.Governors = append(oracle.Governors, feed.UserID)
		}
	}

	if passed := int64(len(pubs)) >= oracle.Threshold && blst.AggregatePublicKeys(pubs).Verify(p.Payload(), &p.Signature.Signature); !passed {
		return nil, errors.New("oracle verify not pass")
	}

	return oracle, nil
}
