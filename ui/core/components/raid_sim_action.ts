import { Encounter as EncounterProto } from '/tbc/core/proto/common.js';
import { Raid as RaidProto } from '/tbc/core/proto/api.js';
import { RaidSimRequest, RaidSimResult, ProgressMetrics } from '/tbc/core/proto/api.js';
import { SimRunData } from '/tbc/core/proto/ui.js';
import { ActionMetrics, SimResult, SimResultFilter } from '/tbc/core/proto_utils/sim_result.js';
import { SimUI } from '/tbc/core/sim_ui.js';
import { EventID, TypedEvent } from '/tbc/core/typed_event.js';
import { formatDeltaTextElem } from '/tbc/core/utils.js';

declare var tippy: any;

export function addRaidSimAction(simUI: SimUI): RaidSimResultsManager {
	simUI.addAction('DPS', 'dps-action', async () => simUI.runSim((progress: ProgressMetrics) => {
		resultsManager.setSimProgress(progress);
	}));

	const resultsManager = new RaidSimResultsManager(simUI);
	simUI.sim.simResultEmitter.on((eventID, simResult) => {
		resultsManager.setSimResult(eventID, simResult);
	});
	return resultsManager;
}

export type ReferenceData = {
	simResult: SimResult,
	settings: any,
	raidProto: RaidProto,
	encounterProto: EncounterProto,
};

export class RaidSimResultsManager {
	readonly currentChangeEmitter: TypedEvent<void> = new TypedEvent<void>();
	readonly referenceChangeEmitter: TypedEvent<void> = new TypedEvent<void>();

	readonly changeEmitter: TypedEvent<void> = new TypedEvent<void>();

	private readonly simUI: SimUI;

	private currentData: ReferenceData | null = null;
	private referenceData: ReferenceData | null = null;

	constructor(simUI: SimUI) {
		this.simUI = simUI;

		[
			this.currentChangeEmitter,
			this.referenceChangeEmitter,
		].forEach(emitter => emitter.on(eventID => this.changeEmitter.emit(eventID)));
	}

	setSimProgress(progress: ProgressMetrics) {
		if (progress.presimRunning) {
			this.simUI.resultsViewer.setContent(`
				<div class="results-sim">
						<div class="results-sim-dps">
							<span class="topline-result-avg">${progress.dps.toFixed(2)}</span>
						</div>
						<div class="">
							presimulations running
						</div>
				</div>
			`);
		} else {
			this.simUI.resultsViewer.setContent(`
				<div class="results-sim">
						<div class="results-sim-dps">
							<span class="topline-result-avg">${progress.dps.toFixed(2)}</span>
						</div>
						<div class="">
							${progress.completedIterations} / ${progress.totalIterations}<br>iterations complete
						</div>
				</div>
			`);
		}
	}

