import { BaseModal } from './components/base_modal.js';
import { Component } from './components/component.js';
import { NumberPicker } from './components/number_picker.js';
import { ResultsViewer } from './components/results_viewer.js';
import { SimHeader } from './components/sim_header';
import { SimTab } from './components/sim_tab.js';
import { SimTitleDropdown } from './components/sim_title_dropdown.js';
import { SocialLinks } from './components/social_links.jsx';
import { LaunchStatus } from './launched_sims.js';
import { Spec } from './proto/common.js';
import { ActionId } from './proto_utils/action_id.js';
import { Sim, SimError } from './sim.js';
import { EventID, TypedEvent } from './typed_event.js';

const URLMAXLEN = 2048;
const globalKnownIssues: Array<string> = [];

// Config for displaying a warning to the user whenever a condition is met.
export interface SimWarning {
	updateOn: TypedEvent<any>;
	getContent: () => string | Array<string>;
}

export interface SimUIConfig {
	// Additional css class to add to the root element.
	cssClass: string;
	// Scheme used for themeing on a per-class Basis or for other sims
	cssScheme: string;
	// The spec, if an individual sim, or null if the raid sim.
	spec: Spec | null;
	launchStatus: LaunchStatus;
	knownIssues?: Array<string>;
	noticeText?: string;
}

// Shared UI for all individual sims and the raid sim.
export abstract class SimUI extends Component {
	readonly sim: Sim;
	readonly cssClass: string;
	readonly cssScheme: string;
	readonly isWithinRaidSim: boolean;

	// Emits when anything from the sim, raid, or encounter changes.
	readonly changeEmitter;

	readonly resultsViewer: ResultsViewer;
	readonly simHeader: SimHeader;

	readonly simContentContainer: HTMLElement;
	readonly simMain: HTMLElement;
	readonly simActionsContainer: HTMLElement;
	readonly iterationsPicker: HTMLElement;
	readonly simTabContentsContainer: HTMLElement;

