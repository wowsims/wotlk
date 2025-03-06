import { TypedEvent } from '../../core/typed_event';
import { REPO_NAME } from '../constants/other';
import { DetailedResultsUpdate, SimRun, SimRunData } from '../proto/ui';
import { SimResult } from '../proto_utils/sim_result';
import { SimUI } from '../sim_ui';
import { Component } from './component';
import { AuraMetricsTable } from './detailed_results/aura_metrics';
import { CastMetricsTable } from './detailed_results/cast_metrics';
import { DpsHistogram } from './detailed_results/dps_histogram';
import { DtpsMeleeMetricsTable } from './detailed_results/dtps_melee_metrics';
import { DtpsSpellMetricsTable } from './detailed_results/dtps_spell_metrics';
import { HealingMetricsTable } from './detailed_results/healing_metrics';
import { LogRunner } from './detailed_results/log_runner';
import { MeleeMetricsTable } from './detailed_results/melee_metrics';
import { PlayerDamageMetricsTable } from './detailed_results/player_damage';
import { ResourceMetricsTable } from './detailed_results/resource_metrics';
import { SimResultData } from './detailed_results/result_component';
import { ResultsFilter } from './detailed_results/results_filter';
import { SpellMetricsTable } from './detailed_results/spell_metrics';
import { Timeline } from './detailed_results/timeline';
import { ToplineResults } from './detailed_results/topline_results';
import { RaidSimResultsManager } from './raid_sim_action';

const layoutHTML = `
<div class="dr-root dr-no-results">
	<div class="dr-toolbar">
		<div class="results-filter"></div>
		<div class="tabs-filler"></div>
		<ul class="nav nav-tabs" role="tablist">
			<li class="nav-item dr-tab-tab damage-metrics" role="presentation">
				<a
					class="nav-link active"
					data-bs-toggle="tab"
					data-bs-target="#damageTab"
					type="button"
					role="tab"
					aria-controls="damageTab"
					aria-selected="true"
				>Damage</a>
			</li>
			<li class="nav-item dr-tab-tab healing-metrics" role="presentation">
				<a
					class="nav-link"
					data-bs-toggle="tab"
					data-bs-target="#healingTab"
					type="button"
					role="tab"
					aria-controls="healingTab"
					aria-selected="false"
				>Healing</a>
			</li>
			<li class="nav-item dr-tab-tab threat-metrics" role="presentation">
				<a
					class="nav-link"
					data-bs-toggle="tab"
					data-bs-target="#damageTakenTab"
					type="button"
					role="tab"
					aria-controls="damageTakenTab"
					aria-selected="false"
				>Damage Taken</a>
			</li>
			<li class="nav-item dr-tab-tab" role="presentation">
				<a
					class="nav-link"
					data-bs-toggle="tab"
					data-bs-target="#buffsTab"
					type="button"
					role="tab"
					aria-controls="buffsTab"
					aria-selected="false"
				>Buffs</a>
			</li>
			<li class="nav-item dr-tab-tab" role="presentation">
				<a
					class="nav-link"
					data-bs-toggle="tab"
					data-bs-target="#debuffsTab"
					type="button"
					role="tab"
					aria-controls="debuffsTab"
					aria-selected="false"
				>Debuffs</a>
			</li>
			<li class="nav-item dr-tab-tab" role="presentation">
				<a
					class="nav-link"
					data-bs-toggle="tab"
					data-bs-target="#castsTab"
					type="button"
					role="tab"
					aria-controls="castsTab"
					aria-selected="false"
				>Casts</a>
			</li>
			<li class="nav-item dr-tab-tab" role="presentation">
				<a
					class="nav-link"
					data-bs-toggle="tab"
					data-bs-target="#resourcesTab"
					type="button"
					role="tab"
					aria-controls="resourcesTab"
					aria-selected="false"
				>Resources</a>
			</li>
			<li class="nav-item dr-tab-tab" role="presentation">
				<a
					id="timelineTabTab"
					class="nav-link"
					data-bs-toggle="tab"
					data-bs-target="#timelineTab"
					type="button"
					role="tab"
					aria-controls="timelineTab"
					aria-selected="false"
				>Timeline</a>
			<li class="nav-item dr-tab-tab" role="presentation">
				<a
					id="logTabTab"
					class="nav-link"
					data-bs-toggle="tab"
					data-bs-target="#logTab"
					type="button"
					role="tab"
					aria-controls="logTab"
					aria-selected="false"
				>Log</a>
			</li>
		</ul>
	</div>
	<div class="tab-content">
		<div id="noResultsTab" class="tab-pane dr-tab-content fade active show">
			Run a simulation to view results
		</div>
		<div id="damageTab" class="tab-pane dr-tab-content damage-content fade active show">
			<div class="dr-row topline-results">
			</div>
			<div class="dr-row all-players-only">
				<div class="player-damage-metrics">
				</div>
			</div>
			<div class="dr-row single-player-only">
				<div class="melee-metrics">
				</div>
			</div>
			<div class="dr-row single-player-only">
				<div class="spell-metrics">
				</div>
			</div>
			<div class="dr-row dps-histogram">
			</div>
		</div>
		<div id="healingTab" class="tab-pane dr-tab-content healing-content fade">
			<div class="dr-row topline-results">
			</div>
			<div class="dr-row single-player-only">
				<div class="healing-spell-metrics">
				</div>
			</div>
			<div class="dr-row hps-histogram">
			</div>
		</div>
		<div id="damageTakenTab" class="tab-pane dr-tab-content damage-taken-content fade">
			<div class="dr-row single-player-only">
				<div class="dtps-melee-metrics">
				</div>
			</div>
			<div class="dr-row single-player-only">
				<div class="dtps-spell-metrics">
				</div>
			</div>
			<div class="dr-row damage-taken-histogram single-player-only">
			</div>
		</div>
		<div id="buffsTab" class="tab-pane dr-tab-content buffs-content fade">
			<div class="dr-row">
				<div class="buff-aura-metrics">
				</div>
			</div>
		</div>
		<div id="debuffsTab" class="tab-pane dr-tab-content debuffs-content fade">
			<div class="dr-row">
				<div class="debuff-aura-metrics">
				</div>
			</div>
		</div>
		<div id="castsTab" class="tab-pane dr-tab-content casts-content fade">
			<div class="dr-row">
				<div class="cast-metrics">
				</div>
			</div>
		</div>
		<div id="resourcesTab" class="tab-pane dr-tab-content resources-content fade">
			<div class="dr-row">
				<div class="resource-metrics">
				</div>
			</div>
		</div>
		<div id="timelineTab" class="tab-pane dr-tab-content timeline-content fade">
			<div class="dr-row">
				<div class="timeline">
				</div>
			</div>
		</div>
		<div id="logTab" class="tab-pane dr-tab-content log-content fade">
			<div class="dr-row">
				<div class="log">
				</div>
			</div>
		</div>
	</div>
</div>
`;

