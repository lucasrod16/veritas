let chart = null;

async function getSBOM() {
    const url = "http://localhost:8080/scan/report?image=rockylinux:8.7";
    try {
        const response = await fetch(url);
        if (!response.ok) {
            throw new Error(`Response status: ${response.status}`);
        }
        const contentType = response.headers.get("content-type");
        if (!contentType || !contentType.includes("application/json")) {
            throw new TypeError("Oops, we haven't got JSON!");
        }

        const sbom = await response.json();
        return sbom;
    } catch (error) {
        console.error(error.message);
        return null;
    }
}

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
                'Negligible',
                'Low',
                'Medium',
                'High',
                'Critical',
                'Unknown'
            ],
            datasets: [{
                label: 'Vulnerabilities',
                data: [
                    sevCount.negligible,
                    sevCount.low,
                    sevCount.medium,
                    sevCount.high,
                    sevCount.critical,
                    sevCount.unknown
                ],
                backgroundColor: [
                    '#0031c0', // Negligible
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
                severity === "negligible"
            }
            if (!severityCount[severity]) {
                severityCount[severity] = 0;
            }
            severityCount[severity]++;
        });
    });
    return severityCount;
}

function loadFromLocalStorage() {
    const savedData = localStorage.getItem('vulnData');
    if (savedData) {
        const vulnData = JSON.parse(savedData);
        renderChart(vulnData);
    }
}
