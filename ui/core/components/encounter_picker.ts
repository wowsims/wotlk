import {
	InputType,
	MobType,
	SpellSchool,
	Stat,
	Target as TargetProto,
 } from '../proto/common.js';
import { Encounter } from '../encounter.js';
import { Raid } from '../raid.js';
import { Target } from '../target.js';
import { EventID, TypedEvent } from '../typed_event.js';
import { BooleanPicker } from '../components/boolean_picker.js';
import { EnumPicker } from '../components/enum_picker.js';
import { ListPicker } from '../components/list_picker.js';
import { NumberPicker } from '../components/number_picker.js';
import { isHealingSpec, isTankSpec } from '../proto_utils/utils.js';
import { statNames } from '../proto_utils/names.js';

import { Component } from './component.js';

import * as Mechanics from '../constants/mechanics.js';
import { IndividualSimUI } from '../individual_sim_ui.js';
import { SimUI } from '../sim_ui.js';
import { Input } from './input.js';
import { BaseModal } from './base_modal.js';
import { TargetInputs } from '../target_inputs.js';

export interface EncounterPickerConfig {
	showExecuteProportion: boolean,
}

export class EncounterPicker extends Component {
	constructor(parent: HTMLElement, modEncounter: Encounter, config: EncounterPickerConfig, simUI: SimUI) {
		super(parent, 'encounter-picker-root');

		addEncounterFieldPickers(this.rootElem, modEncounter, config.showExecuteProportion);

		// Need to wait so that the encounter and target presets will be loaded.
		modEncounter.sim.waitForInit().then(() => {
			const presetTargets = modEncounter.sim.db.getAllPresetTargets();

			new EnumPicker<Encounter>(this.rootElem, modEncounter, {
				extraCssClasses: ['damage-metrics', 'npc-picker'],
				label: 'NPC',
				labelTooltip: 'Selects a preset NPC configuration.',
				values: [
					{ name: 'Custom', value: -1 },
				].concat(presetTargets.map((pe, i) => {
					return {
						name: pe.path,
						value: i,
					};
				})),
				changedEvent: (encounter: Encounter) => encounter.changeEmitter,
				getValue: (encounter: Encounter) => presetTargets.findIndex(pe => encounter.primaryTarget.matchesPreset(pe)),
				setValue: (eventID: EventID, encounter: Encounter, newValue: number) => {
					if (newValue != -1) {
						encounter.primaryTarget.applyPreset(eventID, presetTargets[newValue]);
					}

					EncounterPicker.updatePrimaryTargetInputs(encounter.primaryTarget);
				},
			});

			//new EnumPicker<Encounter>(this.rootElem, modEncounter, {
			//	label: 'Target Level',
			//	values: [
			//		{ name: '83', value: 83 },
			//		{ name: '82', value: 82 },
			//		{ name: '81', value: 81 },
			//		{ name: '80', value: 80 },
			//	],
			//	changedEvent: (encounter: Encounter) => encounter.changeEmitter,
			//	getValue: (encounter: Encounter) => encounter.primaryTarget.getLevel(),
			//	setValue: (eventID: EventID, encounter: Encounter, newValue: number) => {
			//		encounter.primaryTarget.setLevel(eventID, newValue);
			//	},
			//});

			//new EnumPicker(this.rootElem, modEncounter, {
			//	label: 'Mob Type',
			//	values: mobTypeEnumValues,
			//	changedEvent: (encounter: Encounter) => encounter.changeEmitter,
			//	getValue: (encounter: Encounter) => encounter.primaryTarget.getMobType(),
			//	setValue: (eventID: EventID, encounter: Encounter, newValue: number) => {
			//		encounter.primaryTarget.setMobType(eventID, newValue);
			//	},
			//});

			// Leaving this commented in case we want it later. But it takes up a lot of
			// screen space and none of these fields get changed much.
			//if (config.simpleTargetStats) {
			//	config.simpleTargetStats.forEach(stat => {
			//		new NumberPicker(this.rootElem, modEncounter, {
			//			label: statNames[stat],
			//			changedEvent: (encounter: Encounter) => encounter.changeEmitter,
			//			getValue: (encounter: Encounter) => encounter.primaryTarget.getStats().getStat(stat),
			//			setValue: (eventID: EventID, encounter: Encounter, newValue: number) => {
			//				encounter.primaryTarget.setStats(eventID, encounter.primaryTarget.getStats().withStat(stat, newValue));
			//			},
			//		});
			//	});
			//}

			if (simUI.isIndividualSim() && isHealingSpec((simUI as IndividualSimUI<any>).player.spec)) {
				new NumberPicker(this.rootElem, simUI.sim.raid, {
					label: 'Num Allies',
					labelTooltip: 'Number of allied players in the raid.',
					changedEvent: (raid: Raid) => raid.targetDummiesChangeEmitter,
					getValue: (raid: Raid) => raid.getTargetDummies(),
					setValue: (eventID: EventID, raid: Raid, newValue: number) => {
						raid.setTargetDummies(eventID, newValue);
					},
				});
			}

			if (simUI.isIndividualSim() && isTankSpec((simUI as IndividualSimUI<any>).player.spec)) {
				new NumberPicker(this.rootElem, modEncounter, {
					label: 'Min Base Damage',
					labelTooltip: 'Base damage for auto attacks, i.e. lowest roll with 0 AP against a 0-armor Player.',
					changedEvent: (encounter: Encounter) => encounter.changeEmitter,
					getValue: (encounter: Encounter) => encounter.primaryTarget.getMinBaseDamage(),
					setValue: (eventID: EventID, encounter: Encounter, newValue: number) => {
						encounter.primaryTarget.setMinBaseDamage(eventID, newValue);
					},
				});
			}

			const advancedButton = document.createElement('button');
			advancedButton.classList.add('advanced-button', 'btn', 'btn-primary');
			advancedButton.textContent = 'Advanced';
			advancedButton.addEventListener('click', () => new AdvancedEncounterModal(simUI.rootElem, simUI, modEncounter));
			this.rootElem.appendChild(advancedButton);

			// Transfer Target Inputs from target Id if they dont match (possible when custom AI is selected)
			let targetIndex = presetTargets.findIndex(pe => modEncounter.primaryTarget.getId() == pe.target?.id);
			let targetInputs = new TargetInputs(presetTargets[targetIndex]?.target?.targetInputs);

			if (targetInputs.getLength() != modEncounter.primaryTarget.getTargetInputsLength()) {
				modEncounter.primaryTarget.setTargetInputs(TypedEvent.nextEventID(), targetInputs);
			} else {
				let isDiff = false
				for (let i = 0; i < modEncounter.primaryTarget.getTargetInputsLength(); i++) {
					if (modEncounter.primaryTarget.getTargetInputs().getTargetInput(i).label != targetInputs.getTargetInput(i)?.label) {
						isDiff = true
						break;
					}
				}

				if (isDiff) {
					modEncounter.primaryTarget.setTargetInputs(TypedEvent.nextEventID(), targetInputs);
				}
			}

			EncounterPicker.primaryRootElem = this.rootElem;
			EncounterPicker.updatePrimaryTargetInputs(modEncounter.primaryTarget);
		});
	}

