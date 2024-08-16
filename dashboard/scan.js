export async function handleScan(app, event) {
	app.reportAvailable = false;

	if (event.key === "Enter") {
		event.preventDefault();
	}

	const imgRef = app.imageReference.trim();

	if (!imgRef) {
		Swal.fire({
			title: "Oops!",
			text: "Please enter a container image reference to scan.",
			icon: "info",
			confirmButtonText: "Continue",
		});
		return;
	}

	const encodedImageReference = encodeURIComponent(imgRef);

	app.isScanning = true;

	try {
		const [sbomResponse, tableResponse] = await Promise.all([
			fetch(`http://localhost:8080/scan/report?image=${encodedImageReference}`),
			fetch(
				`http://localhost:8080/scan/details?image=${encodedImageReference}`
			),
		]);

		if (!sbomResponse.ok) {
			Swal.fire({
				title: "Error",
				text: await sbomResponse.text(),
				icon: "error",
				confirmButtonText: "Continue",
			});
			return;
		}
		if (!tableResponse.ok) {
			Swal.fire({
				title: "Error",
				text: await tableResponse.text(),
				icon: "error",
				confirmButtonText: "Continue",
			});
			return;
		}

		const sbomData = await sbomResponse.json();
		const tableData = await tableResponse.json();

		if (sbomData.vulnerabilities) {
			app.vulnerabilities = sbomData.vulnerabilities;
			console.log("Vulnerabilities loaded:", app.vulnerabilities);
			app.$nextTick(() => {
				app.renderChart();
			});
			sessionStorage.setItem(
				"sbomData",
				JSON.stringify(sbomData.vulnerabilities)
			);
		}

		if (tableData) {
			app.tableData = tableData;
			console.log("Table data loaded:", app.tableData);
			sessionStorage.setItem("tableData", JSON.stringify(tableData));
		}

		if (sbomData.vulnerabilities.length === 0) {
			Swal.fire({
				title: "😎",
				text: `No vulnerabilities found for ${imgRef}`,
				icon: "success",
				confirmButtonText: "Continue",
			});
		}
		app.reportData = sbomData;
		app.reportAvailable = true;
	} catch (error) {
		Swal.fire({
			title: "Error",
			text: error.message,
			icon: "error",
			confirmButtonText: "Continue",
		});
		console.error("Fetch error:", error.message);
	} finally {
		app.isScanning = false;
	}
}
