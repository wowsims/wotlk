import { ResourceType } from '/tbc/core/proto/api.js';
import { OtherAction } from '/tbc/core/proto/common.js';
import { UnitMetrics, SimResult, SimResultFilter } from '/tbc/core/proto_utils/sim_result.js';
import { ActionId, resourceTypeToIcon } from '/tbc/core/proto_utils/action_id.js';
import { resourceColors, resourceNames } from '/tbc/core/proto_utils/names.js';
import { orderedResourceTypes } from '/tbc/core/proto_utils/utils.js';
import { EventID, TypedEvent } from '/tbc/core/typed_event.js';
import { bucket, distinct, getEnumValues, maxIndex, stringComparator, sum } from '/tbc/core/utils.js';

import {
	AuraUptimeLog,
	CastLog,
	DamageDealtLog,
	ResourceChangedLogGroup,
	DpsLog,
	SimLog,
	ThreatLogGroup,
} from '/tbc/core/proto_utils/logs_parser.js';

import { actionColors } from './color_settings.js';
import { ResultComponent, ResultComponentConfig, SimResultData } from './result_component.js';

declare var $: any;
declare var tippy: any;
declare var ApexCharts: any;

type TooltipHandler = (dataPointIndex: number) => string;

const dpsColor = '#ed5653';
const manaColor = '#2E93fA';
const threatColor = '#b56d07';

export class Timeline extends ResultComponent {
	private readonly dpsResourcesPlotElem: HTMLElement;
	private dpsResourcesPlot: any;

	private readonly rotationPlotElem: HTMLElement;
	private readonly rotationLabels: HTMLElement;
	private readonly rotationTimeline: HTMLElement;
	private readonly rotationHiddenIdsContainer: HTMLElement;
	private readonly chartPicker: HTMLSelectElement;

	private resultData: SimResultData | null;
	private rendered: boolean;

	private hiddenIds: Array<ActionId>;
	private hiddenIdsChangeEmitter;

	constructor(config: ResultComponentConfig) {
		config.rootCssClass = 'timeline-root';
		super(config);
		this.resultData = null;
		this.rendered = false;
		this.hiddenIds = [];
		this.hiddenIdsChangeEmitter = new TypedEvent<void>();

		this.rootElem.innerHTML = `
		<div class="timeline-disclaimer">
			<span class="timeline-warning fa fa-exclamation-triangle"></span>
			<span class="timeline-warning-description">Timeline data visualizes only 1 sim iteration.</span>
			<div class="timeline-run-again-button sim-button">SIM 1 ITERATION</div>
			<select class="timeline-chart-picker">
				<option class="rotation-option" value="rotation">Rotation</option>
				<option class="dps-option" value="dps">DPS</option>
				<option class="threat-option" value="threat">Threat</option>
			</select>
		</div>
		<div class="timeline-plots-container">
			<div class="timeline-plot dps-resources-plot hide"></div>
			<div class="timeline-plot rotation-plot">
				<div class="rotation-container">
					<div class="rotation-labels">
					</div>
					<div class="rotation-timeline">
					</div>
				</div>
				<div class="rotation-hidden-ids">
				</div>
			</div>
		</div>
		`;

		const runAgainButton = this.rootElem.getElementsByClassName('timeline-run-again-button')[0] as HTMLElement;
		runAgainButton.addEventListener('click', event => {
			(window.opener || window.parent)!.postMessage('runOnce', '*');
		});

		this.chartPicker = this.rootElem.getElementsByClassName('timeline-chart-picker')[0] as HTMLSelectElement;
		this.chartPicker.addEventListener('change', event => {
			if (this.chartPicker.value == 'rotation') {
				this.dpsResourcesPlotElem.classList.add('hide');
				this.rotationPlotElem.classList.remove('hide');
			} else {
				this.dpsResourcesPlotElem.classList.remove('hide');
				this.rotationPlotElem.classList.add('hide');
			}
			this.updatePlot();
		});

		this.dpsResourcesPlotElem = this.rootElem.getElementsByClassName('dps-resources-plot')[0] as HTMLElement;
		this.dpsResourcesPlot = new ApexCharts(this.dpsResourcesPlotElem, {
			chart: {
				type: 'line',
				foreColor: 'white',
				id: 'dpsResources',
				animations: {
					enabled: false,
				},
				height: '100%',
			},
			series: [], // Set dynamically
			xaxis: {
				title: {
					text: 'Time (s)',
				},
				type: 'datetime',
			},
			noData: {
				text: 'Waiting for data...',
			},
			stroke: {
				width: 2,
				curve: 'straight',
			},
		});

		this.rotationPlotElem = this.rootElem.getElementsByClassName('rotation-plot')[0] as HTMLElement;
		this.rotationLabels = this.rootElem.getElementsByClassName('rotation-labels')[0] as HTMLElement;
		this.rotationTimeline = this.rootElem.getElementsByClassName('rotation-timeline')[0] as HTMLElement;
		this.rotationHiddenIdsContainer = this.rootElem.getElementsByClassName('rotation-hidden-ids')[0] as HTMLElement;
	}