	static primaryRootElem: HTMLElement;
	static primaryPickers = new Array<Component>();

	static clearTargetInputPickers(targetInputPickers: Array<Component>) {
		targetInputPickers.forEach(picker => {
			picker?.rootElem.remove()
			picker?.dispose()
		})
	}

	static updatePrimaryTargetInputs(target: Target) {
		EncounterPicker.clearTargetInputPickers(EncounterPicker.primaryPickers)
		EncounterPicker.primaryPickers = []
		EncounterPicker.rebuildTargetInputs(EncounterPicker.primaryRootElem, target, EncounterPicker.primaryPickers, true)
	}

	static rebuildTargetInputs(rootElem: HTMLElement, target: Target, pickers: Array<Component>, beforeLast: boolean) {
		if (target.hasTargetInputs()) {
			for (let index = 0; index < target.getTargetInputsLength(); index++) {
				let targetInput = target.getTargetInputs().getTargetInput(index)
				if (targetInput.inputType == InputType.Number) {
					let numberPicker = new NumberPicker(rootElem, target, {
						label: targetInput.label,
						labelTooltip: targetInput.tooltip,
						changedEvent: (target: Target) => target.propChangeEmitter,
						getValue: (target: Target) => target.getTargetInputNumberValue(index),
						setValue: (eventID: EventID, target: Target, newValue: number) => {
							target.setTargetInputNumberValue(eventID, index, newValue)
						},
					});
					if (beforeLast) {
						let parent = numberPicker.rootElem.parentElement;
						parent?.removeChild(numberPicker.rootElem)
						parent?.insertBefore(numberPicker.rootElem, parent.lastChild)
					}
					pickers.push(numberPicker)
				} else if (targetInput.inputType == InputType.Bool) {
					let booleanPicker = new BooleanPicker(rootElem, target, {
						label: targetInput.label,
						labelTooltip: targetInput.tooltip,
						changedEvent: (target: Target) => target.propChangeEmitter,
						getValue: (target: Target) => target.getTargetInputBooleanValue(index),
						setValue: (eventID: EventID, target: Target, newValue: boolean) => {
							target.setTargetInputBooleanValue(eventID, index, newValue);
						},
					});
					if (beforeLast) {
						let parent = booleanPicker.rootElem.parentElement;
						parent?.removeChild(booleanPicker.rootElem)
						parent?.insertBefore(booleanPicker.rootElem, parent.lastChild)
					}
					pickers.push(booleanPicker)
				}
			}
		}
	}
}

