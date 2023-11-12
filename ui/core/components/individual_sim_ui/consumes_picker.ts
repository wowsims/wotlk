import { IndividualSimUI } from "../../individual_sim_ui";
import {
	Flask,
	Food,
	Profession,
	Spec,
	Stat,
	WeaponBuff,
} from "../../proto/common";
import { Component } from "../component";
import { IconEnumPicker } from "../icon_enum_picker";

import * as IconInputs from '../icon_inputs.js';
import { buildIconInput } from "../icon_inputs.js";
import { SettingsTab } from "./settings_tab";

export class ConsumesPicker extends Component {
	protected settingsTab: SettingsTab;
	protected simUI: IndividualSimUI<Spec>;

	constructor(parentElem: HTMLElement, settingsTab: SettingsTab, simUI: IndividualSimUI<Spec>) {
		super(parentElem, 'consumes-picker-root');
		this.settingsTab = settingsTab;
		this.simUI = simUI;

		this.buildElixirsPicker();
		this.buildWeaponBuffPicker();
		this.buildFoodPicker();
		this.buildPhysicalBuffPicker();
		this.buildSpellPowerBuffPicker();
		this.buildEngPicker();
		this.buildPetPicker();
	}

	private buildElixirsPicker() {
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

		const flaskOptions = this.simUI.splitRelevantOptions([
			{ item: Flask.FlaskOfTheTitans, stats: [Stat.StatStamina] },
			{ item: Flask.FlaskOfDistilledWisdom, stats: [Stat.StatMP5, Stat.StatSpellPower] },
			{ item: Flask.FlaskOfSupremePower, stats: [Stat.StatMP5, Stat.StatSpellPower] },
			{ item: Flask.FlaskOfChromaticResistance, stats: [Stat.StatStamina] },
		]);
		if (flaskOptions.length) {
			const elem = this.rootElem.querySelector('.consumes-flasks') as HTMLElement;
			new IconEnumPicker(
				elem,
				this.simUI.player,
				IconInputs.makeFlasksInput(flaskOptions, 'Flask')
			);
		}
	}

	private buildWeaponBuffPicker() {
		let fragment = document.createElement('fragment');
		fragment.innerHTML = `
      <div class="consumes-row input-root input-inline">
        <label class="form-label">Weapon</label>
        <div class="consumes-row-inputs consumes-weapon"></div>
      </div>
    `;

		this.rootElem.appendChild(fragment.children[0] as HTMLElement);

		const weaponOptions = this.simUI.splitRelevantOptions([
			{ item: WeaponBuff.BrillianWizardOil, stats: [Stat.StatSpellPower] },
			{ item: WeaponBuff.BrilliantManaOil, stats: [Stat.StatHealing, Stat.StatSpellPower] },
			{ item: WeaponBuff.DenseSharpeningStone, stats: [Stat.StatAttackPower] },
			{ item: WeaponBuff.ElementalSharpeningStone, stats: [Stat.StatAttackPower] },
		]);
		if (weaponOptions.length) {
			const elem = this.rootElem.querySelector('.consumes-weapon') as HTMLElement;
			new IconEnumPicker(
				elem,
				this.simUI.player,
				IconInputs.makeWeaponBuffsInput(weaponOptions, 'WeaponBuff'),
			);
		}
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

		const foodOptions = this.simUI.splitRelevantOptions([
			{ item: Food.FoodNightfinSoup, stats: [Stat.StatMP5, Stat.StatSpellPower] },
			{ item: Food.FoodGrilledSquid, stats: [Stat.StatAgility] },
			{ item: Food.FoodSmokedDesertDumpling, stats: [Stat.StatStrength] },
			{ item: Food.FoodRunnTumTuberSurprise, stats: [Stat.StatIntellect] },
			{ item: Food.FoodDirgesKickChimaerokChops, stats: [Stat.StatStamina] },
			{ item: Food.FoodBlessSunfruit, stats: [Stat.StatStrength] },
			{ item: Food.FoodBlessedSunfruitJuice, stats: [Stat.StatSpirit] },
		]);
		if (foodOptions.length) {
			const elem = this.rootElem.querySelector('.consumes-food') as HTMLElement;
			new IconEnumPicker(elem, this.simUI.player, IconInputs.makeFoodInput(foodOptions));
		}
	}

	private buildPhysicalBuffPicker() {
		let fragment = document.createElement('fragment');
		fragment.innerHTML = `
      <div class="consumes-row input-root input-inline">
        <label class="form-label">Physical</label>
        <div class="consumes-row-inputs consumes-physical"></div>
      </div>
    `;

		this.rootElem.appendChild(fragment.children[0] as HTMLElement);
		const physicalConsumesElem = this.rootElem.querySelector('.consumes-physical') as HTMLElement;

		buildIconInput(physicalConsumesElem, this.simUI.player, IconInputs.AgilityBuffInput);
		buildIconInput(physicalConsumesElem, this.simUI.player, IconInputs.StrengthBuffInput);
	}

	private buildSpellPowerBuffPicker() {
		let fragment = document.createElement('fragment');
		fragment.innerHTML = `
      <div class="consumes-row input-root input-inline">
        <label class="form-label">Spells</label>
        <div class="consumes-row-inputs consumes-spells"></div>
      </div>
    `;

		this.rootElem.appendChild(fragment.children[0] as HTMLElement);
		const spellsCnsumesElem = this.rootElem.querySelector('.consumes-spells') as HTMLElement;

		buildIconInput(spellsCnsumesElem, this.simUI.player, IconInputs.SpellDamageBuff);
		buildIconInput(spellsCnsumesElem, this.simUI.player, IconInputs.ShadowDamageBuff);
		buildIconInput(spellsCnsumesElem, this.simUI.player, IconInputs.FireDamageBuff);
		buildIconInput(spellsCnsumesElem, this.simUI.player, IconInputs.FrostDamageBuff);
	}

	private buildEngPicker() {
		let fragment = document.createElement('fragment');
		fragment.innerHTML = `
      <div class="consumes-row input-root input-inline">
        <label class="form-label">Engineering</label>
        <div class="consumes-row-inputs consumes-trade"></div>
      </div>
    `;

		this.rootElem.appendChild(fragment.children[0] as HTMLElement);

		const tradeConsumesElem = this.rootElem.querySelector('.consumes-trade') as HTMLElement;

		// TODO Classic
		buildIconInput(tradeConsumesElem, this.simUI.player, IconInputs.Sapper);
		buildIconInput(tradeConsumesElem, this.simUI.player, IconInputs.FillerExplosiveInput);

		const updateProfession = () => {
			if (this.simUI.player.hasProfession(Profession.Engineering))
				tradeConsumesElem.parentElement!.classList.remove('hide');
			else
				tradeConsumesElem.parentElement!.classList.add('hide');
		};
		this.simUI.player.professionChangeEmitter.on(updateProfession);
		updateProfession();
	}

	private buildPetPicker() {
		if (this.simUI.individualConfig.petConsumeInputs?.length) {
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
}
