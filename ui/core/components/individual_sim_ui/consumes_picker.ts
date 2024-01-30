import { IndividualSimUI } from "../../individual_sim_ui";
import { Player } from "../../player";
import {
	Spec,
} from "../../proto/common";
import { TypedEvent } from "../../typed_event";

import { Component } from "../component";
import { buildIconInput } from "../icon_inputs.js";
import { relevantStatOptions } from "../inputs/stat_options";
import { TypedIconEnumPickerConfig, TypedIconPickerConfig } from "../input_helpers";

import { SettingsTab } from "./settings_tab";

import * as ConsumablesInputs from '../inputs/consumables';

export class ConsumesPicker extends Component {
	protected settingsTab: SettingsTab;
	protected simUI: IndividualSimUI<Spec>;

	constructor(parentElem: HTMLElement, settingsTab: SettingsTab, simUI: IndividualSimUI<Spec>) {
		super(parentElem, 'consumes-picker-root');
		this.settingsTab = settingsTab;
		this.simUI = simUI;

		this.buildPotionsPicker();
		this.buildElixirsPicker();
		this.buildFoodPicker();
		this.buildEngPicker();
		this.buildPetPicker();
	}

	private buildPotionsPicker() {
		let fragment = document.createElement('fragment');
		fragment.innerHTML = `
			<div class="consumes-row input-root input-inline">
				<label class="form-label">Potions</label>
				<div class="consumes-row-inputs">
					<div class="consumes-prepot"></div>
          <div class="consumes-potions"></div>
          <div class="consumes-conjured"></div>
				</div>
			</div>
    `;

		const rowElem = this.rootElem.appendChild(fragment.children[0] as HTMLElement);
		const prePotsElem = this.rootElem.querySelector('.consumes-prepot') as HTMLElement;
		const potionsElem = this.rootElem.querySelector('.consumes-potions') as HTMLElement;
		const conjuredElem = this.rootElem.querySelector('.consumes-conjured') as HTMLElement;

		this.buildPickers({
			// GearChangeEmitter for ConjuredMinorRecombobulator
			changeEmitters: [],
			containerElem: rowElem,
			options: [
				{
					getConfig: () => ConsumablesInputs.makePrepopPotionsInput(
						relevantStatOptions(ConsumablesInputs.PRE_POTIONS_CONFIG, this.simUI),
						'Combat Potion',
					),
					parentElem: prePotsElem,
				},
				{
					getConfig: () => ConsumablesInputs.makePotionsInput(
						relevantStatOptions(ConsumablesInputs.POTIONS_CONFIG, this.simUI)
					),
					parentElem: potionsElem,
				},
				{
					getConfig: () => ConsumablesInputs.makeConjuredInput(
						relevantStatOptions(ConsumablesInputs.CONJURED_CONFIG, this.simUI)
					),
					parentElem: conjuredElem,
				}
			],
		})
	}

	private buildElixirsPicker() {
		let fragment = document.createElement('fragment');
		fragment.innerHTML = `
      <div class="consumes-row input-root input-inline">
        <label class="form-label">Elixirs</label>
        <div class="consumes-row-inputs">
					<div class="consumes-flasks"></div>
					<span class="elixir-space">or</span>
					<div class="consumes-battle-elixirs"></div>
					<div class="consumes-guardian-elixirs"></div>
				</div>
      </div>
    `;

		const rowElem = this.rootElem.appendChild(fragment.children[0] as HTMLElement);
		const flasksElem = this.rootElem.querySelector('.consumes-flasks') as HTMLElement;
		const battleElixirsElem = this.rootElem.querySelector('.consumes-battle-elixirs') as HTMLElement;
		const guardianElixirsElem = this.rootElem.querySelector('.consumes-guardian-elixirs') as HTMLElement;

		this.buildPickers({
			changeEmitters: [],
			containerElem: rowElem,
			options: [
				{
					getConfig: () => ConsumablesInputs.makeFlasksInput(
						relevantStatOptions(ConsumablesInputs.FLASKS_CONFIG, this.simUI)
					),
					parentElem: flasksElem,
				},
				{
					getConfig: () => ConsumablesInputs.makeBattleElixirsInput(
						relevantStatOptions(ConsumablesInputs.BATTLE_ELIXIRS_CONFIG, this.simUI)
					),
					parentElem: battleElixirsElem,
				},
				{
					getConfig: () => ConsumablesInputs.makeGuardianElixirsInput(
						relevantStatOptions(ConsumablesInputs.GUARDIAN_ELIXIRS_CONFIG, this.simUI)
					),
					parentElem: guardianElixirsElem,
				}
			],
		})
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

		const rowElem = this.rootElem.appendChild(fragment.children[0] as HTMLElement);
		const foodsElem = this.rootElem.querySelector('.consumes-food') as HTMLElement;

		this.buildPickers({
			changeEmitters: [],
			containerElem: rowElem,
			options: [
				{
					getConfig: () => ConsumablesInputs.makeFoodInput(
						relevantStatOptions(ConsumablesInputs.FOOD_CONFIG, this.simUI),
					),
					parentElem: foodsElem,
				},
			],
		})
	}

