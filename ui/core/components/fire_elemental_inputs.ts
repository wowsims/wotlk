import { IndividualSimUI } from "../individual_sim_ui";
import { Player } from "../player";
import { ShamanTotems } from "../proto/shaman";
import { ShamanSpecs } from "../proto_utils/utils";
import { EventID } from "../typed_event";
import { ContentBlock } from "./content_block";
import { IconPicker } from "./icon_picker";
import * as InputHelpers from '../components/input_helpers.js';
import { ActionId } from '../proto_utils/action_id.js';
import { Input } from "./input";
import { NumberPicker } from "./number_picker";
import { BooleanPicker } from "./boolean_picker";

export function FireElementalSection(parentElem: HTMLElement, simUI: IndividualSimUI<ShamanSpecs>): ContentBlock {
	let contentBlock = new ContentBlock(parentElem, 'fire-elemental-settings', {
		header: { title: 'Fire Elemental' }
	});

	let fireElementalIconContainer = Input.newGroupContainer();
	fireElementalIconContainer.classList.add('fire-elemental-icon-container');

	contentBlock.bodyElement.appendChild(fireElementalIconContainer);

	const fireElementalBooleanIconInput = InputHelpers.makeBooleanIconInput<ShamanSpecs, ShamanTotems, Player<ShamanSpecs>>({
		getModObject: (player: Player<ShamanSpecs>) => player,
		getValue: (player: Player<ShamanSpecs>) => player.getSpecOptions().totems || ShamanTotems.create(),
		setValue: (eventID: EventID, player: Player<ShamanSpecs>, newVal: ShamanTotems) => {
			const newOptions = player.getSpecOptions();
			newOptions.totems = newVal;
			player.setSpecOptions(eventID, newOptions);

			// Hacky fix ItemSwapping is in the Rotation proto, this will let the Rotation know to update showWhen
			// TODO move the ItemSwap enabled to a spec option and have the ItemSwap proto be apart of player.
			player.rotationChangeEmitter.emit(eventID)
		},
		changeEmitter: (player: Player<ShamanSpecs>) => player.specOptionsChangeEmitter,
	}, ActionId.fromSpellId(2894), "useFireElemental");

	new IconPicker(fireElementalIconContainer, simUI.player, fireElementalBooleanIconInput);

	new NumberPicker(contentBlock.bodyElement, simUI.player, {
		positive: true,
		label: "Bonus spell power",
		labelTooltip: "Bonus spell power to snapshot Fire Elemental with. Will prioritize dropping Fire Elemental if greater then 0",
		inline: true,
		getValue: (player: Player<ShamanSpecs>) => player.getSpecOptions().totems?.bonusSpellpower || 0,
		setValue: (eventID: EventID, player: Player<ShamanSpecs>, newVal: number) => {
			const newOptions = player.getSpecOptions();

			if (newOptions.totems) {
				newOptions.totems.bonusSpellpower = newVal
			}

			player.setSpecOptions(eventID, newOptions);
		},
		changedEvent: (player: Player<ShamanSpecs>) => player.specOptionsChangeEmitter,
	})

	new BooleanPicker(contentBlock.bodyElement, simUI.player, {
		label: "Use Tier 10 (4pc)",
		labelTooltip: "Will use Tier 10 (4pc) to snapshot Fire Elemental.",
		inline: true,
		getValue: (player: Player<ShamanSpecs>) => player.getSpecOptions().totems?.enhTierTenBonus || false,
		setValue: (eventID: EventID, player: Player<ShamanSpecs>, newVal: boolean) => {
			const newOptions = player.getSpecOptions();

			if (newOptions.totems) {
				newOptions.totems.enhTierTenBonus = newVal
			}

			player.setSpecOptions(eventID, newOptions);
		},
		changedEvent: (player: Player<ShamanSpecs>) => player.currentStatsEmitter,
		showWhen: (player: Player<ShamanSpecs>) => {
			const hasBonus = player.getCurrentStats().sets.includes('Frost Witch\'s Battlegear (4pc)');
			return hasBonus
		}
	})


	return contentBlock;
}