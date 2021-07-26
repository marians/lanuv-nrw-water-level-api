// Package waterlevel provides functions to fetch waterlevel data from
// the https://luadb.lds.nrw.de website, hosted by LANUV NRW,
// for all stations.
package waterlevel

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/marians/lanuv-nrw-water-level-api/pkg/stations"
)

const (
	// Here we download the gzipped tarball
	url = "https://luadb.lds.nrw.de/LUA/hygon/messwerte/messwerte.tar.gz"

	// The file name we expect to be in the tarball
	fileName = "messwerte.txt"

	// The expected first line content
	firstLine = "Name;Datum_zeit;Datum;Messwert"

	// Data file field separator
	separator = ";"

	// Data file time format
	timeFormat = "2006-01-02 15:04:05"

	// Time zone location used for the data
	timeZone = "Europe/Berlin"
)

// StationMeasurement represents one water level measurement on one station.
type StationMeasurement struct {
	StationName string
	Time        time.Time
	Value       *float32
}

type Measurement struct {
	Time  time.Time
	Value *float32
}

// Get raw data
func Fetch() ([]byte, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	contentCompressed, err := io.ReadAll(res.Body)
	defer res.Body.Close()
	if err != nil {
		return nil, err
	}

	reader := bytes.NewReader(contentCompressed)
	gzipReader, err := gzip.NewReader(reader)
	if err != nil {
		return nil, err
	}

	tarFile, err := ioutil.ReadAll(gzipReader)
	if err != nil {
		return nil, err
	}

	var data []byte

	tarReader := bytes.NewReader(tarFile)
	tr := tar.NewReader(tarReader)
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break // End of archive
		} else if err != nil {
			return nil, err
		}

		if hdr.Name != fileName {
			continue
		}

		data, err = ioutil.ReadAll(tr)
		if err != nil {
			return nil, err
		}
	}

	return data, nil
}

func Parse(data []byte) ([]StationMeasurement, error) {
	s := string(data)
	lines := strings.Split(s, "\n")

	if lines[0] != firstLine {
		return nil, fmt.Errorf("input does not match header line requirements: %s", lines[0])
	}

	timeLocation, _ := time.LoadLocation(timeZone)

	measurements := []StationMeasurement{}

	for i, line := range lines {
		if i == 0 {
			continue
		}

		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		fields := strings.Split(line, separator)
		station := fields[0]

		t, err := time.ParseInLocation(timeFormat, fields[1], timeLocation)
		if err != nil {
			return nil, err
		}

		m := StationMeasurement{
			StationName: station,
			Time:        t,
		}

		if fields[3] != "" {
			v, err := strconv.ParseFloat(fields[3], 32)
			if err != nil {
				return nil, err
			}

			v32 := float32(v)
			m.Value = &v32
		}

		measurements = append(measurements, m)

	}

	return measurements, nil
}

func ParseByLocation(input []StationMeasurement) (map[string]Measurement, error) {
	data := make(map[string]Measurement)

	// Per station, take the latest value and write it into the output structure
	for _, ds := range input {
		key := stations.Normalize(ds.StationName)

		if _, ok := data[key]; !ok {
			data[key] = Measurement{
				Time:  ds.Time,
				Value: ds.Value,
			}
		}

		if ds.Time.After(data[key].Time) {
			data[key] = Measurement{
				Time:  ds.Time,
				Value: ds.Value,
			}
		}
	}

	return data, nil
}