class AdvancedEncounterModal extends BaseModal {
	private readonly encounter: Encounter;

	constructor(parent: HTMLElement, simUI: SimUI, encounter: Encounter) {
		super(parent, 'advanced-encounter-picker-modal');

		this.encounter = encounter;

		this.addHeader();
		this.body.innerHTML = `
			<div class="encounter-header"></div>
			<div class="encounter-targets"></div>
		`;

		const header = this.rootElem.getElementsByClassName('encounter-header')[0] as HTMLElement;
		const targetsElem = this.rootElem.getElementsByClassName('encounter-targets')[0] as HTMLElement;

		addEncounterFieldPickers(header, this.encounter, true);
		if (!simUI.isIndividualSim()) {
			new BooleanPicker<Encounter>(header, encounter, {
				label: 'Use Health',
				labelTooltip: 'Uses a damage limit in place of a duration limit. Damage limit is equal to sum of all targets health.',
				changedEvent: (encounter: Encounter) => encounter.changeEmitter,
				getValue: (encounter: Encounter) => encounter.getUseHealth(),
				setValue: (eventID: EventID, encounter: Encounter, newValue: boolean) => {
					encounter.setUseHealth(eventID, newValue);
				},
			});
		}
		new ListPicker<Encounter, Target>(targetsElem, simUI, this.encounter, {
			extraCssClasses: ['targets-picker', 'mb-0'],
			itemLabel: 'Target',
			changedEvent: (encounter: Encounter) => encounter.targetsChangeEmitter,
			getValue: (encounter: Encounter) => encounter.getTargets(),
			setValue: (eventID: EventID, encounter: Encounter, newValue: Array<Target>) => {
				encounter.setTargets(eventID, newValue);
			},
			newItem: () => Target.fromDefaults(TypedEvent.nextEventID(), this.encounter.sim),
			copyItem: (oldItem: Target) => oldItem.clone(TypedEvent.nextEventID()),
			newItemPicker: (parent: HTMLElement, target: Target) => new TargetPicker(parent, target),
		});
	}

	private addHeader() {
		const presetEncounters = this.encounter.sim.db.getAllPresetEncounters();

		new EnumPicker<Encounter>(this.header as HTMLElement, this.encounter, {
			label: 'Encounter',
			extraCssClasses: ['encounter-picker', 'mb-0', 'pe-2'],
			values: [
				{ name: 'Custom', value: -1 },
			].concat(presetEncounters.map((pe, i) => {
				return {
					name: pe.path,
					value: i,
				};
			})),
			changedEvent: (encounter: Encounter) => encounter.changeEmitter,
			getValue: (encounter: Encounter) => presetEncounters.findIndex(pe => encounter.matchesPreset(pe)),
			setValue: (eventID: EventID, encounter: Encounter, newValue: number) => {
				if (newValue != -1) {
					encounter.applyPreset(eventID, presetEncounters[newValue]);
				}
				EncounterPicker.updatePrimaryTargetInputs(encounter.primaryTarget);
			},
		});
	}
}

