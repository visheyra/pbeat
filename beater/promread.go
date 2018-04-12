package beater

import (
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/elastic/beats/libbeat/common"
	"github.com/golang/protobuf/proto"
	"github.com/golang/snappy"
	"github.com/prometheus/prometheus/prompb"
	"github.com/visheyra/pbeat/config"
)

//PrometheusServer ...
type PrometheusServer struct {
	config config.Config
	events chan common.MapStr
}

//NewServer ...
func NewServer() *PrometheusServer {
	return &PrometheusServer{
		config: config.DefaultConfig,
	}
}

//StartServer ...
func (s *PrometheusServer) StartServer(ch chan common.MapStr) {
	s.events = ch
	http.HandleFunc(s.config.Path, s.writeHandler)
	http.ListenAndServe(s.config.ListenAddr, nil)
}

//writeHandler ...
func (s *PrometheusServer) writeHandler(w http.ResponseWriter, r *http.Request) {
	if compressed, err := ioutil.ReadAll(r.Body); r != nil {
		http.Error(w, err.Error(), 500)
		return
	} else {
		reqBuf, err := snappy.Decode(nil, compressed)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		var remoteReq prompb.WriteRequest
		if err := proto.Unmarshal(reqBuf, &remoteReq); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		} else {
			s.toChan(remoteReq)
		}
	}
}

func (s *PrometheusServer) toChan(r prompb.WriteRequest) {
	for _, sr := range r.GetTimeseries() {
		event := map[string]interface{}{}
		labels := map[string]interface{}{}

		for _, l := range sr.GetLabels() {
			field := strings.Replace(l.Name, "_", "", -1)
			labels[field] = l.Value
		}
		event["labels"] = labels

		for _, s := range sr.GetSamples() {
			event["value"] = s.Value
			event["timestamp"] = common.Time(time.Unix(0, s.Timestamp*1000000))
		}

		s.events <- event
	}
}
