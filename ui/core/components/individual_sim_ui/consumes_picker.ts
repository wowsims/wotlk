import { IndividualSimUI } from "../../individual_sim_ui";
import {
	Profession,
	Spec,
	Stat,
} from "../../proto/common";

import { Component } from "../component";
import { IconEnumPicker } from "../icon_enum_picker";
import { buildIconInput } from "../icon_inputs.js";
import { relevantStatOptions } from "../inputs/stat_options";

import * as ConsumablesInputs from '../inputs/consumables';

export class ConsumesPicker extends Component {
	protected simUI: IndividualSimUI<Spec>;

	constructor(parentElem: HTMLElement, simUI: IndividualSimUI<Spec>) {
		super(parentElem, 'consumes-picker-root');
		this.simUI = simUI;

		this.buildPotionsPicker();
		this.buildFlaskPicker();
		this.buildWeaponImbuePicker();
		this.buildFoodPicker();
		this.buildPhysicalBuffPicker();
		this.buildSpellPowerBuffPicker();
		this.buildEngPicker();
		this.buildPetPicker();
	}

	private buildPotionsPicker() {
		let fragment = document.createElement('fragment');
		fragment.innerHTML = `
      <div class="consumes-row input-root input-inline">
        <label class="form-label">Potions</label>
        <div class="consumes-row-inputs">
          <div class="consumes-potions"></div>
          <div class="consumes-conjured"></div>
        </div>
      </div>
    `;

		this.rootElem.appendChild(fragment.children[0] as HTMLElement);
		const potionsElem = this.rootElem.querySelector('.consumes-potions') as HTMLElement;
		const conjuredElem = this.rootElem.querySelector('.consumes-conjured') as HTMLElement;

		const potionOptions = relevantStatOptions(ConsumablesInputs.POTIONS_CONFIG, this.simUI);
		new IconEnumPicker(
			potionsElem,
			this.simUI.player,
			ConsumablesInputs.makePotionsInput(potionOptions, 'Combat Potion')
		);

		const conjuredOptions = relevantStatOptions(ConsumablesInputs.CONJURED_CONFIG, this.simUI);
		new IconEnumPicker(conjuredElem, this.simUI.player, ConsumablesInputs.makeConjuredInput(conjuredOptions));
	}

	private buildFlaskPicker() {
		let fragment = document.createElement('fragment');
		fragment.innerHTML = `
      <div class="consumes-row input-root input-inline">
        <label class="form-label">Elixirs</label>
        <div class="consumes-row-inputs">
          <div class="consumes-flasks"></div>
        </div>
      </div>
    `;

		this.rootElem.appendChild(fragment.children[0] as HTMLElement);

		const flaskOptions = relevantStatOptions(ConsumablesInputs.FLASKS_CONFIG, this.simUI)
		const elem = this.rootElem.querySelector('.consumes-flasks') as HTMLElement;
		new IconEnumPicker(elem, this.simUI.player, ConsumablesInputs.makeFlasksInput(flaskOptions));
	}

	private buildWeaponImbuePicker() {
		let fragment = document.createElement('fragment');
		fragment.innerHTML = `
    <div class="consumes-row input-root input-inline">
        <label class="form-label">Weapon Imbues</label>
        <div class="consumes-row-inputs">
					<div class="consumes-mainhand"></div>
					<div class="consumes-offhand"></div>
				</div>
    </div>
    `;

		this.rootElem.appendChild(fragment.children[0] as HTMLElement);
		const mhElem = this.rootElem.querySelector('.consumes-mainhand') as HTMLElement;
		const ohElem = this.rootElem.querySelector('.consumes-offhand') as HTMLElement;

		const mhImbueOptions = relevantStatOptions(ConsumablesInputs.WEAPON_IMBUES_MH_CONFIG, this.simUI);
		const ohImbueOptions = relevantStatOptions(ConsumablesInputs.WEAPON_IMBUES_OH_CONFIG, this.simUI);

		new IconEnumPicker(mhElem, this.simUI.player,	ConsumablesInputs.makeMainHandImbuesInput(mhImbueOptions, 'Main Hand Imbue'));
		new IconEnumPicker(ohElem, this.simUI.player,	ConsumablesInputs.makeOffHandImbuesInput(ohImbueOptions, 'Off Hand Imbue'));
	}

	private buildFoodPicker() {
		let fragment = document.createElement('fragment');
		fragment.innerHTML = `
      <div class="consumes-row input-root input-inline">
        <label class="form-label">Food</label>
        <div class="consumes-row-inputs">
          <div class="consumes-food"></div>
        </div>
      </div>
    `;

		this.rootElem.appendChild(fragment.children[0] as HTMLElement);
		const elem = this.rootElem.querySelector('.consumes-food') as HTMLElement;

		const foodOptions = relevantStatOptions(ConsumablesInputs.FOOD_CONFIG, this.simUI)
		new IconEnumPicker(elem, this.simUI.player, ConsumablesInputs.makeFoodInput(foodOptions));
	}

