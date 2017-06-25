package labeler

import (
	"github.com/bivas/rivi/bot"
	"github.com/bivas/rivi/util"
	"github.com/mitchellh/mapstructure"
)

type action struct {
	rule *rule
}

func (a *action) Apply(config bot.Configuration, meta bot.EventData) {
	apply := a.rule.Label
	if meta.HasLabel(apply) {
		util.Logger.Debug("Skipping label '%s' as it already exists", apply)
	} else {
		meta.AddLabel(apply)
	}
}

type factory struct {
}

func (*factory) BuildAction(config map[string]interface{}) bot.Action {
	item := rule{}
	if e := mapstructure.Decode(config, &item); e != nil {
		panic(e)
	}
	return &action{rule: &item}
}

func init() {
	bot.RegisterAction("labeler", &factory{})
}