	setSimResult(eventID: EventID, simResult: SimResult) {
		this.currentData = {
			simResult: simResult,
			settings: {
				'raid': RaidProto.toJson(this.simUI.sim.raid.toProto()),
				'encounter': EncounterProto.toJson(this.simUI.sim.encounter.toProto()),
			},
			raidProto: RaidProto.clone(simResult.request.raid || RaidProto.create()),
			encounterProto: EncounterProto.clone(simResult.request.encounter || EncounterProto.create()),
		};
		this.currentChangeEmitter.emit(eventID);

		const dpsMetrics = simResult.raidMetrics.dps;
		this.simUI.resultsViewer.setContent(`
      <div class="results-sim">
				${RaidSimResultsManager.makeToplineResultsContent(simResult)}
				<div class="results-sim-reference">
					<span class="results-sim-set-reference fa fa-map-pin"></span>
					<div class="results-sim-reference-bar">
						<span class="results-sim-reference-dps-diff"></span>
						<span class="results-sim-reference-diff-separator threat-metrics">/</span>
						<span class="results-sim-reference-tps-diff threat-metrics"></span>
						<span class="results-sim-reference-diff-separator threat-metrics">/</span>
						<span class="results-sim-reference-dtps-diff threat-metrics"></span>
						<span class="results-sim-reference-diff-separator threat-metrics">/</span>
						<span class="results-sim-reference-chanceOfDeath-diff threat-metrics"></span>
						<span class="results-sim-reference-text"> vs. reference</span>
						<span class="results-sim-reference-swap fa fa-retweet"></span>
						<span class="results-sim-reference-delete fa fa-times"></span>
					</div>
				</div>
      </div>
    `);

		const setResultTippy = (cssClass: string, tippyContent: string) => {
			const resultDivElem = this.simUI.resultsViewer.contentElem.getElementsByClassName(cssClass)[0] as HTMLElement | undefined;
			if (resultDivElem) {
				tippy(resultDivElem, {
					'content': tippyContent,
					'allowHTML': true,
					placement: 'right',
				});
			}
		};
		setResultTippy('results-sim-dps', 'Damage Per Second');
		setResultTippy('results-sim-tps', 'Threat Per Second');
		setResultTippy('results-sim-dtps', 'Damage Taken Per Second');
		setResultTippy('results-sim-chanceOfDeath', `
			<p>Chance of Death</p>
			<p>The percentage of iterations in which the player died, based on incoming damage from the enemies and incoming healing (see the <b>Incoming HPS</b> and <b>Healing Cadence</b> options).</p>
			<p>DTPS alone is not a good measure of tankiness because it is not affected by health and ignores damage spikes. Chance of Death attempts to capture overall tankiness.</p>
		`);

		if (!this.simUI.isIndividualSim()) {
			Array.from(this.simUI.resultsViewer.contentElem.getElementsByClassName('results-sim-reference-diff-separator')).forEach(e => e.remove());
			Array.from(this.simUI.resultsViewer.contentElem.getElementsByClassName('results-sim-reference-tps-diff')).forEach(e => e.remove());
			Array.from(this.simUI.resultsViewer.contentElem.getElementsByClassName('results-sim-reference-dtps-diff')).forEach(e => e.remove());
			Array.from(this.simUI.resultsViewer.contentElem.getElementsByClassName('results-sim-reference-chanceOfDeath-diff')).forEach(e => e.remove());
		}

		const simReferenceElem = this.simUI.resultsViewer.contentElem.getElementsByClassName('results-sim-reference')[0] as HTMLDivElement;
		const simReferenceDiffElem = this.simUI.resultsViewer.contentElem.getElementsByClassName('results-sim-reference-diff')[0] as HTMLSpanElement;

		const simReferenceSetButton = this.simUI.resultsViewer.contentElem.getElementsByClassName('results-sim-set-reference')[0] as HTMLSpanElement;
		simReferenceSetButton.addEventListener('click', event => {
			this.referenceData = this.currentData;
			this.referenceChangeEmitter.emit(TypedEvent.nextEventID());
			this.updateReference();
		});
		tippy(simReferenceSetButton, {
			'content': 'Use as reference',
			'allowHTML': true,
		});

		const simReferenceSwapButton = this.simUI.resultsViewer.contentElem.getElementsByClassName('results-sim-reference-swap')[0] as HTMLSpanElement;
		simReferenceSwapButton.addEventListener('click', event => {
			TypedEvent.freezeAllAndDo(() => {
				if (this.currentData && this.referenceData) {
					const swapEventID = TypedEvent.nextEventID();
					const tmpData = this.currentData;
					this.currentData = this.referenceData;
					this.referenceData = tmpData;

					this.simUI.sim.raid.fromProto(swapEventID, this.currentData.raidProto);
					this.simUI.sim.encounter.fromProto(swapEventID, this.currentData.encounterProto);
					this.setSimResult(swapEventID, this.currentData.simResult);

					this.referenceChangeEmitter.emit(swapEventID);
					this.updateReference();
				}
			});
		});
		tippy(simReferenceSwapButton, {
			'content': 'Swap reference with current',
			'allowHTML': true,
		});

		const simReferenceDeleteButton = this.simUI.resultsViewer.contentElem.getElementsByClassName('results-sim-reference-delete')[0] as HTMLSpanElement;
		simReferenceDeleteButton.addEventListener('click', event => {
			this.referenceData = null;
			this.referenceChangeEmitter.emit(TypedEvent.nextEventID());
			this.updateReference();
		});
		tippy(simReferenceDeleteButton, {
			'content': 'Remove reference',
			'allowHTML': true,
		});

		this.updateReference();
	}

	private updateReference() {
		const simReferenceElem = this.simUI.resultsViewer.contentElem.getElementsByClassName('results-sim-reference')[0] as HTMLDivElement;
		const simReferenceDpsDiffElem = this.simUI.resultsViewer.contentElem.getElementsByClassName('results-sim-reference-dps-diff')[0] as HTMLSpanElement;

		if (!this.referenceData || !this.currentData) {
			simReferenceElem.classList.remove('has-reference');
			return;
		}
		simReferenceElem.classList.add('has-reference');

		const currentDpsMetrics = this.currentData.simResult.raidMetrics.dps;
		const referenceDpsMetrics = this.referenceData.simResult.raidMetrics.dps;
		formatDeltaTextElem(simReferenceDpsDiffElem, referenceDpsMetrics.avg, currentDpsMetrics.avg, 2);

		if (this.simUI.isIndividualSim()) {
			const simReferenceTpsDiffElem = this.simUI.resultsViewer.contentElem.getElementsByClassName('results-sim-reference-tps-diff')[0] as HTMLSpanElement;
			const simReferenceDtpsDiffElem = this.simUI.resultsViewer.contentElem.getElementsByClassName('results-sim-reference-dtps-diff')[0] as HTMLSpanElement;
			const simReferenceCodDiffElem = this.simUI.resultsViewer.contentElem.getElementsByClassName('results-sim-reference-chanceOfDeath-diff')[0] as HTMLSpanElement;

			const curPlayerMetrics = this.currentData.simResult.getPlayers()[0]!;
			const refPlayerMetrics = this.referenceData.simResult.getPlayers()[0]!;
			formatDeltaTextElem(simReferenceTpsDiffElem, refPlayerMetrics.tps.avg, curPlayerMetrics.tps.avg, 2);
			formatDeltaTextElem(simReferenceDtpsDiffElem, refPlayerMetrics.dtps.avg, curPlayerMetrics.dtps.avg, 2);
			formatDeltaTextElem(simReferenceCodDiffElem, refPlayerMetrics.chanceOfDeath, curPlayerMetrics.chanceOfDeath, 1);
		}
	}

