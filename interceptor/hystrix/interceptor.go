package easyhttphystrix

import (
	"errors"
	"net/http"
	"time"

	"github.com/afex/hystrix-go/hystrix"
	metricCollector "github.com/afex/hystrix-go/hystrix/metric_collector"
	"github.com/afex/hystrix-go/plugins"
	"github.com/soyacen/easyhttp"
)

var err5xx = errors.New("server returned 5xx status code")

func CircuitBreaker(opts ...Option) easyhttp.Interceptor {
	o := defaultOptions()
	o.apply(opts...)
	if o.statsD != nil {
		c, err := plugins.InitializeStatsdCollector(o.statsD)
		if err != nil {
			panic(err)
		}
		metricCollector.Registry.Register(c.NewStatsdCollector)
	}
	hystrix.ConfigureCommand(o.hystrixCommandName, hystrix.CommandConfig{
		Timeout:                durationToInt(o.hystrixTimeout, time.Millisecond),
		MaxConcurrentRequests:  o.maxConcurrentRequests,
		RequestVolumeThreshold: o.requestVolumeThreshold,
		SleepWindow:            o.sleepWindow,
		ErrorPercentThreshold:  o.errorPercentThreshold,
	})
	return func(cli *easyhttp.Client, req *easyhttp.Request, do easyhttp.Doer) (reply *easyhttp.Reply, err error) {
		err = hystrix.Do(o.hystrixCommandName, func() error {
			reply, err = do(cli, req)
			if err != nil {
				return err
			}
			if reply == nil || reply.RawResponse() == nil {
				return err
			}
			if reply.RawResponse().StatusCode >= http.StatusInternalServerError {
				return err5xx
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
