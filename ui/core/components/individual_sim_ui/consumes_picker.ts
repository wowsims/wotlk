import { IndividualSimUI } from "../../individual_sim_ui";
import { Player } from "../../player";
import {
	Spec,
	Stat,
} from "../../proto/common";
import { TypedEvent } from "../../typed_event";

import { Component } from "../component";
import { buildIconInput } from "../icon_inputs.js";
import { relevantStatOptions } from "../inputs/stat_options";

import * as ConsumablesInputs from '../inputs/consumables';
import { TypedIconEnumPickerConfig, TypedIconPickerConfig } from "../input_helpers";

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
        <div class="consumes-row-inputs consumes-potions"></div>
      </div>
    `;

		const rowElem = this.rootElem.appendChild(fragment.children[0] as HTMLElement);
		const potionsElem = this.rootElem.querySelector('.consumes-potions') as HTMLElement;

		this.buildPickers({
			changeEmitters: [this.simUI.player.levelChangeEmitter],
			containerElem: rowElem,
			options: [
				{
					getConfig: () => ConsumablesInputs.makePotionsInput(
						relevantStatOptions(ConsumablesInputs.POTIONS_CONFIG, this.simUI),
						'Combat Potion',
					),
				},
				{
					getConfig: () => ConsumablesInputs.makeConjuredInput(
						relevantStatOptions(ConsumablesInputs.CONJURED_CONFIG, this.simUI)
					),
				}
			],
			parentElem: potionsElem,
		})
	}

	private buildFlaskPicker() {
		let fragment = document.createElement('fragment');
		fragment.innerHTML = `
      <div class="consumes-row input-root input-inline">
        <label class="form-label">Elixirs</label>
        <div class="consumes-row-inputs consumes-flasks"></div>
      </div>
    `;

		const rowElem = this.rootElem.appendChild(fragment.children[0] as HTMLElement);
		const flasksElem = this.rootElem.querySelector('.consumes-flasks') as HTMLElement;

		this.buildPickers({
			changeEmitters: [this.simUI.player.levelChangeEmitter],
			containerElem: rowElem,
			options: [
				{
					getConfig: () => ConsumablesInputs.makeFlasksInput(
						relevantStatOptions(ConsumablesInputs.FLASKS_CONFIG, this.simUI)
					),
				}
			],
			parentElem: flasksElem,
		})
	}

	private buildWeaponImbuePicker() {
		let fragment = document.createElement('fragment');
		fragment.innerHTML = `
    	<div class="consumes-row input-root input-inline">
        <label class="form-label">Weapon Imbues</label>
        <div class="consumes-row-inputs consumes-weapon-imbues"></div>
    	</div>
    `;

		const rowElem = this.rootElem.appendChild(fragment.children[0] as HTMLElement);
		const imbuesElem = this.rootElem.querySelector('.consumes-weapon-imbues') as HTMLElement;

		this.buildPickers({
			changeEmitters: [this.simUI.player.levelChangeEmitter],
			containerElem: rowElem,
			options: [
				{
					getConfig: () => ConsumablesInputs.makeMainHandImbuesInput(
						relevantStatOptions(ConsumablesInputs.WEAPON_IMBUES_MH_CONFIG, this.simUI),
						'Main-Hand',
					),
				},
				{
					getConfig: () => ConsumablesInputs.makeOffHandImbuesInput(
						relevantStatOptions(ConsumablesInputs.WEAPON_IMBUES_OH_CONFIG, this.simUI),
						'Off-Hand',
					),
				},
			],
			parentElem: imbuesElem,
		})
	}

	private buildFoodPicker() {
		let fragment = document.createElement('fragment');
		fragment.innerHTML = `
      <div class="consumes-row input-root input-inline">
        <label class="form-label">Food</label>
        <div class="consumes-row-inputs consumes-food"></div>
      </div>
    `;

		const rowElem = this.rootElem.appendChild(fragment.children[0] as HTMLElement);
		const foodsElem = this.rootElem.querySelector('.consumes-food') as HTMLElement;

		this.buildPickers({
			changeEmitters: [this.simUI.player.levelChangeEmitter],
			containerElem: rowElem,
			options: [
				{
					getConfig: () => ConsumablesInputs.makeFoodInput(
						relevantStatOptions(ConsumablesInputs.FOOD_CONFIG, this.simUI),
					),
				},
			],
			parentElem: foodsElem,
		})
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

		const rowElem = this.rootElem.appendChild(fragment.children[0] as HTMLElement);
		const physicalConsumesElem = this.rootElem.querySelector('.consumes-physical') as HTMLElement;

		this.buildPickers({
			changeEmitters: [this.simUI.player.levelChangeEmitter],
			containerElem: rowElem,
			options: [
				{
					getConfig: () => ConsumablesInputs.makeAgilityConsumeInput(
						relevantStatOptions(ConsumablesInputs.AGILITY_CONSUMES_CONFIG, this.simUI),
						'Agility',
					),
				},
				{
					getConfig: () => ConsumablesInputs.makeStrengthConsumeInput(
						relevantStatOptions(ConsumablesInputs.STRENGTH_CONSUMES_CONFIG, this.simUI),
						'Strength',
					),
				},
				{
					getConfig: () => ConsumablesInputs.BoglingRootDebuff,
				},
			],
			parentElem: physicalConsumesElem,
		})
	}

	private buildSpellPowerBuffPicker() {
		let fragment = document.createElement('fragment');
		fragment.innerHTML = `
      <div class="consumes-row input-root input-inline">
        <label class="form-label">Spells</label>
        <div class="consumes-row-inputs consumes-spells"></div>
      </div>
    `;

		const rowElem = this.rootElem.appendChild(fragment.children[0] as HTMLElement);
		const spellsCnsumesElem = this.rootElem.querySelector('.consumes-spells') as HTMLElement;

		this.buildPickers({
			changeEmitters: [this.simUI.player.levelChangeEmitter],
			containerElem: rowElem,
			options: [
				{
					getConfig: () => ConsumablesInputs.makeSpellPowerConsumeInput(
						relevantStatOptions(ConsumablesInputs.SPELL_POWER_CONFIG, this.simUI),
						'Arcane',
					),
				},
				{
					getConfig: () => ConsumablesInputs.makeFirePowerConsumeInput(
						relevantStatOptions(ConsumablesInputs.FIRE_POWER_CONFIG, this.simUI),
						'Fire',
					),
				},
				{
					getConfig: () => ConsumablesInputs.makeFrostPowerConsumeInput(
						relevantStatOptions(ConsumablesInputs.FROST_POWER_CONFIG, this.simUI),
						'Frost',
					),
				},
				{
					getConfig: () => ConsumablesInputs.makeshadowPowerConsumeInput(
						relevantStatOptions(ConsumablesInputs.SHADOW_POWER_CONFIG, this.simUI),
						'Shadow',
					),
				},
			],
			parentElem: spellsCnsumesElem,
		})
	}

	private buildEngPicker() {
		let fragment = document.createElement('fragment');
		fragment.innerHTML = `
      <div class="consumes-row input-root input-inline">
        <label class="form-label">Engineering</label>
        <div class="consumes-row-inputs consumes-trade">
				</div>
      </div>
    `;

		const rowElem = this.rootElem.appendChild(fragment.children[0] as HTMLElement);
		const tradeConsumesElem = this.rootElem.querySelector('.consumes-trade') as HTMLElement;

		this.buildPickers({
			changeEmitters: [this.simUI.player.levelChangeEmitter, this.simUI.player.professionChangeEmitter],
			containerElem: rowElem,
			options: [
				{
					getConfig: () => ConsumablesInputs.Sapper,
				},
				{
					getConfig: () => ConsumablesInputs.makeExplosivesInput(
						relevantStatOptions(ConsumablesInputs.EXPLOSIVES_CONFIG, this.simUI),
						'Explosives',
					),
				}
			],
			parentElem: tradeConsumesElem,
		})
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

	private buildPickers({containerElem, changeEmitters, options, parentElem}: {
		containerElem: HTMLElement,
		changeEmitters: TypedEvent<any>[],
		options: {
			getConfig: () => TypedIconPickerConfig<Player<Spec>, boolean> | TypedIconEnumPickerConfig<Player<Spec>, number>
		}[],
		parentElem: HTMLElement,
	}) {
		const buildPickers = () => {
			parentElem.innerHTML = '';

			const anyPickersShown = options.map(optionSet => {
				const config = optionSet.getConfig();

				let isShown: boolean;
				if (config.type == 'icon') {
					isShown = !config.showWhen || config.showWhen(this.simUI.player);
				} else {
					isShown = config.values.filter(value => !value.showWhen || value.showWhen(this.simUI.player)).length > 1;
				}

				if (isShown) buildIconInput(parentElem, this.simUI.player, config);

				return isShown;
			}).filter(isShown => isShown).length > 0;

			if (anyPickersShown)
				containerElem.classList.remove('hide');
			else
				containerElem.classList.add('hide');
		};

		TypedEvent.onAny(changeEmitters).on(buildPickers)
		buildPickers()
	}
}
