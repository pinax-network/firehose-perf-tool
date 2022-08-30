package chain

import (
	"fmt"
	"github.com/streamingfast/eth-go"
	pbtransform "github.com/streamingfast/streamingfast-client/pb/sf/ethereum/transform/v1"
	"google.golang.org/protobuf/types/known/anypb"
	"strings"
)

func ParseMultiLogFilter(in []string) (*anypb.Any, error) {

	mf := &pbtransform.MultiLogFilter{}

	for _, filter := range in {
		parts := strings.Split(filter, ":")
		if len(parts) != 2 {
			return nil, fmt.Errorf("option --log-filter-multi must be of type address_hash+address_hash+address_hash:event_sig_hash+event_sig_hash (repeated, separated by comma)")
		}
		var addrs []eth.Address
		for _, a := range strings.Split(parts[0], "+") {
			if a != "" {
				addr := eth.MustNewAddress(a)
				addrs = append(addrs, addr)
			}
		}
		var sigs []eth.Hash
		for _, s := range strings.Split(parts[1], "+") {
			if s != "" {
				sig := eth.MustNewHash(s)
				sigs = append(sigs, sig)
			}
		}

		mf.LogFilters = append(mf.LogFilters, BasicLogFilter(addrs, sigs))
	}

	t, err := anypb.New(mf)
	if err != nil {
		return nil, err
	}
	return t, nil
}

func BasicLogFilter(addrs []eth.Address, sigs []eth.Hash) *pbtransform.LogFilter {
	var addrBytes [][]byte
	var sigsBytes [][]byte

	for _, addr := range addrs {
		b := addr.Bytes()
		addrBytes = append(addrBytes, b)
	}

	for _, sig := range sigs {
		b := sig.Bytes()
		sigsBytes = append(sigsBytes, b)
	}

	return &pbtransform.LogFilter{
		Addresses:       addrBytes,
		EventSignatures: sigsBytes,
	}
}

func ParseMultiCallToFilter(in []string) (*anypb.Any, error) {

	mf := &pbtransform.MultiCallToFilter{}

	for _, filter := range in {
		parts := strings.Split(filter, ":")
		if len(parts) != 2 {
			return nil, fmt.Errorf("option --log-filter-multi must be of type address_hash+address_hash+address_hash:event_sig_hash+event_sig_hash (repeated, separated by comma)")
		}
		var addrs []eth.Address
		for _, a := range strings.Split(parts[0], "+") {
			if a != "" {
				addr := eth.MustNewAddress(a)
				addrs = append(addrs, addr)
			}
		}
		var sigs []eth.Hash
		for _, s := range strings.Split(parts[1], "+") {
			if s != "" {
				sig := eth.MustNewHash(s)
				sigs = append(sigs, sig)
			}
		}

		mf.CallFilters = append(mf.CallFilters, BasicCallToFilter(addrs, sigs))
	}

	t, err := anypb.New(mf)
	if err != nil {
		return nil, err
	}
	return t, nil

}

func BasicCallToFilter(addrs []eth.Address, sigs []eth.Hash) *pbtransform.CallToFilter {
	var addrBytes [][]byte
	var sigsBytes [][]byte

	for _, addr := range addrs {
		b := addr.Bytes()
		addrBytes = append(addrBytes, b)
	}

	for _, sig := range sigs {
		b := sig.Bytes()
		sigsBytes = append(sigsBytes, b)
	}

	return &pbtransform.CallToFilter{
		Addresses:  addrBytes,
		Signatures: sigsBytes,
	}
}
