package main

import (
	"flag"
	"log"
	"math"
	"strconv"
	"strings"
)

const psiMemoryFile = "/proc/pressure/memory"

type PsiObserver struct {
	tracker   *Tracker
	reader    Reader
	avgMetric string
}

type PsiValues struct {
	someAvg float64
	fullAvg float64
}

func (o *PsiObserver) Initialize(t *Tracker, r Reader) {
	o.tracker = t
	o.reader = r
	o.process()
}

func (o *PsiObserver) SetFlags() {
	flag.StringVar(&o.avgMetric, "psiAvgMetric", "avg10", "metric to use in PSI observer")
}

func (o *PsiObserver) TimerEvent() {
	o.process()
}

func (o *PsiObserver) parsePsiValue(text string, key string) float64 {
	var result float64 = math.NaN()
	values := strings.Split(text, "=")
	if values[0] == key {
		i, err := strconv.ParseFloat(values[1], 64)
		if err != nil {
			log.Print(err)
		} else {
			result = i
		}
	}
	return result
}

func (o *PsiObserver) getPsiValues() (*PsiValues, error) {
	var values PsiValues

	psiSome, err := o.reader.getTextValue(psiMemoryFile, "some")
	if err != nil {
		return nil, err
	}
	values.someAvg = o.parsePsiValue(psiSome, o.avgMetric)

	psiFull, err := o.reader.getTextValue(psiMemoryFile, "full")
	if err != nil {
		return nil, err
	}
	values.fullAvg = o.parsePsiValue(psiFull, o.avgMetric)

	return &values, nil
}

func (o *PsiObserver) process() {
	const someAvgKey string = "psi_some"
	const fullAvgKey string = "psi_full"

	result := make(map[string]interface{})

	values, err := o.getPsiValues()
	if err == nil {
		result[someAvgKey] = values.someAvg
		result[fullAvgKey] = values.fullAvg
	}
	o.tracker.track(&result)
}