	onSimResult(resultData: SimResultData) {
		this.resultData = resultData;

		if (this.rendered) {
			this.updatePlot();
		}
	}

	private updatePlot() {
		if (this.resultData == null) {
			return;
		}

		const duration = this.resultData!.result.result.firstIterationDuration || 1;
		let options: any = {
			series: [],
			colors: [],
			xaxis: {
				min: this.toDatetime(0).getTime(),
				max: this.toDatetime(duration).getTime(),
				type: 'datetime',
				tickAmount: 10,
				decimalsInFloat: 1,
				labels: {
					show: true,
					formatter: (defaultValue: string, timestamp: number) => {
						return (timestamp / 1000).toFixed(1);
					},
				},
				title: {
					text: 'Time (s)',
				},
			},
			yaxis: [],
			chart: {
				events: {
					beforeResetZoom: () => {
						return {
							xaxis: {
								min: this.toDatetime(0),
								max: this.toDatetime(duration),
							},
						};
					},
				},
			},
		};

		let tooltipHandlers: Array<TooltipHandler | null> = [];
		options.tooltip = {
			enabled: true,
			custom: (data: { series: any, seriesIndex: number, dataPointIndex: number, w: any }) => {
				if (tooltipHandlers[data.seriesIndex]) {
					return tooltipHandlers[data.seriesIndex]!(data.dataPointIndex);
				} else {
					throw new Error('No tooltip handler for series ' + data.seriesIndex);
				}
			},
		};

		const players = this.resultData!.result.getPlayers(this.resultData!.filter);
		if (players.length == 1) {
			const player = players[0];

			const rotationOption = this.rootElem.getElementsByClassName('rotation-option')[0] as HTMLElement;
			rotationOption.classList.remove('hide');
			const threatOption = this.rootElem.getElementsByClassName('threat-option')[0] as HTMLElement;
			threatOption.classList.add('hide');

			this.updateRotationChart(player, duration);

			const dpsData = this.addDpsSeries(player, options, '');
			this.addDpsYAxis(dpsData.maxDps, options);
			tooltipHandlers.push(dpsData.tooltipHandler);
			tooltipHandlers.push(this.addManaSeries(player, options));
			tooltipHandlers.push(this.addThreatSeries(player, options, ''));
			tooltipHandlers = tooltipHandlers.filter(handler => handler != null);

			this.addMajorCooldownAnnotations(player, options);
		} else {
			if (this.chartPicker.value == 'rotation') {
				this.chartPicker.value = 'dps';
				return;
			}
			const rotationOption = this.rootElem.getElementsByClassName('rotation-option')[0] as HTMLElement;
			rotationOption.classList.add('hide');
			const threatOption = this.rootElem.getElementsByClassName('threat-option')[0] as HTMLElement;
			threatOption.classList.remove('hide');

			this.clearRotationChart();

			if (this.chartPicker.value == 'dps') {
				let maxDps = 0;
				players.forEach(player => {
					const dpsData = this.addDpsSeries(player, options, player.classColor);
					maxDps = Math.max(maxDps, dpsData.maxDps);
					tooltipHandlers.push(dpsData.tooltipHandler);
				});
				this.addDpsYAxis(maxDps, options);
			} else { // threat
				let maxThreat = 0;
				players.forEach(player => {
					tooltipHandlers.push(this.addThreatSeries(player, options, player.classColor));
					maxThreat = Math.max(maxThreat, player.maxThreat);
				});
				this.addThreatYAxis(maxThreat, options);
			}
		}

		this.dpsResourcesPlot.updateOptions(options);
	}

	private addDpsYAxis(maxDps: number, options: any) {
		const dpsAxisMax = Math.ceil(maxDps / 100) * 100;
		options.yaxis.push({
			color: dpsColor,
			seriesName: 'DPS',
			min: 0,
			max: dpsAxisMax,
			tickAmount: 10,
			decimalsInFloat: 0,
			title: {
				text: 'DPS',
				style: {
					color: dpsColor,
				},
			},
			axisBorder: {
				show: true,
				color: dpsColor,
			},
			axisTicks: {
				color: dpsColor,
			},
			labels: {
				minWidth: 30,
				style: {
					colors: [dpsColor],
				},
			},
		});
	}

	private addThreatYAxis(maxThreat: number, options: any) {
		const axisMax = Math.ceil(maxThreat / 10000) * 10000;
		options.yaxis.push({
			color: threatColor,
			seriesName: 'Threat',
			min: 0,
			max: axisMax,
			tickAmount: 10,
			decimalsInFloat: 0,
			title: {
				text: 'Threat',
				style: {
					color: threatColor,
				},
			},
			axisBorder: {
				show: true,
				color: threatColor,
			},
			axisTicks: {
				color: threatColor,
			},
			labels: {
				minWidth: 30,
				style: {
					colors: [threatColor],
				},
			},
		});
	}

