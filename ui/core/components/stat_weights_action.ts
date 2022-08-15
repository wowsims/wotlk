import { StatWeightsRequest, StatWeightsResult, StatWeightValues, ProgressMetrics } from '../proto/api.js';
import { ItemSlot } from '../proto/common.js';
import { Gem } from '../proto/common.js';
import { GemColor } from '../proto/common.js';
import { Profession } from '../proto/common.js';
import { Stat } from '../proto/common.js';
import { Stats } from '../proto_utils/stats.js';
import { Gear } from '../proto_utils/gear.js';
import { statNames, statOrder } from '../proto_utils/names.js';
import { IndividualSimUI } from '../individual_sim_ui.js';
import { EventID, TypedEvent } from '../typed_event.js';
import { Player } from '../player.js';
import { stDevToConf90 } from '../utils.js';
import { BooleanPicker } from '../components/boolean_picker.js';
import { NumberPicker } from '../components/number_picker.js';
import { ResultsViewer } from '../components/results_viewer.js';
import { combinations, combinationsWithDups, permutations, getEnumValues, maxIndex, sum } from '../utils.js';

import * as Gems from '../proto_utils/gems.js';

import { Popup } from './popup.js';

declare var tippy: any;

export function addStatWeightsAction(simUI: IndividualSimUI<any>, epStats: Array<Stat>, epReferenceStat: Stat) {
	simUI.addAction('STAT WEIGHTS', 'ep-weights-action', () => {
		new EpWeightsMenu(simUI, epStats, epReferenceStat);
	});
}

class EpWeightsMenu extends Popup {
	private readonly simUI: IndividualSimUI<any>;
	private readonly tableContainer: HTMLElement;
	private readonly tableBody: HTMLElement;
	private readonly tableHeader: HTMLElement;
	private readonly resultsViewer: ResultsViewer;

	private statsType: string;
	private epStats: Array<Stat>;
	private epReferenceStat: Stat;

