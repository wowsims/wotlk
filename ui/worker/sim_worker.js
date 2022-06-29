// Wasm binary calls this function when its done loading.
function wasmready() {
	postMessage({
		msg: "ready"
	});
}

const go = new Go();
let mod, inst;

WebAssembly.instantiateStreaming(fetch("lib.wasm"), go.importObject).then(
	async result => {
		mod = result.module;
		inst = result.instance;
		// console.log("loading wasm...")
		await go.run(inst);
	}
);

var workerID = "";

addEventListener('message', async (e) => {
	const msg = e.data.msg;
	const id = e.data.id;

	let handled = false;

	[
		['computeStats', computeStats],
		['gearList', gearList],
		['raidSim', raidSim],
		['raidSimAsync', (data) => {
			return raidSimAsync(data, (result) => {
				postMessage({
					msg: "progress",
					outputData: result,
					id: id+"progress",
				});
			});
		}],
		['statWeights', statWeights],
		['statWeightsAsync', (data) => {
			return statWeightsAsync(data, (result) => {
				postMessage({
					msg: "progress",
					outputData: result,
					id: id+"progress",
				});
			});
		}],
	].forEach(funcData => {
		const funcName = funcData[0];
		const func = funcData[1];

		if (msg == funcName) {
			const outputData = func(e.data.inputData);

			postMessage({
				msg: funcName,
				outputData: outputData,
				id: id,
			});
			handled = true;
		}
	});

	if (handled) {
		return;
	}

	if (msg == "setID") {
		workerID = id;
		postMessage({ msg: "idconfirm" })
	}
}, false);
