var workerID = "";

addEventListener('message', async (e) => {
	const msg = e.data.msg;
	const id = e.data.id;

	if (msg == "setID") {
		workerID = id;
		postMessage({ msg: "idconfirm" })
		return;
	}

	var url = "/" + msg;
	let response = await fetch(url, {
		method: 'POST',
		headers: {
			'Content-Type': 'application/x-protobuf'
		},
		body: e.data.inputData
	});

	var content = await response.arrayBuffer();
	var outputData;
	if (msg == "raidSimAsync" || msg == "statWeightsAsync" || msg == "bulkSimAsync") {
		while (true) {
			let progressResponse = await fetch("/asyncProgress", {
				method: 'POST',
				headers: {
					'Content-Type': 'application/x-protobuf'
				},
				body: content,
			});

			// If no new data available, stop querying.
			if (progressResponse.status == 204) {
				break
			}

			outputData = await progressResponse.arrayBuffer();
			var uint8View = new Uint8Array(outputData);
			postMessage({
				msg: msg,
				outputData: uint8View,
				id: id + "progress",
			});
			await new Promise(resolve => setTimeout(resolve, 500));
		}
	} else {
		outputData = content;
	}

	var uint8View = new Uint8Array(outputData);
	postMessage({
		msg: msg,
		outputData: uint8View,
		id: id,
	});

}, false);

// Let UI know worker is ready.
postMessage({
	msg: "ready"
});