	private buildEngPicker() {
		let fragment = document.createElement('fragment');
		fragment.innerHTML = `
      <div class="consumes-row input-root input-inline">
        <label class="form-label">Engineering</label>
        <div class="consumes-row-inputs">
					<div class="consumes-sapper"></div>
					<div class="consumes-decoy"></div>
					<div class="consumes-explosives"></div>
				</div>
      </div>
    `;

		const rowElem = this.rootElem.appendChild(fragment.children[0] as HTMLElement);
		const sapperElem = this.rootElem.querySelector('.consumes-sapper') as HTMLElement;
		const decoyElem = this.rootElem.querySelector('.consumes-decoy') as HTMLElement;
		const explosivesElem = this.rootElem.querySelector('.consumes-explosives') as HTMLElement;

		this.buildPickers({
			changeEmitters: [this.simUI.player.professionChangeEmitter],
			containerElem: rowElem,
			options: [
				{
					getConfig: () => ConsumablesInputs.ThermalSapper,
					parentElem: sapperElem,
				},
				{
					getConfig: () => ConsumablesInputs.ExplosiveDecoy,
					parentElem: decoyElem,
				},
				{
					getConfig: () => ConsumablesInputs.makeExplosivesInput(
						relevantStatOptions(ConsumablesInputs.EXPLOSIVES_CONFIG, this.simUI),
						'Explosives',
					),
					parentElem: explosivesElem,
				}
			],
		})
	}

	private buildPetPicker() {
		if (this.simUI.individualConfig.petConsumeInputs?.length) {
			let fragment = document.createElement('fragment');
			fragment.innerHTML = `
        <div class="consumes-row input-root input-inline">
          <label class="form-label">Pet</label>
          <div class="consumes-row-inputs">
						<div class="consumes-pet"></div>
					</div>
        </div>
      `;

			this.rootElem.appendChild(fragment.children[0] as HTMLElement);
			const petConsumesElem = this.rootElem.querySelector('.consumes-pet') as HTMLElement;

			this.simUI.individualConfig.petConsumeInputs.map(iconInput => buildIconInput(petConsumesElem, this.simUI.player, iconInput));
		}
	}

	private buildPickers({containerElem, changeEmitters, options}: {
		containerElem: HTMLElement,
		changeEmitters: TypedEvent<any>[],
		options: {
			getConfig: () => TypedIconPickerConfig<Player<Spec>, boolean> | TypedIconEnumPickerConfig<Player<Spec>, number>,
			parentElem: HTMLElement,
		}[],
	}) {
		const buildPickers = () => {
			const anyPickersShown = options.map(optionSet => {
				optionSet.parentElem.innerHTML = '';
				const config = optionSet.getConfig();

				let isShown: boolean;
				if (config.type == 'icon') {
					isShown = !config.showWhen || config.showWhen(this.simUI.player);
				} else {
					isShown =
						(!config.showWhen || config.showWhen(this.simUI.player)) &&
						config.values.filter(value => !value.showWhen || value.showWhen(this.simUI.player)).length > 1;
				}

				if (isShown) buildIconInput(optionSet.parentElem, this.simUI.player, config);

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
