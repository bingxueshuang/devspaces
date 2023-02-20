package db

import "github.com/bingxueshuang/devspaces/core"

type Tag struct {
	Name     string
	Trapdoor []byte
}

type Space struct {
	Name   string
	Owner  string
	Pubkey []byte
	Tags   []*Tag
}

var spaces []*Space

func FindSpace(space string) (ok bool, s Space, err error) {
	for _, v := range spaces {
		if v.Name == space {
			return true, *v, nil
		}
	}
	return
}

func MessageTag(ciphertext []byte, server *core.SKey, sp Space) (string, error) {
	for _, tag := range sp.Tags {
		ok, err := core.Test(ciphertext, tag.Trapdoor, server)
		if err != nil {
			return "", err
		}
		if ok {
			return tag.Name, nil
		}
	}
	return "others", nil
}
