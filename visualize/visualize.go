package main

import (
	"crypto/ed25519"
	"crypto/sha256"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-echarts/go-echarts/charts"
	"github.com/yoseplee/vrf-go"
	"github.com/yoseplee/vrf-go/sortition"
)

func barChartExample(w http.ResponseWriter, _ *http.Request) {
	nameItems := []string{"aa", "bb", "cc", "dd", "ee"}
	bar := charts.NewBar()
	bar.SetGlobalOptions(charts.TitleOpts{Title: "bar-bar"})
	bar.AddXAxis(nameItems).
		AddYAxis("A", []int{20, 30, 40, 10, 36}).
		AddYAxis("B", []int{35, 14, 25, 60, 44, 23})
	f, err := os.Create("bar.html")
	if err != nil {
		log.Println(err)
	}
	bar.Render(w, f)
}

func loadDataset(amount int) []float64 {
	_, privateKey, err := ed25519.GenerateKey(nil)
	if err != nil {
		log.Panic("failed to generate new key pair")
	}
	publicKey := privateKey.Public().(ed25519.PublicKey)

	var ratios []float64
	for i := 0; i < amount; i++ {
		message := sha256.Sum256([]byte(fmt.Sprintf("Hello%d", i)))
		_, vrfOutput, err := vrf.Prove(publicKey, privateKey, message[:])
		if err != nil {
			log.Println("incorrect calculation of vrf output:", err)
		}
		got := sortition.HashRatio(vrfOutput)
		ratios = append(ratios, got)
	}
	return ratios
}

func getMassByRange(data []float64) []float64 {
	var classifiedData [10]int
	for _, entity := range data {
		switch {
		case entity >= 0.9:
			classifiedData[9]++
		case entity >= 0.8:
			classifiedData[8]++
		case entity >= 0.7:
			classifiedData[7]++
		case entity >= 0.6:
			classifiedData[6]++
		case entity >= 0.5:
			classifiedData[5]++
		case entity >= 0.4:
			classifiedData[4]++
		case entity >= 0.3:
			classifiedData[3]++
		case entity >= 0.2:
			classifiedData[2]++
		case entity >= 0.1:
			classifiedData[1]++
		case entity >= 0.0:
			classifiedData[0]++
		}
	}
	mass := []float64{
		float64(classifiedData[0]) / float64(len(data)),
		float64(classifiedData[1]) / float64(len(data)),
		float64(classifiedData[2]) / float64(len(data)),
		float64(classifiedData[3]) / float64(len(data)),
		float64(classifiedData[4]) / float64(len(data)),
		float64(classifiedData[5]) / float64(len(data)),
		float64(classifiedData[6]) / float64(len(data)),
		float64(classifiedData[7]) / float64(len(data)),
		float64(classifiedData[8]) / float64(len(data)),
		float64(classifiedData[9]) / float64(len(data)),
	}
	return mass
}

func scatter(data []float64, n int) {
	scatter := charts.NewScatter()
	scatter.SetGlobalOptions(charts.TitleOpts{Title: fmt.Sprintf("Probability Mass, N = %d", n)})
	scatter.
		AddXAxis([]string{
			"[0.0, 0.1)",
			"[0.1, 0.2)",
			"[0.2, 0.3)",
			"[0.3, 0.4)",
			"[0.4, 0.5)",
			"[0.5, 0.6)",
			"[0.6, 0.7)",
			"[0.7, 0.8)",
			"[0.8, 0.9)",
			"[0.9, 1.0]",
		}).
		AddYAxis("Ratio from hash(hash/2^(hashlen))", data)
	f, err := os.Create(fmt.Sprintf("./probabilityMass(n=%d).html", n))
	if err != nil {
		log.Println(err)
	}
	scatter.Render(f)
}

func main() {
	amountSet := []int{100, 500, 1000, 1500, 2000, 3000, 4000, 5000, 6000, 8000, 10000, 150000}
	for _, amount := range amountSet {
		data := loadDataset(amount)
		classifiedData := getMassByRange(data)
		scatter(classifiedData, amount)
	}
	// http.HandleFunc("/", barChartExample)
	// http.ListenAndServe(":8081", nil)
}
