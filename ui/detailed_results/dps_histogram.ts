import { SimResult, SimResultFilter } from '/tbc/core/proto_utils/sim_result.js';

import { ResultComponent, ResultComponentConfig, SimResultData } from './result_component.js';

declare var Chart: any;

export class DpsHistogram extends ResultComponent {
	constructor(config: ResultComponentConfig) {
		config.rootCssClass = 'dps-histogram-root';
		super(config);
	}

	onSimResult(resultData: SimResultData) {
		const chartBounds = this.rootElem.getBoundingClientRect();

		this.rootElem.textContent = '';
		const chartCanvas = document.createElement("canvas");
		chartCanvas.height = chartBounds.height;
		chartCanvas.width = chartBounds.width;

		const damageMetrics = resultData.result.getDamageMetrics(resultData.filter);

		const min = damageMetrics.avg - damageMetrics.stdev;
		const max = damageMetrics.avg + damageMetrics.stdev;
		const vals: Array<number> = [];
		const colors: Array<string> = [];

		const labels = Object.keys(damageMetrics.hist);
		labels.forEach((k, i) => {
			vals.push(damageMetrics.hist[Number(k)]);
			const val = parseInt(k);
			if (val > min && val < max) {
				colors.push('#1E87F0');
			} else {
				colors.push('#FF6961');
			}
		});

		const ctx = chartCanvas.getContext('2d');
		const chart = new Chart(ctx, {
			type: 'bar',
			data: {
				labels: labels,
				datasets: [{
					data: vals,
					backgroundColor: colors,
				}],
			},
			options: {
				plugins: {
					title: {
						display: true,
						text: 'DPS Histogram',
					},
					legend: {
						display: false,
						labels: {},
					}
				},
				scales: {
					y: {
						beginAtZero: true,
						ticks: {
							display: false
						},
					},
				},
			},
		});
		this.rootElem.appendChild(chartCanvas);
	}
}