	constructor(simUI: IndividualSimUI<any>, epStats: Array<Stat>, epReferenceStat: Stat) {
		super(simUI.rootElem);
		this.simUI = simUI;
		this.statsType = 'ep';
		this.epStats = epStats;
		this.epReferenceStat = epReferenceStat;

		this.rootElem.classList.add('ep-weights-menu');
		this.rootElem.innerHTML = `
			<div class="ep-weights-header">
				<div class="ep-weights-actions">
					<button class="sim-button calc-weights">CALCULATE</button>
				</div>
				<div class="ep-weights-results">
				</div>
			</div>
			<div class="stats-controls-row">
				<div class="ep-weights-options">
					<select class="ep-type-select">
						<option value="ep">EP</option>
						<option value="weight">Weights</option>
					</select>
				</div>
				<div class="show-all-stats-container">
				</div>
				<button class="sim-button optimize-gems">OPTIMIZE GEMS</button>
			</div>
			<div class="ep-weights-table">
				<table class="results-ep-table">
					<tbody id="ep-tbody">
						<tr>
							<th>Stat</th>
							<th class="type-weight"><span>DPS Weight</span><span class="col-action fa fa-copy"></span></th>
							<th class="type-ep"><span>DPS EP</span><span class="col-action fa fa-copy"></span></th>
							<th class="threat-metrics type-weight"><span>TPS Weight</span><span class="col-action fa fa-copy"></span></th>
							<th class="threat-metrics type-ep"><span>TPS EP</span><span class="col-action fa fa-copy"></span></th>
							<th class="threat-metrics type-weight"><span>DTPS Weight</span><span class="col-action fa fa-copy"></span></th>
							<th class="threat-metrics type-ep"><span>DTPS EP</span><span class="col-action fa fa-copy"></span></th>
							<th><span>Current EP</span><span class="col-action fa fa-recycle"></span></th>
						</tr>
					</tbody>
				</table>
			</div>
		`;

		this.tableContainer = this.rootElem.getElementsByClassName('ep-weights-table')[0] as HTMLElement;
		this.tableBody = this.rootElem.querySelector('#ep-tbody') as HTMLElement;
		this.tableHeader = this.rootElem.querySelector('#ep-tbody > tr') as HTMLElement;

		const resultsViewerElem = this.rootElem.getElementsByClassName('ep-weights-results')[0] as HTMLElement;
		this.resultsViewer = new ResultsViewer(resultsViewerElem);

		const updateType = () => {
			if (this.statsType == 'ep') {
				this.tableContainer.classList.remove('stats-type-weight');
				this.tableContainer.classList.add('stats-type-ep');
			} else {
				this.tableContainer.classList.add('stats-type-weight');
				this.tableContainer.classList.remove('stats-type-ep');
			}
		};

		const selectElem = this.rootElem.getElementsByClassName('ep-type-select')[0] as HTMLSelectElement;
		selectElem.addEventListener('input', event => {
			this.statsType = selectElem.value;
			updateType();
		});
		selectElem.value = this.statsType;
		updateType();

		const calcButton = this.rootElem.getElementsByClassName('calc-weights')[0] as HTMLElement;
		calcButton.addEventListener('click', async event => {
			this.resultsViewer.setPending();
			const iterations = this.simUI.sim.getIterations();
			const result = await this.simUI.player.computeStatWeights(TypedEvent.nextEventID(), this.epStats, this.epReferenceStat, (progress: ProgressMetrics) => {
				this.setSimProgress(progress);
			});
			this.resultsViewer.hideAll();
			this.simUI.prevEpIterations = iterations;
			this.simUI.prevEpSimResult = result;
			this.preprocessResults(result);
			this.updateTable(iterations, result);
		});

		const colActionButtons = Array.from(this.rootElem.getElementsByClassName('col-action')) as Array<HTMLSelectElement>;
		const makeUpdateWeights = (button: HTMLElement, labelTooltip: string, tooltip: string, weightsFunc: () => Array<number>) => {
			tippy(button.previousSibling, {
				'content': labelTooltip,
				'allowHTML': true,
			});
			tippy(button, {
				'content': tooltip,
				'allowHTML': true,
			});
			button.addEventListener('click', event => {
				this.simUI.player.setEpWeights(TypedEvent.nextEventID(), new Stats(weightsFunc()));
			});
		};
		makeUpdateWeights(colActionButtons[0], 'Per-point increase in DPS (Damage Per Second) for each stat.', 'Copy to Current EP', () => this.getPrevSimResult().dps!.weights);
		makeUpdateWeights(colActionButtons[1], `EP (Equivalency Points) for DPS (Damage Per Second) for each stat. Normalized by ${statNames[this.epReferenceStat]}.`, 'Copy to Current EP', () => this.getPrevSimResult().dps!.epValues);
		makeUpdateWeights(colActionButtons[2], 'Per-point increase in TPS (Threat Per Second) for each stat.', 'Copy to Current EP', () => this.getPrevSimResult().tps!.weights);
		makeUpdateWeights(colActionButtons[3], `EP (Equivalency Points) for TPS (Threat Per Second) for each stat. Normalized by ${statNames[this.epReferenceStat]}.`, 'Copy to Current EP', () => this.getPrevSimResult().tps!.epValues);
		makeUpdateWeights(colActionButtons[4], 'Per-point increase in DTPS (Damage Taken Per Second) for each stat.', 'Copy to Current EP', () => this.getPrevSimResult().dtps!.weights);
		makeUpdateWeights(colActionButtons[5], `EP (Equivalency Points) for DTPS (Damage Taken Per Second) for each stat. Normalized by ${statNames[Stat.StatArmor]}.`, 'Copy to Current EP', () => this.getPrevSimResult().dtps!.epValues);
		makeUpdateWeights(colActionButtons[6], 'Current EP Weights. Used to sort the gear selector menus.', 'Restore Default EP', () => this.simUI.individualConfig.defaults.epWeights.asArray());

		const showAllStatsContainer = this.rootElem.getElementsByClassName('show-all-stats-container')[0] as HTMLElement;
		new BooleanPicker(showAllStatsContainer, this, {
			label: 'Show All Stats',
			changedEvent: () => new TypedEvent(),
			getValue: () => this.tableContainer.classList.contains('show-all-stats'),
			setValue: (eventID: EventID, menu: EpWeightsMenu, newValue: boolean) => {
				if (newValue) {
					this.tableContainer.classList.add('show-all-stats');
				} else {
					this.tableContainer.classList.remove('show-all-stats');
				}
				this.applyAlternatingColors();
			},
		});

		this.updateTable(this.simUI.prevEpIterations || 1, this.getPrevSimResult());

		const optimizeGemsButton = this.rootElem.getElementsByClassName('optimize-gems')[0] as HTMLElement;
		tippy(optimizeGemsButton, {
			'content': `
				<p>Optimizes equipped gems to maximize EP, based on the values in <b>Current EP</b>.</p>
				<p>Does not change the meta gem, but ensures that its condition is met. Uses JC gems if Jewelcrafting is a selected profession.</p>
			`,
			'allowHTML': true,
		});
		optimizeGemsButton.addEventListener('click', event => this.optimizeGems(TypedEvent.nextEventID()));

		this.addCloseButton();
	}

