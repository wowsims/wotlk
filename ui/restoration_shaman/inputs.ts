import { Spec } from '../core/proto/common.js';
import { Player } from '../core/player.js';
import { ActionId } from '../core/proto_utils/action_id.js';
import { TypedEvent } from '../core/typed_event.js';

import {
	ShamanShield,
} from '../core/proto/shaman.js';

import * as InputHelpers from '../core/components/input_helpers.js';

// Configuration for spec-specific UI elements on the settings tab.
// These don't need to be in a separate file but it keeps things cleaner.

export const ShamanShieldInput = InputHelpers.makeSpecOptionsEnumIconInput<Spec.SpecRestorationShaman, ShamanShield>({
	fieldName: 'shield',
	values: [
		{ value: ShamanShield.NoShield, tooltip: 'No Shield' },
		{ actionId: ActionId.fromSpellId(57960), value: ShamanShield.WaterShield },
		{ actionId: ActionId.fromSpellId(49281), value: ShamanShield.LightningShield },
	],
});

export const TriggerEarthShield = InputHelpers.makeSpecOptionsNumberInput<Spec.SpecRestorationShaman>({
	fieldName: 'earthShieldPPM',
	label: 'Earth Shield PPM',
	labelTooltip: 'How many times Earth Shield should be triggered per minute.',
	showWhen: (player: Player<Spec.SpecRestorationShaman>) => player.getTalents().earthShield,
	changeEmitter: (player: Player<Spec.SpecRestorationShaman>) => TypedEvent.onAny([player.specOptionsChangeEmitter, player.rotationChangeEmitter, player.talentsChangeEmitter]),
});

