async function loadData() {
  const res = await fetch("/api/week")
  return await res.json()
}
function createChart(ctx, label, data, labels) {
  return new Chart(ctx, {
    type: "line",
    data: {
      labels: labels,
      datasets: [{
        label: label,
        data: data,
        borderWidth: 2,
        tension: 0.3
      }]
    },
    options: {
      responsive: true
    }
  })
}

async function initCharts() {
  const data = await loadData()

  const labels = data.map(d => d.time)

  const pm25 = data.map(d => d.pm25)
  const pm10 = data.map(d => d.pm10)
  const no2  = data.map(d => d.no2)
  const co   = data.map(d => d.co)

  createChart(document.getElementById("pm25Chart"), "PM2.5", pm25, labels)
  createChart(document.getElementById("pm10Chart"), "PM10", pm10, labels)
  createChart(document.getElementById("no2Chart"), "NO2", no2, labels)
  createChart(document.getElementById("coChart"), "CO", co, labels)
}

initCharts()