document.addEventListener('DOMContentLoaded', () => {
    loadFromSessionStorage();
    setupEventListeners();
});

async function fetchData(url) {
    try {
        const response = await fetch(url);
        if (!response.ok) {
            throw new Error(`Response status: ${response.status}`);
        }
        const contentType = response.headers.get("content-type");
        if (!contentType || !contentType.includes("application/json")) {
            throw new TypeError("Oops, we haven't got JSON!");
        }
        return await response.json();
    } catch (error) {
        console.error(error.message);
        return null;
    }
}

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

function setupEventListeners() {
    document.getElementById('scanButton').addEventListener('click', async () => {
        const imageReference = document.getElementById('imageInput').value.trim();
        if (imageReference) {
            const encodedImageReference = encodeURIComponent(imageReference);

            try {
                const [sbomResponse, tableResponse] = await Promise.all([
                    fetchData(`http://localhost:8080/scan/report?image=${encodedImageReference}`),
                    fetchData(`http://localhost:8080/scan/details?image=${encodedImageReference}`)
                ]);

                if (sbomResponse) {
                    renderChart(sbomResponse.vulnerabilities);
                    sessionStorage.setItem('sbomData', JSON.stringify(sbomResponse.vulnerabilities));
                }

                if (tableResponse) {
                    renderTable(tableResponse);
                    sessionStorage.setItem('tableData', JSON.stringify(tableResponse));
                }

            } catch (error) {
                console.error('Error fetching data:', error);
                alert('An error occurred while fetching data. Please try again.');
            }
        } else {
            alert('Please enter a container image reference.');
        }
    });
}
