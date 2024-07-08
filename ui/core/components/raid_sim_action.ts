import { Tooltip } from 'bootstrap';
import tippy from 'tippy.js';

import { DistributionMetrics as DistributionMetricsProto , ProgressMetrics,Raid as RaidProto , RaidSimRequest, RaidSimResult } from '../proto/api.js';
import { Encounter as EncounterProto } from '../proto/common.js';
import { SimRunData } from '../proto/ui.js';
import { ActionMetrics, SimResult, SimResultFilter } from '../proto_utils/sim_result.js';
import { SimUI } from '../sim_ui.js';
import { EventID, TypedEvent } from '../typed_event.js';
import { formatDeltaTextElem } from '../utils.js';

export function addRaidSimAction(simUI: SimUI): RaidSimResultsManager {
	simUI.addAction('开始模拟', 'dps-action', async () => simUI.runSim((progress: ProgressMetrics) => {
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

export interface ResultMetrics {
	cod: string,
	dps: string,
	dpasp: string,
	dtps: string,
	tmi: string,
	dur: string,
	hps: string,
	tps: string,
	tto: string,
}

export interface ResultMetricCategories {
	damage: string,
	demo: string,
	healing: string,
	threat: string,
}

export interface ResultsLineArgs {
	average: number,
	stdev?: number,
	classes?: string
}

export class RaidSimResultsManager {
	static resultMetricCategories: { [ResultMetrics: string]: keyof ResultMetricCategories } = {
		dps: 'damage',
		dpasp: 'demo',
		tps: 'threat',
		dtps: 'threat',
		tmi: 'threat',
		cod: 'threat',
		tto: 'healing',
		hps: 'healing',
	}

	static resultMetricClasses: { [ResultMetrics: string]: string } = {
		cod: 'results-sim-cod',
		dps: 'results-sim-dps',
		dpasp: 'results-sim-dpasp',
		dtps: 'results-sim-dtps',
		tmi: 'results-sim-tmi',
		dur: 'results-sim-dur',
		hps: 'results-sim-hps',
		tps: 'results-sim-tps',
		tto: 'results-sim-tto',
	}

	static metricsClasses: { [ResultMetricCategories: string]: string } = {
		damage: 'damage-metrics',
		demo: 'demo-metrics',
		healing: 'healing-metrics',
		threat: 'threat-metrics',
	}

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
		this.simUI.resultsViewer.setContent(`
			<div class="results-sim">
				<div class="results-sim-dps damage-metrics">
					<span class="topline-result-avg">${progress.dps.toFixed(2)}</span>
				</div>
				${!this.simUI.isIndividualSim() ? '' : `<div class="results-sim-hps healing-metrics">
					<span class="topline-result-avg">${progress.hps.toFixed(2)}</span>
				</div>`}
				<div class="">
					${progress.presimRunning ? 'presimulations running' : `${progress.completedIterations} / ${progress.totalIterations}<br>iterations complete`}
				</div>
			</div>
		`);
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
					<a
						href="javascript:void(0)"
						class="results-sim-set-reference"
						role="button"
					>
						<i class="fa fa-map-pin fa-lg text-${this.simUI.cssScheme} me-2"></i>设为参考指标
					</a>
					<div class="results-sim-reference-bar">
						<a href="javascript:void(0)" class="results-sim-reference-swap me-3" role="button">
							<i class="fas fa-arrows-rotate me-1"></i>替换
						</a>
						<a href="javascript:void(0)" class="results-sim-reference-delete" role="button">
							<i class="fa fa-times fa-lg me-1"></i>取消参考指标
						</a>
					</div>
				</div>
      </div>
    `);

		const setResultTooltip = (cssClass: string, tooltip: string) => {
			const resultDivElem = this.simUI.resultsViewer.contentElem.getElementsByClassName(cssClass)[0] as HTMLElement | undefined;
			if (resultDivElem) {
				Tooltip.getOrCreateInstance(resultDivElem, {title: tooltip, html: true, placement: 'right'});
			}
		};
		setResultTooltip('results-sim-dps', '每秒伤害');
		setResultTooltip('results-sim-dpasp', '恶魔契约的平均法术强度');
		setResultTooltip('results-sim-tto', '空蓝耗时');
		setResultTooltip('results-sim-hps', '每秒治疗+护盾，包括过量治疗');
		setResultTooltip('results-sim-tps', '每秒仇恨');
		setResultTooltip('results-sim-dtps', '每秒承受伤害');
		setResultTooltip('results-sim-tmi', `
			<p>Theck-Meloree 指数 (TMI)</p>
			<p>将闪避的优点与有效生命值结合起来以衡量输入伤害平滑度的指标</p>
			<p><b>越低越好</b> 这表示根据战斗设置，在6秒爆发窗口中预期会损失的生命值百分比。</p>
		`);
		setResultTooltip('results-sim-cod', `
			<p>死亡几率</p>
			<p>基于来自敌人的输入伤害和输入治疗（参见<b>输入HPS</b>和<b>治疗节奏</b>选项），表示玩家死亡的迭代次数百分比。</p>
			<p>单靠每秒承受伤害（DTPS）并不能很好地衡量坦克的坚韧度，因为它不受生命值影响并忽略了伤害尖峰。死亡概率试图捕捉整体的坦克硬度。</p>
		`);

		if (!this.simUI.isIndividualSim()) {
			Array.from(this.simUI.resultsViewer.contentElem.getElementsByClassName('results-sim-reference-diff-separator')).forEach(e => e.remove());
			Array.from(this.simUI.resultsViewer.contentElem.getElementsByClassName('results-sim-dpasp')).forEach(e => e.remove());
			Array.from(this.simUI.resultsViewer.contentElem.getElementsByClassName('results-sim-tto')).forEach(e => e.remove());
			Array.from(this.simUI.resultsViewer.contentElem.getElementsByClassName('results-sim-hps')).forEach(e => e.remove());
			Array.from(this.simUI.resultsViewer.contentElem.getElementsByClassName('results-sim-tps')).forEach(e => e.remove());
			Array.from(this.simUI.resultsViewer.contentElem.getElementsByClassName('results-sim-dtps')).forEach(e => e.remove());
			Array.from(this.simUI.resultsViewer.contentElem.getElementsByClassName('results-sim-tmi')).forEach(e => e.remove());
			Array.from(this.simUI.resultsViewer.contentElem.getElementsByClassName('results-sim-cod')).forEach(e => e.remove());
		}

		const simReferenceElem = this.simUI.resultsViewer.contentElem.getElementsByClassName('results-sim-reference')[0] as HTMLDivElement;
		const simReferenceDiffElem = this.simUI.resultsViewer.contentElem.getElementsByClassName('results-sim-reference-diff')[0] as HTMLSpanElement;

		const simReferenceSetButton = this.simUI.resultsViewer.contentElem.getElementsByClassName('results-sim-set-reference')[0] as HTMLSpanElement;
		simReferenceSetButton.addEventListener('click', event => {
			this.referenceData = this.currentData;
			this.referenceChangeEmitter.emit(TypedEvent.nextEventID());
			this.updateReference();
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

		const simReferenceDeleteButton = this.simUI.resultsViewer.contentElem.getElementsByClassName('results-sim-reference-delete')[0] as HTMLSpanElement;
		simReferenceDeleteButton.addEventListener('click', event => {
			this.referenceData = null;
			this.referenceChangeEmitter.emit(TypedEvent.nextEventID());
			this.updateReference();
		});

		this.updateReference();
	}

	private updateReference() {
		if (!this.referenceData || !this.currentData) {
			// Remove references
			this.simUI.resultsViewer.contentElem.querySelector('.results-sim-reference')?.classList.remove('has-reference');
			this.simUI.resultsViewer.contentElem.querySelectorAll('.results-reference').forEach(e => e.classList.add('hide'));
			return;
		} else {
			// Add references references
			this.simUI.resultsViewer.contentElem.querySelector('.results-sim-reference')?.classList.add('has-reference');
			this.simUI.resultsViewer.contentElem.querySelectorAll('.results-reference').forEach(e => e.classList.remove('hide'));
		}

		this.formatToplineResult(`.${RaidSimResultsManager.resultMetricClasses['dps']} .results-reference-diff`, res => res.raidMetrics.dps, 2);
		if (this.simUI.isIndividualSim()) {
			this.formatToplineResult(`.${RaidSimResultsManager.resultMetricClasses['hps']} .results-reference-diff`, res => res.raidMetrics.hps, 2);
			this.formatToplineResult(`.${RaidSimResultsManager.resultMetricClasses['dpasp']} .results-reference-diff`, res => res.getPlayers()[0]!.dpasp, 2);
			this.formatToplineResult(`.${RaidSimResultsManager.resultMetricClasses['tto']} .results-reference-diff`, res => res.getPlayers()[0]!.tto, 2);
			this.formatToplineResult(`.${RaidSimResultsManager.resultMetricClasses['tps']} .results-reference-diff`, res => res.getPlayers()[0]!.tps, 2);
			this.formatToplineResult(`.${RaidSimResultsManager.resultMetricClasses['dtps']} .results-reference-diff`, res => res.getPlayers()[0]!.dtps, 2, true);
			this.formatToplineResult(`.${RaidSimResultsManager.resultMetricClasses['tmi']} .results-reference-diff`, res => res.getPlayers()[0]!.tmi, 2, true);
			this.formatToplineResult(`.${RaidSimResultsManager.resultMetricClasses['cod']} .results-reference-diff`, res => res.getPlayers()[0]!.chanceOfDeath, 1, true);
		}
	}

	private formatToplineResult(querySelector: string, getMetrics: (result: SimResult) => DistributionMetricsProto | number, precision: number, lowerIsBetter?: boolean) {
		const elem = this.simUI.resultsViewer.contentElem.querySelector(querySelector) as HTMLSpanElement;
		if (!elem) {
			return;
		}

		const cur = this.currentData!.simResult;
		const ref = this.referenceData!.simResult;
		const curMetricsTemp = getMetrics(cur);
		const refMetricsTemp = getMetrics(ref);
		if (typeof curMetricsTemp === 'number') {
			const curMetrics = curMetricsTemp as number;
			const refMetrics = refMetricsTemp as number;
			formatDeltaTextElem(elem, refMetrics, curMetrics, precision, lowerIsBetter);
		} else {
			const curMetrics = curMetricsTemp as DistributionMetricsProto;
			const refMetrics = refMetricsTemp as DistributionMetricsProto;
			const isDiff = this.applyZTestTooltip(elem, ref.iterations, refMetrics.avg, refMetrics.stdev, cur.iterations, curMetrics.avg, curMetrics.stdev);
			formatDeltaTextElem(elem, refMetrics.avg, curMetrics.avg, precision, lowerIsBetter, !isDiff);
		}
	}

	private applyZTestTooltip(elem: HTMLElement, n1: number, avg1: number, stdev1: number, n2: number, avg2: number, stdev2: number): boolean {
		const delta = avg1 - avg2;
		const err1 = stdev1 / Math.sqrt(n1);
		const err2 = stdev2 / Math.sqrt(n2);
		const denom = Math.sqrt(Math.pow(err1, 2) + Math.pow(err2, 2));
		const z = Math.abs(delta / denom);
		const isDiff = z > 1.96;

		let significance_str = '';
		if (isDiff) {
			significance_str = `Difference is significantly different (Z = ${z.toFixed(3)}).`;
		} else {
			significance_str = `Difference is not significantly different (Z = ${z.toFixed(3)}).`;
		}
		tippy(elem, {
			'content': significance_str,
			ignoreAttributes: true,
		});

		return isDiff;
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
		let content = '';

		if (players.length == 1) {
			const playerMetrics = players[0];
			if (playerMetrics.getTargetIndex(filter) == null) {
				const dpsMetrics = playerMetrics.dps;
				const dpaspMetrics = playerMetrics.dpasp;
				const tpsMetrics = playerMetrics.tps;
				const dtpsMetrics = playerMetrics.dtps;
				const tmiMetrics = playerMetrics.tmi;
				content += this.buildResultsLine({
					average: dpsMetrics.avg,
					stdev: dpsMetrics.stdev,
					classes: this.getResultsLineClasses('dps'),
				}).outerHTML;

				// Hide dpasp if it's zero.
				const dpaspContent = this.buildResultsLine({
					average: dpaspMetrics.avg,
					stdev: dpaspMetrics.stdev,
					classes: this.getResultsLineClasses('dpasp'),
				});
				if (dpaspMetrics.avg == 0) {
					dpaspContent.classList.add('hide');
				}
				content += dpaspContent.outerHTML;

				content += this.buildResultsLine({
					average: tpsMetrics.avg,
					stdev: tpsMetrics.stdev,
					classes: this.getResultsLineClasses('tps'),
				}).outerHTML;
				content += this.buildResultsLine({
					average: dtpsMetrics.avg,
					stdev: dtpsMetrics.stdev,
					classes: this.getResultsLineClasses('dtps'),
				}).outerHTML;
				content += this.buildResultsLine({
					average: tmiMetrics.avg,
					stdev: tmiMetrics.stdev,
					classes: this.getResultsLineClasses('tmi'),
				}).outerHTML;
				content += this.buildResultsLine({
					average: playerMetrics.chanceOfDeath,
					classes: this.getResultsLineClasses('cod'),
				}).outerHTML;
			} else {
				const actions = simResult.getActionMetrics(filter);
				if (actions.length > 0) {
					const mergedActions = ActionMetrics.merge(actions);
					content += this.buildResultsLine({
						average: mergedActions.dps,
						classes: this.getResultsLineClasses('dps'),
					}).outerHTML;
					content += this.buildResultsLine({
						average: mergedActions.tps,
						classes: this.getResultsLineClasses('tps'),
					}).outerHTML;
				}

				const targetActions = simResult.getTargets(filter)[0].actions.map(action => action.forTarget(filter));
				if (targetActions.length > 0) {
					const mergedTargetActions = ActionMetrics.merge(targetActions);
					content += this.buildResultsLine({
						average: mergedTargetActions.dps,
						classes: this.getResultsLineClasses('dtps'),
					}).outerHTML;
				}
			}

			content += this.buildResultsLine({
				average: playerMetrics.tto.avg,
				stdev: playerMetrics.tto.stdev,
				classes: this.getResultsLineClasses('tto'),
			}).outerHTML;
			content += this.buildResultsLine({
				average: playerMetrics.hps.avg,
				stdev: playerMetrics.hps.stdev,
				classes: this.getResultsLineClasses('hps'),
			}).outerHTML;
		} else {
			const dpsMetrics = simResult.raidMetrics.dps;
			content += this.buildResultsLine({
				average: dpsMetrics.avg,
				stdev: dpsMetrics.stdev,
				classes: this.getResultsLineClasses('dps'),
			}).outerHTML;
			//const hpsMetrics = simResult.raidMetrics.hps;
			//content += this.buildResultsLine({
			//	average: hpsMetrics.avg,
			//	stdev: hpsMetrics.stdev,
			//	classes: this.getResultsLineClasses('hps'),
			//}).outerHTML;
		}

		if (simResult.request.encounter?.useHealth) {
			content += this.buildResultsLine({
				average: simResult.result.avgIterationDuration,
				classes: this.getResultsLineClasses('dur'),
			});
		}

		return content;
	}

	private static getResultsLineClasses(metric: keyof ResultMetrics): string {
		const classes = [this.resultMetricClasses[metric]];
		if (this.resultMetricCategories[metric])
			classes.push(this.metricsClasses[this.resultMetricCategories[metric]]);

		return classes.join(' ');
	}

	private static buildResultsLine(args: ResultsLineArgs): HTMLElement {
		const resultsFragment = document.createElement('fragment');
		resultsFragment.innerHTML = `
			<div class="results-metric ${args.classes}">
				<span class="topline-result-avg">${args.average.toFixed(2)}</span>
				${args.stdev ? `
					<span class="topline-result-stdev">
						(<i class="fas fa-plus-minus fa-xs"></i>${args.stdev.toFixed()})
					</span>` : ''
			}
				<div class="results-reference hide">
					<span class="results-reference-diff"></span> 对比参考指标
				</div>
			</div>
		`;

		return resultsFragment.children[0] as HTMLElement;
	}

}
