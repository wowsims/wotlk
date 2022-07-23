import { BalanceDruid_Rotation_PrimarySpell as PrimarySpell } from '/wotlk/core/proto/druid.js';
import { BalanceDruid_Options as DruidOptions } from '/wotlk/core/proto/druid.js';
import { RaidTarget } from '/wotlk/core/proto/common.js';
import { Spec } from '/wotlk/core/proto/common.js';
import { NO_TARGET } from '/wotlk/core/proto_utils/utils.js';
import { ActionId } from '/wotlk/core/proto_utils/action_id.js';
import { Player } from '/wotlk/core/player.js';
import { Sim } from '/wotlk/core/sim.js';
import { EventID, TypedEvent } from '/wotlk/core/typed_event.js';
import { IndividualSimUI } from '/wotlk/core/individual_sim_ui.js';
import { Target } from '/wotlk/core/target.js';

import * as InputHelpers from '/wotlk/core/components/input_helpers.js';

// Configuration for spec-specific UI elements on the settings tab.
// These don't need to be in a separate file but it keeps things cleaner.

export const SelfInnervate = InputHelpers.makeSpecOptionsBooleanIconInput<Spec.SpecBalanceDruid>({
	fieldName: 'innervateTarget',
	id: ActionId.fromSpellId(29166),
	extraCssClasses: [
		'within-raid-sim-hide',
	],
	getValue: (player: Player<Spec.SpecBalanceDruid>) => player.getSpecOptions().innervateTarget?.targetIndex != NO_TARGET,
	setValue: (eventID: EventID, player: Player<Spec.SpecBalanceDruid>, newValue: boolean) => {
		const newOptions = player.getSpecOptions();
		newOptions.innervateTarget = RaidTarget.create({
			targetIndex: newValue ? 0 : NO_TARGET,
		});
		player.setSpecOptions(eventID, newOptions);
	},
});

export const BalanceDruidRotationConfig = {
	inputs: [
		InputHelpers.makeRotationEnumInput<Spec.SpecBalanceDruid, PrimarySpell>({
			fieldName: 'primarySpell',
			label: 'Primary Spell',
			labelTooltip: 'If set to \'Adaptive\', will dynamically adjust rotation based on available mana.',
			values: [
				{ name: 'Adaptive', value: PrimarySpell.Adaptive },
				{ name: 'Starfire', value: PrimarySpell.Starfire },
				{ name: 'Starfire R6', value: PrimarySpell.Starfire6 },
				{ name: 'Wrath', value: PrimarySpell.Wrath },
			],
		}),
		InputHelpers.makeRotationBooleanInput<Spec.SpecBalanceDruid>({
			fieldName: 'moonfire',
			label: 'Use Moonfire',
			labelTooltip: 'Use Moonfire as the next cast after the dot expires.',
			enableWhen: (player: Player<Spec.SpecBalanceDruid>) => player.getRotation().primarySpell != PrimarySpell.Adaptive,
		}),
		InputHelpers.makeRotationBooleanInput<Spec.SpecBalanceDruid>({
			fieldName: 'faerieFire',
			label: 'Use Faerie Fire',
			labelTooltip: 'Keep Faerie Fire active on the primary target.',
		}),
		InputHelpers.makeRotationBooleanInput<Spec.SpecBalanceDruid>({
			fieldName: 'insectSwarm',
			label: 'Use Insect Swarm',
			labelTooltip: 'Keep Insect Swarm active on the primary target.',
			enableWhen: (player: Player<Spec.SpecBalanceDruid>) => player.getTalents().insectSwarm,
		}),
		InputHelpers.makeRotationBooleanInput<Spec.SpecBalanceDruid>({
			fieldName: 'hurricane',
			label: 'Use Hurricane',
			labelTooltip: 'Casts Hurricane on cooldown.',
		}),
		InputHelpers.makeSpecOptionsBooleanInput<Spec.SpecBalanceDruid>({
			fieldName: 'battleRes',
			label: 'Use Battle Res',
			labelTooltip: 'Cast Battle Res on an ally sometime during the encounter.',
		}),
	],
};