class TargetPicker extends Input<Target, Target> {
	constructor(parent: HTMLElement, modTarget: Target) {
		super(parent, 'target-picker-root', modTarget, {
			changedEvent: (target: Target) => target.changeEmitter,
			getValue: (target: Target) => target,
			setValue: (eventID: EventID, target: Target, newValue: Target) => {
				target.fromProto(eventID, newValue.toProto());
			},
		});
		this.rootElem.innerHTML = `
			<div class="target-picker-section target-picker-section1"></div>
			<div class="target-picker-section target-picker-section2"></div>
			<div class="target-picker-section target-picker-section3 threat-metrics"></div>
		`;

		const encounter = modTarget.sim.encounter;
		const section1 = this.rootElem.getElementsByClassName('target-picker-section1')[0] as HTMLElement;
		const section2 = this.rootElem.getElementsByClassName('target-picker-section2')[0] as HTMLElement;
		const section3 = this.rootElem.getElementsByClassName('target-picker-section3')[0] as HTMLElement;

		let targetInputPickers = new Array<Component>();

		const presetTargets = modTarget.sim.db.getAllPresetTargets();
		new EnumPicker<Target>(section1, modTarget, {
			extraCssClasses: ['npc-picker'],
			label: 'NPC',
			labelTooltip: 'Selects a preset NPC configuration.',
			values: [
				{ name: 'Custom', value: -1 },
			].concat(presetTargets.map((pe, i) => {
				return {
					name: pe.path,
					value: i,
				};
			})),
			changedEvent: (target: Target) => target.changeEmitter,
			getValue: (target: Target) => presetTargets.findIndex(pe => target.matchesPreset(pe)),
			setValue: (eventID: EventID, target: Target, newValue: number) => {
				if (newValue != -1) {
					target.applyPreset(eventID, presetTargets[newValue]);
				}
				
				EncounterPicker.clearTargetInputPickers(targetInputPickers);
				targetInputPickers = [];
				EncounterPicker.rebuildTargetInputs(section1, target, targetInputPickers, false);

				if (target == encounter.primaryTarget) {
					EncounterPicker.updatePrimaryTargetInputs(encounter.primaryTarget);
				}
			},
		});

		new EnumPicker<Target>(section1, modTarget, {
			extraCssClasses: ['ai-picker'],
			label: 'AI',
			labelTooltip: `
				<p>Determines the target\'s ability rotation.</p>
				<p>Note that most rotations are not yet implemented.</p>
			`,
			values: [
				{ name: 'None', value: 0 },
			].concat(presetTargets.map(pe => {
				return {
					name: pe.path,
					value: pe.target!.id,
				};
			})),
			changedEvent: (target: Target) => target.changeEmitter,
			getValue: (target: Target) => target.getId(),
			setValue: (eventID: EventID, target: Target, newValue: number) => {
				target.setId(eventID, newValue);

				// Transfer Target Inputs from the AI of the selected target
				let newTargetIndex = presetTargets.findIndex(pe => target.getId() == pe.target?.id);
				let newTargetInputs = new TargetInputs(presetTargets[newTargetIndex]?.target?.targetInputs);
				target.setTargetInputs(eventID, newTargetInputs);
				
				// Update picker elements
				EncounterPicker.clearTargetInputPickers(targetInputPickers);
				targetInputPickers = [];
				EncounterPicker.rebuildTargetInputs(section1, target, targetInputPickers, false);

				if (target == encounter.primaryTarget) {
					EncounterPicker.updatePrimaryTargetInputs(encounter.primaryTarget);
				}
			},
		});

		new EnumPicker<Target>(section1, modTarget, {
			label: 'Level',
			values: [
				{ name: '83', value: 83 },
				{ name: '82', value: 82 },
				{ name: '81', value: 81 },
				{ name: '80', value: 80 },
			],
			changedEvent: (target: Target) => target.levelChangeEmitter,
			getValue: (target: Target) => target.getLevel(),
			setValue: (eventID: EventID, target: Target, newValue: number) => {
				target.setLevel(eventID, newValue);
			},
		});
		new EnumPicker(section1, modTarget, {
			label: 'Mob Type',
			values: mobTypeEnumValues,
			changedEvent: (target: Target) => target.mobTypeChangeEmitter,
			getValue: (target: Target) => target.getMobType(),
			setValue: (eventID: EventID, target: Target, newValue: number) => {
				target.setMobType(eventID, newValue);
			},
		});
		new EnumPicker<Target>(section1, modTarget, {
			extraCssClasses: ['threat-metrics'],
			label: 'Tanked By',
			labelTooltip: 'Determines which player in the raid this enemy will attack. If no player is assigned to the specified tank slot, this enemy will not attack.',
			values: [
				{ name: 'None', value: -1 },
				{ name: 'Main Tank', value: 0 },
				{ name: 'Tank 2', value: 1 },
				{ name: 'Tank 3', value: 2 },
				{ name: 'Tank 4', value: 3 },
			],
			changedEvent: (target: Target) => target.propChangeEmitter,
			getValue: (target: Target) => target.getTankIndex(),
			setValue: (eventID: EventID, target: Target, newValue: number) => {
				target.setTankIndex(eventID, newValue);
			},
		});

		EncounterPicker.rebuildTargetInputs(section1, modTarget, targetInputPickers, false);

		ALL_TARGET_STATS.forEach(statData => {
			const stat = statData.stat;
			new NumberPicker(section2, modTarget, {
				inline: true,
				extraCssClasses: statData.extraCssClasses,
				label: statNames[stat],
				labelTooltip: statData.tooltip,
				changedEvent: (target: Target) => target.statsChangeEmitter,
				getValue: (target: Target) => target.getStats().getStat(stat),
				setValue: (eventID: EventID, target: Target, newValue: number) => {
					target.setStats(eventID, target.getStats().withStat(stat, newValue));
				},
			});
		});

		new NumberPicker(section3, modTarget, {
			label: 'Swing Speed',
			labelTooltip: 'Time in seconds between auto attacks. Set to 0 to disable auto attacks.',
			float: true,
			changedEvent: (target: Target) => target.propChangeEmitter,
			getValue: (target: Target) => target.getSwingSpeed(),
			setValue: (eventID: EventID, target: Target, newValue: number) => {
				target.setSwingSpeed(eventID, newValue);
			},
		});
		new NumberPicker(section3, modTarget, {
			label: 'Min Base Damage',
			labelTooltip: 'Base damage for auto attacks, i.e. lowest roll with 0 AP against a 0-armor Player.',
			changedEvent: (target: Target) => target.propChangeEmitter,
			getValue: (target: Target) => target.getMinBaseDamage(),
			setValue: (eventID: EventID, target: Target, newValue: number) => {
				target.setMinBaseDamage(eventID, newValue);
			},
		});
		new BooleanPicker(section3, modTarget, {
			label: 'Dual Wield',
			labelTooltip: 'Uses 2 separate weapons to attack.',
			inline: true,
			changedEvent: (target: Target) => target.propChangeEmitter,
			getValue: (target: Target) => target.getDualWield(),
			setValue: (eventID: EventID, target: Target, newValue: boolean) => {
				target.setDualWield(eventID, newValue);
			},
		});
		new BooleanPicker(section3, modTarget, {
			label: 'DW Miss Penalty',
			labelTooltip: 'Enables the Dual Wield Miss Penalty (+19% chance to miss) if dual wielding. Bosses in Hyjal/BT/SWP usually have this disabled to stop tanks from avoidance stacking.',
			inline: true,
			changedEvent: (target: Target) => target.changeEmitter,
			getValue: (target: Target) => target.getDualWieldPenalty(),
			setValue: (eventID: EventID, target: Target, newValue: boolean) => {
				target.setDualWieldPenalty(eventID, newValue);
			},
			enableWhen: (target: Target) => target.getDualWield(),
		});
		new BooleanPicker(section3, modTarget, {
			label: 'Parry Haste',
			labelTooltip: 'Whether this enemy will gain parry haste when parrying attacks.',
			inline: true,
			changedEvent: (target: Target) => target.propChangeEmitter,
			getValue: (target: Target) => target.getParryHaste(),
			setValue: (eventID: EventID, target: Target, newValue: boolean) => {
				target.setParryHaste(eventID, newValue);
			},
		});
		new EnumPicker<Target>(section3, modTarget, {
			label: 'Spell School',
			labelTooltip: 'Type of damage caused by auto attacks. This is usually Physical, but some enemies have elemental attacks.',
			values: [
				{ name: 'Physical', value: SpellSchool.SpellSchoolPhysical },
				{ name: 'Arcane', value: SpellSchool.SpellSchoolArcane },
				{ name: 'Fire', value: SpellSchool.SpellSchoolFire },
				{ name: 'Frost', value: SpellSchool.SpellSchoolFrost },
				{ name: 'Holy', value: SpellSchool.SpellSchoolHoly },
				{ name: 'Nature', value: SpellSchool.SpellSchoolNature },
				{ name: 'Shadow', value: SpellSchool.SpellSchoolShadow },
			],
			changedEvent: (target: Target) => target.propChangeEmitter,
			getValue: (target: Target) => target.getSpellSchool(),
			setValue: (eventID: EventID, target: Target, newValue: number) => {
				target.setSpellSchool(eventID, newValue);
			},
		});
		new BooleanPicker(section3, modTarget, {
			label: 'Chill of the Throne',
			labelTooltip: 'Reduces the chance for this enemy\'s attacks to be dodged by 20%. Active in Icecrown Citadel.',
			inline: true,
			changedEvent: (target: Target) => target.changeEmitter,
			getValue: (target: Target) => target.getSuppressDodge(),
			setValue: (eventID: EventID, target: Target, newValue: boolean) => {
				target.setSuppressDodge(eventID, newValue);
			},
			enableWhen: (target: Target) => target.getLevel() == Mechanics.BOSS_LEVEL,
		});
		new BooleanPicker(section3, modTarget, {
			label: 'Tightened Damage Range',
			labelTooltip: 'Reduces the damage range of this enemy\'s auto-attacks. Observed behavior for Patchwerk.',
			inline: true,
			changedEvent: (target: Target) => target.changeEmitter,
			getValue: (target: Target) => target.getTightEnemyDamage(),
			setValue: (eventID: EventID, target: Target, newValue: boolean) => {
				target.setTightEnemyDamage(eventID, newValue);
			},
			enableWhen: (target: Target) => target.getLevel() == Mechanics.BOSS_LEVEL,
		});
	}

