package gubrak

import (
	"github.com/DefinitelyMod/gocsv"
	// "github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"testing"
	"time"
)

var (
	dataCSV   = make([]Data, 0)
	tldSearch = "com"
)

type Data struct {
	GlobalRank     string `csv:"GlobalRank"`
	TldRank        string `csv:"TldRank"`
	Domain         string `csv:"Domain"`
	TLD            string `csv:"TLD"`
	RefSubNets     string `csv:"RefSubNets"`
	RefIPs         string `csv:"RefIPs"`
	IDNDomain      string `csv:"IDN_Domain"`
	IDNTLD         string `csv:"IDN_TLD"`
	PrevGlobalRank string `csv:"PrevGlobalRank"`
	PrevTldRank    string `csv:"PrevTldRank"`
	PrevRefSubNets string `csv:"PrevRefSubNets"`
	PrevRefIPs     string `csv:"PrevRefIPs"`
}

func TimeBenchmarker(t *testing.T, start time.Time) {
	duration := time.Since(start)
	t.Log("required time to completed:", duration.Seconds())
}

func TestLoadData(t *testing.T) {
	defer TimeBenchmarker(t, time.Now())

	basePath, _ := os.Getwd()
	fileLocation := filepath.Join(basePath, "majestic_million.csv")

	if _, err := os.Stat(fileLocation); err != nil {
		if os.IsNotExist(err) {

			t.Fatal("To perform performance test, you have to download sample csv data from https://blog.majestic.com/development/majestic-million-csv-daily/")
			return
		}
	}

	f, err := os.OpenFile(fileLocation, os.O_RDWR|os.O_CREATE, os.ModePerm)
	if f != nil {
		defer f.Close()
	}
	if err != nil {
		t.Fatal("error", err.Error())
		return
	}

	rows := make([]Data, 0)

	err = gocsv.UnmarshalFile(f, &rows)
	if err != nil {
		t.Fatal("error", err.Error())
		return
	}

	dataCSV = rows
	t.Log("Loading csv file done. Total", len(dataCSV), "data found")
}

func TestPerformanceFilterUsingForRange(t *testing.T) {
	if len(dataCSV) == 0 {
		t.Skip()
		return
	}

	time.Sleep(time.Second * 2)
	defer TimeBenchmarker(t, time.Now())

	result := make([]Data, 0)

	for _, row := range dataCSV {
		if row.TLD == tldSearch {
			result = append(result, row)
		}
	}

	t.Log("found", len(result))
}

func TestPerformanceFilterUsingOurLibrary(t *testing.T) {
	if len(dataCSV) == 0 {
		t.Skip()
		return
	}

	time.Sleep(time.Second * 2)
	defer TimeBenchmarker(t, time.Now())

	result, err := Filter(dataCSV, func(row Data) bool {
		return row.TLD == tldSearch
	})
	if err != nil {
		t.Fatal("error", err.Error())
		return
	}

	resultParsed := result.([]Data)

	t.Log("found", len(resultParsed))
}
