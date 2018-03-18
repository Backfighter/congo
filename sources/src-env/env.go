package src_env

import "congo"

func New() congo.Source {
	return &source{}
}

type source struct {
}

func (s *source) Init(map[string]*congo.Setting) error {
	// Do nothing
	return nil
}

func (s *source) Load(settings map[string]*congo.Setting) error {

	return nil
}
