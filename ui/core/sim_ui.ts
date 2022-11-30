import { Component } from './components/component.js';
import { NumberPicker } from './components/number_picker.js';
import { ResultsViewer } from './components/results_viewer.js';
import { SimTitleDropdown } from './components/sim_title_dropdown.js';
import { Spec } from './proto/common.js';
import { SimOptions } from './proto/api.js';
import { LaunchStatus } from './launched_sims.js';
import { specToLocalStorageKey } from './proto_utils/utils.js';

import { Sim, SimError } from './sim.js';
import { Target } from './target.js';
import { EventID, TypedEvent } from './typed_event.js';

declare var tippy: any;
declare var pako: any;

const URLMAXLEN = 2048;
const noticeText = '';
//const noticeText = 'We are looking for help migrating our sims to Wrath of the Lich King. If you\'d like to participate in a fun side project working with an open-source community please <a href="https://discord.gg/jJMPr9JWwx" target="_blank">join our discord!</a>';

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

	protected simActionsContainer: HTMLElement;
	protected iterationsPicker: HTMLElement;

	protected simTabsContainer: HTMLElement;
	protected simTabContentsContainer: HTMLElement;
	protected importExportContainer: HTMLElement;
	protected simToolbar: HTMLElement;

	private warnings: Array<SimWarning>;
	private warningsTippy: any;

	constructor(parentElem: HTMLElement, sim: Sim, config: SimUIConfig) {
		super(parentElem, 'sim-ui');
		this.sim = sim;
		this.rootElem.innerHTML = simHTML;
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

		const noticesElem = document.getElementsByClassName('notices')[0] as HTMLElement;
		if (noticeText) {
			tippy(noticesElem, {
				content: noticeText,
				allowHTML: true,
				interactive: true,
			});
		} else {
			noticesElem.remove();
		}

		this.warnings = [];
		const warningsElem = document.getElementsByClassName('warnings')[0] as HTMLElement;
		this.warningsTippy = tippy(warningsElem, {
			content: '',
			allowHTML: true,
		});
		this.updateWarnings();

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
			const knownIssuesContainer = document.getElementsByClassName('known-issues')[0] as HTMLElement;
			knownIssuesContainer.style.display = 'initial';
			tippy(knownIssuesContainer, {
				content: `
				<ul class="known-issues-tooltip">
					${config.knownIssues.map(issue => '<li>' + issue + '</li>').join('')}
				</ul>
				`,
				allowHTML: true,
				interactive: true,
			});
		}

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

		this.simTabsContainer = this.rootElem.querySelector('#simHeader .sim-tabs') as HTMLElement;
		this.simTabContentsContainer = this.rootElem.querySelector('#simMain.tab-content') as HTMLElement;
		this.importExportContainer = this.rootElem.querySelector('#simHeader .import-export') as HTMLElement;
		this.simToolbar = this.rootElem.querySelector('#simHeader .sim-toolbar') as HTMLElement;

		const reportBug = document.createElement('span');
		reportBug.classList.add('report-bug', 'fa', 'fa-bug');
		tippy(reportBug, {
			'content': 'Report a bug / request a feature',
			'allowHTML': true,
		});
		reportBug.addEventListener('click', event => {
			window.open('https://github.com/wowsims/wotlk/issues/new/choose', '_blank');
		});
		this.addToolbarItem(reportBug);

		if (!this.isWithinRaidSim) {
			window.addEventListener('message', async event => {
				if (event.data == 'runOnce') {
					this.runSimOnce();
				}
			});
		}

		const patreon = document.createElement('span');
		patreon.classList.add('patreon-link', 'fa', 'fa-brands', 'fa-patreon');
		tippy(patreon, {
			'content': 'Support our devs on Patreon',
			'allowHTML': true,
		});
		patreon.addEventListener('click', event => {
			window.open('https://patreon.com/wowsims', '_blank');
		});
		this.addToolbarItem(patreon);

		const downloadBinary = document.createElement('span');
		// downloadBinary.src = "/wotlk/assets/img/gauge.svg"
		downloadBinary.classList.add('downbin', 'hide');
		downloadBinary.addEventListener('click', event => {
			window.open('https://github.com/wowsims/wotlk/releases', '_blank');
		});

		if (document.location.href.includes("localhost")) {
			fetch(document.location.protocol + "//" + document.location.host + "/version").then(resp => {
				resp.json()
					.then((versionInfo) => {
						if (versionInfo.outdated == 2) {
							tippy(downloadBinary, {
								'content': 'Newer version of simulator available for download',
								'allowHTML': true,
							});
							downloadBinary.classList.add('downbinalert');
							this.addToolbarItem(downloadBinary);
						}
					})
					.catch(error => {
						console.warn('No version info found!');
					});
			});
		} else {
			tippy(downloadBinary, {
				'content': 'Download simulator for faster simulating',
				'allowHTML': true,
			});
			downloadBinary.classList.add('downbinnorm');
			this.addToolbarItem(downloadBinary);
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
		const isFirstTab = this.simTabsContainer.children.length == 0;

		const tabFragment = document.createElement('fragment');
		tabFragment.innerHTML = `
			<li class="${cssClass}-tab nav-item" role="presentation">
				<a
					class="nav-link ${isFirstTab ? 'active' : ''}"
					data-bs-toggle="tab"
					data-bs-target="#${contentId}"
					type="button"
					role="tab"
					aria-controls="${contentId}"
					aria-selected="${isFirstTab}"
				>${title}</a>
			</li>
		`;

		const tabContentFragment = document.createElement('fragment');
		tabContentFragment.innerHTML = `
			<div
				id="${contentId}"
				class="tab-pane fade ${isFirstTab ? 'active show' : ''}"
			>${innerHTML}</div>
		`;

		this.simTabsContainer.appendChild(tabFragment.children[0] as HTMLElement);
		this.simTabContentsContainer.appendChild(tabContentFragment.children[0] as HTMLElement);
	}

	addToolbarItem(elem: HTMLElement) {
		const toolbarItem = document.createElement('div');
		toolbarItem.appendChild(elem);
		toolbarItem.classList.add('sim-toolbar-item');
		this.simToolbar.appendChild(toolbarItem);
	}

	private updateWarnings() {
		const activeWarnings = this.warnings.map(warning => warning.getContent()).flat().filter(content => content != '');

		const warningsElem = document.getElementsByClassName('warnings')[0] as HTMLElement;
		if (activeWarnings.length == 0) {
			warningsElem.style.display = 'none';
		} else {
			warningsElem.style.display = 'initial';
			this.warningsTippy.setContent(`
				<ul class="known-issues-tooltip">
					${activeWarnings.map(content => '<li>' + content + '</li>').join('')}
				</ul>`
			);
		}
	}

	addWarning(warning: SimWarning) {
		this.warnings.push(warning);
		warning.updateOn.on(() => this.updateWarnings());
		this.updateWarnings();
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
  <aside id="simSidebar"">
    <div id="simTitle"></div>
		<div id="simSidebarContent">
			<div class="sim-sidebar-actions within-raid-sim-hide"></div>
			<div class="sim-sidebar-results within-raid-sim-hide"></div>
			<div class="sim-sidebar-footer"></div>
		</div>
  </aside>
  <div id="simContent" class="container-fluid">
		<header id="simHeader">
			<ul class="sim-tabs nav nav-tabs" role="tablist"></ul>
			<div class="import-export"></div>
			<div class="sim-toolbar">
				<div class="sim-toolbar-item">
					<span class="notices fas fa-exclamation-circle"></span>
				</div>
				<div class="sim-toolbar-item">
					<span class="warnings fa fa-exclamation-triangle"></span>
				</div>
				<div class="sim-toolbar-item">
					<div class="known-issues">Known Issues</div>
				</div>
			</div>
    </header>
    <main id="simMain" class="tab-content">
    </main>
  </section>
</div>
`;
