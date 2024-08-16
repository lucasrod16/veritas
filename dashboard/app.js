import { handleScan } from "./scan.js";
import { renderChart } from "./chart.js";
import { loadFromSessionStorage } from "./localStorage.js";
import { downloadReport } from "./downloadReport.js";

PetiteVue.createApp({
	imageReference: "",
	vulnerabilities: [],
	tableData: [],
	isScanning: false,
	reportAvailable: false,
	reportData: {},

	async handleScan(event) {
		await handleScan(this, event);
	},

	renderChart() {
		renderChart(this);
	},

	downloadReport() {
		downloadReport(this);
	},

	mounted() {
		loadFromSessionStorage(this);
		renderChart(this);
	},
}).mount();
