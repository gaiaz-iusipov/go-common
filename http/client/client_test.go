package httpclient_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/alecthomas/assert/v2"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
	"go.opentelemetry.io/otel/sdk/metric/metricdata/metricdatatest"

	httpclient "github.com/gaiaz-iusipov/go-common/http/client"
)

func TestNew(t *testing.T) {
	t.Run("basic", func(t *testing.T) {
		ctx := t.Context()

		server := httptest.NewTestServer(t, http.HandlerFunc(func(rw http.ResponseWriter, _ *http.Request) {
			_, err := rw.Write([]byte("Hello, world!"))
			assert.NoError(t, err)
		}))

		reader := metric.NewManualReader()

		client := httpclient.New(
			server.Client().Transport,
			httpclient.WithOTELOptions(
				otelhttp.WithMeterProvider(metric.NewMeterProvider(metric.WithReader(reader))),
			),
		)

		req, err := http.NewRequestWithContext(ctx, http.MethodGet, server.URL, strings.NewReader("john"))
		assert.NoError(t, err)

		resp, err := client.Do(req)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Equal(t, 13, resp.ContentLength)
		assert.NoError(t, resp.Body.Close())

		rm := metricdata.ResourceMetrics{}
		assert.NoError(t, reader.Collect(ctx, &rm))
		assert.Equal(t, 1, len(rm.ScopeMetrics))

		attrs := attribute.NewSet(
			attribute.String("http.request.method", http.MethodGet),
			attribute.Int("http.response.status_code", http.StatusOK),
			attribute.String("server.address", req.Host),
			attribute.String("url.scheme", "http"),
			attribute.String("network.protocol.name", "http"),
			attribute.String("network.protocol.version", "1.1"),
		)
		assertClientScopeMetrics(t, rm.ScopeMetrics[0], attrs)
	})

	t.Run("with request name", func(t *testing.T) {
		ctx := t.Context()
		server := httptest.NewTestServer(t, http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {}))
		reader := metric.NewManualReader()
		client := httpclient.New(
			server.Client().Transport,
			httpclient.WithOTELOptions(
				otelhttp.WithMeterProvider(metric.NewMeterProvider(metric.WithReader(reader))),
			),
		)

		req, err := http.NewRequestWithContext(
			httpclient.WithRequestName(ctx, "test_request_name"),
			http.MethodGet,
			server.URL,
			strings.NewReader("john"),
		)
		assert.NoError(t, err)

		resp, err := client.Do(req)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Zero(t, resp.ContentLength)
		assert.NoError(t, resp.Body.Close())

		rm := &metricdata.ResourceMetrics{}
		assert.NoError(t, reader.Collect(ctx, rm))
		assert.Equal(t, 1, len(rm.ScopeMetrics))

		attrs := attribute.NewSet(
			attribute.String("http.request.method", http.MethodGet),
			attribute.String("http.request_name", "test_request_name"),
			attribute.Int("http.response.status_code", http.StatusOK),
			attribute.String("server.address", req.Host),
			attribute.String("url.scheme", "http"),
			attribute.String("network.protocol.name", "http"),
			attribute.String("network.protocol.version", "1.1"),
		)
		assertClientScopeMetrics(t, rm.ScopeMetrics[0], attrs)
	})
}

func assertClientScopeMetrics(t *testing.T, sm metricdata.ScopeMetrics, attrs attribute.Set) {
	assert.Equal(t, instrumentation.Scope{
		Name:    "go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp",
		Version: otelhttp.Version,
	}, sm.Scope)

	assert.Equal(t, 2, len(sm.Metrics))

	want := metricdata.ScopeMetrics{
		Scope: instrumentation.Scope{
			Name:    otelhttp.ScopeName,
			Version: otelhttp.Version,
		},
		Metrics: []metricdata.Metrics{
			{
				Name:        "http.client.request.body.size",
				Description: "Size of HTTP client request bodies.",
				Unit:        "By",
				Data: metricdata.Histogram[int64]{
					Temporality: metricdata.CumulativeTemporality,
					DataPoints:  []metricdata.HistogramDataPoint[int64]{{Attributes: attrs}},
				},
			},
			{
				Name:        "http.client.request.duration",
				Description: "Duration of HTTP client requests.",
				Unit:        "s",
				Data: metricdata.Histogram[float64]{
					Temporality: metricdata.CumulativeTemporality,
					DataPoints:  []metricdata.HistogramDataPoint[float64]{{Attributes: attrs}},
				},
			},
		},
	}
	metricdatatest.AssertEqual(t, want, sm, metricdatatest.IgnoreTimestamp(), metricdatatest.IgnoreValue(), metricdatatest.IgnoreExemplars())
}
