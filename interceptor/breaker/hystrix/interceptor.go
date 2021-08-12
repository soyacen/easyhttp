package easyhttphystrix

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/afex/hystrix-go/hystrix"
	metricCollector "github.com/afex/hystrix-go/hystrix/metric_collector"
	"github.com/afex/hystrix-go/plugins"

	"github.com/soyacen/easyhttp"
)

func Interceptor(commandName string, opts ...Option) easyhttp.Interceptor {
	o := defaultOptions()
	o.apply(opts...)
	if o.statsD != nil {
		c, err := plugins.InitializeStatsdCollector(o.statsD)
		if err != nil {
			panic(err)
		}
		metricCollector.Registry.Register(c.NewStatsdCollector)
	}
	hystrix.ConfigureCommand(commandName, hystrix.CommandConfig{
		Timeout:                durationToInt(o.timeout, time.Millisecond),
		MaxConcurrentRequests:  o.maxConcurrentRequests,
		RequestVolumeThreshold: o.requestVolumeThreshold,
		SleepWindow:            durationToInt(o.sleepWindow, time.Millisecond),
		ErrorPercentThreshold:  o.errorPercentThreshold,
	})
	return func(cli *easyhttp.Client, req *easyhttp.Request, do easyhttp.Doer) (reply *easyhttp.Reply, err error) {
		err = hystrix.DoC(req.Context(), commandName, func(ctx context.Context) error {
			reply, err = do(cli, req)
			if err != nil {
				return err
			}
			if reply.RawResponse() == nil {
				return errors.New("http response is nil")
			}
			if reply.RawResponse().StatusCode >= http.StatusInternalServerError {
				return fmt.Errorf("server returned %d status code", reply.RawResponse().StatusCode)
			}
			return nil
		}, o.fallbackFunc)
		return reply, err
	}
}

func durationToInt(duration, unit time.Duration) int {
	durationAsNumber := duration / unit

	if int64(durationAsNumber) > int64(maxInt) {
		// Returning max possible value seems like best possible solution here
		// the alternative is to panic as there is no way of returning an error
		// without changing the NewClient API
		return maxInt
	}
	return int(durationAsNumber)
}
