export function loadFromSessionStorage(app) {
	const sbomData = sessionStorage.getItem("sbomData");
	if (sbomData) {
		app.vulnerabilities = JSON.parse(sbomData);
		app.$nextTick(() => {
			app.renderChart();
		});
	}

	const tableData = sessionStorage.getItem("tableData");
	if (tableData) {
		app.tableData = JSON.parse(tableData);
	}
}