	constructor(parentElem: HTMLElement, sim: Sim, config: SimUIConfig) {
		super(parentElem, 'sim-ui');
		this.sim = sim;
		this.cssClass = config.cssClass;
		this.cssScheme = config.cssScheme;
		this.isWithinRaidSim = this.rootElem.closest('.within-raid-sim') != null;
		// config.noticeText = `Cataclysm sim development has begun! <a href="https://discord.gg/p3DgvmnDCS" target="_blank">Join our Discord</a> for more details or to contribute to the project.`;
		this.rootElem.innerHTML = `
			<div class="sim-root">
				<div class="sim-bg"></div>
				${
					config.noticeText
						? `<div class="notices-banner alert border-bottom mb-0 text-center">${config.noticeText}</div>`
						: ''
				}
				<div class="sim-container">
					<aside class="sim-sidebar">
						<div class="sim-title"></div>
						<div class="sim-sidebar-content ${config.noticeText ? 'smaller' : ''}">
							<div class="sim-sidebar-actions within-raid-sim-hide"></div>
							<div class="sim-sidebar-results within-raid-sim-hide"></div>
							<div class="sim-sidebar-stats"></div>
							<div class="sim-sidebar-socials"></div>
						</div>
					</aside>
					<div class="sim-content container-fluid"></div>
				</div>
			</div>
		`;
		this.simContentContainer = this.rootElem.querySelector('.sim-content') as HTMLElement;
		this.simHeader = new SimHeader(this.simContentContainer, this);
		this.simMain = document.createElement('main');
		this.simMain.classList.add('sim-main', 'tab-content');
		this.simContentContainer.appendChild(this.simMain);

		this.rootElem.classList.add(this.cssClass);

		if (!this.isWithinRaidSim) {
			this.rootElem.classList.add('not-within-raid-sim');
		}

		this.changeEmitter = TypedEvent.onAny([this.sim.changeEmitter], 'SimUIChange');

		this.sim.crashEmitter.on((eventID: EventID, error: SimError) => this.handleCrash(error));

		const updateShowDamageMetrics = () => {
			if (this.sim.getShowDamageMetrics())
				this.rootElem.classList.remove('hide-damage-metrics');
			else this.rootElem.classList.add('hide-damage-metrics');
		};
		updateShowDamageMetrics();
		this.sim.showDamageMetricsChangeEmitter.on(updateShowDamageMetrics);

		const updateShowThreatMetrics = () => {
			if (this.sim.getShowThreatMetrics())
				this.rootElem.classList.remove('hide-threat-metrics');
			else this.rootElem.classList.add('hide-threat-metrics');
		};
		updateShowThreatMetrics();
		this.sim.showThreatMetricsChangeEmitter.on(updateShowThreatMetrics);

		const updateShowHealingMetrics = () => {
			if (this.sim.getShowHealingMetrics())
				this.rootElem.classList.remove('hide-healing-metrics');
			else this.rootElem.classList.add('hide-healing-metrics');
		};
		updateShowHealingMetrics();
		this.sim.showHealingMetricsChangeEmitter.on(updateShowHealingMetrics);

		const updateShowEpRatios = () => {
			// Threat metrics *always* shows multiple columns, so
			// always show ratios when they are shown
			if (this.sim.getShowThreatMetrics()) {
				this.rootElem.classList.remove('hide-ep-ratios');
				// This case doesn't currently happen, but who knows
				// what the future holds...
			} else if (this.sim.getShowDamageMetrics() && this.sim.getShowHealingMetrics()) {
				this.rootElem.classList.remove('hide-ep-ratios');
			} else {
				this.rootElem.classList.add('hide-ep-ratios');
			}
		};

		updateShowEpRatios();
		this.sim.showDamageMetricsChangeEmitter.on(updateShowEpRatios);
		this.sim.showHealingMetricsChangeEmitter.on(updateShowEpRatios);
		this.sim.showThreatMetricsChangeEmitter.on(updateShowEpRatios);

		const updateShowExperimental = () => {
			if (this.sim.getShowExperimental()) this.rootElem.classList.remove('hide-experimental');
			else this.rootElem.classList.add('hide-experimental');
		};
		updateShowExperimental();
		this.sim.showExperimentalChangeEmitter.on(updateShowExperimental);

		this.addKnownIssues(config);

		// Sidebar Contents

		const titleElem = this.rootElem.querySelector('.sim-title') as HTMLElement;
		new SimTitleDropdown(titleElem, config.spec, { noDropdown: this.isWithinRaidSim });

		this.simActionsContainer = this.rootElem.querySelector(
			'.sim-sidebar-actions',
		) as HTMLElement;
		this.iterationsPicker = new NumberPicker(this.simActionsContainer, this.sim, {
			label: 'Iterations',
			extraCssClasses: ['iterations-picker', 'within-raid-sim-hide'],
			changedEvent: (sim: Sim) => sim.iterationsChangeEmitter,
			getValue: (sim: Sim) => sim.getIterations(),
			setValue: (eventID: EventID, sim: Sim, newValue: number) => {
				sim.setIterations(eventID, newValue);
			},
		}).rootElem;

		const resultsViewerElem = this.rootElem.querySelector(
			'.sim-sidebar-results',
		) as HTMLElement;
		this.resultsViewer = new ResultsViewer(resultsViewerElem);

		const socialsContainer = this.rootElem.querySelector('.sim-sidebar-socials') as HTMLElement;
		// socialsContainer.appendChild(SocialLinks.buildDiscordLink());
		// socialsContainer.appendChild(SocialLinks.buildGitHubLink());
		// socialsContainer.appendChild(SocialLinks.buildPatreonLink());

		this.simTabContentsContainer = this.rootElem.querySelector(
			'.sim-main.tab-content',
		) as HTMLElement;

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
		button.classList.add('btn', 'btn-primary', 'w-100', cssClass);
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

	addSimTab(tab: SimTab) {
		this.simHeader.addSimTabLink(tab);
	}

	addWarning(warning: SimWarning) {
		this.resultsViewer.addWarning(warning);
	}

	private addKnownIssues(config: SimUIConfig) {
		let statusStr = '';
		if (config.launchStatus == LaunchStatus.Unlaunched) {
			statusStr =
				'This sim is a WORK IN PROGRESS. It is not fully developed and should not be used for general purposes.';
		} else if (config.launchStatus == LaunchStatus.Alpha) {
			statusStr =
				'This sim is in ALPHA. Bugs are expected; please let us know if you find one!';
		} else if (config.launchStatus == LaunchStatus.Beta) {
			statusStr =
				'This sim is in BETA. There may still be a few bugs; please let us know if you find one!';
		}
		if (statusStr) {
			config.knownIssues = [statusStr].concat(config.knownIssues || []);
		}
		if (config.knownIssues && config.knownIssues.length) {
			config.knownIssues.forEach(issue => this.simHeader.addKnownIssue(issue));
		}
		globalKnownIssues.forEach(issue => this.simHeader.addKnownIssue(issue));
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

	async runSim(onProgress: (_?: any) => void) {
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

	async handleCrash(error: any): Promise<void> {
		if (!(error instanceof SimError)) {
			alert(error);
			return;
		}

		const errorStr = (error as SimError).errorStr;
		if (errorStr.startsWith('[USER_ERROR] ')) {
			let alertStr = errorStr.substring('[USER_ERROR] '.length);
			alertStr = await ActionId.replaceAllInString(alertStr);
			alert(alertStr);
			return;
		}

		if (
			window.confirm('Simulation Failure:\n' + errorStr + '\nPress Ok to file crash report')
		) {
			// Splice out just the line numbers
			const hash = this.hashCode(errorStr);
			const link = this.toLink();
			const rngSeed = this.sim.getLastUsedRngSeed();
			fetch(
				'https://api.github.com/search/issues?q=is:issue+is:open+repo:wowsims/wotlk+' +
					hash,
			)
				.then(resp => {
					resp.json().then(issues => {
						if (issues.total_count > 0) {
							window.open(issues.items[0].html_url, '_blank');
						} else {
							const base_url =
								'https://github.com/wowsims/wotlk/issues/new?assignees=&labels=&title=Crash%20Report%20';
							const base = `${base_url}${hash}&body=`;
							const maxBodyLength = URLMAXLEN - base.length;
							let issueBody = encodeURIComponent(
								`Link:\n${link}\n\nRNG Seed: ${rngSeed}\n\n${errorStr}`,
							);
							if (link.includes('/raid/')) {
								// Move the actual error before the link, as it will likely get truncated.
								issueBody = encodeURIComponent(
									`${errorStr}\nRNG Seed: ${rngSeed}\nLink:\n${link}`,
								);
							}
							let truncated = false;
							while (issueBody.length > maxBodyLength - (truncated ? 3 : 0)) {
								issueBody = issueBody.slice(0, issueBody.lastIndexOf('%')); // Avoid truncating in the middle of a URLencoded segment.
								truncated = true;
							}
							if (truncated) {
								issueBody += '...';
								// The raid links are too large and will always cause truncation.
								// Prompt the user to add more information to the issue.
								new CrashModal(this.rootElem, link);
							}
							window.open(base + issueBody, '_blank');
						}
					});
				})
				.catch(fetchErr => {
					alert('Failed to file report... try again another time:' + fetchErr);
				});
		}
	}

	hashCode(str: string): number {
		let hash = 0;
		for (let i = 0, len = str.length; i < len; i++) {
			const chr = str.charCodeAt(i);
			hash = (hash << 5) - hash + chr;
			hash |= 0; // Convert to 32bit integer
		}
		return hash;
	}

	abstract applyDefaults(eventID: EventID): void;
	abstract toLink(): string;
}

class CrashModal extends BaseModal {
	constructor(parent: HTMLElement, link: string) {
		super(parent, 'crash', { title: 'Extra Crash Information' });
		this.body.innerHTML = `
			<div class="sim-crash-report">
				<h3 class="sim-crash-report-header">Please append the following complete link to the issue you just created. This will simplify debugging the issue.</h3>
				<textarea class="sim-crash-report-text form-control"></textarea>
			</div>
		`;
		const text = document.createTextNode(link);
		this.body.querySelector('textarea')?.appendChild(text);
	}
}
