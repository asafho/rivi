package runner

import (
	"net/http"

	"github.com/bivas/rivi/config/client"
	"github.com/bivas/rivi/runner/types"
	"github.com/bivas/rivi/types/builder"
	"github.com/bivas/rivi/util/log"
	"github.com/spf13/viper"
)

type hookListener struct {
	conf  client.ClientConfig
	queue types.HookListenerQueue

	logger log.Logger
}

func (h *hookListener) HandleEvent(r *http.Request) *HandledEventResult {
	data, ok := builder.BuildFromHook(h.conf, r)
	if !ok {
		return &HandledEventResult{Message: "Skipping hook processing"}
	}
	h.queue.Enqueue(types.NewMessage(h.conf, data))
	return &HandledEventResult{Message: "Processing " + data.GetShortName()}
}

func NewHookListener(clientConfiguration string) (Runner, error) {
	var conf client.ClientConfig
	if clientConfiguration == "" {
		conf = client.NewClientConfig(viper.New())
	} else {
		conf = client.NewClientConfigFromFile(clientConfiguration)
	}
	return NewHookListenerWithClientConfig(conf), nil
}

func NewHookListenerWithClientConfig(config client.ClientConfig) Runner {
	logger := runnerLog.Get("hook.listener")
	return &hookListener{
		conf:   config,
		queue:  CreateHookListenerQueue(),
		logger: logger,
	}
}