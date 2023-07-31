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
		getValue: (player: Player<ShamanSpecs>) => player.getRotation().totems || ShamanTotems.create(),
		setValue: (eventID: EventID, player: Player<ShamanSpecs>, newVal: ShamanTotems) => {
			const newRotation = player.getRotation();
			newRotation.totems = newVal;
			player.setRotation(eventID, newRotation);
		},
		changeEmitter: (player: Player<ShamanSpecs>) => player.rotationChangeEmitter,
	}, ActionId.fromSpellId(2894), "useFireElemental");

	new IconPicker(fireElementalIconContainer, simUI.player, fireElementalBooleanIconInput);

	new NumberPicker(contentBlock.bodyElement, simUI.player, {
		positive: true,
		label: "Bonus spell power",
		labelTooltip: "Bonus spell power to snapshot Fire Elemental with. Will prioritize dropping Fire Elemental if greater then 0",
		inline: true,
		getValue: (player: Player<ShamanSpecs>) => player.getRotation().totems?.bonusSpellpower || 0,
		setValue: (eventID: EventID, player: Player<ShamanSpecs>, newVal: number) => {
			const newRotation = player.getRotation();

			if (newRotation.totems) {
				newRotation.totems.bonusSpellpower = newVal
			}

			player.setRotation(eventID, newRotation);
		},
		changedEvent: (player: Player<ShamanSpecs>) => player.rotationChangeEmitter,
	})

	new BooleanPicker(contentBlock.bodyElement, simUI.player, {
		label: "Use Tier 10 (4pc)",
		labelTooltip: "Will use Tier 10 (4pc) to snapshot Fire Elemental.",
		inline: true,
		getValue: (player: Player<ShamanSpecs>) => player.getRotation().totems?.enhTierTenBonus || false,
		setValue: (eventID: EventID, player: Player<ShamanSpecs>, newVal: boolean) => {
			const newRotation = player.getRotation();

			if (newRotation.totems) {
				newRotation.totems.enhTierTenBonus = newVal
			}

			player.setRotation(eventID, newRotation);
		},
		changedEvent: (player: Player<ShamanSpecs>) => player.currentStatsEmitter,
		showWhen: (player: Player<ShamanSpecs>) => {
			const hasBonus = player.getCurrentStats().sets.includes('Frost Witch\'s Battlegear (4pc)');
			return hasBonus
		}
	})


	return contentBlock;
}