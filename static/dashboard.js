// dashboard.js
document.addEventListener("DOMContentLoaded", function () {
    const metricsDiv = document.getElementById("metrics");
    const cpuMetricSpan = document.getElementById("cpuMetric");
    const memMetricSpan = document.getElementById("memMetric");
    const diskMetricSpan = document.getElementById("diskMetric");
    const messageTextSpan = document.getElementById("messageText");

    function updateMetrics() {
        fetch("/get_metrics")
            .then((response) => response.json())
            .then((data) => {
                cpuMetricSpan.textContent = data.cpu_metric.toFixed(2);
                memMetricSpan.textContent = data.mem_metric.toFixed(2);
                diskMetricSpan.textContent = data.disk_metric.toFixed(2);
                messageTextSpan.textContent = data.message;
            })
            .catch((error) => {
                console.error("Error fetching metrics:", error);
            });
    }

    // Update metrics initially and then every 5 seconds
    updateMetrics();
    setInterval(updateMetrics, 5000);
});