	setSimProgress(progress: ProgressMetrics) {
		this.resultsViewer.setContent(`
  <div class="results-sim">
  			<div class=""> ${progress.completedSims} / ${progress.totalSims}<br>simulations complete</div>
  			<div class="">
				${progress.completedIterations} / ${progress.totalIterations}<br>iterations complete
			</div>
  </div>
`);
	}

	private preprocessResults(result: StatWeightsResult) {
		// Values for a school's power should never exceed the value for regular spell power.
		result.dps!.epValues.forEach((value, index) => {
			if (index == Stat.StatArcaneSpellPower ||
				index == Stat.StatFireSpellPower ||
				index == Stat.StatFrostSpellPower ||
				index == Stat.StatHolySpellPower ||
				index == Stat.StatNatureSpellPower ||
				index == Stat.StatShadowSpellPower) {
				if (value > result.dps!.epValues[Stat.StatSpellPower]) {
					const diff = value - result.dps!.epValues[Stat.StatSpellPower];
					result.dps!.epValues[index] = result.dps!.epValues[Stat.StatSpellPower];
					result.dps!.epValuesStdev[index] -= diff;
					const wdiff = result.dps!.weights[index] - result.dps!.weights[Stat.StatSpellPower];
					result.dps!.weights[index] = result.dps!.weights[Stat.StatSpellPower];
					result.dps!.weightsStdev[index] -= wdiff;
				}
			}
		});
	}

	private updateTable(iterations: number, result: StatWeightsResult) {
		this.tableHeader.remove();
		this.tableBody.innerHTML = '';
		this.tableBody.appendChild(this.tableHeader);

		const allStats = statOrder.filter(stat => ![Stat.StatMana, Stat.StatEnergy, Stat.StatRage].includes(stat));
		allStats.forEach(stat => {
			const row = this.makeTableRow(stat, iterations, result);
			if (!this.epStats.includes(stat)) {
				row.classList.add('non-ep-stat');
			}
			this.tableBody.appendChild(row);
		});

		this.applyAlternatingColors();
	}

	private makeTableRow(stat: Stat, iterations: number, result: StatWeightsResult): HTMLElement {
		const row = document.createElement('tr');
		row.innerHTML = `
			<td>${statNames[stat]}</td>
			<td class="stdev-cell type-weight"><span>${result.dps!.weights[stat].toFixed(2)}</span><span>${stDevToConf90(result.dps!.weightsStdev[stat], iterations).toFixed(2)}</span></td>
			<td class="stdev-cell type-ep"><span>${result.dps!.epValues[stat].toFixed(2)}</span><span>${stDevToConf90(result.dps!.epValuesStdev[stat], iterations).toFixed(2)}</span></td>
			<td class="stdev-cell threat-metrics type-weight"><span>${result.tps!.weights[stat].toFixed(2)}</span><span>${stDevToConf90(result.tps!.weightsStdev[stat], iterations).toFixed(2)}</span></td>
			<td class="stdev-cell threat-metrics type-ep"><span>${result.tps!.epValues[stat].toFixed(2)}</span><span>${stDevToConf90(result.tps!.epValuesStdev[stat], iterations).toFixed(2)}</span></td>
			<td class="stdev-cell threat-metrics type-weight"><span>${result.dtps!.weights[stat].toFixed(2)}</span><span>${stDevToConf90(result.dtps!.weightsStdev[stat], iterations).toFixed(2)}</span></td>
			<td class="stdev-cell threat-metrics type-ep"><span>${result.dtps!.epValues[stat].toFixed(2)}</span><span>${stDevToConf90(result.dtps!.epValuesStdev[stat], iterations).toFixed(2)}</span></td>
			<td class="current-ep"></td>
		`;

		const currentEpCell = row.querySelector('.current-ep') as HTMLElement;
		new NumberPicker(currentEpCell, this.simUI.player, {
			float: true,
			changedEvent: (player: Player<any>) => player.epWeightsChangeEmitter,
			getValue: (player: Player<any>) => player.getEpWeights().getStat(stat),
			setValue: (eventID: EventID, player: Player<any>, newValue: number) => {
				const epWeights = player.getEpWeights().withStat(stat, newValue);
				player.setEpWeights(eventID, epWeights);
			},
		});

		return row;
	}

