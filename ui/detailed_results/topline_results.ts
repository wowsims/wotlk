import { Spec } from '../core/proto/common.js';
import { SimResult, SimResultFilter } from '../core/proto_utils/sim_result.js';

import { ResultComponent, ResultComponentConfig, SimResultData } from './result_component.js';
import { RaidSimResultsManager } from '../core/components/raid_sim_action.js';

export class ToplineResults extends ResultComponent {
	constructor(config: ResultComponentConfig) {
		config.rootCssClass = 'topline-results-root';
		super(config);

		this.rootElem.classList.add('results-sim');
	}

	onSimResult(resultData: SimResultData) {
		let content = RaidSimResultsManager.makeToplineResultsContent(resultData.result, resultData.filter);

		const noManaSpecs = [
			Spec.SpecFeralTankDruid,
			Spec.SpecRogue,
			Spec.SpecWarrior,
			Spec.SpecProtectionWarrior,
		];

		const players = resultData.result.getPlayers(resultData.filter);
		if (players.length == 1 && !noManaSpecs.includes(players[0].spec)) {
			const player = players[0];
			const secondsOOM = player.secondsOomAvg;
			const percentOOM = secondsOOM / resultData.result.encounterMetrics.durationSeconds;
			const dangerLevel = percentOOM < 0.01 ? 'safe' : (percentOOM < 0.05 ? 'warning' : 'danger');

			content += `
				<div class="results-sim-percent-oom ${dangerLevel} damage-metrics">
					<span class="topline-result-avg">${secondsOOM.toFixed(1)}s</span>
				</div>
			`;
		}

		this.rootElem.innerHTML = content;
	}
}