export abstract class DetailedResults extends Component {
	protected readonly simUI: SimUI | null;
	protected latestRun: SimRunData | null = null;

	private currentSimResult: SimResult | null = null;
	private resultsEmitter: TypedEvent<SimResultData | null> = new TypedEvent<SimResultData | null>();
	private resultsFilter: ResultsFilter;

	constructor(parent: HTMLElement, simUI: SimUI | null, cssScheme: string) {
		super(parent, 'detailed-results-manager-root');
		this.rootElem.innerHTML = layoutHTML;
		this.simUI = simUI;

		this.simUI?.sim.settingsChangeEmitter.on(async () => await this.updateSettings());

		// Allow styling the sticky toolbar
		const toolbar = document.querySelector('.dr-toolbar') as HTMLElement;
		new IntersectionObserver(
			([e]) => {
				//console.log(e.intersectionRatio)
				e.target.classList.toggle('stuck', e.intersectionRatio < 1);
			},
			{
				// Intersect with the sim header or top of the separate tab
				rootMargin: this.simUI ? `-${this.simUI.simHeader.rootElem.offsetHeight + 1}px 0px 0px 0px` : '0px',
				threshold: [1],
			},
		).observe(toolbar);

		this.resultsFilter = new ResultsFilter({
			parent: this.rootElem.getElementsByClassName('results-filter')[0] as HTMLElement,
			resultsEmitter: this.resultsEmitter,
		});

		(Array.from(this.rootElem.getElementsByClassName('topline-results')) as Array<HTMLElement>).forEach(toplineResultsDiv => {
			new ToplineResults({ parent: toplineResultsDiv, resultsEmitter: this.resultsEmitter });
		});

		const castMetrics = new CastMetricsTable({
			parent: this.rootElem.getElementsByClassName('cast-metrics')[0] as HTMLElement,
			resultsEmitter: this.resultsEmitter,
		});
		const meleeMetrics = new MeleeMetricsTable({
			parent: this.rootElem.getElementsByClassName('melee-metrics')[0] as HTMLElement,
			resultsEmitter: this.resultsEmitter,
		});
		const spellMetrics = new SpellMetricsTable({
			parent: this.rootElem.getElementsByClassName('spell-metrics')[0] as HTMLElement,
			resultsEmitter: this.resultsEmitter,
		});
		const healingMetrics = new HealingMetricsTable({
			parent: this.rootElem.getElementsByClassName('healing-spell-metrics')[0] as HTMLElement,
			resultsEmitter: this.resultsEmitter,
		});
		const resourceMetrics = new ResourceMetricsTable({
			parent: this.rootElem.getElementsByClassName('resource-metrics')[0] as HTMLElement,
			resultsEmitter: this.resultsEmitter,
		});
		const playerDamageMetrics = new PlayerDamageMetricsTable(
			{ parent: this.rootElem.getElementsByClassName('player-damage-metrics')[0] as HTMLElement, resultsEmitter: this.resultsEmitter },
			this.resultsFilter,
		);
		const buffAuraMetrics = new AuraMetricsTable(
			{
				parent: this.rootElem.getElementsByClassName('buff-aura-metrics')[0] as HTMLElement,
				resultsEmitter: this.resultsEmitter,
			},
			false,
		);
		const debuffAuraMetrics = new AuraMetricsTable(
			{
				parent: this.rootElem.getElementsByClassName('debuff-aura-metrics')[0] as HTMLElement,
				resultsEmitter: this.resultsEmitter,
			},
			true,
		);
		const dpsHistogram = new DpsHistogram({
			parent: this.rootElem.getElementsByClassName('dps-histogram')[0] as HTMLElement,
			resultsEmitter: this.resultsEmitter,
		});

		const dtpsMeleeMetrics = new DtpsMeleeMetricsTable({
			parent: this.rootElem.getElementsByClassName('dtps-melee-metrics')[0] as HTMLElement,
			resultsEmitter: this.resultsEmitter,
		});
		const dtpsSpellMetrics = new DtpsSpellMetricsTable({
			parent: this.rootElem.getElementsByClassName('dtps-spell-metrics')[0] as HTMLElement,
			resultsEmitter: this.resultsEmitter,
		});

		const timeline = new Timeline({
			parent: this.rootElem.getElementsByClassName('timeline')[0] as HTMLElement,
			cssScheme: cssScheme,
			resultsEmitter: this.resultsEmitter,
		});
		document.getElementById('timelineTabTab')?.addEventListener('click', event => timeline.render());

		const log = new LogRunner({
			parent: this.rootElem.getElementsByClassName('log')[0] as HTMLElement,
			cssScheme: cssScheme,
			resultsEmitter: this.resultsEmitter,
		});

		this.rootElem.classList.add('hide-threat-metrics');
		this.rootElem.classList.add('hide-healing-metrics');

		this.resultsFilter.changeEmitter.on(() => this.updateResults());

		const rootDiv = this.rootElem.getElementsByClassName('dr-root')[0] as HTMLElement;
		this.resultsEmitter.on((eventID, resultData) => {
			if (resultData?.filter.player || resultData?.filter.player === 0) {
				rootDiv.classList.remove('all-players');
				rootDiv.classList.add('single-player');
			} else {
				rootDiv.classList.add('all-players');
				rootDiv.classList.remove('single-player');
			}
		});
	}