	getRunData(): SimRunData | null {
		if (this.currentData == null) {
			return null;
		}

		return SimRunData.create({
			run: this.currentData.simResult.toProto(),
			referenceRun: this.referenceData?.simResult.toProto(),
		});
	}

	getCurrentData(): ReferenceData | null {
		if (this.currentData == null) {
			return null;
		}

		// Defensive copy.
		return {
			simResult: this.currentData.simResult,
			settings: JSON.parse(JSON.stringify(this.currentData.settings)),
			raidProto: this.currentData.raidProto,
			encounterProto: this.currentData.encounterProto,
		};
	}

	getReferenceData(): ReferenceData | null {
		if (this.referenceData == null) {
			return null;
		}

		// Defensive copy.
		return {
			simResult: this.referenceData.simResult,
			settings: JSON.parse(JSON.stringify(this.referenceData.settings)),
			raidProto: this.referenceData.raidProto,
			encounterProto: this.referenceData.encounterProto,
		};
	}

	static makeToplineResultsContent(simResult: SimResult, filter?: SimResultFilter): string {
		const players = simResult.getPlayers(filter);
		const playerMetrics = players.length == 1 ? players[0] : null;
		let content = '';

		if (playerMetrics) {
			if (playerMetrics.getTargetIndex(filter) == null) {
				const dpsMetrics = simResult.raidMetrics.dps;
				const tpsMetrics = playerMetrics.tps;
				const dtpsMetrics = playerMetrics.dtps;
				content += `
					<div class="results-sim-dps">
						<span class="topline-result-avg">${dpsMetrics.avg.toFixed(2)}</span>
						<span class="topline-result-stdev">${dpsMetrics.stdev.toFixed(2)}</span>
					</div>
					<div class="results-sim-tps threat-metrics">
						<span class="topline-result-avg">${tpsMetrics.avg.toFixed(2)}</span>
						<span class="topline-result-stdev">${tpsMetrics.stdev.toFixed(2)}</span>
					</div>
					<div class="results-sim-dtps threat-metrics">
						<span class="topline-result-avg">${dtpsMetrics.avg.toFixed(2)}</span>
						<span class="topline-result-stdev">${dtpsMetrics.stdev.toFixed(2)}</span>
					</div>
					<div class="results-sim-chanceOfDeath threat-metrics">
						<span class="topline-result-avg">${playerMetrics.chanceOfDeath.toFixed(2)}</span>
					</div>
				`;
			} else {
				const actions = simResult.getActionMetrics(filter);
				const targetActions = simResult.getTargets(filter)[0].actions.map(action => action.forTarget(filter));
				if (actions.length > 0) {
					const mergedActions = ActionMetrics.merge(actions);
					content += `
						<div class="results-sim-dps">
							<span class="topline-result-avg">${mergedActions.dps.toFixed(2)}</span>
						</div>
						<div class="results-sim-tps threat-metrics">
							<span class="topline-result-avg">${mergedActions.tps.toFixed(2)}</span>
						</div>
					`;
				}
				if (targetActions.length > 0) {
					const mergedTargetActions = ActionMetrics.merge(targetActions);
					content += `
						<div class="results-sim-dtps threat-metrics">
							<span class="topline-result-avg">${mergedTargetActions.dps.toFixed(2)}</span>
						</div>
					`;
				}
			}
		} else {
			const dpsMetrics = simResult.raidMetrics.dps;
			content = `
				<div class="results-sim-dps">
					<span class="topline-result-avg">${dpsMetrics.avg.toFixed(2)}</span>
					<span class="topline-result-stdev">${dpsMetrics.stdev.toFixed(2)}</span>
				</div>
			`;
		}

		if (simResult.request.encounter?.useHealth) {
			content += `<div class="results-sim-dur"><span class="topline-result-avg">${simResult.result.avgIterationDuration.toFixed(2)}s</span></div>`;
		}

		return content;
	}
}
