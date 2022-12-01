import { Component } from './components/component.js';
import { NumberPicker } from './components/number_picker.js';
import { ResultsViewer } from './components/results_viewer.js';
import { SimTitleDropdown } from './components/sim_title_dropdown.js';
import { SimHeader } from './components/sim_header';
import { Spec } from './proto/common.js';
import { SimOptions } from './proto/api.js';
import { LaunchStatus } from './launched_sims.js';
import { specToLocalStorageKey } from './proto_utils/utils.js';

import { Sim, SimError } from './sim.js';
import { Target } from './target.js';
import { EventID, TypedEvent } from './typed_event.js';

import { Tooltip } from 'bootstrap';

declare var tippy: any;
declare var pako: any;

const URLMAXLEN = 2048;
const noticeText = '';
// const noticeText = "We're looking for help migrating our sims to Wrath of the Lich King. If you\'d like to participate in a fun side project working with an open-source community please, <a href='https://discord.gg/jJMPr9JWwx' target='_blank'>join our Discord!</a>";

// Config for displaying a warning to the user whenever a condition is met.
export interface SimWarning {
	updateOn: TypedEvent<any>,
	getContent: () => string | Array<string>,
}

export interface SimUIConfig {
	// The spec, if an individual sim, or null if the raid sim.
	spec: Spec | null,
	launchStatus: LaunchStatus,
	knownIssues?: Array<string>,
}

// Shared UI for all individual sims and the raid sim.
export abstract class SimUI extends Component {
	readonly sim: Sim;
	readonly isWithinRaidSim: boolean;

	// Emits when anything from the sim, raid, or encounter changes.
	readonly changeEmitter;

	readonly resultsViewer: ResultsViewer

	protected readonly simHeader: SimHeader;

	protected readonly simContentContainer: HTMLElement;
	protected readonly simMain: HTMLElement;
	protected readonly simActionsContainer: HTMLElement;
	protected readonly iterationsPicker: HTMLElement;
	protected readonly simTabContentsContainer: HTMLElement;

	private warningsTippy: any;

	constructor(parentElem: HTMLElement, sim: Sim, config: SimUIConfig) {
		super(parentElem, 'sim-ui');
		this.sim = sim;
		this.rootElem.innerHTML = simHTML;
		this.simContentContainer = this.rootElem.querySelector('#simContent') as HTMLElement;
		this.simHeader = new SimHeader(this.simContentContainer, this);
		this.simMain = simMain;
		this.simContentContainer.appendChild(simMain);
		this.isWithinRaidSim = this.rootElem.closest('.within-raid-sim') != null;

		if (!this.isWithinRaidSim) {
			this.rootElem.classList.add('not-within-raid-sim');
		}

		this.changeEmitter = TypedEvent.onAny([
			this.sim.changeEmitter,
		], 'SimUIChange');

		this.sim.crashEmitter.on((eventID: EventID, error: SimError) => this.handleCrash(error));

		const updateShowDamageMetrics = () => {
			if (this.sim.getShowDamageMetrics()) {
				this.rootElem.classList.remove('hide-damage-metrics');
			} else {
				this.rootElem.classList.add('hide-damage-metrics');
			}
		};
		updateShowDamageMetrics();
		this.sim.showDamageMetricsChangeEmitter.on(updateShowDamageMetrics);

		const updateShowThreatMetrics = () => {
			if (this.sim.getShowThreatMetrics()) {
				this.rootElem.classList.remove('hide-threat-metrics');
			} else {
				this.rootElem.classList.add('hide-threat-metrics');
			}
		};
		updateShowThreatMetrics();
		this.sim.showThreatMetricsChangeEmitter.on(updateShowThreatMetrics);

		const updateShowHealingMetrics = () => {
			if (this.sim.getShowHealingMetrics()) {
				this.rootElem.classList.remove('hide-healing-metrics');
			} else {
				this.rootElem.classList.add('hide-healing-metrics');
			}
		};
		updateShowHealingMetrics();
		this.sim.showHealingMetricsChangeEmitter.on(updateShowHealingMetrics);

		const updateShowExperimental = () => {
			if (this.sim.getShowExperimental()) {
				this.rootElem.classList.remove('hide-experimental');
			} else {
				this.rootElem.classList.add('hide-experimental');
			}
		};
		updateShowExperimental();
		this.sim.showExperimentalChangeEmitter.on(updateShowExperimental);

		this.addNoticeBanner();
		this.addKnownIssues(config);

		const titleElem = this.rootElem.querySelector('#simTitle') as HTMLElement;
		new SimTitleDropdown(titleElem, config.spec);

		const resultsViewerElem = this.rootElem.getElementsByClassName('sim-sidebar-results')[0] as HTMLElement;
		this.resultsViewer = new ResultsViewer(resultsViewerElem);

		this.simActionsContainer = this.rootElem.getElementsByClassName('sim-sidebar-actions')[0] as HTMLElement;

		new NumberPicker(this.simActionsContainer, this.sim, {
			label: 'Iterations',
			extraCssClasses: [
				'iterations-picker',
				'within-raid-sim-hide',
			],
			changedEvent: (sim: Sim) => sim.iterationsChangeEmitter,
			getValue: (sim: Sim) => sim.getIterations(),
			setValue: (eventID: EventID, sim: Sim, newValue: number) => {
				sim.setIterations(eventID, newValue);
			},
		});

		this.iterationsPicker = this.rootElem.getElementsByClassName('iterations-picker')[0] as HTMLElement;
		this.simTabContentsContainer = this.rootElem.querySelector('#simMain.tab-content') as HTMLElement;

		if (!this.isWithinRaidSim) {
			window.addEventListener('message', async event => {
				if (event.data == 'runOnce') {
					this.runSimOnce();
				}
			});
		}
	}