	// Returns a function for drawing the tooltip, or null if no series was added.
	private addDpsSeries(unit: UnitMetrics, options: any, colorOverride: string): { maxDps: number, tooltipHandler: TooltipHandler } {
		const dpsLogs = unit.dpsLogs;

		options.colors.push(colorOverride || dpsColor);
		options.series.push({
			name: 'DPS',
			type: 'line',
			data: dpsLogs.map(log => {
				return {
					x: this.toDatetime(log.timestamp),
					y: log.dps,
				};
			}),
		});

		return {
			maxDps: dpsLogs[maxIndex(dpsLogs.map(l => l.dps))!].dps,
			tooltipHandler: (dataPointIndex: number) => {
				const log = dpsLogs[dataPointIndex];
				return this.dpsTooltip(log, true, unit, colorOverride);
			},
		};
	}

	// Returns a function for drawing the tooltip, or null if no series was added.
	private addManaSeries(unit: UnitMetrics, options: any): TooltipHandler | null {
		const manaLogs = unit.groupedResourceLogs[ResourceType.ResourceTypeMana];
		if (manaLogs.length == 0) {
			return null;
		}
		const maxMana = manaLogs[0].valueBefore;

		options.colors.push(manaColor);
		options.series.push({
			name: 'Mana',
			type: 'line',
			data: manaLogs.map(log => {
				return {
					x: this.toDatetime(log.timestamp),
					y: log.valueAfter,
				};
			}),
		});
		options.yaxis.push({
			seriesName: 'Mana',
			opposite: true, // Appear on right side
			min: 0,
			max: maxMana,
			tickAmount: 10,
			title: {
				text: 'Mana',
				style: {
					color: manaColor,
				},
			},
			axisBorder: {
				show: true,
				color: manaColor,
			},
			axisTicks: {
				color: manaColor,
			},
			labels: {
				minWidth: 30,
				style: {
					colors: [manaColor],
				},
				formatter: (val: string) => {
					const v = parseFloat(val);
					return `${v.toFixed(0)} (${(v / maxMana * 100).toFixed(0)}%)`;
				},
			},
		} as any);

		return (dataPointIndex: number) => {
			const log = manaLogs[dataPointIndex];
			return this.resourceTooltip(log, maxMana, true);
		};
	}

	// Returns a function for drawing the tooltip, or null if no series was added.
	private addThreatSeries(unit: UnitMetrics, options: any, colorOverride: string): TooltipHandler | null {
		options.colors.push(colorOverride || threatColor);
		options.series.push({
			name: 'Threat',
			type: 'line',
			data: unit.threatLogs.map(log => {
				return {
					x: this.toDatetime(log.timestamp),
					y: log.threatAfter,
				};
			}),
		});

		return (dataPointIndex: number) => {
			const log = unit.threatLogs[dataPointIndex];
			return this.threatTooltip(log, true, unit, colorOverride);
		};
	}

	private addMajorCooldownAnnotations(unit: UnitMetrics, options: any) {
		const mcdLogs = unit.majorCooldownLogs;
		const mcdAuraLogs = unit.majorCooldownAuraUptimeLogs;

		// Figure out how much to vertically offset cooldown icons, for cooldowns
		// used very close to each other. This is so the icons don't overlap.
		const MAX_ALLOWED_DIST = 10;
		const cooldownIconOffsets = mcdLogs.map((mcdLog, mcdIdx) => mcdLogs.filter((cdLog, cdIdx) => (cdIdx < mcdIdx) && (cdLog.timestamp > mcdLog.timestamp - MAX_ALLOWED_DIST)).length);

		const distinctMcdAuras = distinct(mcdAuraLogs, (a, b) => a.actionId!.equalsIgnoringTag(b.actionId!));
		// Sort by name so auras keep their same colors even if timings change.
		distinctMcdAuras.sort((a, b) => stringComparator(a.actionId!.name, b.actionId!.name));
		const mcdAuraColors = mcdAuraLogs.map(mcdAuraLog => actionColors[distinctMcdAuras.findIndex(dAura => dAura.actionId!.equalsIgnoringTag(mcdAuraLog.actionId!))]);

		options.annotations = {
			position: 'back',
			xaxis: mcdAuraLogs.map((log, i) => {
				return {
					x: this.toDatetime(log.gainedAt).getTime(),
					x2: this.toDatetime(log.fadedAt).getTime(),
					fillColor: mcdAuraColors[i],
				};
			}),
			points: mcdLogs.map((log, i) => {
				return {
					x: this.toDatetime(log.timestamp).getTime(),
					y: 0,
					image: {
						path: log.actionId!.iconUrl,
						width: 20,
						height: 20,
						offsetY: cooldownIconOffsets[i] * -25,
					},
				};
			}),
		};
	}

	private clearRotationChart() {
		this.rotationLabels.innerHTML = `
			<div class="rotation-label-header"></div>
		`;
		this.rotationTimeline.innerHTML = `
			<div class="rotation-timeline-header">
				<canvas class="rotation-timeline-canvas"></canvas>
			</div>
		`;
		this.rotationHiddenIdsContainer.innerHTML = '';
		this.hiddenIdsChangeEmitter = new TypedEvent<void>();
	}

