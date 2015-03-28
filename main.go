package main

import (
	"time"

	ui "github.com/gizak/termui"
	tm "github.com/nsf/termbox-go"

	"github.com/Depado/test-termui/utils"
)

func GetCPUPercentage() float64 {
	idle0, total0 := utils.GetCPUSample()
	time.Sleep(500 * time.Millisecond)
	idle1, total1 := utils.GetCPUSample()

	idleTicks := float64(idle1 - idle0)
	totalTicks := float64(total1 - total0)
	return 100 * (totalTicks - idleTicks) / totalTicks
}

func GetRAMPercentage() float64 {
	time.Sleep(500 * time.Millisecond)
	total, _, available := utils.GetRAMUsage()
	return float64(100 - 100*available/total)
}

func UpdateGenericGauge(g *ui.Gauge, updateFunc func() float64) {
	for {
		g.Percent = int(updateFunc())
		time.Sleep(500 * time.Millisecond)
	}
}

func UpdateGenericChart(l *ui.LineChart, updateFunc func() float64) {
	data := make([]float64, 150)
	for {
		for i := len(data) - 1; i > 0; i-- {
			data[i] = data[i-1]
		}
		data[0] = updateFunc()
		l.Data = data
		time.Sleep(500 * time.Millisecond)
	}
}

func main() {
	// Init
	err := ui.Init()
	if err != nil {
		panic(err)
	}
	defer ui.Close()

	// Theme Setting
	ui.UseTheme("helloworld")

	// Setup the CPU Gauge
	cpuGauge := ui.NewGauge()
	cpuGauge.Height = 2
	cpuGauge.BarColor = ui.ColorRed
	cpuGauge.HasBorder = false
	cpuGauge.PaddingBottom = 1
	go UpdateGenericGauge(cpuGauge, GetCPUPercentage)

	// Setup the RAM Gauge
	ramGauge := ui.NewGauge()
	ramGauge.Height = 2
	ramGauge.BarColor = ui.ColorGreen
	ramGauge.HasBorder = false
	ramGauge.PaddingBottom = 1
	go UpdateGenericGauge(ramGauge, GetRAMPercentage)

	// Setup the Label list
	ls := ui.NewList()
	ls.HasBorder = false
	ls.Items = []string{
		"CPU",
		"",
		"RAM",
	}
	ls.Height = 5

	// Setup the CPU Line Chart
	cpuLineChart := ui.NewLineChart()
	cpuLineChart.Width = 50
	cpuLineChart.Height = 11
	cpuLineChart.Border.Label = "CPU Usage"
	cpuLineChart.AxesColor = ui.ColorWhite
	cpuLineChart.LineColor = ui.ColorGreen | ui.AttrBold
	go UpdateGenericChart(cpuLineChart, GetCPUPercentage)

	// Setup the RAM Line Chart
	ramLineChart := ui.NewLineChart()
	ramLineChart.Width = 50
	ramLineChart.Height = 11
	ramLineChart.Border.Label = "RAM Usage"
	ramLineChart.AxesColor = ui.ColorWhite
	ramLineChart.LineColor = ui.ColorGreen | ui.AttrBold
	go UpdateGenericChart(ramLineChart, GetRAMPercentage)

	// Setup the layout
	ui.Body.AddRows(
		ui.NewRow(
			ui.NewCol(3, 0, cpuGauge, ramGauge),
			ui.NewCol(3, 0, ls),
		),
		ui.NewRow(
			ui.NewCol(6, 0, cpuLineChart),
			ui.NewCol(6, 0, ramLineChart),
		),
	)

	// Align
	ui.Body.Align()

	// Create the event polling system
	evt := make(chan tm.Event)
	go func() {
		for {
			evt <- tm.PollEvent()
		}
	}()

	for {
		select {
		case e := <-evt:
			if e.Type == tm.EventKey && e.Ch == 'q' {
				return
			}
			if e.Type == tm.EventResize {
				ui.Body.Width = ui.TermWidth()
				ui.Body.Align()
			}
		default:
			ui.Render(ui.Body)
			time.Sleep(time.Second / 2)
		}
	}
}