	private buildPhysicalBuffPicker() {
		const includeAgi = this.simUI.individualConfig.epStats.includes(Stat.StatAgility)
		const includeStr = this.simUI.individualConfig.epStats.includes(Stat.StatStrength)

		if (!includeAgi && !includeStr) return;

		let fragment = document.createElement('fragment');
		fragment.innerHTML = `
      <div class="consumes-row input-root input-inline">
        <label class="form-label">Physical</label>
        <div class="consumes-row-inputs consumes-physical"></div>
      </div>
    `;

		this.rootElem.appendChild(fragment.children[0] as HTMLElement);
		const physicalConsumesElem = this.rootElem.querySelector('.consumes-physical') as HTMLElement;

		if (includeAgi) {
			const agilityConsumesOptions = ConsumablesInputs.AGILITY_CONSUMES_CONFIG;
			buildIconInput(physicalConsumesElem, this.simUI.player, ConsumablesInputs.makeAgilityConsumeInput(agilityConsumesOptions, 'Agility'));
		}
		if (includeStr) {
			const strengthConsumesOptions = ConsumablesInputs.STRENGTH_CONSUMES_CONFIG;
			buildIconInput(physicalConsumesElem, this.simUI.player, ConsumablesInputs.makeStrengthConsumeInput(strengthConsumesOptions, 'Strength'));
		}

		buildIconInput(physicalConsumesElem, this.simUI.player, ConsumablesInputs.BoglingRootDebuff);
	}

	private buildSpellPowerBuffPicker() {
		const spellPowerOptions = relevantStatOptions(ConsumablesInputs.SPELL_POWER_CONFIG, this.simUI);
		if (!spellPowerOptions.length) return;

		const firePowerOptions = relevantStatOptions(ConsumablesInputs.FIRE_POWER_CONFIG, this.simUI);
		const frostPowerOptions = relevantStatOptions(ConsumablesInputs.FROST_POWER_CONFIG, this.simUI);
		const shadowPowerOptions = relevantStatOptions(ConsumablesInputs.SHADOW_POWER_CONFIG, this.simUI);

		let fragment = document.createElement('fragment');
		fragment.innerHTML = `
      <div class="consumes-row input-root input-inline">
        <label class="form-label">Spells</label>
        <div class="consumes-row-inputs consumes-spells"></div>
      </div>
    `;

		this.rootElem.appendChild(fragment.children[0] as HTMLElement);
		const spellsCnsumesElem = this.rootElem.querySelector('.consumes-spells') as HTMLElement;

		buildIconInput(spellsCnsumesElem, this.simUI.player, ConsumablesInputs.makeSpellPowerConsumeInput(spellPowerOptions, 'Arcane'));
		buildIconInput(spellsCnsumesElem, this.simUI.player, ConsumablesInputs.makeFirePowerConsumeInput(firePowerOptions, 'Fire'));
		buildIconInput(spellsCnsumesElem, this.simUI.player, ConsumablesInputs.makeFrostPowerConsumeInput(frostPowerOptions, 'Frost'));
		buildIconInput(spellsCnsumesElem, this.simUI.player, ConsumablesInputs.makeshadowPowerConsumeInput(shadowPowerOptions, 'Shadow'));
	}

	private buildEngPicker() {
		let fragment = document.createElement('fragment');
		fragment.innerHTML = `
      <div class="consumes-row input-root input-inline">
        <label class="form-label">Engineering</label>
        <div class="consumes-row-inputs consumes-trade"></div>
      </div>
    `;

		const rowElem = this.rootElem.appendChild(fragment.children[0] as HTMLElement);
		const tradeConsumesElem = this.rootElem.querySelector('.consumes-trade') as HTMLElement;

		buildIconInput(tradeConsumesElem, this.simUI.player, ConsumablesInputs.Sapper);

		const explosivesOptions = relevantStatOptions(ConsumablesInputs.EXPLOSIVES_CONFIG, this.simUI);
		buildIconInput(tradeConsumesElem, this.simUI.player, ConsumablesInputs.makeExplosivesInput(explosivesOptions, 'Explosives'));

		const updateProfession = () => {
			if (this.simUI.player.hasProfession(Profession.Engineering))
				rowElem.classList.remove('hide');
			else
				rowElem.classList.add('hide');
		};
		this.simUI.player.professionChangeEmitter.on(updateProfession);
		updateProfession();
	}

	private buildPetPicker() {
		if (!this.simUI.individualConfig.petConsumeInputs?.length) return

		let fragment = document.createElement('fragment');
		fragment.innerHTML = `
			<div class="consumes-row input-root input-inline">
				<label class="form-label">Pet</label>
				<div class="consumes-row-inputs consumes-pet"></div>
			</div>
		`;

		this.rootElem.appendChild(fragment.children[0] as HTMLElement);
		const petConsumesElem = this.rootElem.querySelector('.consumes-pet') as HTMLElement;

		this.simUI.individualConfig.petConsumeInputs.map(iconInput => buildIconInput(petConsumesElem, this.simUI.player, iconInput));
	}
}