	private updateRotationChart(player: UnitMetrics, duration: number) {
		const targets = this.resultData!.result.getTargets(this.resultData!.filter);
		if (targets.length == 0) {
			return;
		}
		const target = targets[0];

		this.clearRotationChart();
		this.drawRotationTimeRuler(this.rotationTimeline.getElementsByClassName('rotation-timeline-canvas')[0] as HTMLCanvasElement, duration);

		orderedResourceTypes.forEach(resourceType => this.addResourceRow(resourceType, player.groupedResourceLogs[resourceType], duration));

		const buffsById = Object.values(bucket(player.auraUptimeLogs, log => log.actionId!.toString()));
		buffsById.sort((a, b) => stringComparator(a[0].actionId!.name, b[0].actionId!.name));
		const debuffsById = Object.values(bucket(target.auraUptimeLogs, log => log.actionId!.toString()));
		debuffsById.sort((a, b) => stringComparator(a[0].actionId!.name, b[0].actionId!.name));
		const buffsAndDebuffsById = buffsById.concat(debuffsById);

		const playerCastsByAbility = this.getSortedCastsByAbility(player);
		playerCastsByAbility.forEach(castLogs => this.addCastRow(castLogs, buffsAndDebuffsById, duration));

		// Don't add a row for buffs that were already visualized in a cast row.
		const buffsToShow = buffsById.filter(auraUptimeLogs => playerCastsByAbility.findIndex(casts => casts[0].actionId!.equalsIgnoringTag(auraUptimeLogs[0].actionId!)));
		if (buffsToShow.length > 0) {
			this.addSeparatorRow(duration);
			buffsToShow.forEach(auraUptimeLogs => this.addAuraRow(auraUptimeLogs, duration));
		}

		const targetCastsByAbility = this.getSortedCastsByAbility(target);
		if (targetCastsByAbility.length > 0) {
			this.addSeparatorRow(duration);
			targetCastsByAbility.forEach(castLogs => this.addCastRow(castLogs, buffsAndDebuffsById, duration));
		}

		// Add a row for all debuffs, even those which have already been visualized in a cast row.
		const debuffsToShow = debuffsById;
		if (debuffsToShow.length > 0) {
			this.addSeparatorRow(duration);
			debuffsToShow.forEach(auraUptimeLogs => this.addAuraRow(auraUptimeLogs, duration));
		}
	}

	private getSortedCastsByAbility(player: UnitMetrics): Array<Array<CastLog>> {
		const meleeActionIds = player.getMeleeActions().map(action => action.actionId);
		const spellActionIds = player.getSpellActions().map(action => action.actionId);
		const getActionCategory = (actionId: ActionId): number => {
			const fixedCategory = idToCategoryMap[actionId.anyId()];
			if (fixedCategory != null) {
				return fixedCategory;
			} else if (meleeActionIds.find(meleeActionId => meleeActionId.equals(actionId))) {
				return MELEE_ACTION_CATEGORY;
			} else if (spellActionIds.find(spellActionId => spellActionId.equals(actionId))) {
				return SPELL_ACTION_CATEGORY;
			} else {
				return DEFAULT_ACTION_CATEGORY;
			}
		};

		const castsByAbility = Object.values(bucket(player.castLogs, log => {
			if (idsToGroupForRotation.includes(log.actionId!.spellId)) {
				return log.actionId!.toStringIgnoringTag();
			} else {
				return log.actionId!.toString();
			}
		}));

		castsByAbility.sort((a, b) => {
			const categoryA = getActionCategory(a[0].actionId!);
			const categoryB = getActionCategory(b[0].actionId!);
			if (categoryA != categoryB) {
				return categoryA - categoryB;
			} else if (a[0].actionId!.anyId() == b[0].actionId!.anyId()) {
				return a[0].actionId!.tag - b[0].actionId!.tag;
			} else {
				return stringComparator(a[0].actionId!.name, b[0].actionId!.name);
			}
		});

		return castsByAbility;
	}

