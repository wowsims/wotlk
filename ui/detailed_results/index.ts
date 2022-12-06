import { DetailedResultsUpdate, SimRun, SimRunData } from '../core/proto/ui.js';
import { EventID, TypedEvent } from '../core/typed_event.js';
import { Database } from '../core/proto_utils/database.js';
import { SimResult, SimResultFilter } from '../core/proto_utils/sim_result.js';

import { SimResultData } from './result_component.js';
import { ResultsFilter } from './results_filter.js';
import { CastMetricsTable } from './cast_metrics.js';
import { DtpsMeleeMetricsTable } from './dtps_melee_metrics.js';
import { HealingMetricsTable } from './healing_metrics.js';
import { MeleeMetricsTable } from './melee_metrics.js';
import { SpellMetricsTable } from './spell_metrics.js';
import { ResourceMetricsTable } from './resource_metrics.js';
import { PlayerDamageMetricsTable } from './player_damage.js';
import { AuraMetricsTable } from './aura_metrics.js'
import { DpsHistogram } from './dps_histogram.js';
import { SourceChart } from './source_chart.js';
import { Timeline } from './timeline.js';
import { ToplineResults } from './topline_results.js';

declare var Chart: any;

Database.get();

const urlParams = new URLSearchParams(window.location.search);
if (urlParams.has('mainTextColor')) {
	document.body.style.setProperty('--main-text-color', urlParams.get('mainTextColor')!);
}
if (urlParams.has('themeColorPrimary')) {
	document.body.style.setProperty('--theme-color-primary', urlParams.get('themeColorPrimary')!);
}
if (urlParams.has('themeColorBackground')) {
	document.body.style.setProperty('--theme-color-background', urlParams.get('themeColorBackground')!);
}
if (urlParams.has('themeColorBackgroundRaw')) {
	document.body.style.setProperty('--theme-color-background-raw', urlParams.get('themeColorBackgroundRaw')!);
}
if (urlParams.has('themeBackgroundImage')) {
	document.body.style.setProperty('--theme-background-image', urlParams.get('themeBackgroundImage')!);
}
if (urlParams.has('themeBackgroundOpacity')) {
	document.body.style.setProperty('--theme-background-opacity', urlParams.get('themeBackgroundOpacity')!);
}

const isIndividualSim = urlParams.has('isIndividualSim');
if (isIndividualSim) {
	document.body.classList.add('individual-sim');
}

const isInIframe = Boolean(window.frameElement);
if (isInIframe) {
	// Causes links opened from the iframe to be opened as tabs in the parent window instead.
	const base = document.createElement('base');
	base.target = '_parent';
	document.head.appendChild(base);
} else {
	document.body.classList.add('new-tab');
}

const colorSettings = {
	mainTextColor: document.body.style.getPropertyValue('--main-text-color'),
};

Chart.defaults.color = colorSettings.mainTextColor;

const layoutHTML = `
<div class="dr-root">
	<div class="dr-toolbar">
		<div class="results-filter"></div>
		<div class="tabs-filler"></div>
		<ul class="dr-toolbar nav nav-tabs" role="tablist">
			<li class="nav-item dr-tab-tab damage-metrics" role="presentation">
				<a
					class="nav-link active"
					data-bs-toggle="tab"
					data-bs-target="#damageTab"
					type="button"
					role="tab"
					aria-controls="damageTab"
					aria-selected="true"
				>DAMAGE</a>
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
				>HEALING</a>
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
				>DAMAGE TAKEN</a>
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
				>BUFFS</a>
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
				>DEBUFFS</a>
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
				>CASTS</a>
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
				>RESOURCES</a>
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
				>TIMELINE</a>
			</li>
		</ul>
	</div>
	<div class="tab-content">
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
	</div>
</div>
`;

document.body.innerHTML = layoutHTML;
const resultsEmitter = new TypedEvent<SimResultData | null>();

const resultsFilter = new ResultsFilter({
	parent: document.body.getElementsByClassName('results-filter')[0] as HTMLElement,
	resultsEmitter: resultsEmitter,
	colorSettings: colorSettings,
});

(Array.from(document.body.getElementsByClassName('topline-results')) as Array<HTMLElement>).forEach(toplineResultsDiv => {
	new ToplineResults({ parent: toplineResultsDiv, resultsEmitter: resultsEmitter, colorSettings: colorSettings });
});

