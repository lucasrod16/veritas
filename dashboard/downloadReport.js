export function downloadReport(app) {
	const blob = new Blob([JSON.stringify(app.reportData, null, 2)], {
		type: "application/json",
	});

	const url = URL.createObjectURL(blob);
	const a = document.createElement("a");
	a.href = url;
	a.download = "scan-report.json";
	a.click();
	URL.revokeObjectURL(url);
}