	private makeLabelElem(actionId: ActionId, isHiddenLabel: boolean): HTMLElement {
		const labelElem = document.createElement('div');
		labelElem.classList.add('rotation-label', 'rotation-row');
		if (isHiddenLabel) {
			labelElem.classList.add('rotation-label-hidden');
		}
		const labelText = idsToGroupForRotation.includes(actionId.spellId) ? actionId.baseName : actionId.name;
		labelElem.innerHTML = `
			<span class="fas fa-eye${isHiddenLabel ? '' : '-slash'}"></span>
			<a class="rotation-label-icon"></a>
			<span class="rotation-label-text">${labelText}</span>
		`;
		const hideElem = labelElem.getElementsByClassName('fas')[0] as HTMLElement;
		hideElem.addEventListener('click', event => {
			if (isHiddenLabel) {
				const index = this.hiddenIds.findIndex(hiddenId => hiddenId.equals(actionId));
				if (index != -1) {
					this.hiddenIds.splice(index, 1);
				}
			} else {
				this.hiddenIds.push(actionId);
			}
			this.hiddenIdsChangeEmitter.emit(TypedEvent.nextEventID());
		});
		tippy(hideElem, {
			content: isHiddenLabel ? 'Show Row' : 'Hide Row',
			allowHTML: true,
		});
		const updateHidden = () => {
			if (isHiddenLabel == Boolean(this.hiddenIds.find(hiddenId => hiddenId.equals(actionId)))) {
				labelElem.classList.remove('hide');
			} else {
				labelElem.classList.add('hide');
			}
		};
		this.hiddenIdsChangeEmitter.on(updateHidden);
		updateHidden();
		const labelIcon = labelElem.getElementsByClassName('rotation-label-icon')[0] as HTMLAnchorElement;
		actionId.setBackgroundAndHref(labelIcon);
		return labelElem;
	}

	private makeRowElem(actionId: ActionId, duration: number): HTMLElement {
		const rowElem = document.createElement('div');
		rowElem.classList.add('rotation-timeline-row', 'rotation-row');
		rowElem.style.width = this.timeToPx(duration);

		const updateHidden = () => {
			if (this.hiddenIds.find(hiddenId => hiddenId.equals(actionId))) {
				rowElem.classList.add('hide');
			} else {
				rowElem.classList.remove('hide');
			}
		};
		this.hiddenIdsChangeEmitter.on(updateHidden);
		updateHidden();
		return rowElem;
	}

	private addSeparatorRow(duration: number) {
		let separatorElem = document.createElement('div');
		separatorElem.classList.add('rotation-timeline-separator');
		this.rotationLabels.appendChild(separatorElem);
		separatorElem = document.createElement('div');
		separatorElem.classList.add('rotation-timeline-separator');
		separatorElem.style.width = this.timeToPx(duration);
		this.rotationTimeline.appendChild(separatorElem);
	}

	private addResourceRow(resourceType: ResourceType, resourceLogs: Array<ResourceChangedLogGroup>, duration: number) {
		if (resourceLogs.length == 0) {
			return;
		}
		const startValue = resourceLogs[0].valueBefore;

		const labelElem = document.createElement('div');
		labelElem.classList.add('rotation-label', 'rotation-row');
		labelElem.innerHTML = `
			<a class="rotation-label-icon" style="background-image:url('${resourceTypeToIcon[resourceType]}')"></a>
			<span class="rotation-label-text">${resourceNames[resourceType]}</span>
		`;
		this.rotationLabels.appendChild(labelElem);

		const rowElem = document.createElement('div');
		rowElem.classList.add('rotation-timeline-row', 'rotation-row');
		rowElem.style.width = this.timeToPx(duration);

		resourceLogs.forEach((resourceLogGroup, i) => {
			const resourceElem = document.createElement('div');
			resourceElem.classList.add('rotation-timeline-resource', 'series-color', resourceNames[resourceType].toLowerCase().replaceAll(' ', '-'));
			resourceElem.style.left = this.timeToPx(resourceLogGroup.timestamp);
			resourceElem.style.width = this.timeToPx((resourceLogs[i + 1]?.timestamp || duration) - resourceLogGroup.timestamp);
			if (percentageResources.includes(resourceType)) {
				resourceElem.textContent = (resourceLogGroup.valueAfter / startValue * 100).toFixed(0) + '%';
			} else {
				resourceElem.textContent = Math.floor(resourceLogGroup.valueAfter).toFixed(0);
			}
			rowElem.appendChild(resourceElem);

			tippy(resourceElem, {
				content: this.resourceTooltip(resourceLogGroup, startValue, false),
				allowHTML: true,
				placement: 'bottom',
			});
		});
		this.rotationTimeline.appendChild(rowElem);
	}

