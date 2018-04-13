package beater

import (
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/elastic/beats/libbeat/beat"
	"github.com/elastic/beats/libbeat/common"
	"github.com/golang/protobuf/proto"
	"github.com/golang/snappy"
	"github.com/prometheus/prometheus/prompb"
	"github.com/visheyra/pbeat/config"
)

//PrometheusServer ...
type PrometheusServer struct {
	config config.Config
	events chan beat.Event
}

//NewServer ...
func NewServer() *PrometheusServer {
	return &PrometheusServer{
		config: config.DefaultConfig,
	}
}

//StartServer ...
func (s *PrometheusServer) StartServer(ch chan beat.Event) {
	s.events = ch
	http.HandleFunc(s.config.Path, s.writeHandler)
	http.ListenAndServe(s.config.ListenAddr, nil)
}

//writeHandler ...
func (s *PrometheusServer) writeHandler(w http.ResponseWriter, r *http.Request) {

	// Read prometheus compressed protobuf
	compressed, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// Inflate protobuf
	reqBuf, err := snappy.Decode(nil, compressed)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Unpack protobuf
	var remoteReq prompb.WriteRequest
	if err := proto.Unmarshal(reqBuf, &remoteReq); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Pack data in ES understandable format (aka json) and returns it
	s.toChan(remoteReq)
}

func (s *PrometheusServer) toChan(r prompb.WriteRequest) {

	// Iterate over TS
	for _, sr := range r.GetTimeseries() {
		event := map[string]interface{}{}
		labels := map[string]interface{}{}

		// Find labels for each point of TS
		for _, l := range sr.GetLabels() {
			field := strings.Replace(l.Name, "_", "", -1)
			labels[field] = l.Value
		}
		event["labels"] = labels

		// Add sample to object
		for _, s := range sr.GetSamples() {
			event["value"] = s.Value
			event["timestamp"] = common.Time(time.Unix(0, s.Timestamp*1000000))
		}

		final := beat.Event{
			Timestamp: time.Now(),
			Fields:    event,
		}
		// Send object
		s.events <- final
	}
}