	abstract postMessage(update: DetailedResultsUpdate): Promise<void>;

	protected async setSimRunData(simRunData: SimRunData) {
		this.latestRun = simRunData;
		await this.postMessage(
			DetailedResultsUpdate.create({
				data: {
					oneofKind: 'runData',
					runData: simRunData,
				},
			}),
		);
	}

	protected async updateSettings() {
		if (!this.simUI) return;
		await this.postMessage(
			DetailedResultsUpdate.create({
				data: {
					oneofKind: 'settings',
					settings: this.simUI.sim.toProto(),
				},
			}),
		);
	}

	private updateResults() {
		const eventID = TypedEvent.nextEventID();
		if (this.currentSimResult == null) {
			this.rootElem.querySelector('.dr-root')?.classList.add('dr-no-results');
			this.resultsEmitter.emit(eventID, null);
		} else {
			this.rootElem.querySelector('.dr-root')?.classList.remove('dr-no-results');
			this.resultsEmitter.emit(eventID, {
				eventID: eventID,
				result: this.currentSimResult,
				filter: this.resultsFilter.getFilter(),
			});
		}
	}

	protected async handleMessage(data: DetailedResultsUpdate) {
		switch (data.data.oneofKind) {
			case 'runData':
				const runData = data.data.runData;
				this.currentSimResult = await SimResult.fromProto(runData.run || SimRun.create());
				this.updateResults();
				break;
			case 'settings':
				const settings = data.data.settings;
				if (settings.showDamageMetrics) {
					this.rootElem.classList.remove('hide-damage-metrics');
				} else {
					this.rootElem.classList.add('hide-damage-metrics');
					if (document.getElementById('damageTab')!.classList.contains('active')) {
						document.getElementById('damageTab')!.classList.remove('active', 'show');
						document.getElementById('healingTab')!.classList.add('active', 'show');

						const toolbar = document.getElementsByClassName('dr-toolbar')[0] as HTMLElement;
						(toolbar.getElementsByClassName('damage-metrics')[0] as HTMLElement).children[0]!.classList.remove('active');
						(toolbar.getElementsByClassName('healing-metrics')[0] as HTMLElement).children[0]!.classList.add('active');
					}
				}
				if (settings.showThreatMetrics) {
					this.rootElem.classList.remove('hide-threat-metrics');
				} else {
					this.rootElem.classList.add('hide-threat-metrics');
				}
				if (settings.showHealingMetrics) {
					this.rootElem.classList.remove('hide-healing-metrics');
				} else {
					this.rootElem.classList.add('hide-healing-metrics');
				}
				if (settings.showExperimental) {
					this.rootElem.classList.remove('hide-experimental');
				} else {
					this.rootElem.classList.add('hide-experimental');
				}
				break;
		}
	}
}