	private addCastRow(castLogs: Array<CastLog>, aurasById: Array<Array<AuraUptimeLog>>, duration: number) {
		const actionId = castLogs[0].actionId!;

		this.rotationLabels.appendChild(this.makeLabelElem(actionId, false));
		this.rotationHiddenIdsContainer.appendChild(this.makeLabelElem(actionId, true));

		const rowElem = this.makeRowElem(actionId, duration);
		castLogs.forEach(castLog => {
			const castElem = document.createElement('div');
			castElem.classList.add('rotation-timeline-cast');
			castElem.style.left = this.timeToPx(castLog.timestamp);
			castElem.style.minWidth = this.timeToPx(castLog.castTime);
			rowElem.appendChild(castElem);

			if (castLog.damageDealtLogs.length > 0) {
				const ddl = castLog.damageDealtLogs[0];
				if (ddl.miss || ddl.dodge || ddl.parry) {
					castElem.classList.add('outcome-miss');
				} else if (ddl.glance || ddl.block || ddl.partialResist1_4 || ddl.partialResist2_4 || ddl.partialResist3_4) {
					castElem.classList.add('outcome-partial');
				} else if (ddl.crit) {
					castElem.classList.add('outcome-crit');
				} else {
					castElem.classList.add('outcome-hit');
				}
			}

			const iconElem = document.createElement('a');
			iconElem.classList.add('rotation-timeline-cast-icon');
			actionId.setBackground(iconElem);
			castElem.appendChild(iconElem);
			tippy(castElem, {
				content: `
					<span>${castLog.actionId!.name} from ${castLog.timestamp.toFixed(2)}s to ${(castLog.timestamp + castLog.castTime).toFixed(2)}s (${castLog.castTime.toFixed(2)}s)</span>
					<ul class="rotation-timeline-cast-damage-list">
						${castLog.damageDealtLogs.map(ddl => `
								<li>
									<span>${ddl.timestamp.toFixed(2)}s - ${ddl.resultString()}</span>
									${ddl.source?.isTarget ? '' : `<span class="threat-metrics"> (${ddl.threat.toFixed(1)} Threat)</span>`}
								</li>`)
						.join('')}
					</ul>
				`,
				allowHTML: true,
				placement: 'bottom',
			});

			castLog.damageDealtLogs.filter(ddl => ddl.tick).forEach(ddl => {
				const tickElem = document.createElement('div');
				tickElem.classList.add('rotation-timeline-tick');
				tickElem.style.left = this.timeToPx(ddl.timestamp);
				rowElem.appendChild(tickElem);

				tippy(tickElem, {
					content: `
						<span>${ddl.timestamp.toFixed(2)}s - ${ddl.actionId!.name} ${ddl.resultString()}</span>
						${ddl.source?.isTarget ? '' : `<span class="threat-metrics"> (${ddl.threat.toFixed(1)} Threat)</span>`}
					`,
					allowHTML: true,
					placement: 'bottom',
				});
			});
		});

		// If there are any auras that correspond to this cast, visualize them in the same row.
		aurasById
			.filter(auraUptimeLogs => auraUptimeLogs[0].actionId!.equalsIgnoringTag(actionId))
			.forEach(auraUptimeLogs => this.applyAuraUptimeLogsToRow(auraUptimeLogs, rowElem, duration));

		this.rotationTimeline.appendChild(rowElem);
	}

	private addAuraRow(auraUptimeLogs: Array<AuraUptimeLog>, duration: number) {
		const actionId = auraUptimeLogs[0].actionId!;

		let rowElem = this.makeRowElem(actionId, duration);
		this.rotationLabels.appendChild(this.makeLabelElem(actionId, false));
		this.rotationHiddenIdsContainer.appendChild(this.makeLabelElem(actionId, true));
		this.rotationTimeline.appendChild(rowElem);

		this.applyAuraUptimeLogsToRow(auraUptimeLogs, rowElem, duration);
	}

	private applyAuraUptimeLogsToRow(auraUptimeLogs: Array<AuraUptimeLog>, rowElem: HTMLElement, duration: number) {
		auraUptimeLogs.forEach(aul => {
			const auraElem = document.createElement('div');
			auraElem.classList.add('rotation-timeline-aura');
			auraElem.style.left = this.timeToPx(aul.gainedAt);
			auraElem.style.minWidth = this.timeToPx((aul.fadedAt || duration) - aul.gainedAt);
			rowElem.appendChild(auraElem);

			tippy(auraElem, {
				content: `
					<span>${aul.actionId!.name}: ${aul.gainedAt.toFixed(2)}s - ${(aul.fadedAt || duration).toFixed(2)}s</span>
				`,
				allowHTML: true,
			});
		});
	}

	private timeToPxValue(time: number): number {
		return time * 100;
	}
	private timeToPx(time: number): string {
		return this.timeToPxValue(time) + 'px';
	}

	private drawRotationTimeRuler(canvas: HTMLCanvasElement, duration: number) {
		const height = 30;
		canvas.width = this.timeToPxValue(duration);
		canvas.height = height;

		const ctx = canvas.getContext('2d')!;
		ctx.strokeStyle = 'white'

		ctx.font = 'bold 14px SimDefaultFont';
		ctx.fillStyle = 'white';
		ctx.lineWidth = 2;
		ctx.beginPath();

		// Bottom border line
		ctx.moveTo(0, height);
		ctx.lineTo(canvas.width, height);

		// Tick lines
		const numTicks = 1 + Math.floor(duration * 10);
		for (let i = 0; i <= numTicks; i++) {
			const time = i * 0.1;
			let x = this.timeToPxValue(time);
			if (i == 0) {
				ctx.textAlign = 'left';
				x++;
			} else if (i % 10 == 0 && time + 1 > duration) {
				ctx.textAlign = 'right';
				x--;
			} else {
				ctx.textAlign = 'center';
			}

			let lineHeight = 0;
			if (i % 10 == 0) {
				lineHeight = height * 0.5;
				ctx.fillText(time + 's', x, height - height * 0.6);
			} else if (i % 5 == 0) {
				lineHeight = height * 0.25;
			} else {
				lineHeight = height * 0.125;
			}
			ctx.moveTo(x, height);
			ctx.lineTo(x, height - lineHeight);
		}
		ctx.stroke();
	}