	private applyAlternatingColors() {
		(Array.from(this.tableBody.childNodes) as Array<HTMLElement>)
			.filter(row => window.getComputedStyle(row).getPropertyValue('display') != 'none')
			.forEach((row, i) => {
				if (i % 2 == 0) {
					row.classList.remove('odd');
				} else {
					row.classList.add('odd');
				}
			});
	}

	private getPrevSimResult(): StatWeightsResult {
		return this.simUI.prevEpSimResult || StatWeightsResult.create({
			dps: {
				weights: new Stats().asArray(),
				weightsStdev: new Stats().asArray(),
				epValues: new Stats().asArray(),
				epValuesStdev: new Stats().asArray(),
			},
			tps: {
				weights: new Stats().asArray(),
				weightsStdev: new Stats().asArray(),
				epValues: new Stats().asArray(),
				epValuesStdev: new Stats().asArray(),
			},
			dtps: {
				weights: new Stats().asArray(),
				weightsStdev: new Stats().asArray(),
				epValues: new Stats().asArray(),
				epValuesStdev: new Stats().asArray(),
			},
		});
	}

	private optimizeGems(eventID: EventID) {
		// Replace 0 weights with a very tiny value, so we always prefer to take free stats even if the user gave a 0 weight.
		let epWeights = this.simUI.player.getEpWeights();
		epWeights = new Stats(epWeights.asArray().map(w => w == 0 ? 1e-8 : w));

		const gear = this.simUI.player.getGear();
		const allGems = this.simUI.sim.getGems();
		const phase = this.simUI.sim.getPhase();
		const isBlacksmithing = this.simUI.player.isBlacksmithing();
		const isJewelcrafting = this.simUI.player.hasProfession(Profession.Jewelcrafting);

		const optimizedGear = EpWeightsMenu.optimizeGemsForWeights(epWeights, gear, allGems, phase, isBlacksmithing, isJewelcrafting);
		this.simUI.player.setGear(eventID, optimizedGear);
	}