	getInputElem(): HTMLElement|null {
		return null;
	}
	getInputValue(): Target {
		return this.value;
	}
	setInputValue(newValue: Target) {

	}
}

function addEncounterFieldPickers(rootElem: HTMLElement, encounter: Encounter, showExecuteProportion: boolean) {
	let durationGroup = Input.newGroupContainer();
	rootElem.appendChild(durationGroup);

	new NumberPicker(durationGroup, encounter, {
		label: 'Duration',
		labelTooltip: 'The fight length for each sim iteration, in seconds.',
		changedEvent: (encounter: Encounter) => encounter.changeEmitter,
		getValue: (encounter: Encounter) => encounter.getDuration(),
		setValue: (eventID: EventID, encounter: Encounter, newValue: number) => {
			encounter.setDuration(eventID, newValue);
		},
		enableWhen: (obj) => { return !encounter.getUseHealth() },
	});
	new NumberPicker(durationGroup, encounter, {
		label: 'Duration +/-',
		labelTooltip: 'Adds a random amount of time, in seconds, between [value, -1 * value] to each sim iteration. For example, setting Duration to 180 and Duration +/- to 10 will result in random durations between 170s and 190s.',
		changedEvent: (encounter: Encounter) => encounter.changeEmitter,
		getValue: (encounter: Encounter) => encounter.getDurationVariation(),
		setValue: (eventID: EventID, encounter: Encounter, newValue: number) => {
			encounter.setDurationVariation(eventID, newValue);
		},
		enableWhen: (obj) => { return !encounter.getUseHealth() },
	});

	if (showExecuteProportion) {
		let executeGroup = Input.newGroupContainer();
		executeGroup.classList.add('execute-group');
		rootElem.appendChild(executeGroup);

		new NumberPicker(executeGroup, encounter, {
			label: 'Execute Duration 20 (%)',
			labelTooltip: 'Percentage of the total encounter duration, for which the targets will be considered to be in execute range (< 20% HP) for the purpose of effects like Warrior Execute or Mage Molten Fury.',
			changedEvent: (encounter: Encounter) => encounter.changeEmitter,
			getValue: (encounter: Encounter) => encounter.getExecuteProportion20() * 100,
			setValue: (eventID: EventID, encounter: Encounter, newValue: number) => {
				encounter.setExecuteProportion20(eventID, newValue / 100);
			},
			enableWhen: (obj) => { return !encounter.getUseHealth() },
		});
		new NumberPicker(executeGroup, encounter, {
			label: 'Execute Duration 25 (%)',
			labelTooltip: 'Percentage of the total encounter duration, for which the targets will be considered to be in execute range (< 25% HP) for the purpose of effects like Warlock\'s Drain Soul.',
			changedEvent: (encounter: Encounter) => encounter.changeEmitter,
			getValue: (encounter: Encounter) => encounter.getExecuteProportion25() * 100,
			setValue: (eventID: EventID, encounter: Encounter, newValue: number) => {
				encounter.setExecuteProportion25(eventID, newValue / 100);
			},
			enableWhen: (obj) => { return !encounter.getUseHealth() },
		});
		new NumberPicker(executeGroup, encounter, {
			label: 'Execute Duration 35 (%)',
			labelTooltip: 'Percentage of the total encounter duration, for which the targets will be considered to be in execute range (< 35% HP) for the purpose of effects like Warrior Execute or Mage Molten Fury.',
			changedEvent: (encounter: Encounter) => encounter.changeEmitter,
			getValue: (encounter: Encounter) => encounter.getExecuteProportion35() * 100,
			setValue: (eventID: EventID, encounter: Encounter, newValue: number) => {
				encounter.setExecuteProportion35(eventID, newValue / 100);
			},
			enableWhen: (obj) => { return !encounter.getUseHealth() },
		});
	}
}