	private dpsTooltip(log: DpsLog, includeAuras: boolean, player: UnitMetrics, colorOverride: string): string {
		const showPlayerLabel = colorOverride != '';
		return `<div class="timeline-tooltip dps">
			<div class="timeline-tooltip-header">
				${showPlayerLabel ? `
				<img class="timeline-tooltip-icon" src="${player.iconUrl}">
				<span class="" style="color: ${colorOverride}">${player.label}</span><span> - </span>
				` : ''}
				<span class="bold">${log.timestamp.toFixed(2)}s</span>
			</div>
			<div class="timeline-tooltip-body">
				<ul class="timeline-dps-events">
					${log.damageLogs.map(damageLog => this.tooltipLogItem(damageLog, damageLog.resultString())).join('')}
				</ul>
				<div class="timeline-tooltip-body-row">
					<span class="series-color">DPS: ${log.dps.toFixed(2)}</span>
				</div>
			</div>
			${this.tooltipAurasSection(log)}
		</div>`;
	}

	private threatTooltip(log: ThreatLogGroup, includeAuras: boolean, player: UnitMetrics, colorOverride: string): string {
		const showPlayerLabel = colorOverride != '';
		return `<div class="timeline-tooltip threat">
			<div class="timeline-tooltip-header">
				${showPlayerLabel ? `
				<img class="timeline-tooltip-icon" src="${player.iconUrl}">
				<span class="" style="color: ${colorOverride}">${player.label}</span></span> - </span>
				` : ''}
				<span class="bold">${log.timestamp.toFixed(2)}s</span>
			</div>
			<div class="timeline-tooltip-body">
				<div class="timeline-tooltip-body-row">
					<span class="series-color">Before: ${log.threatBefore.toFixed(1)}</span>
				</div>
				<ul class="timeline-threat-events">
					${log.logs.map(log => this.tooltipLogItem(log, `${log.threat.toFixed(1)} Threat`)).join('')}
				</ul>
				<div class="timeline-tooltip-body-row">
					<span class="series-color">After: ${log.threatAfter.toFixed(1)}</span>
				</div>
			</div>
			${includeAuras ? this.tooltipAurasSection(log) : ''}
		</div>`;
	}

	private resourceTooltip(log: ResourceChangedLogGroup, maxValue: number, includeAuras: boolean): string {
		const valToDisplayString = percentageResources.includes(log.resourceType)
			? (val: number) => `${val.toFixed(1)} (${(val / maxValue * 100).toFixed(0)}%)`
			: (val: number) => `${val.toFixed(1)}`;

		return `<div class="timeline-tooltip ${resourceNames[log.resourceType].toLowerCase().replaceAll(' ', '-')}">
			<div class="timeline-tooltip-header">
				<span class="bold">${log.timestamp.toFixed(2)}s</span>
			</div>
			<div class="timeline-tooltip-body">
				<div class="timeline-tooltip-body-row">
					<span class="series-color">Before: ${valToDisplayString(log.valueBefore)}</span>
				</div>
				<ul class="timeline-mana-events">
					${log.logs.map(manaChangedLog => this.tooltipLogItem(manaChangedLog, manaChangedLog.resultString())).join('')}
				</ul>
				<div class="timeline-tooltip-body-row">
					<span class="series-color">After: ${valToDisplayString(log.valueAfter)}</span>
				</div>
			</div>
			${includeAuras ? this.tooltipAurasSection(log) : ''}
		</div>`;
	}

	private tooltipLogItem(log: SimLog, value: string): string {
		const valueElem = `<span class="series-color">${value}</span>`;
		let actionElem = '';
		if (log.actionId) {
			let iconElem = '';
			if (log.actionId.iconUrl) {
				iconElem = `<img class="timeline-tooltip-icon" src="${log.actionId.iconUrl}">`;
			}
			actionElem = `
			${iconElem}
			<span>${log.actionId.name}:</span>
			`;
		}
		return `
		<li>
			${actionElem}
			${valueElem}
		</li>`;
	}

	private tooltipAurasSection(log: SimLog): string {
		if (log.activeAuras.length == 0) {
			return '';
		}

		return `
		<div class="timeline-tooltip-auras">
			<div class="timeline-tooltip-body-row">
				<span class="bold">Active Auras</span>
			</div>
			<ul class="timeline-active-auras">
				${log.activeAuras.map(auraLog => {
			let iconElem = '';
			if (auraLog.actionId!.iconUrl) {
				iconElem = `<img class="timeline-tooltip-icon" src="${auraLog.actionId!.iconUrl}">`;
			}
			return `
					<li>
						${iconElem}
						<span>${auraLog.actionId!.name}</span>
					</li>`;
		}).join('')}
			</ul>
		</div>
		`;
	}