	private static optimizeGemsForWeights(epWeights: Stats, gear: Gear, allGems: Array<Gem>, phase: number, isBlacksmithing: boolean, isJewelcrafting: boolean): Gear {
		const unrestrictedGems = allGems.filter(gem => Gems.isUnrestrictedGem(gem, phase));

		const {
			bestGemForColor: bestGemForColor,
			bestGemForColorEP: bestGemForColorEP,
			bestGemForSocket: bestGemForSocket,
			bestGemForSocketEP: bestGemForSocketEP,
			bestGem: bestGem,
			bestGemEP: bestGemEP,
		} = EpWeightsMenu.findBestGems(unrestrictedGems, epWeights);

		const items = gear.asMap();
		const socketBonusEPs = Object.values(items).map(item => item != null ? new Stats(item.item.socketBonus).computeEP(epWeights) : 0);

		// Start by optimally filling all items, ignoring meta condition.
		Object.entries(items).forEach(([itemSlot, equippedItem], i) => {
			if (equippedItem == null) {
				return;
			}
			const item = equippedItem.item;
			const socketColors = equippedItem.curSocketColors(isBlacksmithing);

			// Compare whether its better to match sockets + get socket bonus, or just use best gems.
			const bestGemEPNotMatchingSockets = sum(socketColors.map(socketColor => socketColor == GemColor.GemColorMeta ? 0 : bestGemEP));
			const bestGemEPMatchingSockets = socketBonusEPs[i] + sum(socketColors.map(socketColor => socketColor == GemColor.GemColorMeta ? 0 : bestGemForSocketEP[socketColor]));

			if (bestGemEPNotMatchingSockets > bestGemEPMatchingSockets) {
				socketColors.forEach((socketColor, i) => {
					if (socketColor != GemColor.GemColorMeta) {
						equippedItem = equippedItem!.withGem(bestGem, i);
					}
				});
			} else {
				socketColors.forEach((socketColor, i) => {
					if (socketColor != GemColor.GemColorMeta) {
						equippedItem = equippedItem!.withGem(bestGemForSocket[socketColor], i);
					}
				});
			}

			items[Number(itemSlot) as ItemSlot] = equippedItem;
		});
		gear = new Gear(items);

		const allSockets: Array<{itemSlot: ItemSlot, socketIdx: number}> = Object.keys(items).map((itemSlotStr) => {
			const itemSlot = parseInt(itemSlotStr) as ItemSlot;
			const item = items[itemSlot];
			if (!item) {
				return [];
			}

			const numSockets = item.numSockets(isBlacksmithing);
			return [...Array(numSockets).keys()]
			.filter(socketIdx => item.item.gemSockets[socketIdx] != GemColor.GemColorMeta)
			.map(socketIdx => {
				return {
					itemSlot: itemSlot,
					socketIdx: socketIdx,
				};
			});
		}).flat();
		const threeSocketCombos = permutations(allSockets, 3);
		const calculateGearGemsEP = (gear: Gear): number => gear.statsFromGems(isBlacksmithing).computeEP(epWeights);

		// Now make adjustments to satisfy meta condition.
		// Use a wrapper function so we can return for readability.
		gear = ((gear: Gear): Gear => {
			const metaGem = gear.getMetaGem();
			if (!metaGem) {
				return gear;
			}

			const condition = Gems.getMetaGemCondition(metaGem.id);
			// Only TBC gems use compare color conditions, so just ignore them.
			if (!condition || condition.isCompareColorCondition()) {
				return gear;
			}

			// If there are very few non-meta gem slots, just skip because it's annoying to deal with.
			if (gear.getAllGems(isBlacksmithing).length - 1 < 3) {
				return gear;
			}

			// In wrath, all meta gems use min colors condition (numRed >= r && numYellow >= y && numBlue >= b)
			// All conditions require 3 gems, e.g. 3 of a single color, 2 of one color and 1 of another, or 1 of each.
			// So the maximum number of gems that ever need to change is 3.

			const colorCombos = EpWeightsMenu.getColorCombosToSatisfyCondition(condition);

			let bestGear = gear;
			let bestGearEP = calculateGearGemsEP(gear);

			// Use brute-force to try every possibility.
			colorCombos.forEach(colorCombo => {
				threeSocketCombos.forEach(socketCombo => {
					const curItems = gear.asMap();
					for (let i = 0; i < colorCombo.length; i++) {
						const gemColor = colorCombo[i];
						const { itemSlot, socketIdx } = socketCombo[i];
						curItems[itemSlot] = curItems[itemSlot]!.withGem(bestGemForColor[gemColor], socketIdx);
					}
					const curGear = new Gear(curItems);
					if (curGear.hasActiveMetaGem(isBlacksmithing)) {
						const curGearEP = calculateGearGemsEP(curGear);
						if (curGearEP > bestGearEP) {
							bestGear = curGear;
							bestGearEP = curGearEP;
						}
					}
				});
			});

			return bestGear;
		})(gear);

		// Now insert 3 JC gems, if Jewelcrafting is selected.
		// Use a wrapper function so we can return for readability.
		gear = ((gear: Gear): Gear => {
			if (!isJewelcrafting) {
				return gear;
			}

			const jcGems = allGems.filter(gem => gem.requiredProfession == Profession.Jewelcrafting);

			const {
				bestGemForColor: bestJcGemForColor,
				bestGemForColorEP: bestJcGemForColorEP,
				bestGemForSocket: bestJcGemForSocket,
				bestGemForSocketEP: bestJcGemForSocketEP,
				bestGem: bestJcGem,
				bestGemEP: bestJcGemEP,
			} = EpWeightsMenu.findBestGems(jcGems, epWeights);

			let bestGear = gear;
			let bestGearEP = calculateGearGemsEP(gear);

			threeSocketCombos.forEach(socketCombo => {
				const curItems = gear.asMap();
				for (let i = 0; i < socketCombo.length; i++) {
					const { itemSlot, socketIdx } = socketCombo[i];
					const ei = curItems[itemSlot]!;
					const gemColor = ei.gems[socketIdx]!.color;
					curItems[itemSlot] = ei.withGem(bestJcGemForColor[gemColor], socketIdx);
				}

				const curGear = new Gear(curItems);
				if (curGear.hasActiveMetaGem(isBlacksmithing)) {
					const curGearEP = calculateGearGemsEP(curGear);
					if (curGearEP > bestGearEP) {
						bestGear = curGear;
						bestGearEP = curGearEP;
					}
				}
			});

			return bestGear;
		})(gear);

		return gear;
	}

