import { Spec } from '../core/proto/common.js';
import { SimResult, SimResultFilter } from '../core/proto_utils/sim_result.js';

import { ResultComponent, ResultComponentConfig, SimResultData } from './result_component.js';
import { RaidSimResultsManager } from '../core/components/raid_sim_action.js';

export class ToplineResults extends ResultComponent {
	constructor(config: ResultComponentConfig) {
		config.rootCssClass = 'topline-results-root';
		super(config);
	}

	onSimResult(resultData: SimResultData) {
		let content = RaidSimResultsManager.makeToplineResultsContent(resultData.result, resultData.filter);

		const noManaSpecs = [
			Spec.SpecFeralTankDruid,
			Spec.SpecRogue,
			Spec.SpecWarrior,
			Spec.SpecProtectionWarrior,
		];

		const demonicPactSpecs = [
			Spec.SpecWarlock
		];

		const players = resultData.result.getPlayers(resultData.filter);
		if (players.length == 1) {
			const player = players[0];

			if (demonicPactSpecs.includes(player.spec) && player.dpasp.avg > 0) {
				content += `
					<div class="dpasp damage-metrics">
						<span class="topline-result-avg">${player.dpasp.avg.toFixed(1)} </span>
						<span class="topline-result-stdev">(${"\u00B1"}${player.dpasp.stdev.toFixed(1)})</span>
						<span class="topline-result-label"> DP Avg SP</span>
					</div>
				`;
			}

			if (!noManaSpecs.includes(player.spec)) {
				const secondsOOM = player.secondsOomAvg;
				const percentOOM = secondsOOM / resultData.result.encounterMetrics.durationSeconds;
				const dangerLevel = percentOOM < 0.01 ? 'safe' : (percentOOM < 0.05 ? 'warning' : 'danger');

				content += `
					<div class="percent-oom ${dangerLevel} damage-metrics">
						<span class="topline-result-avg">${secondsOOM.toFixed(1)}s</span>
						<span class="topline-result-label"> spent OOM</span>
					</div>
				`;
			}
		}

		this.rootElem.innerHTML = content;
	}
}