	addAction(name: string, cssClass: string, actFn: () => void) {
		const button = document.createElement('button');
		button.classList.add('sim-sidebar-actions-button', 'btn', 'btn-outline-primary', cssClass);
		button.textContent = name;
		button.addEventListener('click', actFn);
		this.simActionsContainer.appendChild(button);
	}

	addTab(title: string, cssClass: string, innerHTML: string) {
		const contentId = cssClass.replace(/\s+/g, '-') + '-tab';
		const isFirstTab = this.simTabContentsContainer.children.length == 0;

		this.simHeader.addTab(title, contentId);

		const tabContentFragment = document.createElement('fragment');
		tabContentFragment.innerHTML = `
			<div
				id="${contentId}"
				class="tab-pane fade ${isFirstTab ? 'active show' : ''}"
			>${innerHTML}</div>
		`;
		this.simTabContentsContainer.appendChild(tabContentFragment.children[0] as HTMLElement);
	}

	addWarning(warning: SimWarning) {
		this.simHeader.addWarning(warning);
	}

	protected addImportLink(link: HTMLElement) {
		this.simHeader.addImportLink(link);
	}

	protected addExportLink(link: HTMLElement) {
		this.simHeader.addExportLink(link);
	}

	private addNoticeBanner() {
		const noticesElem = document.querySelector('#noticesBanner') as HTMLElement;

		if (!noticeText) {
			noticesElem.remove();
		}
	}

	private addKnownIssues(config: SimUIConfig) {
		let statusStr = '';
		if (config.launchStatus == LaunchStatus.Unlaunched) {
			statusStr = 'This sim is a WORK IN PROGRESS. It is not fully developed and should not be used for general purposes.';
		} else if (config.launchStatus == LaunchStatus.Alpha) {
			statusStr = 'This sim is in ALPHA. Bugs are expected; please let us know if you find one!';
		} else if (config.launchStatus == LaunchStatus.Beta) {
			statusStr = 'This sim is in BETA. There may still be a few bugs; please let us know if you find one!';
		}
		if (statusStr) {
			config.knownIssues = [statusStr].concat(config.knownIssues || []);
		}
		if (config.knownIssues && config.knownIssues.length) {
			config.knownIssues.forEach(issue => this.simHeader.addKnownIssue(issue));
		}
	}

