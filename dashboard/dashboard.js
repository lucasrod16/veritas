let chartInstance = null;

const app = PetiteVue.createApp({
    imageReference: '',
    errorMessage: '',
    vulnerabilities: [],
    tableData: [],

    async handleScan() {
        this.errorMessage = '';

        if (this.imageReference.trim()) {
            const encodedImageReference = encodeURIComponent(this.imageReference);

            try {
                const [sbomResponse, tableResponse] = await Promise.all([
                    fetch(`http://localhost:8080/scan/report?image=${encodedImageReference}`),
                    fetch(`http://localhost:8080/scan/details?image=${encodedImageReference}`)
                ]);

                if (!sbomResponse.ok) {
                    this.errorMessage = await sbomResponse.text();
                    return;
                }
                if (!tableResponse.ok) {
                    this.errorMessage = await tableResponse.text();
                    return;
                }

                const sbomData = await sbomResponse.json();
                const tableData = await tableResponse.json();

                if (sbomData.vulnerabilities) {
                    this.vulnerabilities = sbomData.vulnerabilities;
                    console.log('Vulnerabilities loaded:', this.vulnerabilities);
                    this.$nextTick(() => {
                        this.renderChart();
                    });
                    sessionStorage.setItem('sbomData', JSON.stringify(sbomData.vulnerabilities));
                }

                if (tableData) {
                    this.tableData = tableData;
                    console.log('Table data loaded:', this.tableData);
                    sessionStorage.setItem('tableData', JSON.stringify(tableData));
                }

            } catch (error) {
                this.errorMessage = error.message;
                console.error('Fetch error:', error.message);
            }
        } else {
            this.errorMessage = 'Please enter a container image reference.';
        }
    },

    renderChart() {
        const canvas = document.getElementById('myChart');
        if (!canvas) {
            console.error('Canvas element not found.');
            return;
        }

        const ctx = canvas.getContext('2d');
        if (!ctx) {
            console.error('Failed to get context for the chart.');
            return;
        }

        if (chartInstance) {
            chartInstance.destroy();
        }

        const sevCount = this.getSeverityCount();
        console.log('Severity Count:', sevCount);

        chartInstance = new Chart(ctx, {
            type: 'pie',
            data: {
                labels: ['Low', 'Medium', 'High', 'Critical', 'Unknown'],
                datasets: [{
                    label: 'Vulnerabilities',
                    data: [
                        sevCount.low || 0,
                        sevCount.medium || 0,
                        sevCount.high || 0,
                        sevCount.critical || 0,
                        sevCount.unknown || 0
                    ],
                    backgroundColor: [
                        '#0bb400', // Low
                        '#ffff00', // Medium
                        '#ff7000', // High
                        '#ff0000', // Critical
                        '#747474'  // Unknown
                    ],
                }]
            }
        });
    },

    getSeverityCount() {
        const severityCount = {};
        this.vulnerabilities.forEach(vuln => {
            const ratings = vuln.ratings || [];
            ratings.forEach(rating => {
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
    },

    mounted() {
        this.loadFromSessionStorage();
    },

    loadFromSessionStorage() {
        const sbomData = sessionStorage.getItem('sbomData');
        if (sbomData) {
            this.vulnerabilities = JSON.parse(sbomData);
            this.$nextTick(() => {
                this.renderChart();
            });
        }

        const tableData = sessionStorage.getItem('tableData');
        if (tableData) {
            this.tableData = JSON.parse(tableData);
        }
    }
});

app.mount();