export class WindowedDetailedResults extends DetailedResults {
	constructor(parent: HTMLElement) {
		super(parent, null, new URLSearchParams(window.location.search).get('cssScheme') ?? '');

		window.addEventListener('message', async event => await this.handleMessage(DetailedResultsUpdate.fromJson(event.data)));

		this.rootElem.insertAdjacentHTML(
			'beforeend',
			`
			<div class="sim-bg"></div>
		`,
		);
	}

	async postMessage(update: DetailedResultsUpdate): Promise<void> {
		await this.handleMessage(update);
	}
}

export class EmbeddedDetailedResults extends DetailedResults {
	private tabWindow: Window | null = null;

	constructor(parent: HTMLElement, simUI: SimUI, simResultsManager: RaidSimResultsManager) {
		super(parent, simUI, simUI.cssScheme);

		const newTabBtn = document.createElement('div');
		newTabBtn.classList.add('detailed-results-controls-div');
		newTabBtn.innerHTML = `
			<button class="detailed-results-new-tab-button btn btn-primary">View in Separate Tab</button>
			<button class="detailed-results-1-iteration-button btn btn-primary">Sim 1 Iteration</button>
		`;

		this.rootElem.prepend(newTabBtn);

		const computedStyles = window.getComputedStyle(this.rootElem);

		const url = new URL(`${window.location.protocol}//${window.location.host}/${REPO_NAME}/detailed_results/index.html`);
		url.searchParams.append('cssClass', simUI.cssClass);

		if (simUI.isIndividualSim()) {
			url.searchParams.append('isIndividualSim', '');
			this.rootElem.classList.add('individual-sim');
		}

		const newTabButton = this.rootElem.querySelector('.detailed-results-new-tab-button') as HTMLButtonElement;
		newTabButton.addEventListener('click', event => {
			if (this.tabWindow == null || this.tabWindow.closed) {
				this.tabWindow = window.open(url.href, 'Detailed Results');
				this.tabWindow!.addEventListener('load', async event => {
					if (this.latestRun) {
						await this.updateSettings();
						await this.setSimRunData(this.latestRun);
					}
				});
			} else {
				this.tabWindow.focus();
			}
		});

		const simButton = this.rootElem.querySelector('.detailed-results-1-iteration-button') as HTMLButtonElement;
		simButton.addEventListener('click', () => {
			(window.opener || window.parent)!.postMessage('runOnce', '*');
		});

		simResultsManager.currentChangeEmitter.on(async () => {
			const runData = simResultsManager.getRunData();
			if (runData) {
				await this.updateSettings();
				await this.setSimRunData(runData);
			}
		});
	}

	async postMessage(update: DetailedResultsUpdate) {
		if (this.tabWindow) {
			this.tabWindow.postMessage(DetailedResultsUpdate.toJson(update), '*');
		}
		await this.handleMessage(update);
	}
}
