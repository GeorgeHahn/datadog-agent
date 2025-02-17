// Copyright The OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package translator

import (
	"context"
	"encoding"
	"fmt"

	"github.com/DataDog/datadog-agent/pkg/quantile"
	"github.com/DataDog/datadog-agent/pkg/trace/pb"
)

// MetricDataType is a timeseries-style metric type.
type MetricDataType int

var _ encoding.TextUnmarshaler = (*MetricDataType)(nil)
var _ encoding.TextMarshaler = (MetricDataType)(Gauge)

const (
	// Gauge is the Datadog Gauge metric type.
	Gauge MetricDataType = iota
	// Count is the Datadog Count metric type.
	Count
)

// UnmarshalText implements encoding.TextUnmarshaler.
func (t *MetricDataType) UnmarshalText(text []byte) error {
	switch string(text) {
	case "gauge":
		*t = Gauge
	case "count":
		*t = Count
	default:
		return fmt.Errorf("invalid metric data type %q", text)
	}
	return nil
}

// MarshalText implements encoding.TextMarshaler.
func (t MetricDataType) MarshalText() ([]byte, error) {
	switch t {
	case Gauge:
		return []byte("gauge"), nil
	case Count:
		return []byte("count"), nil
	}

	return nil, fmt.Errorf("invalid metric data type %d", t)
}

// TimeSeriesConsumer is timeseries consumer.
type TimeSeriesConsumer interface {
	// ConsumeTimeSeries consumes a timeseries-style metric.
	ConsumeTimeSeries(
		ctx context.Context,
		dimensions *Dimensions,
		typ MetricDataType,
		timestamp uint64,
		value float64,
	)
}

// SketchConsumer is a pkg/quantile sketch consumer.
type SketchConsumer interface {
	// ConsumeSketch consumes a pkg/quantile-style sketch.
	ConsumeSketch(
		ctx context.Context,
		dimensions *Dimensions,
		timestamp uint64,
		sketch *quantile.Sketch,
	)
}

// Consumer is a metrics consumer.
type Consumer interface {
	TimeSeriesConsumer
	SketchConsumer
	APMStatsConsumer
}

// APMStatsConsumer implementations are able to consume APM Stats generated by
// a Translator.
type APMStatsConsumer interface {
	// ConsumeAPMStats consumes the given StatsPayload.
	ConsumeAPMStats(pb.ClientStatsPayload)
}

// HostConsumer is a hostname consumer.
// It is an optional interface that can be implemented by a Consumer.
type HostConsumer interface {
	// ConsumeHost consumes a hostname.
	ConsumeHost(host string)
}

// TagsConsumer is a tags consumer.
// It is an optional interface that can be implemented by a Consumer.
// Consumed tags are used for running metrics, and should represent
// some resource running a Collector (e.g. Fargate task).
type TagsConsumer interface {
	// ConsumeTag consumes a tag
	ConsumeTag(tag string)
}
