document.addEventListener('DOMContentLoaded', async () => {
    loadFromSessionStorage();
    await setupEventListeners();
});

let chart = null;
function renderChart(vulns) {
    const ctx = document.getElementById('myChart').getContext('2d');

    if (chart) {
        chart.destroy();
    }

    const sevCount = getSeverityCount(vulns);

    chart = new Chart(ctx, {
        type: 'pie',
        data: {
            labels: [
                'Low',
                'Medium',
                'High',
                'Critical',
                'Unknown'
            ],
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
                    '#747474', // Unknown
                ],
            }]
        }
    });
}

function getSeverityCount(vulns) {
    const severityCount = {};

    vulns.forEach(vuln => {
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
}

function renderTable(data) {
    const templateData = {
        details: data.map(item => ({
            name: item.package.Name,
            installed: item.package.Version,
            fixedIn: item.vulnerability.Fix.Versions || 'N/A',
            type: item.package.Type,
            vulnerabilityId: item.vulnerability.ID,
            severity: item.severity
        }))
    };
    const template = document.getElementById('table-template').innerHTML;
    const rendered = Mustache.render(template, templateData);
    document.getElementById('table-container').innerHTML = rendered;
}

function loadFromSessionStorage() {
    const sbomData = sessionStorage.getItem('sbomData');
    if (sbomData) {
        const parsedSbomData = JSON.parse(sbomData);
        renderChart(parsedSbomData);
    }

    const tableData = sessionStorage.getItem('tableData');
    if (tableData) {
        const parsedTableData = JSON.parse(tableData);
        renderTable(parsedTableData);
    }
}

async function setupEventListeners() {
    document.getElementById('scanButton').addEventListener('click', async () => {
        clearError();

        const imageReference = document.getElementById('imageInput').value.trim();

        if (imageReference) {
            const encodedImageReference = encodeURIComponent(imageReference);

            try {
                const [sbomResponse, tableResponse] = await Promise.all([
                    fetch(`http://localhost:8080/scan/report?image=${encodedImageReference}`),
                    fetch(`http://localhost:8080/scan/details?image=${encodedImageReference}`)
                ]);

                if (!sbomResponse.ok) {
                    showError(await sbomResponse.text())
                    return;
                }
                if (!tableResponse.ok) {
                    showError(await tableResponse.text())
                    return;
                }

                const sbomData = await sbomResponse.json();
                const tableData = await tableResponse.json();

                if (sbomData.vulnerabilities) {
                    renderChart(sbomData.vulnerabilities);
                    sessionStorage.setItem('sbomData', JSON.stringify(sbomData.vulnerabilities));
                }

                if (tableData) {
                    renderTable(tableData);
                    sessionStorage.setItem('tableData', JSON.stringify(tableData));
                }

            } catch (error) {
                console.error(error.message);
                showError(error.message);
            }
        } else {
            showError('Please enter a container image reference.');
        }
    });
}

function showError(message) {
    const errorElement = document.getElementById('error-message');
    if (errorElement) {
        errorElement.textContent = message;
        errorElement.style.display = 'block';
    }
}

function clearError() {
    const errorElement = document.getElementById('error-message');
    if (errorElement) {
        errorElement.textContent = '';
        errorElement.style.display = 'none';
    }
}