const castMetrics = new CastMetricsTable({ parent: document.body.getElementsByClassName('cast-metrics')[0] as HTMLElement, resultsEmitter: resultsEmitter, colorSettings: colorSettings });
const meleeMetrics = new MeleeMetricsTable({ parent: document.body.getElementsByClassName('melee-metrics')[0] as HTMLElement, resultsEmitter: resultsEmitter, colorSettings: colorSettings });
const spellMetrics = new SpellMetricsTable({ parent: document.body.getElementsByClassName('spell-metrics')[0] as HTMLElement, resultsEmitter: resultsEmitter, colorSettings: colorSettings });
const healingMetrics = new HealingMetricsTable({ parent: document.body.getElementsByClassName('healing-spell-metrics')[0] as HTMLElement, resultsEmitter: resultsEmitter, colorSettings: colorSettings });
const resourceMetrics = new ResourceMetricsTable({ parent: document.body.getElementsByClassName('resource-metrics')[0] as HTMLElement, resultsEmitter: resultsEmitter, colorSettings: colorSettings });
const playerDamageMetrics = new PlayerDamageMetricsTable({ parent: document.body.getElementsByClassName('player-damage-metrics')[0] as HTMLElement, resultsEmitter: resultsEmitter, colorSettings: colorSettings }, resultsFilter);
const buffAuraMetrics = new AuraMetricsTable({
	parent: document.body.getElementsByClassName('buff-aura-metrics')[0] as HTMLElement,
	resultsEmitter: resultsEmitter,
	colorSettings: colorSettings,
}, false);
const debuffAuraMetrics = new AuraMetricsTable({
	parent: document.body.getElementsByClassName('debuff-aura-metrics')[0] as HTMLElement,
	resultsEmitter: resultsEmitter,
	colorSettings: colorSettings,
}, true);
const dpsHistogram = new DpsHistogram({ parent: document.body.getElementsByClassName('dps-histogram')[0] as HTMLElement, resultsEmitter: resultsEmitter, colorSettings: colorSettings });

const dtpsMeleeMetrics = new DtpsMeleeMetricsTable({ parent: document.body.getElementsByClassName('dtps-melee-metrics')[0] as HTMLElement, resultsEmitter: resultsEmitter, colorSettings: colorSettings });

const timeline = new Timeline({
	parent: document.body.getElementsByClassName('timeline')[0] as HTMLElement,
	resultsEmitter: resultsEmitter,
	colorSettings: colorSettings,
});
(document.getElementById('timelineTabTab') as HTMLElement).addEventListener('click', event => timeline.render());

let currentSimResult: SimResult | null = null;
function updateResults() {
	const eventID = TypedEvent.nextEventID();
	if (currentSimResult == null) {
		resultsEmitter.emit(eventID, null);
	} else {
		resultsEmitter.emit(eventID, {
			eventID: eventID,
			result: currentSimResult,
			filter: resultsFilter.getFilter(),
		});
	}
}

document.body.classList.add('hide-threat-metrics');
document.body.classList.add('hide-healing-metrics');
window.addEventListener('message', async event => {
	const data = DetailedResultsUpdate.fromJson(event.data);
	switch (data.data.oneofKind) {
		case 'runData':
			const runData = data.data.runData;
			currentSimResult = await SimResult.fromProto(runData.run || SimRun.create());
			updateResults();
			break;
		case 'settings':
			const settings = data.data.settings;
			if (settings.showDamageMetrics) {
				document.body.classList.remove('hide-damage-metrics');
			} else {
				document.body.classList.add('hide-damage-metrics');
				if (document.getElementById('damageTab')!.classList.contains('active')) {
					document.getElementById('damageTab')!.classList.remove('active', 'show');
					document.getElementById('healingTab')!.classList.add('active', 'show');

					const toolbar = document.getElementsByClassName('dr-toolbar')[0] as HTMLElement;
					(toolbar.getElementsByClassName('damage-metrics')[0] as HTMLElement).children[0]!.classList.remove('active');
					(toolbar.getElementsByClassName('healing-metrics')[0] as HTMLElement).children[0]!.classList.add('active');
				}
			}
			if (settings.showThreatMetrics) {
				document.body.classList.remove('hide-threat-metrics');
			} else {
				document.body.classList.add('hide-threat-metrics');
			}
			if (settings.showHealingMetrics) {
				document.body.classList.remove('hide-healing-metrics');
			} else {
				document.body.classList.add('hide-healing-metrics');
			}
			if (settings.showExperimental) {
				document.body.classList.remove('hide-experimental');
			} else {
				document.body.classList.add('hide-experimental');
			}
			break;
	}
});

resultsFilter.changeEmitter.on(() => updateResults());

const rootDiv = document.body.getElementsByClassName('dr-root')[0] as HTMLElement;
resultsEmitter.on((eventID, resultData) => {
	if (resultData?.filter.player || resultData?.filter.player === 0) {
		rootDiv.classList.remove('all-players');
		rootDiv.classList.add('single-player');
	} else {
		rootDiv.classList.add('all-players');
		rootDiv.classList.remove('single-player');
	}
});
