package src_flag

import (
	"congo"
	"flag"
	"os"
)

func New() congo.Source {
	return FromFlagSet(flag.CommandLine, standardLoader)
}

type ArgLoader func() []string

func FromFlagSet(set *flag.FlagSet, loader ArgLoader) congo.Source {
	return &source{set, loader}
}

func standardLoader() []string {
	return os.Args[1:]
}

type source struct {
	set *flag.FlagSet
	ArgLoader
}

func (s *source) Init(settings map[string]*congo.Setting) error {
	for key, setting := range settings {
		s.set.Var(setting.Value, key, setting.Usage)
	}
	return nil
}

func (s *source) Load(settings map[string]*congo.Setting) error {
	s.set.Parse(s.ArgLoader())
	return nil
}
