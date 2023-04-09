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

func AddSpace(s *Space) (bool, error) {
	spaces = append(spaces, s)
	return true, nil
}

func AddTag(space string, tag *Tag) (ok bool, err error) {
	for _, s := range spaces {
		if s.Name == space {
			s.Tags = append(s.Tags, tag)
			return true, nil
		}
	}
	return
}

func ListTags(sp string) ([]*Tag, error) {
	ok, space, err := FindSpace(sp)
	if !ok || err != nil {
		return nil, err
	}
	tags := make([]*Tag, len(space.Tags))
	copy(tags, space.Tags)
	return tags, nil
}

func ListSpaces(owner string) ([]Space, error) {
	sp := make([]Space, 0, len(spaces))
	for _, s := range spaces {
		if s.Owner == owner {
			sp = append(sp, *s)
		}
	}
	return sp, nil
}

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
