let chartInstance = null;

export function renderChart(app) {
	const canvas = document.getElementById("myChart");
	if (!canvas) {
		console.error("Canvas element not found.");
		return;
	}

	const ctx = canvas.getContext("2d");
	if (!ctx) {
		console.error("Failed to get context for the chart.");
		return;
	}

	if (chartInstance) {
		chartInstance.destroy();
	}

	const sevCount = getSeverityCount(app);
	console.log("Severity Count:", sevCount);

	chartInstance = new Chart(ctx, {
		type: "pie",
		data: {
			labels: ["Low", "Medium", "High", "Critical", "Unknown"],
			datasets: [
				{
					label: "Vulnerabilities",
					data: [
						sevCount.low || 0,
						sevCount.medium || 0,
						sevCount.high || 0,
						sevCount.critical || 0,
						sevCount.unknown || 0,
					],
					backgroundColor: [
						"#9de0a8", // Low
						"#fdfd96", // Medium
						"#f8b400", // High
						"#e57373", // Critical
						"#c0c0c0", // Unknown
					],
				},
			],
		},
		options: {
			maintainAspectRatio: false,
		},
	});
}

function getSeverityCount(app) {
	const severityCount = {};
	app.vulnerabilities.forEach((vuln) => {
		const ratings = vuln.ratings || [];
		ratings.forEach((rating) => {
			let severity = rating.severity;
			if (severity === "none") {
				severity = "negligible";
			}
			if (!severityCount[severity]) {
				severityCount[severity] = 0;
			}
			severityCount[severity]++;
		});
	});
	return severityCount;
}
