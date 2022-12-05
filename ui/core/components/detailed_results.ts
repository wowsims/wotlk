import { REPO_NAME } from '../constants/other.js'
import { RaidSimRequest, RaidSimResult } from '../proto/api.js';
import { DetailedResultsUpdate, SimRun, SimRunData } from '../proto/ui.js';
import { SimResult } from '../proto_utils/sim_result.js';
import { SimUI } from '../sim_ui.js';

import { Component } from './component.js';
import { RaidSimResultsManager } from './raid_sim_action.js';

export class DetailedResults extends Component {
	private readonly simUI: SimUI;
	private readonly iframeElem: HTMLIFrameElement;
	private tabWindow: Window | null;
	private latestRun: SimRunData | null;

	constructor(parent: HTMLElement, simUI: SimUI, simResultsManager: RaidSimResultsManager) {
		super(parent, 'detailed-results-manager-root');
		this.simUI = simUI;
		this.tabWindow = null;
		this.latestRun = null;

		this.simUI.sim.settingsChangeEmitter.on(() => this.updateSettings());

		const computedStyles = window.getComputedStyle(this.rootElem);

		const url = new URL(`${window.location.protocol}//${window.location.host}/${REPO_NAME}/detailed_results/index.html`);
		url.searchParams.append('mainTextColor', computedStyles.getPropertyValue('--main-text-color').trim());
		url.searchParams.append('themeColorPrimary', computedStyles.getPropertyValue('--theme-color-primary').trim());
		url.searchParams.append('themeColorBackground', computedStyles.getPropertyValue('--theme-color-background').trim());
		url.searchParams.append('themeColorBackgroundRaw', computedStyles.getPropertyValue('--theme-color-background-raw').trim());
		url.searchParams.append('themeBackgroundImage', computedStyles.getPropertyValue('--theme-background-image').trim());
		url.searchParams.append('themeBackgroundOpacity', computedStyles.getPropertyValue('--theme-background-opacity').trim());
		if (simUI.isIndividualSim()) {
			url.searchParams.append('isIndividualSim', '');
		}

		this.rootElem.innerHTML = `
		<div class="detailed-results-controls-div">
			<button class="detailed-results-new-tab-button sim-button">VIEW IN SEPARATE TAB</button>
		</div>
		<iframe class="detailed-results-iframe" src="${url.href}" allowtransparency="true"></iframe>
		`;

		this.iframeElem = this.rootElem.getElementsByClassName('detailed-results-iframe')[0] as HTMLIFrameElement;

		const newTabButton = this.rootElem.getElementsByClassName('detailed-results-new-tab-button')[0] as HTMLButtonElement;
		newTabButton.addEventListener('click', event => {
			if (this.tabWindow == null || this.tabWindow.closed) {
				this.tabWindow = window.open(url.href, 'Detailed Results');
				this.tabWindow!.addEventListener('load', event => {
					if (this.latestRun) {
						this.updateSettings();
						this.setSimRunData(this.latestRun);
					}
				});
			} else {
				this.tabWindow.focus();
			}
		});

		simResultsManager.currentChangeEmitter.on(() => {
			const runData = simResultsManager.getRunData();
			if (runData) {
				this.updateSettings();
				this.setSimRunData(runData);
			}
		});
	}

	private setSimRunData(simRunData: SimRunData) {
		this.latestRun = simRunData;
		this.postMessage(DetailedResultsUpdate.create({
			data: {
				oneofKind: 'runData',
				runData: simRunData,
			},
		}));
	}

	private updateSettings() {
		this.postMessage(DetailedResultsUpdate.create({
			data: {
				oneofKind: 'settings',
				settings: this.simUI.sim.toProto(),
			},
		}));
	}

	private postMessage(update: DetailedResultsUpdate) {
		this.iframeElem.contentWindow!.postMessage(DetailedResultsUpdate.toJson(update), '*');
		if (this.tabWindow) {
			this.tabWindow.postMessage(DetailedResultsUpdate.toJson(update), '*');
		}
	}
}