	render() {
		setTimeout(() => {
			this.dpsResourcesPlot.render();
			this.rendered = true;
			this.updatePlot();
		}, 300);
	}

	private toDatetime(timestamp: number): Date {
		return new Date(timestamp * 1000);
	}
}

const MELEE_ACTION_CATEGORY = 1;
const SPELL_ACTION_CATEGORY = 2;
const DEFAULT_ACTION_CATEGORY = 3;

// Hard-coded spell categories for controlling rotation ordering.
const idToCategoryMap: Record<number, number> = {
	[OtherAction.OtherActionAttack]: 0,
	[OtherAction.OtherActionShoot]: 0.5,

	// Druid
	[26996]: 0.1, // Maul
	[33987]: MELEE_ACTION_CATEGORY + 0.1, // Mangle (Bear)
	[33745]: MELEE_ACTION_CATEGORY + 0.2, // Lacerate
	[26997]: MELEE_ACTION_CATEGORY + 0.3, // Swipe

	[33983]: MELEE_ACTION_CATEGORY + 0.1, // Mangle (Cat)
	[27002]: MELEE_ACTION_CATEGORY + 0.2, // Shred
	[27008]: MELEE_ACTION_CATEGORY + 0.3, // Rip
	[24248]: MELEE_ACTION_CATEGORY + 0.4, // Ferocious Bite

	// Hunter
	[27014]: 0.1, // Raptor Strike

	// Paladin
	[27156]: 0.1, // Seal of Righteousness Proc
	[20182]: 0.2, // Reckoning
	[27179]: SPELL_ACTION_CATEGORY + 0.1, // Holy Shield

	// Rogue
	[6774]: MELEE_ACTION_CATEGORY + 0.1, // Slice and Dice
	[26866]: MELEE_ACTION_CATEGORY + 0.2, // Expose Armor
	[26865]: MELEE_ACTION_CATEGORY + 0.3, // Eviscerate
	[26867]: MELEE_ACTION_CATEGORY + 0.3, // Rupture

	// Shaman
	[17364]: MELEE_ACTION_CATEGORY + 0.1, // Stormstrike
	[25454]: MELEE_ACTION_CATEGORY + 0.2, // Earth Shock
	[25457]: MELEE_ACTION_CATEGORY + 0.2, // Flame Shock
	[25464]: MELEE_ACTION_CATEGORY + 0.2, // Frost Shock
	[25533]: SPELL_ACTION_CATEGORY + 0.2, // Searing Totem
	[25552]: SPELL_ACTION_CATEGORY + 0.2, // Magma Totem
	[25537]: SPELL_ACTION_CATEGORY + 0.2, // Fire Nova Totem
	[25359]: SPELL_ACTION_CATEGORY + 0.3, // Grace of Air Totem
	[8512]: SPELL_ACTION_CATEGORY + 0.3, // Windfury Totem r1
	[10613]: SPELL_ACTION_CATEGORY + 0.3, // Windfury Totem r2
	[10614]: SPELL_ACTION_CATEGORY + 0.3, // Windfury Totem r3
	[25585]: SPELL_ACTION_CATEGORY + 0.3, // Windfury Totem r4
	[25587]: SPELL_ACTION_CATEGORY + 0.3, // Windfury Totem r5
	[2825]: DEFAULT_ACTION_CATEGORY + 0.1, // Bloodlust

	// Warrior
	[25231]: 0.1, // Cleave
	[29707]: 0.1, // Heroic Strike
	[25242]: MELEE_ACTION_CATEGORY + 0.05, // Slam
	[30335]: MELEE_ACTION_CATEGORY + 0.1, // Bloodthirst
	[30330]: MELEE_ACTION_CATEGORY + 0.1, // Mortal Strike
	[30356]: MELEE_ACTION_CATEGORY + 0.1, // Shield Slam
	[1680]: MELEE_ACTION_CATEGORY + 0.2, // Whirlwind
	[11585]: MELEE_ACTION_CATEGORY + 0.3, // Overpower
	[25212]: MELEE_ACTION_CATEGORY + 0.4, // Hamstring
	[25236]: MELEE_ACTION_CATEGORY + 0.5, // Execute
	[71]: DEFAULT_ACTION_CATEGORY + 0.1, // Defensive Stance
	[2457]: DEFAULT_ACTION_CATEGORY + 0.1, // Battle Stance
	[2458]: DEFAULT_ACTION_CATEGORY + 0.1, // Berserker Stance

	// Generic
	[26992]: SPELL_ACTION_CATEGORY + 0.91, // Thorns
	[27150]: SPELL_ACTION_CATEGORY + 0.92, // Retribution Aura
	[27169]: SPELL_ACTION_CATEGORY + 0.93, // Blessing of Sanctuary
};

const idsToGroupForRotation: Array<number> = [
	6774,  // Slice and Dice
	26866, // Expose Armor
	26865, // Eviscerate
	26867, // Rupture
];

const percentageResources: Array<ResourceType> = [
	ResourceType.ResourceTypeHealth,
	ResourceType.ResourceTypeMana,
];