	// Returns every possible way we could satisfy the gem condition.
	private static getColorCombosToSatisfyCondition(condition: Gems.MetaGemCondition): Array<Array<GemColor>> {
		if (condition.isOneOfEach()) {
			return [
				Gems.PRIMARY_COLORS,
				[GemColor.GemColorPrismatic],
			].concat(
				Gems.SECONDARY_COLORS.map((secondaryColor, i) => {
					const remainingColor = Gems.PRIMARY_COLORS[i];
					return Gems.socketToMatchingColors.get(remainingColor)!.map(matchingColor => [matchingColor, secondaryColor]);
				}).flat()
			);
		} else if (condition.isTwoAndOne()) {
			const oneColor = Gems.PRIMARY_COLORS[[condition.minRed, condition.minYellow, condition.minBlue].indexOf(1)];
			const twoColor = Gems.PRIMARY_COLORS[[condition.minRed, condition.minYellow, condition.minBlue].indexOf(2)];
			const secondaryColor = Gems.SECONDARY_COLORS.find(color => Gems.gemColorMatchesSocket(color, oneColor) && Gems.gemColorMatchesSocket(color, twoColor))!;

			return [
				// All the ways to get 1 point in both colors. These are partial combos,
				// which still need 1 more gem in the 2-color.
				[GemColor.GemColorPrismatic],
				[secondaryColor],
				[oneColor, twoColor],
			].map(partialCombo => {
					return Gems.socketToMatchingColors.get(twoColor)!.map(matchingColor => partialCombo.concat([matchingColor]));
			}).flat();
		} else if (condition.isThreeOfAColor()) {
			const threeColor = Gems.PRIMARY_COLORS[[condition.minRed, condition.minYellow, condition.minBlue].indexOf(3)];
			const matchingColors = Gems.socketToMatchingColors.get(threeColor)!;
			return combinationsWithDups(matchingColors, 3);
		} else {
			return [];
		}
	}

	private static findBestGems(gemList: Array<Gem>, epWeights: Stats): BestGemsResult {
		// Best gem when we need a gem of a specific color.
		const bestGemForColor: Array<Gem> = Gems.GEM_COLORS.map(color => null as unknown as Gem);
		const bestGemForColorEP: Array<number> = Gems.GEM_COLORS.map(color => 0);
		// Best gem when we need to match a socket to activate a bonus.
		const bestGemForSocket: Array<Gem> = bestGemForColor.slice();
		const bestGemForSocketEP: Array<number> = bestGemForColorEP.slice();
		// The single best gem, when color doesn't matter.
		let bestGem = gemList[0];
		let bestGemEP = 0;
		gemList.forEach(gem => {
			const gemEP = new Stats(gem.stats).computeEP(epWeights);
			if (gemEP > bestGemForColorEP[gem.color]) {
				bestGemForColorEP[gem.color] = gemEP;
				bestGemForColor[gem.color] = gem;

				if (gem.color != GemColor.GemColorMeta && gemEP > bestGemEP) {
					bestGemEP = gemEP;
					bestGem = gem;
				}
			}

			Gems.GEM_COLORS.forEach(socketColor => {
				if (Gems.gemMatchesSocket(gem, socketColor) && gemEP > bestGemForSocketEP[socketColor]) {
					bestGemForSocketEP[socketColor] = gemEP;
					bestGemForSocket[socketColor] = gem;
				}
			});
		});

		return {
			bestGemForColor: bestGemForColor,
			bestGemForColorEP: bestGemForColorEP,
			bestGemForSocket: bestGemForSocket,
			bestGemForSocketEP: bestGemForSocketEP,
			bestGem: bestGem,
			bestGemEP: bestGemEP,
		};
	}
}

interface BestGemsResult {
	bestGemForColor: Array<Gem>,
	bestGemForColorEP: Array<number>,
	bestGemForSocket: Array<Gem>,
	bestGemForSocketEP: Array<number>,
	bestGem: Gem,
	bestGemEP: number,
}
