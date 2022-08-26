import { Component } from './components/component.js';
import { NumberPicker } from './components/number_picker.js';
import { ResultsViewer } from './components/results_viewer.js';
import { Title } from './components/title.js';
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

        const updateShowThreatMetrics = () => {
            if (this.sim.getShowThreatMetrics()) {
                this.rootElem.classList.remove('hide-threat-metrics');
            } else {
                this.rootElem.classList.add('hide-threat-metrics');
            }
        };
        updateShowThreatMetrics();
        this.sim.showThreatMetricsChangeEmitter.on(updateShowThreatMetrics);

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

        const resultsViewerElem = this.rootElem.getElementsByClassName('sim-sidebar-results')[0] as HTMLElement;
        this.resultsViewer = new ResultsViewer(resultsViewerElem);

        const titleElem = this.rootElem.getElementsByClassName('sim-sidebar-title')[0] as HTMLElement;
        const title = new Title(titleElem, config.spec);

        const simActionsContainer = this.rootElem.getElementsByClassName('sim-sidebar-actions')[0] as HTMLElement;
        const iterationsPicker = new NumberPicker(simActionsContainer, this.sim, {
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
        const simActionsContainer = this.rootElem.getElementsByClassName('sim-sidebar-actions')[0] as HTMLElement;
        const iterationsPicker = this.rootElem.getElementsByClassName('iterations-picker')[0] as HTMLElement;

        const button = document.createElement('button');
        button.classList.add('sim-sidebar-actions-button', cssClass);
        button.textContent = name;
        button.addEventListener('click', actFn);
        simActionsContainer.insertBefore(button, iterationsPicker);
    }

    addTab(title: string, cssClass: string, innerHTML: string) {
        const simTabsContainer = this.rootElem.getElementsByClassName('sim-tabs')[0] as HTMLElement;
        const simTabContentsContainer = this.rootElem.getElementsByClassName('tab-content')[0] as HTMLElement;
        const topBar = simTabsContainer.getElementsByClassName('sim-top-bar')[0] as HTMLElement;

        const contentId = cssClass.replace(/\s+/g, '-') + '-tab';
        const isFirstTab = simTabsContainer.children.length == 1;

        const newTab = document.createElement('li');
        newTab.innerHTML = `<a data-toggle="tab" href="#${contentId}">${title}</a>`;
        newTab.classList.add(cssClass + '-tab');

        const newContent = document.createElement('div');
        newContent.id = contentId;
        newContent.classList.add(cssClass, 'tab-pane', 'fade');
        newContent.innerHTML = innerHTML;

        if (isFirstTab) {
            newTab.classList.add('active');
            newContent.classList.add('active', 'in');
        }

        simTabsContainer.insertBefore(newTab, topBar);
        simTabContentsContainer.appendChild(newContent);
    }

    addToolbarItem(elem: HTMLElement) {
        const topBar = this.rootElem.getElementsByClassName('sim-top-bar')[0] as HTMLElement;
        elem.classList.add('sim-top-bar-item');
        topBar.appendChild(elem);
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
  <section class="sim-sidebar">
    <div class="sim-sidebar-title"></div>
    <div class="sim-sidebar-actions within-raid-sim-hide"></div>
    <div class="sim-sidebar-results within-raid-sim-hide"></div>
    <div class="sim-sidebar-footer"></div>
  </section>
  <section class="sim-main">
		<div class="sim-toolbar">
			<ul class="sim-tabs nav nav-tabs">
				<li class="sim-top-bar">
					<span class="notices fas fa-exclamation-circle"></span>
					<span class="warnings fa fa-exclamation-triangle"></span>
					<div class="known-issues">Known Issues</div>
				</li>
			</ul>
    </div>
    <div class="tab-content">
    </div>
  </section>
</div>
`;