	// Returns a key suitable for the browser's localStorage feature.
	abstract getStorageKey(postfix: string): string;

	getSettingsStorageKey(): string {
		return this.getStorageKey('__currentSettings__');
	}

	getSavedEncounterStorageKey(): string {
		// By skipping the call to this.getStorageKey(), saved encounters will be
		// shared across all sims.
		return 'sharedData__savedEncounter__';
	}

	isIndividualSim(): boolean {
		return this.rootElem.classList.contains('individual-sim-ui');
	}

	async runSim(onProgress: Function) {
		this.resultsViewer.setPending();
		try {
			await this.sim.runRaidSim(TypedEvent.nextEventID(), onProgress);
		} catch (e) {
			this.resultsViewer.hideAll();
			this.handleCrash(e);
		}
	}

	async runSimOnce() {
		this.resultsViewer.setPending();
		try {
			await this.sim.runRaidSimWithLogs(TypedEvent.nextEventID());
		} catch (e) {
			this.resultsViewer.hideAll();
			this.handleCrash(e);
		}
	}

	handleCrash(error: any) {
		if (!(error instanceof SimError)) {
			alert(error);
			return;
		}

		const errorStr = (error as SimError).errorStr;
		if (window.confirm('Simulation Failure:\n' + errorStr + '\nPress Ok to file crash report')) {
			// Splice out just the line numbers
			const hash = this.hashCode(errorStr);
			const link = this.toLink();
			const rngSeed = this.sim.getLastUsedRngSeed();
			fetch('https://api.github.com/search/issues?q=is:issue+is:open+repo:wowsims/wotlk+' + hash).then(resp => {
				resp.json().then((issues) => {
					if (issues.total_count > 0) {
						window.open(issues.items[0].html_url, '_blank');
					} else {
						const base_url = 'https://github.com/wowsims/wotlk/issues/new?assignees=&labels=&title=Crash%20Report%20'
						const base = `${base_url}${hash}&body=`;
						const maxBodyLength = URLMAXLEN - base.length;
						let issueBody = encodeURIComponent(`Link:\n${link}\n\nRNG Seed: ${rngSeed}\n\n${errorStr}`)
						while (issueBody.length > maxBodyLength) {
							issueBody = issueBody.slice(0, issueBody.lastIndexOf('%')) // Avoid truncating in the middle of a URLencoded segment
						}
						window.open(base + issueBody, '_blank');
					}
				});
			}).catch(fetchErr => {
				alert('Failed to file report... try again another time:' + fetchErr);
			});
		}
		return;
	}

	hashCode(str: string): number {
		let hash = 0;
		for (let i = 0, len = str.length; i < len; i++) {
			let chr = str.charCodeAt(i);
			hash = (hash << 5) - hash + chr;
			hash |= 0; // Convert to 32bit integer
		}
		return hash;
	}

	abstract applyDefaults(eventID: EventID): void;
	abstract toLink(): string;
}

const simHTML = `
<div class="sim-root">
	<div id="noticesBanner" class="alert border-bottom mb-0 text-center">${noticeText}</div>
  <aside id="simSidebar">
    <div id="simTitle"></div>
		<div id="simSidebarContent">
			<div class="sim-sidebar-actions within-raid-sim-hide"></div>
			<div class="sim-sidebar-results within-raid-sim-hide"></div>
			<div class="sim-sidebar-footer"></div>
		</div>
  </aside>
  <div id="simContent" class="container-fluid">
	</div>
  </section>
</div>
`;

// TODO: Build SimMain component instead
let simMainFragment = document.createElement('fragment');
simMainFragment.innerHTML = `<main id="simMain" class="tab-content"></main>`;
const simMain = simMainFragment.children[0] as HTMLElement;