const ALL_TARGET_STATS: Array<{ stat: Stat, tooltip: string, extraCssClasses: Array<string> }> = [
	{ stat: Stat.StatHealth, tooltip: '', extraCssClasses: [] },
	{ stat: Stat.StatArmor, tooltip: '', extraCssClasses: [] },
	{ stat: Stat.StatArcaneResistance, tooltip: '', extraCssClasses: [] },
	{ stat: Stat.StatFireResistance, tooltip: '', extraCssClasses: [] },
	{ stat: Stat.StatFrostResistance, tooltip: '', extraCssClasses: [] },
	{ stat: Stat.StatNatureResistance, tooltip: '', extraCssClasses: [] },
	{ stat: Stat.StatShadowResistance, tooltip: '', extraCssClasses: [] },
	{ stat: Stat.StatAttackPower, tooltip: '', extraCssClasses: ['threat-metrics'] },
	{ stat: Stat.StatBlockValue, tooltip: '', extraCssClasses: ['threat-metrics'] },
];

const mobTypeEnumValues = [
	{ name: 'None', value: MobType.MobTypeUnknown },
	{ name: 'Beast', value: MobType.MobTypeBeast },
	{ name: 'Demon', value: MobType.MobTypeDemon },
	{ name: 'Dragonkin', value: MobType.MobTypeDragonkin },
	{ name: 'Elemental', value: MobType.MobTypeElemental },
	{ name: 'Giant', value: MobType.MobTypeGiant },
	{ name: 'Humanoid', value: MobType.MobTypeHumanoid },
	{ name: 'Mechanical', value: MobType.MobTypeMechanical },
	{ name: 'Undead', value: MobType.MobTypeUndead },
];
