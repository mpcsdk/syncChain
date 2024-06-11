package event

import "syncChain/internal/conf"

func token2Native(chainId int64, contract string) bool {
	if addr, ok := conf.Config.Token2NativeChain[chainId]; ok {
		if addr == contract {
			return true
		}
	}
	return false
}
func skipToAddr(chainId int64, toaddr string) bool {
	if addrs, ok := conf.Config.SkipToAddrChain[chainId]; ok {
		if _, ok := addrs[toaddr]; ok {
			return true
		}
	}
	return false
}
